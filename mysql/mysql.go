package mysql

import (
	"context"
	"database/sql"

	"github.com/cntechpower/utils/tracing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/opentracing/opentracing-go/ext"
)

const (
	dbTypeMySQL = "mysql"
)

type DB struct {
	*sql.DB
}

func New(dsn string) (db *DB, err error) {
	db = &DB{}
	m, err := sql.Open(dbTypeMySQL, dsn)
	if err != nil {
		return
	}
	db.DB = m
	return
}

func (d *DB) Query(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error) {
	span, _ := tracing.New(ctx, "mysql.Query")
	ext.DBStatement.Set(span, query)
	ext.DBType.Set(span, dbTypeMySQL)
	rows, err = d.DB.Query(query, args...)
	if err != nil {
		ext.LogError(span, err)
		ext.Error.Set(span, true)
	}
	span.Finish()
	return
}

func (d *DB) QueryRow(ctx context.Context, query string, args ...interface{}) (rows *sql.Row) {
	span, _ := tracing.New(ctx, "mysql.QueryRow")
	ext.DBStatement.Set(span, query)
	ext.DBType.Set(span, dbTypeMySQL)
	rows = d.DB.QueryRowContext(ctx, query, args...)
	if rows.Err() != nil {
		ext.LogError(span, rows.Err())
		ext.Error.Set(span, true)
	}
	span.Finish()
	return
}

func (d *DB) Exec(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	span, _ := tracing.New(ctx, "mysql.Exec")
	ext.DBStatement.Set(span, query)
	ext.DBType.Set(span, dbTypeMySQL)
	res, err = d.DB.Exec(query, args...)
	if err != nil {
		ext.LogError(span, err)
		ext.Error.Set(span, true)
	}
	span.Finish()
	return
}
