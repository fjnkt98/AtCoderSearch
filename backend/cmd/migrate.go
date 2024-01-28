package cmd

import (
	"errors"
	"net/url"
	"os"
	"strconv"

	"log/slog"

	"github.com/goark/errs"
	"github.com/k0kubun/sqldef"
	"github.com/k0kubun/sqldef/database"
	"github.com/k0kubun/sqldef/database/postgres"
	"github.com/k0kubun/sqldef/parser"
	"github.com/k0kubun/sqldef/schema"
	"github.com/spf13/cobra"
)

func Migrate(dsn, schemaFile string) error {
	u, err := url.Parse(dsn)
	if err != nil {
		return errs.New(
			"invalid database url was given",
			errs.WithCause(err),
			errs.WithContext("database url", dsn),
			errs.WithContext("schema", schemaFile),
		)
	}
	password, ok := u.User.Password()
	if !ok {
		return errs.New(
			"failed to get password from the database url",
			errs.WithCause(err),
			errs.WithContext("database url", dsn),
			errs.WithContext("schema", schemaFile),
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
				errs.WithContext("database url", dsn),
				errs.WithContext("schema", schemaFile),
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
			errs.WithContext("database url", dsn),
			errs.WithContext("schema", schemaFile),
		)
	}

	sqlParser := database.NewParser(parser.ParserModePostgres)
	desiredDDLs, err := sqldef.ReadFile(schemaFile)
	if err != nil {
		return errs.New(
			"failed to read schema file",
			errs.WithCause(err),
			errs.WithContext("database url", dsn),
			errs.WithContext("schema", schemaFile),
		)
	}
	options := &sqldef.Options{DesiredDDLs: desiredDDLs}
	if u.Hostname() == "localhost" {
		os.Setenv("PGSSLMODE", "disable")
	}
	sqldef.Run(schema.GeneratorModePostgres, db, sqlParser, options)

	return nil
}

func MustMigrate(dsn, schemaFile string) {
	if err := Migrate(dsn, schemaFile); err != nil {
		slog.Error("failed to migrate database schema", slog.Any("error", err))
		panic("failed to migrate database schema")
	}
	slog.Info("finished migrating database schema successfully.")
}

func newMigrateCmd(args []string, config *RootConfig, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database table schema",
		Long:  "Migrate database table schema",
		PreRun: func(cmd *cobra.Command, args []string) {
			MustLoadConfig(cmd.Flags(), config)
		},
		Run: func(cmd *cobra.Command, args []string) {
			MustMigrate(config.DataBaseURL, config.TableSchema)
		},
	}

	migrateCmd.SetArgs(args)
	if runFunc != nil {
		migrateCmd.Run = runFunc
	}

	return migrateCmd
}
