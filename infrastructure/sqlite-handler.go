package infrastructure

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteHandler struct {
	db *sql.DB
}

func NewSqliteHandler(filePath string) abstruct.SqlHandler {
	// file exists check
	s, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}
	if s.IsDir() {
		panic(fmt.Sprintf("path is directory: %s\n", filePath))
	}

	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		panic(err)
	}

	handler := &SqliteHandler{db: db}

	if _, err := handler.db.Query("SELECT 1+1;"); err != nil {
		panic(err)
	}

	return handler
}

func (handler *SqliteHandler) Prepare(t string) (abstruct.SqlStmt, error) {
	stmt, err := handler.db.Prepare(t)
	if err != nil {
		return nil, err
	}
	sqliteStmt := &SqliteStmt{stmt: stmt}
	return sqliteStmt, nil
}

func (handler *SqliteHandler) Query(t string) (abstruct.SqlRows, error) {
	rows, err := handler.db.Query(t)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (handler *SqliteHandler) Exec(t string) (abstruct.SqlResult, error) {
	result, err := handler.db.Exec(t)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (handler *SqliteHandler) Close() error {
	return handler.db.Close()
}

type SqliteStmt struct {
	stmt *sql.Stmt
}

func (stmt *SqliteStmt) Query(values ...any) (abstruct.SqlRows, error) {
	rows, err := stmt.stmt.Query(values...)
	if err != nil {
		return nil, utilerror.New(err.Error(), utilerror.ERR_SQL_QUERY)
	}
	return rows, nil
}

func (stmt *SqliteStmt) Exec(values ...any) (abstruct.SqlResult, error) {
	result, err := stmt.stmt.Exec(values...)
	if err != nil {
		return nil, utilerror.New(err.Error(), utilerror.ERR_SQL_QUERY)
	}
	return result, nil
}

func (stmt *SqliteStmt) Close() error {
	return stmt.stmt.Close()
}
