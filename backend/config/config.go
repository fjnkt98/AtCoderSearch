package config

import (
	"strings"

	"github.com/goark/errs"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var Config RootConfig

func Load(configFile string) error {
	godotenv.Load()

	viper.SetConfigFile(configFile)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return errs.New(
			"failed to read config file",
			errs.WithCause(err),
			errs.WithContext("file", configFile),
		)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		return errs.New(
			"failed to unmarshal config file",
			errs.WithCause(err),
			errs.WithContext("file", configFile),
		)
	}
	return nil
}

type RootConfig struct {
	CommonConfig `mapstructure:",squash"`
	Problem      ProblemConfig    `mapstructure:"problem" json:"problem"`
	User         UserConfig       `mapstructure:"user" json:"user"`
	Submission   SubmissionConfig `mapstructure:"submission" json:"submission"`
	Recommend    RecommendConfig  `mapstructure:"recommend" json:"recommend"`
}

type CommonConfig struct {
	DataBaseURL     string `mapstructure:"database_url" json:"database_url"`
	TableSchema     string `mapstructure:"table_schema" json:"table_schema"`
	SolrHost        string `mapstructure:"solr_host" json:"solr_host"`
	AtCoderUserName string `mapstructure:"atcoder_username" json:"atcoder_username"`
	AtCoderPassword string `mapstructure:"atcoder_password" json:"atcoder_password"`
	DoMigrate       bool   `mapstructure:"do_migrate" json:"do_migrate"`
}

type ProblemConfig struct {
	CoreName string                `mapstructure:"core_name" json:"core_name"`
	Crawl    CrawlProblemConfig    `mapstructure:"crawl" json:"crawl"`
	Generate GenerateProblemConfig `mapstructure:"generate" json:"generate"`
	Upload   UploadProblemConfig   `mapstructure:"upload" json:"upload"`
	Update   UpdateProblemConfig   `mapstructure:"update" json:"update"`
}

type UserConfig struct {
	CoreName string             `mapstructure:"core_name" json:"core_name"`
	Crawl    CrawlUserConfig    `mapstructure:"crawl" json:"crawl"`
	Generate GenerateUserConfig `mapstructure:"generate" json:"generate"`
	Upload   UploadUserConfig   `mapstructure:"upload" json:"upload"`
	Update   UpdateUserConfig   `mapstructure:"update" json:"update"`
}

type SubmissionConfig struct {
	CoreName string                   `mapstructure:"core_name" json:"core_name"`
	Crawl    CrawlSubmissionConfig    `mapstructure:"crawl" json:"crawl"`
	Generate GenerateSubmissionConfig `mapstructure:"generate" json:"generate"`
	Read     ReadSubmissionConfig     `mapstructure:"read" json:"read"`
	Upload   UploadSubmissionConfig   `mapstructure:"upload" json:"upload"`
	Update   UpdateSubmissionConfig   `mapstructure:"update" json:"update"`
}

type RecommendConfig struct {
}

type CrawlCommonConfig struct {
	Duration int `mapstructure:"duration" json:"duration"`
}

type CrawlProblemConfig struct {
	CrawlCommonConfig `mapstructure:",squash"`
	All               bool `mapstructure:"all" json:"all"`
}

type CrawlUserConfig struct {
	CrawlCommonConfig `mapstructure:",squash"`
}

type CrawlSubmissionConfig struct {
	CrawlCommonConfig `mapstructure:",squash"`
	Retry             int    `mapstructure:"retry" json:"retry"`
	Targets           string `mapstructure:"targets" json:"targets"`
}

type GenerateCommonConfig struct {
	SaveDir    string `mapstructure:"save_dir" json:"save_dir"`
	Concurrent int    `mapstructure:"concurrent" json:"concurrent"`
	ChunkSize  int    `mapstructure:"chunk_size" json:"chunk_size"`
}

type GenerateUserConfig struct {
	GenerateCommonConfig `mapstructure:",squash"`
}

type GenerateSubmissionConfig struct {
	GenerateCommonConfig `mapstructure:",squash"`
}

type ReadSubmissionConfig struct {
	Interval int  `mapstructure:"interval" json:"interval"`
	All      bool `mapstructure:"all" json:"all"`
}

type GenerateRecommendConfig struct {
	GenerateCommonConfig `mapstructure:",squash"`
}

type GenerateProblemConfig struct {
	GenerateCommonConfig `mapstructure:",squash"`
}

type UploadCommonConfig struct {
	SaveDir    string `mapstructure:"save_dir" json:"save_dir"`
	Concurrent int    `mapstructure:"concurrent" json:"concurrent"`
	Optimize   bool   `mapstructure:"optimize" json:"optimize"`
	Truncate   bool   `mapstructure:"truncate" json:"truncate"`
}

type UploadProblemConfig struct {
	UploadCommonConfig `mapstructure:",squash"`
}

type UploadUserConfig struct {
	UploadCommonConfig `mapstructure:",squash"`
}

type UploadSubmissionConfig struct {
	UploadCommonConfig `mapstructure:",squash"`
}

type UploadRecommendConfig struct {
	UploadCommonConfig `mapstructure:",squash"`
}

type UpdateProblemConfig struct {
	SkipFetch bool `mapstructure:"skip_fetch" json:"skip_fetch"`
}

type UpdateUserConfig struct {
	SkipFetch bool `mapstructure:"skip_fetch" json:"skip_fetch"`
}

type UpdateSubmissionConfig struct {
}
