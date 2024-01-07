package cmd

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func newDummyUploadCmd(args []string, config *RootConfig) *cobra.Command {
	dummyFunc := func(cmd *cobra.Command, args []string) {}

	rootCmd := newRootCmd(args)
	rootCmd.AddCommand(
		newUploadCmd(
			args,
			newUploadProblemCmd(args, config, dummyFunc),
			newUploadUserCmd(args, config, dummyFunc),
			newUploadSubmissionCmd(args, config, dummyFunc),
		),
	)

	return rootCmd
}

func TestUploadProblemCmd(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want UploadProblemConfig
	}{
		{name: "default", args: []string{"upload", "problem", "--config=config.yaml"}, want: UploadProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 4, Optimize: false, Truncate: false}},
		{name: "save_dir", args: []string{"upload", "problem", "--save-dir=/tmp/problem", "--config=config.yaml"}, want: UploadProblemConfig{SaveDir: "/tmp/problem", Concurrent: 4, Optimize: false, Truncate: false}},
		{name: "concurrent", args: []string{"upload", "problem", "--concurrent=1", "--config=config.yaml"}, want: UploadProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 1, Optimize: false, Truncate: false}},
		{name: "optimize", args: []string{"upload", "problem", "--optimize", "--config=config.yaml"}, want: UploadProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 4, Optimize: true, Truncate: false}},
		{name: "truncate", args: []string{"upload", "problem", "--truncate", "--config=config.yaml"}, want: UploadProblemConfig{SaveDir: "/var/tmp/atcoder/problem", Concurrent: 4, Optimize: false, Truncate: true}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var config RootConfig
			newDummyUploadCmd(tt.args, &config).Execute()

			result := config.Upload.Problem
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("expected %+v, but got %+v", tt.want, result)
			}
		})
	}
}

func TestUploadUserCmd(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want UploadUserConfig
	}{
		{name: "default", args: []string{"upload", "user", "--config=config.yaml"}, want: UploadUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 4, Optimize: false, Truncate: false}},
		{name: "save_dir", args: []string{"upload", "user", "--save-dir=/tmp/user", "--config=config.yaml"}, want: UploadUserConfig{SaveDir: "/tmp/user", Concurrent: 4, Optimize: false, Truncate: false}},
		{name: "concurrent", args: []string{"upload", "user", "--concurrent=1", "--config=config.yaml"}, want: UploadUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 1, Optimize: false, Truncate: false}},
		{name: "optimize", args: []string{"upload", "user", "--optimize", "--config=config.yaml"}, want: UploadUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 4, Optimize: true, Truncate: false}},
		{name: "truncate", args: []string{"upload", "user", "--truncate", "--config=config.yaml"}, want: UploadUserConfig{SaveDir: "/var/tmp/atcoder/user", Concurrent: 4, Optimize: false, Truncate: true}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var config RootConfig
			newDummyUploadCmd(tt.args, &config).Execute()

			result := config.Upload.User
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("expected %+v, but got %+v", tt.want, result)
			}
		})
	}
}

func TestUploadSubmissionCmd(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want UploadSubmissionConfig
	}{
		{name: "default", args: []string{"upload", "submission", "--config=config.yaml"}, want: UploadSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 4, Optimize: false, Truncate: false}},
		{name: "save_dir", args: []string{"upload", "submission", "--save-dir=/tmp/submission", "--config=config.yaml"}, want: UploadSubmissionConfig{SaveDir: "/tmp/submission", Concurrent: 4, Optimize: false, Truncate: false}},
		{name: "concurrent", args: []string{"upload", "submission", "--concurrent=1", "--config=config.yaml"}, want: UploadSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 1, Optimize: false, Truncate: false}},
		{name: "optimize", args: []string{"upload", "submission", "--optimize", "--config=config.yaml"}, want: UploadSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 4, Optimize: true, Truncate: false}},
		{name: "truncate", args: []string{"upload", "submission", "--truncate", "--config=config.yaml"}, want: UploadSubmissionConfig{SaveDir: "/var/tmp/atcoder/submission", Concurrent: 4, Optimize: false, Truncate: true}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var config RootConfig
			newDummyUploadCmd(tt.args, &config).Execute()

			result := config.Upload.Submission
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("expected %+v, but got %+v", tt.want, result)
			}
		})
	}
}
