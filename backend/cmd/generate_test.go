package cmd

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func newDummyGenerateCmd(args []string, config *RootConfig) *cobra.Command {
	dummyFunc := func(cmd *cobra.Command, args []string) {}

	rootCmd := newRootCmd(args)
	rootCmd.AddCommand(
		newGenerateCmd(
			args,
			newGenerateProblemCmd(args, config, dummyFunc),
			newGenerateUserCmd(args, config, dummyFunc),
			newGenerateSubmissionCmd(args, config, dummyFunc),
		),
	)

	return rootCmd
}

func TestGenerateProblemCmd(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want GenerateProblemConfig
	}{
		{name: "default", args: []string{"generate", "problem", "--config=config.yaml"}, want: GenerateProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 6, ChunkSize: 1000}},
		{name: "save_dir", args: []string{"generate", "problem", "--save-dir=/tmp/problem", "--config=config.yaml"}, want: GenerateProblemConfig{SaveDir: "/tmp/problem", Concurrent: 6, ChunkSize: 1000}},
		{name: "concurrent", args: []string{"generate", "problem", "--concurrent=1", "--config=config.yaml"}, want: GenerateProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 1, ChunkSize: 1000}},
		{name: "chunk_size", args: []string{"generate", "problem", "--chunk-size=2000", "--config=config.yaml"}, want: GenerateProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 6, ChunkSize: 2000}}}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var config RootConfig
			newDummyGenerateCmd(tt.args, &config).Execute()

			result := config.Generate.Problem
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("expected %+v, but got %+v", tt.want, result)
			}
		})
	}
}

func TestGenerateUserCmd(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want GenerateUserConfig
	}{
		{name: "default", args: []string{"generate", "user", "--config=config.yaml"}, want: GenerateUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 6, ChunkSize: 10000}},
		{name: "save_dir", args: []string{"generate", "user", "--save-dir=/tmp/user", "--config=config.yaml"}, want: GenerateUserConfig{SaveDir: "/tmp/user", Concurrent: 6, ChunkSize: 10000}},
		{name: "concurrent", args: []string{"generate", "user", "--concurrent=1", "--config=config.yaml"}, want: GenerateUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 1, ChunkSize: 10000}},
		{name: "chunk_size", args: []string{"generate", "user", "--chunk-size=2000", "--config=config.yaml"}, want: GenerateUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 6, ChunkSize: 2000}}}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var config RootConfig
			newDummyGenerateCmd(tt.args, &config).Execute()

			result := config.Generate.User
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("expected %+v, but got %+v", tt.want, result)
			}
		})
	}
}

func TestGenerateSubmissionCmd(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want GenerateSubmissionConfig
	}{
		{name: "default", args: []string{"generate", "submission", "--config=config.yaml"}, want: GenerateSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 6, ChunkSize: 10000, Interval: 90, All: false}},
		{name: "save_dir", args: []string{"generate", "submission", "--save-dir=/tmp/submission", "--config=config.yaml"}, want: GenerateSubmissionConfig{SaveDir: "/tmp/submission", Concurrent: 6, ChunkSize: 10000, Interval: 90, All: false}},
		{name: "concurrent", args: []string{"generate", "submission", "--concurrent=1", "--config=config.yaml"}, want: GenerateSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 1, ChunkSize: 10000, Interval: 90, All: false}},
		{name: "chunk_size", args: []string{"generate", "submission", "--chunk-size=2000", "--config=config.yaml"}, want: GenerateSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 6, ChunkSize: 2000, Interval: 90, All: false}},
		{name: "interval", args: []string{"generate", "submission", "--interval=30", "--config=config.yaml"}, want: GenerateSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 6, ChunkSize: 10000, Interval: 30, All: false}},
		{name: "all", args: []string{"generate", "submission", "--all", "--config=config.yaml"}, want: GenerateSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 6, ChunkSize: 10000, Interval: 90, All: true}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var config RootConfig
			newDummyGenerateCmd(tt.args, &config).Execute()

			result := config.Generate.Submission
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("expected %+v, but got %+v", tt.want, result)
			}
		})
	}
}
