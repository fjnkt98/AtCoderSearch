package cmd

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func newDummyUpdateCmd(args []string, config *RootConfig) *cobra.Command {
	dummyFunc := func(cmd *cobra.Command, args []string) {}

	rootCmd := newRootCmd(args)
	rootCmd.AddCommand(
		newUpdateCmd(
			args,
			newUpdateProblemCmd(args, config, dummyFunc),
			newUpdateUserCmd(args, config, dummyFunc),
			newUpdateSubmissionCmd(args, config, dummyFunc),
			newUpdateLanguageCmd(args, config, dummyFunc),
		),
	)

	return rootCmd
}

func TestUpdateProblemCmd(t *testing.T) {
	type UpdateConfig struct {
		crawl    CrawlProblemConfig
		generate GenerateProblemConfig
		upload   UploadProblemConfig
		update   UpdateProblemConfig
	}
	cases := []struct {
		name string
		args []string
		want UpdateConfig
	}{
		{name: "default", args: []string{"update", "problem", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlProblemConfig{Duration: 3000, All: false}, generate: GenerateProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 2, ChunkSize: 1000}, upload: UploadProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 2, Optimize: false, Truncate: false}, update: UpdateProblemConfig{SkipFetch: false}}},
		{name: "duration", args: []string{"update", "problem", "--duration=500", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlProblemConfig{Duration: 500, All: false}, generate: GenerateProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 2, ChunkSize: 1000}, upload: UploadProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 2, Optimize: false, Truncate: false}, update: UpdateProblemConfig{SkipFetch: false}}},
		{name: "all", args: []string{"update", "problem", "--all", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlProblemConfig{Duration: 3000, All: true}, generate: GenerateProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 2, ChunkSize: 1000}, upload: UploadProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 2, Optimize: false, Truncate: false}, update: UpdateProblemConfig{SkipFetch: false}}},
		{name: "save-dir", args: []string{"update", "problem", "--save-dir=/tmp/problem", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlProblemConfig{Duration: 3000, All: false}, generate: GenerateProblemConfig{SaveDir: "/tmp/problem", Concurrent: 2, ChunkSize: 1000}, upload: UploadProblemConfig{SaveDir: "/tmp/problem", Concurrent: 2, Optimize: false, Truncate: false}, update: UpdateProblemConfig{SkipFetch: false}}},
		{name: "chunk-size", args: []string{"update", "problem", "--chunk-size=2000", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlProblemConfig{Duration: 3000, All: false}, generate: GenerateProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 2, ChunkSize: 2000}, upload: UploadProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 2, Optimize: false, Truncate: false}, update: UpdateProblemConfig{SkipFetch: false}}},
		{name: "generate-concurrent", args: []string{"update", "problem", "--generate-concurrent=1", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlProblemConfig{Duration: 3000, All: false}, generate: GenerateProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 1, ChunkSize: 1000}, upload: UploadProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 2, Optimize: false, Truncate: false}, update: UpdateProblemConfig{SkipFetch: false}}},
		{name: "optimize", args: []string{"update", "problem", "--optimize", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlProblemConfig{Duration: 3000, All: false}, generate: GenerateProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 2, ChunkSize: 1000}, upload: UploadProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 2, Optimize: true, Truncate: false}, update: UpdateProblemConfig{SkipFetch: false}}},
		{name: "truncate", args: []string{"update", "problem", "--truncate", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlProblemConfig{Duration: 3000, All: false}, generate: GenerateProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 2, ChunkSize: 1000}, upload: UploadProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 2, Optimize: false, Truncate: true}, update: UpdateProblemConfig{SkipFetch: false}}},
		{name: "upload-concurrent", args: []string{"update", "problem", "--upload-concurrent=1", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlProblemConfig{Duration: 3000, All: false}, generate: GenerateProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 2, ChunkSize: 1000}, upload: UploadProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 1, Optimize: false, Truncate: false}, update: UpdateProblemConfig{SkipFetch: false}}},
		{name: "skip-fetch", args: []string{"update", "problem", "--skip-fetch", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlProblemConfig{Duration: 3000, All: false}, generate: GenerateProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 2, ChunkSize: 1000}, upload: UploadProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 2, Optimize: false, Truncate: false}, update: UpdateProblemConfig{SkipFetch: true}}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var config RootConfig
			newDummyUpdateCmd(tt.args, &config).Execute()

			result := UpdateConfig{
				crawl:    config.Crawl.Problem,
				generate: config.Generate.Problem,
				upload:   config.Upload.Problem,
				update:   config.Update.Problem,
			}
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("expected %+v, but got %+v", tt.want, result)
			}
		})
	}
}

func TestUpdateUserCmd(t *testing.T) {
	type UpdateConfig struct {
		crawl    CrawlUserConfig
		generate GenerateUserConfig
		upload   UploadUserConfig
		update   UpdateUserConfig
	}
	cases := []struct {
		name string
		args []string
		want UpdateConfig
	}{
		{name: "default", args: []string{"update", "user", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlUserConfig{Duration: 1000}, generate: GenerateUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 2, ChunkSize: 10000}, upload: UploadUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 2, Optimize: false, Truncate: false}, update: UpdateUserConfig{SkipFetch: false}}},
		{name: "duration", args: []string{"update", "user", "--duration=500", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlUserConfig{Duration: 500}, generate: GenerateUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 2, ChunkSize: 10000}, upload: UploadUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 2, Optimize: false, Truncate: false}, update: UpdateUserConfig{SkipFetch: false}}},
		{name: "save-dir", args: []string{"update", "user", "--save-dir=/tmp/user", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlUserConfig{Duration: 1000}, generate: GenerateUserConfig{SaveDir: "/tmp/user", Concurrent: 2, ChunkSize: 10000}, upload: UploadUserConfig{SaveDir: "/tmp/user", Concurrent: 2, Optimize: false, Truncate: false}, update: UpdateUserConfig{SkipFetch: false}}},
		{name: "chunk-size", args: []string{"update", "user", "--chunk-size=2000", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlUserConfig{Duration: 1000}, generate: GenerateUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 2, ChunkSize: 2000}, upload: UploadUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 2, Optimize: false, Truncate: false}, update: UpdateUserConfig{SkipFetch: false}}},
		{name: "generate-concurrent", args: []string{"update", "user", "--generate-concurrent=1", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlUserConfig{Duration: 1000}, generate: GenerateUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 1, ChunkSize: 10000}, upload: UploadUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 2, Optimize: false, Truncate: false}, update: UpdateUserConfig{SkipFetch: false}}},
		{name: "optimize", args: []string{"update", "user", "--optimize", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlUserConfig{Duration: 1000}, generate: GenerateUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 2, ChunkSize: 10000}, upload: UploadUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 2, Optimize: true, Truncate: false}, update: UpdateUserConfig{SkipFetch: false}}},
		{name: "truncate", args: []string{"update", "user", "--truncate", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlUserConfig{Duration: 1000}, generate: GenerateUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 2, ChunkSize: 10000}, upload: UploadUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 2, Optimize: false, Truncate: true}, update: UpdateUserConfig{SkipFetch: false}}},
		{name: "upload-concurrent", args: []string{"update", "user", "--upload-concurrent=1", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlUserConfig{Duration: 1000}, generate: GenerateUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 2, ChunkSize: 10000}, upload: UploadUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 1, Optimize: false, Truncate: false}, update: UpdateUserConfig{SkipFetch: false}}},
		{name: "skip-fetch", args: []string{"update", "user", "--skip-fetch", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlUserConfig{Duration: 1000}, generate: GenerateUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 2, ChunkSize: 10000}, upload: UploadUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 2, Optimize: false, Truncate: false}, update: UpdateUserConfig{SkipFetch: true}}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var config RootConfig
			newDummyUpdateCmd(tt.args, &config).Execute()

			result := UpdateConfig{
				crawl:    config.Crawl.User,
				generate: config.Generate.User,
				upload:   config.Upload.User,
				update:   config.Update.User,
			}
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("expected %+v, but got %+v", tt.want, result)
			}
		})
	}
}

func TestUpdateSubmissionCmd(t *testing.T) {
	type UpdateConfig struct {
		crawl    CrawlSubmissionConfig
		generate GenerateSubmissionConfig
		upload   UploadSubmissionConfig
	}
	cases := []struct {
		name string
		args []string
		want UpdateConfig
	}{
		{name: "default", args: []string{"update", "submission", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlSubmissionConfig{Duration: 3000, Retry: 5, Targets: []string{"ABC", "ABC-Like", "ARC", "ARC-Like", "AGC", "AGC-Like", "JOI", "Other Sponsored", "PAST"}}, generate: GenerateSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, ChunkSize: 10000, Interval: 90, All: false}, upload: UploadSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, Optimize: false, Truncate: false}}},
		{name: "duration", args: []string{"update", "submission", "--duration=500", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlSubmissionConfig{Duration: 500, Retry: 5, Targets: []string{"ABC", "ABC-Like", "ARC", "ARC-Like", "AGC", "AGC-Like", "JOI", "Other Sponsored", "PAST"}}, generate: GenerateSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, ChunkSize: 10000, Interval: 90, All: false}, upload: UploadSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, Optimize: false, Truncate: false}}},
		{name: "retry", args: []string{"update", "submission", "--retry=1", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlSubmissionConfig{Duration: 3000, Retry: 1, Targets: []string{"ABC", "ABC-Like", "ARC", "ARC-Like", "AGC", "AGC-Like", "JOI", "Other Sponsored", "PAST"}}, generate: GenerateSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, ChunkSize: 10000, Interval: 90, All: false}, upload: UploadSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, Optimize: false, Truncate: false}}},
		{name: "target", args: []string{"update", "submission", "--target=ABC", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlSubmissionConfig{Duration: 3000, Retry: 5, Targets: []string{"ABC"}}, generate: GenerateSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, ChunkSize: 10000, Interval: 90, All: false}, upload: UploadSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, Optimize: false, Truncate: false}}},
		{name: "save-dir", args: []string{"update", "submission", "--save-dir=/tmp/submission", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlSubmissionConfig{Duration: 3000, Retry: 5, Targets: []string{"ABC", "ABC-Like", "ARC", "ARC-Like", "AGC", "AGC-Like", "JOI", "Other Sponsored", "PAST"}}, generate: GenerateSubmissionConfig{SaveDir: "/tmp/submission", Concurrent: 2, ChunkSize: 10000, Interval: 90, All: false}, upload: UploadSubmissionConfig{SaveDir: "/tmp/submission", Concurrent: 2, Optimize: false, Truncate: false}}},
		{name: "chunk-size", args: []string{"update", "submission", "--chunk-size=2000", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlSubmissionConfig{Duration: 3000, Retry: 5, Targets: []string{"ABC", "ABC-Like", "ARC", "ARC-Like", "AGC", "AGC-Like", "JOI", "Other Sponsored", "PAST"}}, generate: GenerateSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, ChunkSize: 2000, Interval: 90, All: false}, upload: UploadSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, Optimize: false, Truncate: false}}},
		{name: "generate-concurrent", args: []string{"update", "submission", "--generate-concurrent=1", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlSubmissionConfig{Duration: 3000, Retry: 5, Targets: []string{"ABC", "ABC-Like", "ARC", "ARC-Like", "AGC", "AGC-Like", "JOI", "Other Sponsored", "PAST"}}, generate: GenerateSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 1, ChunkSize: 10000, Interval: 90, All: false}, upload: UploadSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, Optimize: false, Truncate: false}}},
		{name: "interval", args: []string{"update", "submission", "--interval=30", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlSubmissionConfig{Duration: 3000, Retry: 5, Targets: []string{"ABC", "ABC-Like", "ARC", "ARC-Like", "AGC", "AGC-Like", "JOI", "Other Sponsored", "PAST"}}, generate: GenerateSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, ChunkSize: 10000, Interval: 30, All: false}, upload: UploadSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, Optimize: false, Truncate: false}}},
		{name: "all", args: []string{"update", "submission", "--all", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlSubmissionConfig{Duration: 3000, Retry: 5, Targets: []string{"ABC", "ABC-Like", "ARC", "ARC-Like", "AGC", "AGC-Like", "JOI", "Other Sponsored", "PAST"}}, generate: GenerateSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, ChunkSize: 10000, Interval: 90, All: true}, upload: UploadSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, Optimize: false, Truncate: false}}},
		{name: "optimize", args: []string{"update", "submission", "--optimize", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlSubmissionConfig{Duration: 3000, Retry: 5, Targets: []string{"ABC", "ABC-Like", "ARC", "ARC-Like", "AGC", "AGC-Like", "JOI", "Other Sponsored", "PAST"}}, generate: GenerateSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, ChunkSize: 10000, Interval: 90, All: false}, upload: UploadSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, Optimize: true, Truncate: false}}},
		{name: "truncate", args: []string{"update", "submission", "--truncate", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlSubmissionConfig{Duration: 3000, Retry: 5, Targets: []string{"ABC", "ABC-Like", "ARC", "ARC-Like", "AGC", "AGC-Like", "JOI", "Other Sponsored", "PAST"}}, generate: GenerateSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, ChunkSize: 10000, Interval: 90, All: false}, upload: UploadSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, Optimize: false, Truncate: true}}},
		{name: "upload-concurrent", args: []string{"update", "submission", "--upload-concurrent=1", "--config=config.yaml"}, want: UpdateConfig{crawl: CrawlSubmissionConfig{Duration: 3000, Retry: 5, Targets: []string{"ABC", "ABC-Like", "ARC", "ARC-Like", "AGC", "AGC-Like", "JOI", "Other Sponsored", "PAST"}}, generate: GenerateSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 2, ChunkSize: 10000, Interval: 90, All: false}, upload: UploadSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 1, Optimize: false, Truncate: false}}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var config RootConfig
			newDummyUpdateCmd(tt.args, &config).Execute()

			result := UpdateConfig{
				crawl:    config.Crawl.Submission,
				generate: config.Generate.Submission,
				upload:   config.Upload.Submission,
			}
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("expected %+v, but got %+v", tt.want, result)
			}
		})
	}
}
