package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/core-go/auth"
	"reflect"
	"strings"
)

const (
	driverPostgres   = "postgres"
	driverMysql      = "mysql"
	driverMssql      = "mssql"
	driverOracle     = "oracle"
	driverSqlite3    = "sqlite3"
	driverNotSupport = "no support"
)

type SqlPrivilegesReader struct {
	DB         *sql.DB
	Query      string
	NoSequence bool
	Driver     string
}

func NewPrivilegesReader(db *sql.DB, query string, options ...bool) *SqlPrivilegesReader {
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
	return &SqlPrivilegesReader{DB: db, Query: query, NoSequence: noSequence, Driver: driver}
}
func (l SqlPrivilegesReader) Privileges(ctx context.Context) ([]auth.Privilege, error) {
	models := make([]auth.Module, 0)
	p0 := make([]auth.Privilege, 0)
	rows, er1 := l.DB.QueryContext(ctx, l.Query)
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
	modelType := reflect.TypeOf(auth.Module{})
	indexes, er3 := getColumnIndexes(modelType, columns, l.Driver)
	if er3 != nil {
		return p0, er3
	}
	tb, er4 := scanType(rows, modelTypes, indexes)
	if er4 != nil {
		return p0, er4
	}
	for _, v := range tb {
		if c, ok := v.(*auth.Module); ok {
			models = append(models, *c)
		}
	}
	var p []auth.Privilege
	if l.NoSequence == true {
		p = auth.ToPrivilegesWithNoSequence(models)
	} else {
		p = auth.ToPrivileges(models)
	}
	return p, nil
}
func getDriver(db *sql.DB) string {
	if db == nil {
		return driverNotSupport
	}
	driver := reflect.TypeOf(db.Driver()).String()
	switch driver {
	case "*pq.Driver":
		return driverPostgres
	case "*godror.drv":
		return driverOracle
	case "*mysql.MySQLDriver":
		return driverMysql
	case "*mssql.Driver":
		return driverMssql
	case "*sqlite3.SQLiteDriver":
		return driverSqlite3
	default:
		return driverNotSupport
	}
}

func replaceQueryArgs(driver string, query string) string {
	if driver == driverOracle || driver == driverPostgres || driver == driverMssql {
		var x string
		if driver == driverOracle {
			x = ":val"
		} else if driver == driverPostgres {
			x = "$"
		} else if driver == driverMssql {
			x = "@p"
		}
		i := 1
		k := strings.Index(query, "?")
		if k >= 0 {
			for {
				query = strings.Replace(query, "?", x+fmt.Sprintf("%v", i), 1)
				i = i + 1
				k := strings.Index(query, "?")
				if k < 0 {
					return query
				}
			}
		}
	}
	return query
}
