package repository

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/goark/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thanhpk/randstr"
)

func Columns(model any) []string {
	ty := reflect.TypeOf(model)
	if ty.Kind() != reflect.Pointer {
		return nil
	}

	ty = ty.Elem()
	if ty.Kind() != reflect.Struct {
		return nil
	}

	columns := make([]string, 0, ty.NumField())
	for i := 0; i < ty.NumField(); i++ {
		f := ty.Field(i)

		if f.Name == "UpdatedAt" {
			continue
		}

		var column string
		if tag, ok := f.Tag.Lookup("db"); ok {
			column = tag
		} else {
			column = f.Name
		}
		columns = append(columns, column)
	}
	return columns
}

func UniqueKey(model any) []string {
	ty := reflect.TypeOf(model)
	if ty.Kind() != reflect.Pointer {
		return nil
	}

	ty = ty.Elem()
	if ty.Kind() != reflect.Struct {
		return nil
	}

	uniqueKeys := make([]string, 0)
	for i := 0; i < ty.NumField(); i++ {
		f := ty.Field(i)
		var key string
		if tag, ok := f.Tag.Lookup("bulk"); ok && tag == "unique" {
			if tag, ok := f.Tag.Lookup("db"); ok {
				key = tag
			} else {
				key = f.Name
			}
			uniqueKeys = append(uniqueKeys, key)
		}
	}
	return uniqueKeys
}

func Rows(data any) [][]any {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return nil
	}

	rows := make([][]any, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		d := v.Index(i)
		ty := d.Type()
		row := make([]any, 0, d.NumField())
		for j := 0; j < d.NumField(); j++ {
			value := d.Field(j)
			if ty.Field(j).Name == "UpdatedAt" {
				continue
			}
			if value.Kind() == reflect.Ptr && value.IsNil() {
				row = append(row, nil)
			} else {
				row = append(row, value.Interface())
			}
		}
		rows = append(rows, row)
	}
	return rows
}

func BulkUpdate[T any](ctx context.Context, pool *pgxpool.Pool, realTable string, data []T) (int64, error) {
	// トランザクション作成
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return 0, errs.New("failed to acquire connection from pool", errs.WithCause(err))
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return 0, errs.New("failed to start transaction", errs.WithCause(err))
	}
	defer tx.Rollback(ctx)

	// 一時テーブル作成
	tempTable := randstr.String(32)
	sql := fmt.Sprintf(
		`CREATE TEMPORARY TABLE %s (LIKE %s INCLUDING ALL);`,
		pgx.Identifier{tempTable}.Sanitize(),
		pgx.Identifier{realTable}.Sanitize(),
	)
	slog.Debug("create temporary table", slog.String("sql", sql))
	if _, err := tx.Exec(ctx, sql); err != nil {
		return 0, errs.New("failed to create temporary table", errs.WithCause(err))
	}
	defer func() {
		sql := fmt.Sprintf(`DROP TEMPORARY TABLE IF EXISTS %s;`, pgx.Identifier{tempTable}.Sanitize())
		slog.Debug("drop temporary table", slog.String("sql", sql))
		tx.Exec(ctx, sql)
	}()

	// データを一時テーブル内に投入
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{tempTable},
		Columns(new(T)),
		pgx.CopyFromRows(Rows(data)),
	)
	if err != nil {
		return 0, errs.New("failed to copy data into temporary table", errs.WithCause(err))
	}

	// 一時テーブルに存在するデータを実テーブルから削除
	uniqueKeys := UniqueKey(new(T))
	if len(uniqueKeys) == 0 {
		return 0, errs.New("no unique keys available")
	}

	exprs := make([]string, len(uniqueKeys))
	for i, key := range uniqueKeys {
		expr := fmt.Sprintf(
			`%s = %s`,
			pgx.Identifier{tempTable, key}.Sanitize(),
			pgx.Identifier{realTable, key}.Sanitize(),
		)
		exprs[i] = expr
	}
	sql = fmt.Sprintf(
		`DELETE FROM %s WHERE EXISTS (SELECT 1 FROM %s WHERE %s);`,
		pgx.Identifier{realTable}.Sanitize(),
		pgx.Identifier{tempTable}.Sanitize(),
		strings.Join(exprs, " AND "),
	)
	slog.Debug("delete from data", slog.String("sql", sql))
	if _, err := tx.Exec(ctx, sql); err != nil {
		return 0, errs.New("failed to delete existing rows from real table", errs.WithCause(err))
	}

	// 一時テーブル内のデータを実テーブルに投入
	sql = fmt.Sprintf(
		`INSERT INTO %s SELECT * FROM %s;`,
		pgx.Identifier{realTable}.Sanitize(),
		pgx.Identifier{tempTable}.Sanitize(),
	)
	slog.Debug("insert data into real table", slog.String("sql", sql))
	res, err := tx.Exec(ctx, sql)
	if err != nil {
		return 0, errs.New("failed to insert data of temporary table into real table", errs.WithCause(err))
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, errs.New("failed to commit transaction", errs.WithCause(err))
	}
	return res.RowsAffected(), nil
}
