package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/openshift/ci-search/jira"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"sigs.k8s.io/prow/prow/flagutil"
)

type options struct {
	jira       flagutil.JiraOptions
	JiraSearch string
}

func main() {
	original := flag.CommandLine
	klog.InitFlags(original)
	original.Set("alsologtostderr", "true")
	original.Set("v", "2")

	opt := &options{
		JiraSearch: "project=OCPBUGS&created>='-14d'&status!='CLOSED'&affectedVersion IN (4.14,4.13,4.12,4.11,4.10,4.9,4.8,4.7,4.6,4.5,4.4,4.3,4.2)",
	}

	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, arguments []string) {
			if err := opt.Run(); err != nil {
				klog.Fatalf("error: %v", err)
			}
		},
	}
	flagset := cmd.Flags()

	goFlagSet := flag.NewFlagSet("prowflags", flag.ContinueOnError)
	opt.jira.AddFlags(goFlagSet)
	flagset.AddGoFlagSet(goFlagSet)

	flagset.AddGoFlag(original.Lookup("v"))

	flagset.StringVar(&opt.JiraSearch, "jira-search", opt.JiraSearch, "A JQL query to search for issues to index.")

	if err := cmd.Execute(); err != nil {
		klog.Exitf("error: %v", err)
	}
}

func (o *options) Run() error {
	err := o.jira.Validate(true)
	if err != nil {
		klog.Fatalf("Invalid Jira options specified: %w", err)
		return err
	}

	jc, err := o.jira.Client()
	if err != nil {
		klog.Fatalf("Unable to create jira client: %w", err)
	}

	c := &jira.Client{
		Client: jc,
	}

	ctx := context.Background()

	issues, err := c.SearchIssues(ctx, jira.SearchIssuesArgs{
		Jql: o.JiraSearch,
	})
	if err != nil {
		return err
	}

	for _, issue := range issues {
		fmt.Println("Found issue: ", issue.ID)
	}
	return nil
}
