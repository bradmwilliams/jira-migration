package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	jiraBaseClient "github.com/andygrunwald/go-jira"
	bigquery2 "github.com/bradmwilliams/jira-migration/pkg/bigquery"
	"github.com/openshift/ci-search/jira"
	"github.com/openshift/ci-search/pkg/bigquery"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
	"os"
	"sigs.k8s.io/prow/prow/flagutil"
	"time"
)

type Options struct {
	DryRun bool

	jira flagutil.JiraOptions

	// BigQuery Options
	GoogleProjectID                    string
	GoogleServiceAccountCredentialFile string
	BigQueryRefreshInterval            time.Duration
}

func main() {
	original := flag.CommandLine
	klog.InitFlags(original)
	original.Set("alsologtostderr", "true")
	original.Set("v", "2")

	opt := &Options{
		BigQueryRefreshInterval: 1 * time.Minute,
	}
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, arguments []string) {
			err := opt.Validate(context.TODO())
			if err != nil {
				klog.Exitf("error: %v", err)
			}

			if err := opt.Run(); err != nil {
				klog.Exitf("error: %v", err)
			}
		},
	}
	flagset := cmd.Flags()

	goFlagSet := flag.NewFlagSet("prowflags", flag.ContinueOnError)
	opt.jira.AddFlags(goFlagSet)
	flagset.AddGoFlagSet(goFlagSet)

	flagset.AddGoFlag(original.Lookup("v"))

	opt.AddFlags(flagset)

	if err := cmd.Execute(); err != nil {
		klog.Exitf("error: %v", err)
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.GoogleProjectID, "google-project-id", os.Getenv("GOOGLE_PROJECT_ID"), "Google project name.")
	fs.StringVar(&o.GoogleServiceAccountCredentialFile, "google-service-account-credential-file", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"), "location of a credential file described by https://cloud.google.com/docs/authentication/production")
	fs.DurationVar(&o.BigQueryRefreshInterval, "bigquery-refresh-interval", o.BigQueryRefreshInterval, "How often to push comments into BigQuery. Defaults to 15 minutes.")
	fs.BoolVar(&o.DryRun, "dry-run", o.DryRun, "Perform no actions.")
}

func (o *Options) Validate(ctx context.Context) error {
	if len(o.GoogleProjectID) == 0 {
		return errors.New("--google-project-id flag must be set")
	}
	if len(o.GoogleServiceAccountCredentialFile) == 0 {
		return errors.New("--google-service-account-credential-file flag must be set")
	}
	return nil
}

func (o *Options) Run() error {
	c, err := o.jira.Client()
	if err != nil {
		klog.Fatalf("Unable to create jira client: %w", err)
	}

	issue, err := c.GetIssue("OCPBUGS-35865")
	if err != nil {
		klog.Errorf("Unable to get jira issue: %w", err)
		return err
	}

	b, err := json.MarshalIndent(issue, "", "    ")
	if err != nil {
		klog.Errorf("unable to marshal Jira Issue: %v", err)
		return nil
	}
	klog.V(2).Infof("Retrieved issue:\n%s", string(b))

	bqc, err := bigquery.NewBigQueryClient(o.GoogleProjectID, o.GoogleServiceAccountCredentialFile)
	if err != nil {
		klog.Fatalf("Unable to configure bigquery client: %v", err)
	}

	var tickets []*bigquery2.Ticket
	timestamp := time.Now()

	updated := jira.NewIssueComments(issue.ID, issue.Fields.Comments)
	updated.Info = jiraBaseClient.Issue{
		ID:     issue.ID,
		Key:    issue.Key,
		Fields: issue.Fields,
	}
	updated.RefreshTime = timestamp

	tickets = append(tickets, bigquery2.ConvertToTicket(updated, timestamp))

	b, err = json.MarshalIndent(tickets, "", "    ")
	if err != nil {
		klog.Errorf("unable to marshal tickets: %v", err)
		return nil
	}
	klog.V(2).Infof("Tickets:\n%s", string(b))

	if len(tickets) > 0 {
		if o.DryRun {
			klog.Infof("[Dry Run] Syncing %d issues to bigquery", len(tickets))
		} else {
			klog.V(5).Infof("Syncing %d issues to bigquery", len(tickets))
			err := bqc.WriteRows(context.TODO(), bigquery2.BigqueryDatasetId, bigquery2.BigqueryTableId, tickets)
			if err != nil {
				return fmt.Errorf("unable to write to bigquery: %v", err)
			}
		}
	}

	return nil
}
