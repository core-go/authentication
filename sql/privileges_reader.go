package sql

import (
	"context"
	"database/sql"
	auth "github.com/core-go/authentication"
	"reflect"
)

type PrivilegesReader struct {
	DB           *sql.DB
	Query        string
	NoSequence   bool
	Driver       string
	moduleFields map[string]int
}

func NewPrivilegesReader(db *sql.DB, query string, options ...bool) (*PrivilegesReader, error) {
	var handleDriver, noSequence bool
	if len(options) >= 1 {
		handleDriver = options[0]
	} else {
		handleDriver = true
	}
	if len(options) >= 2 {
		noSequence = options[1]
	} else {
		noSequence = true
	}
	driver := getDriver(db)
	if handleDriver {
		query = replaceQueryArgs(driver, query)
	}
	var module auth.Module
	moduleType := reflect.TypeOf(module)
	moduleFields, err := getColumnIndexes(moduleType)
	if err != nil {
		return nil, err
	}
	return &PrivilegesReader{DB: db, Query: query, NoSequence: noSequence, Driver: driver, moduleFields: moduleFields}, nil
}
func (l PrivilegesReader) Privileges(ctx context.Context) ([]auth.Privilege, error) {
	var models []auth.Module
	p0 := make([]auth.Privilege, 0)
	_, er1 := queryWithMap(ctx, l.DB, l.moduleFields, &models, l.Query)
	if er1 != nil {
		return p0, er1
	}
	var p []auth.Privilege
	if l.NoSequence == true {
		p = auth.ToPrivilegesWithNoSequence(models)
	} else {
		p = auth.ToPrivileges(models)
	}
	return p, nil
}
