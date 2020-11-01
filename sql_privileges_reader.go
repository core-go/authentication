package auth

import (
	"context"
	"database/sql"
	"reflect"
)

type SqlPrivilegesReader struct {
	DB    *sql.DB
	Query string
}

func NewSqlPrivilegesReader(db *sql.DB, query string) *SqlPrivilegesReader {
	return &SqlPrivilegesReader{DB: db, Query: query}
}
func (l SqlPrivilegesReader) Privileges(ctx context.Context) ([]Privilege, error) {
	models := make([]Module, 0)
	p0 := make([]Privilege, 0)
	rows, er1 := l.DB.Query(l.Query)
	if er1 != nil {
		return p0, er1
	}
	defer rows.Close()
	columns, er2 := rows.Columns()
	if er2 != nil {
		return p0, er2
	}
	// get list indexes column
	modelTypes := reflect.TypeOf(models).Elem()
	modelType := reflect.TypeOf(Module{})
	indexes, er3 := getColumnIndexes(modelType, columns)
	if er3 != nil {
		return p0, er3
	}
	tb, er4 := ScanType(rows, modelTypes, indexes)
	if er4 != nil {
		return p0, er4
	}
	for _, v := range tb {
		if c, ok := v.(*Module); ok {
			models = append(models, *c)
		}
	}
	p := ToPrivileges(models)
	return p, nil
}
