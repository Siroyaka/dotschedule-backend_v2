package infrastructure

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	_ "modernc.org/sqlite"
)

func NewSqliteHandlerCGOLess(filePath string) abstruct.SqlHandler {
	// file exists check
	s, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}
	if s.IsDir() {
		panic(fmt.Sprintf("path is directory: %s\n", filePath))
	}

	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		panic(err)
	}

	handler := &SqliteHandler{db: db}

	if _, err := handler.db.Query("SELECT 1+1;"); err != nil {
		panic(err)
	}

	return handler
}
