package cmd

import (
	"fmt"
	"strings"

	"log/slog"

	"github.com/goark/errs"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func LoadConfig(file string, config *RootConfig) error {
	godotenv.Load()

	viper.SetConfigFile(file)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return errs.New(
			"failed to read config file",
			errs.WithCause(err),
			errs.WithContext("file", file),
		)
	}

	if err := viper.Unmarshal(config); err != nil {
		return errs.New(
			"failed to unmarshal config file",
			errs.WithCause(err),
			errs.WithContext("file", file),
		)
	}
	return nil
}

func MustLoadConfig(file string, config *RootConfig) {
	if err := LoadConfig(file, config); err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		panic("failed to load config.")
	} else {
		slog.Info(fmt.Sprintf("using config file: %s", viper.ConfigFileUsed()))
	}
}

func MustLoadConfigFromFlags(flags *pflag.FlagSet, config *RootConfig) {
	file, err := flags.GetString("config")
	if err != nil {
		panic("failed to get config file path from flags")
	}

	if err := LoadConfig(file, config); err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		panic("failed to load config.")
	} else {
		slog.Info(fmt.Sprintf("using config file: %s", viper.ConfigFileUsed()))
	}
}

type RootConfig struct {
	CommonConfig `mapstructure:",squash"`
	Crawl        CrawlConfig    `mapstructure:"crawl" json:"crawl"`
	Generate     GenerateConfig `mapstructure:"generate" json:"generate"`
	Upload       UploadConfig   `mapstructure:"upload" json:"upload"`
	Update       UpdateConfig   `mapstructure:"update" json:"update"`
}

type CommonConfig struct {
	DataBaseURL        string   `mapstructure:"database_url" json:"database_url"`
	TableSchema        string   `mapstructure:"table_schema" json:"table_schema"`
	SolrHost           string   `mapstructure:"solr_host" json:"solr_host"`
	AtCoderUserName    string   `mapstructure:"atcoder_username" json:"atcoder_username"`
	AtCoderPassword    string   `mapstructure:"atcoder_password" json:"atcoder_password"`
	TrustedProxies     []string `mapstructure:"trusted_proxies" json:"trusted_proxies"`
	DoMigrate          bool     `mapstructure:"do_migrate" json:"do_migrate"`
	ProblemCoreName    string   `mapstructure:"problem_core_name" json:"problem_core_name"`
	UserCoreName       string   `mapstructure:"user_core_name" json:"user_core_name"`
	SubmissionCoreName string   `mapstructure:"submission_core_name" json:"submission_core_name"`
}

type CrawlConfig struct {
	Problem    CrawlProblemConfig    `mapstructure:"problem" json:"problem"`
	User       CrawlUserConfig       `mapstructure:"user" json:"user"`
	Submission CrawlSubmissionConfig `mapstructure:"submission" json:"submission"`
}

type CrawlProblemConfig struct {
	Duration int  `mapstructure:"duration" json:"duration"`
	All      bool `mapstructure:"all" json:"all"`
}

type CrawlUserConfig struct {
	Duration int `mapstructure:"duration" json:"duration"`
}

type CrawlSubmissionConfig struct {
	Duration int      `mapstructure:"duration" json:"duration"`
	Retry    int      `mapstructure:"retry" json:"retry"`
	Targets  []string `mapstructure:"targets" json:"targets"`
}

type GenerateConfig struct {
	Problem    GenerateProblemConfig    `mapstructure:"problem" json:"problem"`
	User       GenerateUserConfig       `mapstructure:"user" json:"user"`
	Submission GenerateSubmissionConfig `mapstructure:"submission" json:"submission"`
}

type GenerateProblemConfig struct {
	SaveDir    string `mapstructure:"save_dir" json:"save_dir"`
	Concurrent int    `mapstructure:"concurrent" json:"concurrent"`
	ChunkSize  int    `mapstructure:"chunk_size" json:"chunk_size"`
}

type GenerateUserConfig struct {
	SaveDir    string `mapstructure:"save_dir" json:"save_dir"`
	Concurrent int    `mapstructure:"concurrent" json:"concurrent"`
	ChunkSize  int    `mapstructure:"chunk_size" json:"chunk_size"`
}

type GenerateSubmissionConfig struct {
	SaveDir    string `mapstructure:"save_dir" json:"save_dir"`
	Concurrent int    `mapstructure:"concurrent" json:"concurrent"`
	ChunkSize  int    `mapstructure:"chunk_size" json:"chunk_size"`
	Interval   int    `mapstructure:"interval" json:"interval"`
	All        bool   `mapstructure:"all" json:"all"`
}

type UploadConfig struct {
	Problem    UploadProblemConfig    `mapstructure:"problem" json:"problem"`
	User       UploadUserConfig       `mapstructure:"user" json:"user"`
	Submission UploadSubmissionConfig `mapstructure:"submission" json:"submission"`
}

type UploadProblemConfig struct {
	SaveDir    string `mapstructure:"save_dir" json:"save_dir"`
	Concurrent int    `mapstructure:"concurrent" json:"concurrent"`
	Optimize   bool   `mapstructure:"optimize" json:"optimize"`
	Truncate   bool   `mapstructure:"truncate" json:"truncate"`
}

type UploadUserConfig struct {
	SaveDir    string `mapstructure:"save_dir" json:"save_dir"`
	Concurrent int    `mapstructure:"concurrent" json:"concurrent"`
	Optimize   bool   `mapstructure:"optimize" json:"optimize"`
	Truncate   bool   `mapstructure:"truncate" json:"truncate"`
}

type UploadSubmissionConfig struct {
	SaveDir    string `mapstructure:"save_dir" json:"save_dir"`
	Concurrent int    `mapstructure:"concurrent" json:"concurrent"`
	Optimize   bool   `mapstructure:"optimize" json:"optimize"`
	Truncate   bool   `mapstructure:"truncate" json:"truncate"`
}

type UpdateConfig struct {
	Problem UpdateProblemConfig `mapstructure:"problem" json:"problem"`
	User    UpdateUserConfig    `mapstructure:"user" json:"user"`
}

type UpdateProblemConfig struct {
	SkipFetch bool `mapstructure:"skip_fetch" json:"skip_fetch"`
}

type UpdateUserConfig struct {
	SkipFetch bool `mapstructure:"skip_fetch" json:"skip_fetch"`
}
