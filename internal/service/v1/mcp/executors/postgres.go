package executors

import (
	"context"
	"database/sql"
	"fmt"

	"go.redsock.ru/rerrors"
	"google.golang.org/grpc/codes"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/utils"
)

const (
	ToolPgListTables    = "pg_list_tables"
	ToolPgDescribeTable = "pg_describe_table"
	ToolPgQuery         = "pg_query"
	ToolPgExecute       = "pg_execute"
)

// maxQueryRows caps how many rows pg_query returns in one call — analogous in spirit to
// http.go's httpMaxResponseBytes cap on HttpExecutor responses. Truncation is reported via the
// result's "truncated" field rather than silently dropping rows.
const maxQueryRows = 1000

// pgTableNameRequired/pgSqlRequired validate this executor's own required params. They are
// defined here rather than in internal/service/user_errors (McpTableNameRequired/McpSqlRequired
// do not exist there yet) because user_errors.go is owned by a parallel in-flight task per this
// track's instructions — see the Track C report for a request to promote these into
// internal/service/user_errors mirroring McpPathRequired's shape.
var (
	pgTableNameRequired = rerrors.New("table_name is required and must be a string", codes.InvalidArgument)
	pgSqlRequired       = rerrors.New("sql is required and must be a string", codes.InvalidArgument)
)

// PostgresExecutor implements the pg_list_tables/pg_describe_table/pg_query/pg_execute builtin
// tools. It holds no state of its own — the vault's tenant *sql.DB connection is resolved by
// internal/service/v1/mcp.ExecuteTool (via the shared vaultpg.ConnPool) and passed in per call,
// mirroring how VaultExecutor receives its couchdb/S3 clients per call.
type PostgresExecutor struct{}

func NewPostgresExecutor() *PostgresExecutor {
	return &PostgresExecutor{}
}

// Execute dispatches a postgres tool call against db, the caller vault's dedicated tenant
// database connection.
func (e *PostgresExecutor) Execute(
	ctx context.Context,
	toolName string,
	db *sql.DB,
	params map[string]interface{},
) (domain.ToolExecResult, error) {
	switch toolName {
	case ToolPgListTables:
		return e.listTables(ctx, db)
	case ToolPgDescribeTable:
		return e.describeTable(ctx, db, params)
	case ToolPgQuery:
		return e.query(ctx, db, params)
	case ToolPgExecute:
		return e.execute(ctx, db, params)
	default:
		return domain.ToolExecResult{}, user_errors.McpToolNotFound
	}
}

// -- pg_list_tables --

func (e *PostgresExecutor) listTables(ctx context.Context, db *sql.DB) (domain.ToolExecResult, error) {
	const q = "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_name"

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error listing tables")
	}
	defer utils.CloseWithLog(rows, "pg_list_tables rows")

	var tables []string

	for rows.Next() {
		var tableName string

		err = rows.Scan(&tableName)
		if err != nil {
			return domain.ToolExecResult{}, rerrors.Wrap(err, "error scanning table name")
		}

		tables = append(tables, tableName)
	}

	err = rows.Err()
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error iterating table rows")
	}

	text, err := marshalResult(tables)
	if err != nil {
		return domain.ToolExecResult{}, err
	}

	return domain.ToolExecResult{Text: text}, nil
}

// -- pg_describe_table --

type tableColumnRow struct {
	ColumnName    string `json:"column_name"`
	DataType      string `json:"data_type"`
	IsNullable    bool   `json:"is_nullable"`
	ColumnDefault string `json:"column_default,omitempty"`
}

func (e *PostgresExecutor) describeTable(
	ctx context.Context,
	db *sql.DB,
	params map[string]interface{},
) (domain.ToolExecResult, error) {
	tableName, ok := paramString(params, fieldTableName)
	if !ok || tableName == "" {
		return domain.ToolExecResult{}, pgTableNameRequired
	}

	const q = `SELECT column_name, data_type, is_nullable, column_default
FROM information_schema.columns
WHERE table_schema = 'public' AND table_name = $1
ORDER BY ordinal_position`

	rows, err := db.QueryContext(ctx, q, tableName)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error describing table")
	}
	defer utils.CloseWithLog(rows, "pg_describe_table rows")

	var columns []tableColumnRow

	for rows.Next() {
		var (
			columnName    string
			dataType      string
			isNullable    string
			columnDefault sql.NullString
		)

		err = rows.Scan(&columnName, &dataType, &isNullable, &columnDefault)
		if err != nil {
			return domain.ToolExecResult{}, rerrors.Wrap(err, "error scanning column row")
		}

		column := tableColumnRow{
			ColumnName:    columnName,
			DataType:      dataType,
			IsNullable:    isNullable == "YES",
			ColumnDefault: columnDefault.String,
		}

		columns = append(columns, column)
	}

	err = rows.Err()
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error iterating column rows")
	}

	text, err := marshalResult(columns)
	if err != nil {
		return domain.ToolExecResult{}, err
	}

	return domain.ToolExecResult{Text: text}, nil
}

// -- pg_query --

func (e *PostgresExecutor) query(
	ctx context.Context,
	db *sql.DB,
	params map[string]interface{},
) (domain.ToolExecResult, error) {
	sqlText, ok := paramString(params, fieldSql)
	if !ok || sqlText == "" {
		return domain.ToolExecResult{}, pgSqlRequired
	}

	// No bind params: sqlText is the whole agent-authored statement, per the "full access, role
	// isolation only" decision — there is no separate params list to parameterize against.
	rows, err := db.QueryContext(ctx, sqlText)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error executing query")
	}
	defer utils.CloseWithLog(rows, "pg_query rows")

	columns, err := rows.Columns()
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error reading result columns")
	}

	resultRows, truncated, err := scanRowsToMaps(rows, columns)
	if err != nil {
		return domain.ToolExecResult{}, err
	}

	err = rows.Err()
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error iterating query rows")
	}

	result := map[string]interface{}{
		fieldRows:      resultRows,
		fieldTruncated: truncated,
	}

	text, err := marshalResult(result)
	if err != nil {
		return domain.ToolExecResult{}, err
	}

	return domain.ToolExecResult{Text: text}, nil
}

// scanRowsToMaps drains rows into a slice of column-name-keyed maps, stopping and reporting
// truncated=true once maxQueryRows have been read rather than reading the entire result set.
func scanRowsToMaps(rows *sql.Rows, columns []string) ([]map[string]interface{}, bool, error) {
	var results []map[string]interface{}

	truncated := false

	for rows.Next() {
		if len(results) >= maxQueryRows {
			truncated = true

			break
		}

		values := make([]interface{}, len(columns))
		scanTargets := make([]interface{}, len(columns))

		for i := range values {
			scanTargets[i] = &values[i]
		}

		err := rows.Scan(scanTargets...)
		if err != nil {
			return nil, false, rerrors.Wrap(err, "error scanning row")
		}

		row := make(map[string]interface{}, len(columns))
		for i, col := range columns {
			row[col] = normalizeScanValue(values[i])
		}

		results = append(results, row)
	}

	return results, truncated, nil
}

// normalizeScanValue converts driver-returned []byte (lib/pq returns raw bytes for several text
// types when scanned into interface{}) into a string so marshalResult produces JSON strings
// instead of base64-encoded byte arrays.
func normalizeScanValue(v interface{}) interface{} {
	b, ok := v.([]byte)
	if ok {
		return string(b)
	}

	return v
}

// -- pg_execute --

func (e *PostgresExecutor) execute(
	ctx context.Context,
	db *sql.DB,
	params map[string]interface{},
) (domain.ToolExecResult, error) {
	sqlText, ok := paramString(params, fieldSql)
	if !ok || sqlText == "" {
		return domain.ToolExecResult{}, pgSqlRequired
	}

	result, err := db.ExecContext(ctx, sqlText)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error executing statement")
	}

	// RowsAffected is not meaningfully supported by every statement (e.g. CREATE TABLE/CREATE
	// FUNCTION) — tolerate that by falling back to 0 rather than treating it as an error.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		rowsAffected = 0
	}

	text := fmt.Sprintf("OK, %d row(s) affected", rowsAffected)

	return domain.ToolExecResult{Text: text}, nil
}

// -- tool catalog --

func PostgresToolDefinitions() []domain.McpToolDef {
	queryDescription := fmt.Sprintf(
		"Run a read query (e.g. SELECT) against the vault's dedicated Postgres database and return rows as "+
			"JSON. The sql string is executed as-is, exactly as written, against the vault's own isolated "+
			"database and role — there is no parsing, rewriting, or statement allowlist, so write valid "+
			"Postgres SQL and be mindful of what you run. Results are capped at %d rows; truncated is set "+
			"true when the cap was hit.",
		maxQueryRows,
	)

	executeDescription := "Run a write/DDL statement (INSERT/UPDATE/DELETE/CREATE TABLE/CREATE FUNCTION/etc) " +
		"against the vault's dedicated Postgres database. The sql string is executed as-is against the " +
		"vault's own isolated database and role — there is no statement allowlist, so the agent is fully " +
		"responsible for what it runs."

	return []domain.McpToolDef{
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolPgListTables,
				Description: "List all tables in the vault's dedicated Postgres database (public schema)",
				Properties:  map[string]domain.ToolProperty{},
				Required:    []string{},
			},
			OutputSchema: domain.ToolSchema{
				IsArray: true,
				Items:   &domain.ToolProperty{Type: schemaTypeString, Description: "Table name"},
			},
		},
		{
			ApiDescription: domain.ToolApiDescription{
				Name: ToolPgDescribeTable,
				Description: "Describe a table's columns (name, data type, nullability, default) in the " +
					"vault's dedicated Postgres database",
				Properties: map[string]domain.ToolProperty{
					fieldTableName: {Type: schemaTypeString, Description: "Table name, public schema"},
				},
				Required: []string{fieldTableName},
			},
			OutputSchema: domain.ToolSchema{
				IsArray: true,
				Items: &domain.ToolProperty{
					Type: schemaTypeObject,
					Properties: map[string]domain.ToolProperty{
						"column_name":    {Type: schemaTypeString},
						"data_type":      {Type: schemaTypeString},
						"is_nullable":    {Type: schemaTypeBoolean},
						"column_default": {Type: schemaTypeString, Description: "Empty if the column has no default"},
					},
					Required: []string{"column_name", "data_type", "is_nullable"},
				},
			},
		},
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolPgQuery,
				Description: queryDescription,
				Properties: map[string]domain.ToolProperty{
					fieldSql: {Type: schemaTypeString, Description: "Full SQL statement to execute, e.g. a SELECT"},
				},
				Required: []string{fieldSql},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					fieldRows: {
						Type:        schemaTypeArray,
						Description: "Result rows, each an object keyed by column name",
						Items:       &domain.ToolProperty{Type: schemaTypeObject},
					},
					fieldTruncated: {
						Type:        schemaTypeBoolean,
						Description: "True if the result was capped at the row limit",
					},
				},
				Required: []string{fieldRows, fieldTruncated},
			},
		},
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolPgExecute,
				Description: executeDescription,
				Properties: map[string]domain.ToolProperty{
					fieldSql: {Type: schemaTypeString, Description: "Full SQL statement to execute"},
				},
				Required: []string{fieldSql},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					fieldMessage: {Type: schemaTypeString, Description: "e.g. \"OK, 3 row(s) affected\""},
				},
				Required: []string{fieldMessage},
			},
		},
	}
}
