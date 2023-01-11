package abstruct

type SqlHandler interface {
	Query(string) (SqlRows, error)
	Prepare(string) (SqlStmt, error)
	Exec(string) (SqlResult, error)
	Close() error
}

type SqlResult interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type SqlRows interface {
	Scan(...interface{}) error
	Next() bool
	Close() error
}

type SqlStmt interface {
	Query(...any) (SqlRows, error)
	Exec(values ...any) (SqlResult, error)
	Close() error
}
