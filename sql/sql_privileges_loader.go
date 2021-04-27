package sql

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"strings"

	"github.com/common-go/auth"
)

type SqlPrivilegesLoader struct {
	DB             *sql.DB
	Query          string
	ParameterCount int
	NoSequence     bool
	HandleDriver   bool
	Driver         string
	Or             bool
}

func NewPrivilegesLoader(db *sql.DB, query string, options ...int) *SqlPrivilegesLoader {
	var parameterCount int
	if len(options) >= 1 && options[0] > 0 {
		parameterCount = options[0]
	} else {
		parameterCount = 0
	}
	return NewSqlPrivilegesLoader(db, query, parameterCount, false, true, true)
}
func NewSqlPrivilegesLoader(db *sql.DB, query string, parameterCount int, options ...bool) *SqlPrivilegesLoader {
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
	return &SqlPrivilegesLoader{DB: db, Query: query, ParameterCount: parameterCount, Or: or, NoSequence: noSequence, HandleDriver: handleDriver, Driver: driver}
}
func (l SqlPrivilegesLoader) Load(ctx context.Context, id string) ([]auth.Privilege, error) {
	models := make([]auth.Module, 0)
	p0 := make([]auth.Privilege, 0)
	params := make([]interface{}, 0)
	params = append(params, id)
	if l.ParameterCount > 1 {
		for i := 2; i <= l.ParameterCount; i++ {
			params = append(params, id)
		}
	}
	driver := l.Driver
	rows, er1 := l.DB.Query(l.Query, params...)
	if er1 != nil {
		return p0, er1
	}
	defer rows.Close()
	columns, er2 := rows.Columns()
	hasPermission := hasPermissions(columns)
	if er2 != nil {
		return p0, er2
	}
	// get list indexes column
	modelTypes := reflect.TypeOf(models).Elem()
	modelType := reflect.TypeOf(auth.Module{})
	indexes, er3 := getColumnIndexes(modelType, columns, driver)
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
func scanType(rows *sql.Rows, modelTypes reflect.Type, indexes []int) (t []interface{}, err error) {
	for rows.Next() {
		initArray := reflect.New(modelTypes).Interface()
		if err = rows.Scan(structScan(initArray, indexes)...); err == nil {
			t = append(t, initArray)
		}
	}
	return
}
func structScan(s interface{}, indexColumns []int) (r []interface{}) {
	if s != nil {
		maps := reflect.Indirect(reflect.ValueOf(s))
		for _, index := range indexColumns {
			r = append(r, maps.Field(index).Addr().Interface())
		}
	}
	return
}

func getColumnIndex(modelType reflect.Type, columnsName string, driver string) (index int, err error) {
	if modelType.Kind() != reflect.Struct {
		return -1, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		column, ok := findTag(ormTag, "column")
		if driver == driverOracle {
			column = strings.ToUpper(column)
		} else {
			column = strings.ToLower(column)
		}
		if ok {
			if columnsName == column {
				return i, nil
			}
		}
	}
	return -1, errors.New("col " + columnsName + "not found")
}

func getColumnIndexes(modelType reflect.Type, columnsNames []string, driver string) (indexes []int, err error) {
	if modelType.Kind() != reflect.Struct {
		return nil, errors.New("bad type")
	}
	for i := 0; i < len(columnsNames); i++ {
		index, err := getColumnIndex(modelType, columnsNames[i], driver)
		if err != nil{
			return nil, err
		}
		indexes = append(indexes, index)
	}
	/*for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		column, ok := FindTag(ormTag, "column")
		if driver == DriverOracle {
			column = strings.ToUpper(column)
		}
		if ok {
			if contains(columnsNames, column) {
				indexes = append(indexes, i)
			}
		}
	}*/
	return
}
func findTag(tag string, key string) (string, bool) {
	if has := strings.Contains(tag, key); has {
		str1 := strings.Split(tag, ";")
		num := len(str1)
		for i := 0; i < num; i++ {
			str2 := strings.Split(str1[i], ":")
			for j := 0; j < len(str2); j++ {
				if str2[j] == key {
					return str2[j+1], true
				}
			}
		}
	}
	return "", false
}
func contains(array []string, v string) bool {
	for _, s := range array {
		if s == v {
			return true
		}
	}
	return false
}
