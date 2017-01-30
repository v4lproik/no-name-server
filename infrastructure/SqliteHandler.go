package infrastructure

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yinkozi/no-name-server/interfaces/repositories"
)

type SqliteHandler struct {
	Conn *sql.DB
}

func NewSqliteHandler(filepath string) *SqliteHandler {
	db, err := sql.Open("sqlite3", filepath + "?_busy_timeout=5000")
	if err != nil { panic(err) }
	if db == nil { panic("db nil") }

	sqliteHandler := new(SqliteHandler)
	sqliteHandler.Conn = db

	return sqliteHandler
}

func (handler *SqliteHandler) Execute(statement string) error {
	_, err := handler.Conn.Exec(statement)
	return err
}

func (handler *SqliteHandler) ExecuteWithParam(statement string, args ...interface{}) error {
	stmt, err := handler.Conn.Prepare(statement)
	if err != nil { return err }
	defer stmt.Close()

	_, err2 := stmt.Exec(args...)
	if err2 != nil { return err2 }

	return nil
}

func (handler *SqliteHandler) Query(statement string) (repositories.Row, error) {
	rows, err := handler.Conn.Query(statement)
	if err != nil {
		return new(SqliteRow), err
	}
	row := new(SqliteRow)
	row.Rows = rows
	return row, nil
}

func (handler *SqliteHandler) QueryWithParam(statement string, args ...interface{}) (repositories.Row, error) {
	rows, err := handler.Conn.Query(statement, args...)
	if err != nil {
		return new(SqliteRow), err
	}
	row := new(SqliteRow)
	row.Rows = rows
	return row, nil
}

type SqliteRow struct {
	Rows *sql.Rows
}

func (r SqliteRow) Scan(dest ...interface{}) {
	r.Rows.Scan(dest...)
}

func (r SqliteRow) Next() bool {
	return r.Rows.Next()
}

func (r SqliteRow) Close() {
	r.Rows.Close()
}