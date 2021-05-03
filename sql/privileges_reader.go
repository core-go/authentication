package sql

import (
	"context"
	"database/sql"
	"github.com/core-go/auth"
)

type PrivilegesReader struct {
	DB         *sql.DB
	Query      string
	NoSequence bool
	Driver     string
}

func NewPrivilegesReader(db *sql.DB, query string, options ...bool) *PrivilegesReader {
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
	return &PrivilegesReader{DB: db, Query: query, NoSequence: noSequence, Driver: driver}
}
func (l PrivilegesReader) Privileges(ctx context.Context) ([]auth.Privilege, error) {
	models := make([]auth.Module, 0)
	p0 := make([]auth.Privilege, 0)
	_, er1 := query(ctx, l.DB, &models, l.Query)
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
