package cmd

import (
	"errors"
	"fjnkt98/atcodersearch/config"
	"net/url"
	"os"
	"strconv"

	"github.com/goark/errs"
	"github.com/k0kubun/sqldef"
	"github.com/k0kubun/sqldef/database"
	"github.com/k0kubun/sqldef/database/postgres"
	"github.com/k0kubun/sqldef/parser"
	"github.com/k0kubun/sqldef/schema"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

func Migrate() error {
	u, err := url.Parse(config.Config.DataBaseURL)
	if err != nil {
		return errs.New(
			"invalid database url was given",
			errs.WithCause(err),
			errs.WithContext("database url", config.Config.DataBaseURL),
			errs.WithContext("schema", config.Config.TableSchema),
		)
	}
	password, ok := u.User.Password()
	if !ok {
		return errs.New(
			"failed to get password from the database url",
			errs.WithCause(err),
			errs.WithContext("database url", config.Config.DataBaseURL),
			errs.WithContext("schema", config.Config.TableSchema),
		)
	}

	port, err := strconv.Atoi(u.Port())
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			port = 5432
		} else {
			return errs.New(
				"failed to convert port number from given database url",
				errs.WithCause(err),
				errs.WithContext("database url", config.Config.DataBaseURL),
				errs.WithContext("schema", config.Config.TableSchema),
			)
		}
	}

	db, err := postgres.NewDatabase(database.Config{
		DbName:   u.Path[1:],
		User:     u.User.Username(),
		Password: password,
		Host:     u.Hostname(),
		Port:     port,
	})
	if err != nil {
		return errs.New(
			"failed to create a database adapter",
			errs.WithCause(err),
			errs.WithContext("database url", config.Config.DataBaseURL),
			errs.WithContext("schema", config.Config.TableSchema),
		)
	}

	sqlParser := database.NewParser(parser.ParserModePostgres)
	desiredDDLs, err := sqldef.ReadFile(config.Config.TableSchema)
	if err != nil {
		return errs.New(
			"failed to read schema file",
			errs.WithCause(err),
			errs.WithContext("database url", config.Config.DataBaseURL),
			errs.WithContext("schema", config.Config.TableSchema),
		)
	}
	options := &sqldef.Options{DesiredDDLs: desiredDDLs}
	if u.Hostname() == "localhost" {
		os.Setenv("PGSSLMODE", "disable")
	}
	sqldef.Run(schema.GeneratorModePostgres, db, sqlParser, options)

	return nil
}

func MustMigrate() {
	if err := Migrate(); err != nil {
		slog.Error("failed to migrate database schema", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("finished migrating database schema successfully.")
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database table schema",
	Long:  "Migrate database table schema",
	Run: func(cmd *cobra.Command, args []string) {
		MustMigrate()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
