package cmd

import (
	"errors"
	"fjnkt98/atcodersearch/acs"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/k0kubun/sqldef"
	"github.com/k0kubun/sqldef/database"
	"github.com/k0kubun/sqldef/database/postgres"
	"github.com/k0kubun/sqldef/parser"
	"github.com/k0kubun/sqldef/schema"
	"github.com/morikuni/failure"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

func migrate(databaseURL string, schemaFile string) error {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return failure.Translate(err, acs.InvalidURL, failure.Context{"databaseURL": databaseURL, "schema": schemaFile}, failure.Message("invalid database url was given"))
	}
	password, ok := u.User.Password()
	if !ok {
		return failure.New(acs.InvalidURL, failure.Context{"databaseURL": databaseURL, "schema": schemaFile}, failure.Message("failed to get password from the database url"))
	}

	port, err := strconv.Atoi(u.Port())
	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			port = 5432
		} else {
			return failure.Translate(err, acs.InvalidURL, failure.Context{"databaseURL": databaseURL, "schema": schemaFile}, failure.Message("failed to convert port number from given database url"))
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
		return failure.Translate(err, acs.DBError, failure.Context{"databaseURL": databaseURL, "schema": schemaFile}, failure.Message("failed to create a database adapter"))
	}

	sqlParser := database.NewParser(parser.ParserModePostgres)
	desiredDDLs, err := sqldef.ReadFile(schemaFile)
	if err != nil {
		return failure.Translate(err, acs.FileOperationError, failure.Context{"databaseURL": databaseURL, "schema": schemaFile}, failure.Message("failed to read schema file"))
	}
	options := &sqldef.Options{DesiredDDLs: desiredDDLs}
	if u.Hostname() == "localhost" {
		os.Setenv("PGSSLMODE", "disable")
	}
	sqldef.Run(schema.GeneratorModePostgres, db, sqlParser, options)

	return nil
}

func DoMigrate() {
	url := os.Getenv("DATABASE_URL")
	schema := os.Getenv("DB_SCHEMA_FILE")
	if err := migrate(url, schema); err != nil {
		slog.Error("failed to migrate database schema", slog.String("error", fmt.Sprintf("%+v", err)))
		os.Exit(1)
	}
	slog.Info("finished migrating database schema successfully.")
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database table schema",
	Long:  "Migrate database table schema",
	Run: func(cmd *cobra.Command, args []string) {
		url := os.Getenv("DATABASE_URL")
		schema := os.Getenv("DB_SCHEMA_FILE")
		if err := migrate(url, schema); err != nil {
			slog.Error("failed to migrate database schema", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
		slog.Info("finished migrating database schema successfully.")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
