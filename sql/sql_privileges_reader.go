package sql

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/common-go/auth"
)

const (
	DriverPostgres         = "postgres"
	DriverMysql            = "mysql"
	DriverMssql            = "mssql"
	DriverOracle           = "oracle"
	DriverSqlite3          = "sqlite3"
	DriverNotSupport       = "no support"
)

type SqlPrivilegesReader struct {
	DB         *sql.DB
	Query      string
	NoSequence bool
	Driver     string
}
func NewPrivilegesReader(db *sql.DB, query string, options...bool) *SqlPrivilegesReader {
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
		return DriverNotSupport
	}
	driver := reflect.TypeOf(db.Driver()).String()
	switch driver {
	case "*pq.Driver":
		return DriverPostgres
	case "*godror.drv":
		return DriverOracle
	case "*mysql.MySQLDriver":
		return DriverMysql
	case "*mssql.Driver":
		return DriverMssql
	case "*sqlite3.SQLiteDriver":
		return DriverSqlite3
	default:
		return DriverNotSupport
	}
}

func replaceQueryArgs(driver string, query string) string {
	if driver == DriverOracle || driver == DriverPostgres || driver == DriverMssql {
		var x string
		if driver == DriverOracle {
			x = ":val"
		} else if driver == DriverPostgres {
			x = "$"
		} else if driver == DriverMssql {
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
