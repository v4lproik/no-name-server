package infrastructure

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
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

func (handler *SqliteHandler) Execute(statement string) {
	handler.Conn.Exec(statement)
}

func (handler *SqliteHandler) ExecuteWithParam(statement string, args ...interface{}) {
	stmt, err := handler.Conn.Prepare(statement)
	if err != nil { panic(err) }
	defer stmt.Close()

	_, err2 := stmt.Exec(args...)
	if err2 != nil { panic(err2) }
}

func (handler *SqliteHandler) Query(statement string) repositories.Row {
	rows, err := handler.Conn.Query(statement)
	if err != nil {
		fmt.Println(err)
		return new(SqliteRow)
	}
	row := new(SqliteRow)
	row.Rows = rows
	return row
}

func (handler *SqliteHandler) QueryWithParam(statement string, args ...interface{}) repositories.Row {
	rows, err := handler.Conn.Query(statement, args...)
	if err != nil {
		fmt.Println(err)
		return new(SqliteRow)
	}
	row := new(SqliteRow)
	row.Rows = rows
	return row
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