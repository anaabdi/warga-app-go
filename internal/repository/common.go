package repository

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/anaabdi/warga-app-go/pkg/postgres"
)

type QueryAdapter interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type database struct {
	*postgres.DB
}

func NewDatabase(db *postgres.DB) Database {
	return &database{
		DB: db,
	}
}

type Database interface {
	GetDB() *sql.DB
	Begin(ctx context.Context, txOpts *sql.TxOptions) (*sql.Tx, error)
}

func (gw *database) GetDB() *sql.DB {
	return gw.DB.Pool
}

func (gw *database) Begin(ctx context.Context, txOpts *sql.TxOptions) (*sql.Tx, error) {
	return gw.DB.Pool.BeginTx(ctx, txOpts)
}

func closeSilently(rows io.Closer) {
	_ = rows.Close()
}

func like(columnCondition, keyword string) string {
	return fmt.Sprintf("%s LIKE '%%%v%%'", columnCondition, keyword)
}

func ilike(columnCondition, keyword string) string {
	return fmt.Sprintf("%s ILIKE '%%%v%%'", columnCondition, keyword)
}

func appendPagination(query string, page, limit int) string {
	offset := (page - 1) * limit

	return fmt.Sprintf("%s limit %d offset %d", query, limit, offset)
}

func appendOrderBy(query string, column string, isDecremented bool) string {
	if isDecremented {
		return fmt.Sprintf("%s ORDER BY %s DESC", query, column)
	}

	return fmt.Sprintf("%s ORDER BY %s", query, column)

}

func appendKeywordSearch(query string, keyword string, columms []string) string {
	var buf bytes.Buffer

	buf.WriteString(query)
	buf.WriteString(" AND ")

	buf.WriteString("(")

	for i, columnName := range columms {
		if i != 0 {
			buf.WriteString(" OR ")
		}

		buf.WriteString(ilike(columnName, keyword))
	}

	buf.WriteString(")")

	return buf.String()
}

func removeBracketFromCoordinate(coord string) string {
	return strings.ReplaceAll(strings.ReplaceAll(coord, "(", ""), ")", "")
}
