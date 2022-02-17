package main

import (
	"flag"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"k8s.io/test-infra/prow/flagutil"
)

type options struct {
	jira flagutil.JiraOptions
}

func main() {
	original := flag.CommandLine
	klog.InitFlags(original)
	original.Set("alsologtostderr", "true")
	original.Set("v", "2")

	opt := &options{}
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, arguments []string) {
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

	if err := cmd.Execute(); err != nil {
		klog.Exitf("error: %v", err)
	}
}

func (o *options) Run() error {
	c, err := o.jira.Client()
	if err != nil {
		klog.Fatalf("Unable to create jira client: %w", err)
	}

	issue, err := c.GetIssue("OCPBUGS-10")
	if err != nil {
		klog.Errorf("Unable to get jira issue: %w", err)
		return err
	}

	klog.V(2).Infof("Found issue ID: %s", issue.ID)
	klog.V(2).Infof("Found issue Key: %s", issue.Key)
	klog.V(2).Infof("Found issue Self: %s", issue.Self)
	klog.V(2).Infof("Found issue Expand: %s", issue.Expand)

	projects, err := c.ListProjects()
	if err != nil {
		klog.Errorf("Unable to get jira projects: %w", err)
		return err
	}

	for _, project := range *projects {
		klog.V(2).Infof("Found project: [%s] %s", project.Key, project.Name)
	}
	return nil
}
