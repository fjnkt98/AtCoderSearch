package cmd

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func newDummyCrawlCmd(args []string, config *RootConfig) *cobra.Command {
	dummyFunc := func(cmd *cobra.Command, args []string) {}

	rootCmd := newRootCmd(args)
	rootCmd.AddCommand(
		newCrawlCmd(
			args,
			newCrawlProblemCmd(args, config, dummyFunc),
			newCrawlUserCmd(args, config, dummyFunc),
			newCrawlSubmissionCmd(args, config, dummyFunc),
		),
	)

	return rootCmd
}

func TestCrawlProblemCmd(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want CrawlProblemConfig
	}{
		{name: "default", args: []string{"crawl", "problem", "--config=config.yaml"}, want: CrawlProblemConfig{Duration: 3000, All: false}},
		{name: "duration", args: []string{"crawl", "problem", "--duration=500", "--config=config.yaml"}, want: CrawlProblemConfig{Duration: 500, All: false}},
		{name: "all", args: []string{"crawl", "problem", "--all", "--config=config.yaml"}, want: CrawlProblemConfig{Duration: 3000, All: true}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var config RootConfig
			newDummyCrawlCmd(tt.args, &config).Execute()

			result := config.Crawl.Problem
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("expected %+v, but got %+v", tt.want, result)
			}
		})
	}
}

func TestCrawlUserCmd(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want CrawlUserConfig
	}{
		{name: "default", args: []string{"crawl", "user", "--config=config.yaml"}, want: CrawlUserConfig{Duration: 1000}},
		{name: "duration", args: []string{"crawl", "user", "--duration=500", "--config=config.yaml"}, want: CrawlUserConfig{Duration: 500}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var config RootConfig
			newDummyCrawlCmd(tt.args, &config).Execute()

			result := config.Crawl.User
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("expected %+v, but got %+v", tt.want, result)
			}
		})
	}
}

func TestCrawlSubmissionCmd(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want CrawlSubmissionConfig
	}{
		{name: "default", args: []string{"crawl", "submission", "--config=config.yaml"}, want: CrawlSubmissionConfig{Duration: 3000, Retry: 5, Targets: []string{"ABC", "ABC-Like", "ARC", "ARC-Like", "AGC", "AGC-Like", "JOI", "Other Sponsored", "PAST"}}},
		{name: "duration", args: []string{"crawl", "submission", "--duration=500", "--config=config.yaml"}, want: CrawlSubmissionConfig{Duration: 500, Retry: 5, Targets: []string{"ABC", "ABC-Like", "ARC", "ARC-Like", "AGC", "AGC-Like", "JOI", "Other Sponsored", "PAST"}}},
		{name: "retry", args: []string{"crawl", "submission", "--retry=1", "--config=config.yaml"}, want: CrawlSubmissionConfig{Duration: 3000, Retry: 1, Targets: []string{"ABC", "ABC-Like", "ARC", "ARC-Like", "AGC", "AGC-Like", "JOI", "Other Sponsored", "PAST"}}},
		{name: "target", args: []string{"crawl", "submission", "--target=ABC,Other Contests", "--config=config.yaml"}, want: CrawlSubmissionConfig{Duration: 3000, Retry: 5, Targets: []string{"ABC", "Other Contests"}}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var config RootConfig
			newDummyCrawlCmd(tt.args, &config).Execute()

			result := config.Crawl.Submission
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("expected %+v, but got %+v", tt.want, result)
			}
		})
	}
}
