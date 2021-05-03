package sql

import (
	"context"
	"database/sql"
	"github.com/core-go/auth"
	"strings"
)

type PrivilegesLoader struct {
	DB             *sql.DB
	Query          string
	ParameterCount int
	NoSequence     bool
	HandleDriver   bool
	Driver         string
	Or             bool
}

func NewPrivilegesLoader(db *sql.DB, query string, options ...int) *PrivilegesLoader {
	var parameterCount int
	if len(options) >= 1 && options[0] > 0 {
		parameterCount = options[0]
	} else {
		parameterCount = 0
	}
	return NewSqlPrivilegesLoader(db, query, parameterCount, false, true, true)
}
func NewSqlPrivilegesLoader(db *sql.DB, query string, parameterCount int, options ...bool) *PrivilegesLoader {
	var or, handleDriver, noSequence bool
	if len(options) >= 1 {
		or = options[0]
	} else {
		or = false
	}
	if len(options) >= 2 {
		handleDriver = options[1]
	} else {
		handleDriver = true
	}
	if len(options) >= 3 {
		noSequence = options[2]
	} else {
		noSequence = true
	}
	driver := getDriver(db)
	if handleDriver {
		query = replaceQueryArgs(driver, query)
	}
	return &PrivilegesLoader{DB: db, Query: query, ParameterCount: parameterCount, Or: or, NoSequence: noSequence, HandleDriver: handleDriver, Driver: driver}
}
func (l PrivilegesLoader) Load(ctx context.Context, id string) ([]auth.Privilege, error) {
	models := make([]auth.Module, 0)
	p0 := make([]auth.Privilege, 0)
	params := make([]interface{}, 0)
	params = append(params, id)
	if l.ParameterCount > 1 {
		for i := 2; i <= l.ParameterCount; i++ {
			params = append(params, id)
		}
	}
	columns, er1 := query(ctx, l.DB, &models, l.Query, params...)
	if er1 != nil {
		return p0, er1
	}
	hasPermission := hasPermissions(columns)
	if hasPermission && l.Or {
		models = auth.OrPermissions(models)
	}
	var p []auth.Privilege
	if l.NoSequence == true {
		p = auth.ToPrivilegesWithNoSequence(models)
	} else {
		p = auth.ToPrivileges(models)
	}
	return p, nil
}
func hasPermissions(cols []string) bool {
	for _, col := range cols {
		lcol := strings.ToLower(col)
		if lcol == "permissions" {
			return true
		}
	}
	return false
}
