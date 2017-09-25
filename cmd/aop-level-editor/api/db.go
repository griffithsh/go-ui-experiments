package api

import (
	"database/sql"
	"fmt"
)

type database struct {
	db *sql.DB
}

func newDatabase(db *sql.DB) *database {
	return &database{
		db: db,
	}
}

// Result of a database query to be encoded into JSON and returned in a http response.
type Result struct {
	Columns []string        `json:"Columns"`
	Rows    [][]interface{} `json:"Rows"`
	Error   error           `json:"Error,omitempty"`
}

func (db *database) execute(sql string) (*Result, error) {
	rows, err := db.db.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("query \"%s\": %s", sql, err)
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("query \"%s\": %s", sql, err)
	}
	var result [][]interface{}
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		rows.Scan(vals...)
		result = append(result, vals)
	}
	return &Result{
		Columns: cols,
		Rows:    result,
	}, rows.Err()
}
