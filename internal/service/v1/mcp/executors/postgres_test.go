package executors

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// newMockDB opens a sqlmock-backed *sql.DB for one test. There is no pre-existing *sql.DB unit
// test double convention in this repo (the *sql.DB-consuming tests under tests/ are e2e suites
// against a real Postgres instance) — go-sqlmock was added as a new go.mod dependency for this
// package specifically because it implements database/sql/driver directly, so PostgresExecutor
// can be exercised against a real *sql.DB without a live database.
func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db, mock
}

func TestPostgresExecutor_Execute_UnknownTool(t *testing.T) {
	db, _ := newMockDB(t)
	e := NewPostgresExecutor()

	_, err := e.Execute(context.Background(), "unknown_tool", db, nil)
	if !errors.Is(err, user_errors.McpToolNotFound) {
		t.Fatalf("expected McpToolNotFound, got %v", err)
	}
}

func TestPostgresExecutor_ListTables(t *testing.T) {
	db, mock := newMockDB(t)
	e := NewPostgresExecutor()

	rows := sqlmock.NewRows([]string{"table_name"}).AddRow("users").AddRow("posts")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT table_name FROM information_schema.tables")).WillReturnRows(rows)

	result, err := e.Execute(context.Background(), ToolPgListTables, db, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var tables []string

	err = json.Unmarshal([]byte(result.Text), &tables)
	if err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	want := []string{"users", "posts"}
	if !reflect.DeepEqual(tables, want) {
		t.Fatalf("tables = %v, want %v", tables, want)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestPostgresExecutor_ListTables_QueryError(t *testing.T) {
	db, mock := newMockDB(t)
	e := NewPostgresExecutor()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT table_name FROM information_schema.tables")).
		WillReturnError(errors.New("boom"))

	_, err := e.Execute(context.Background(), ToolPgListTables, db, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPostgresExecutor_DescribeTable_MissingTableName(t *testing.T) {
	db, _ := newMockDB(t)
	e := NewPostgresExecutor()

	_, err := e.Execute(context.Background(), ToolPgDescribeTable, db, map[string]interface{}{})
	if !errors.Is(err, pgTableNameRequired) {
		t.Fatalf("expected pgTableNameRequired, got %v", err)
	}
}

func TestPostgresExecutor_DescribeTable(t *testing.T) {
	db, mock := newMockDB(t)
	e := NewPostgresExecutor()

	rows := sqlmock.NewRows([]string{"column_name", "data_type", "is_nullable", "column_default"}).
		AddRow("id", "integer", "NO", nil).
		AddRow("name", "text", "YES", "'unnamed'::text")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT column_name, data_type, is_nullable, column_default")).
		WithArgs("users").
		WillReturnRows(rows)

	params := map[string]interface{}{fieldTableName: "users"}

	result, err := e.Execute(context.Background(), ToolPgDescribeTable, db, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var columns []tableColumnRow

	err = json.Unmarshal([]byte(result.Text), &columns)
	if err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	if len(columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(columns))
	}

	if columns[0].ColumnName != "id" || columns[0].IsNullable {
		t.Fatalf("unexpected first column: %+v", columns[0])
	}

	if columns[1].ColumnName != "name" || !columns[1].IsNullable || columns[1].ColumnDefault != "'unnamed'::text" {
		t.Fatalf("unexpected second column: %+v", columns[1])
	}
}

func TestPostgresExecutor_Query_MissingSql(t *testing.T) {
	db, _ := newMockDB(t)
	e := NewPostgresExecutor()

	_, err := e.Execute(context.Background(), ToolPgQuery, db, map[string]interface{}{})
	if !errors.Is(err, pgSqlRequired) {
		t.Fatalf("expected pgSqlRequired, got %v", err)
	}
}

func TestPostgresExecutor_Query(t *testing.T) {
	db, mock := newMockDB(t)
	e := NewPostgresExecutor()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, []byte("alice")).
		AddRow(2, []byte("bob"))

	sqlText := "SELECT id, name FROM users"
	mock.ExpectQuery(regexp.QuoteMeta(sqlText)).WillReturnRows(rows)

	params := map[string]interface{}{fieldSql: sqlText}

	result, err := e.Execute(context.Background(), ToolPgQuery, db, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded struct {
		Rows      []map[string]interface{} `json:"rows"`
		Truncated bool                     `json:"truncated"`
	}

	err = json.Unmarshal([]byte(result.Text), &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	if decoded.Truncated {
		t.Fatal("expected truncated=false")
	}

	if len(decoded.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(decoded.Rows))
	}

	// normalizeScanValue must have converted the []byte("alice") scan target into a plain
	// JSON string, not a base64-encoded byte array.
	if decoded.Rows[0]["name"] != "alice" {
		t.Fatalf("expected name to normalize to \"alice\", got %#v", decoded.Rows[0]["name"])
	}
}

func TestPostgresExecutor_Query_Truncation(t *testing.T) {
	db, mock := newMockDB(t)
	e := NewPostgresExecutor()

	rows := sqlmock.NewRows([]string{"n"})
	for i := 0; i < maxQueryRows+1; i++ {
		rows.AddRow(i)
	}

	sqlText := "SELECT n FROM generate_series(1, 1001) n"
	mock.ExpectQuery(regexp.QuoteMeta(sqlText)).WillReturnRows(rows)

	params := map[string]interface{}{fieldSql: sqlText}

	result, err := e.Execute(context.Background(), ToolPgQuery, db, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded struct {
		Rows      []map[string]interface{} `json:"rows"`
		Truncated bool                     `json:"truncated"`
	}

	err = json.Unmarshal([]byte(result.Text), &decoded)
	if err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}

	if !decoded.Truncated {
		t.Fatal("expected truncated=true")
	}

	if len(decoded.Rows) != maxQueryRows {
		t.Fatalf("expected %d rows, got %d", maxQueryRows, len(decoded.Rows))
	}
}

func TestPostgresExecutor_Execute_MissingSql(t *testing.T) {
	db, _ := newMockDB(t)
	e := NewPostgresExecutor()

	_, err := e.Execute(context.Background(), ToolPgExecute, db, map[string]interface{}{})
	if !errors.Is(err, pgSqlRequired) {
		t.Fatalf("expected pgSqlRequired, got %v", err)
	}
}

func TestPostgresExecutor_Execute(t *testing.T) {
	db, mock := newMockDB(t)
	e := NewPostgresExecutor()

	sqlText := "UPDATE users SET active = true"
	mock.ExpectExec(regexp.QuoteMeta(sqlText)).WillReturnResult(sqlmock.NewResult(0, 3))

	params := map[string]interface{}{fieldSql: sqlText}

	result, err := e.Execute(context.Background(), ToolPgExecute, db, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "OK, 3 row(s) affected"
	if result.Text != want {
		t.Fatalf("result.Text = %q, want %q", result.Text, want)
	}
}

func TestPostgresExecutor_Execute_ExecError(t *testing.T) {
	db, mock := newMockDB(t)
	e := NewPostgresExecutor()

	sqlText := "DROP TABLE nonexistent"
	mock.ExpectExec(regexp.QuoteMeta(sqlText)).WillReturnError(errors.New("boom"))

	params := map[string]interface{}{fieldSql: sqlText}

	_, err := e.Execute(context.Background(), ToolPgExecute, db, params)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNormalizeScanValue(t *testing.T) {
	got := normalizeScanValue([]byte("hello"))
	if got != "hello" {
		t.Fatalf("normalizeScanValue([]byte) = %#v, want %q", got, "hello")
	}

	got = normalizeScanValue(int64(42))
	if got != int64(42) {
		t.Fatalf("normalizeScanValue(int64) = %#v, want 42", got)
	}

	got = normalizeScanValue(nil)
	if got != nil {
		t.Fatalf("normalizeScanValue(nil) = %#v, want nil", got)
	}
}

func TestPostgresToolDefinitions(t *testing.T) {
	defs := PostgresToolDefinitions()

	if len(defs) != 4 {
		t.Fatalf("expected 4 tool definitions, got %d", len(defs))
	}

	names := make(map[string]bool, len(defs))
	for _, d := range defs {
		names[d.ApiDescription.Name] = true
	}

	wantNames := []string{ToolPgListTables, ToolPgDescribeTable, ToolPgQuery, ToolPgExecute}
	for _, want := range wantNames {
		if !names[want] {
			t.Fatalf("missing tool definition for %s", want)
		}
	}
}
