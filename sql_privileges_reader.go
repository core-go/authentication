package auth

import (
	"context"
	"database/sql"
	"reflect"
)

const (
	DriverPostgres   = "postgres"
	DriverMysql      = "mysql"
	DriverMssql      = "mssql"
	DriverOracle     = "oracle"
	DriverNotSupport = "no support"
)

type SqlPrivilegesReader struct {
	DB         *sql.DB
	Query      string
	NoSequence bool
	Driver     string
}
func NewPrivilegesReader(db *sql.DB, query string) *SqlPrivilegesReader {
	return NewSqlPrivilegesReader(db, query, true)
}
func NewSqlPrivilegesReader(db *sql.DB, query string, noSequence bool) *SqlPrivilegesReader {
	driver := GetDriver(db)
	return &SqlPrivilegesReader{DB: db, Query: query, NoSequence: noSequence, Driver: driver}
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
	indexes, er3 := GetColumnIndexes(modelType, columns, l.Driver)
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
	var p []Privilege
	if l.NoSequence == true {
		p = ToPrivilegesWithNoSequence(models)
	} else {
		p = ToPrivileges(models)
	}
	return p, nil
}

func GetDriver(db *sql.DB) string {
	driver := reflect.TypeOf(db.Driver()).String()
	switch driver {
	case "*postgres.Driver":
		return DriverPostgres
	case "*mysql.MySQLDriver":
		return DriverMysql
	case "*mssql.Driver":
		return DriverMssql
	case "*godror.drv":
		return DriverOracle
	default:
		return DriverNotSupport
	}
}
