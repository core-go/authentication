package auth

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type SqlPrivilegesLoader struct {
	DB             *sql.DB
	Query          string
	ParameterCount int
	NoSequence     bool
	HandleDriver   bool
	Driver         string
}

func NewSqlPrivilegesLoader(db *sql.DB, query string, parameterCount int, noSequence bool, handleDriver bool) *SqlPrivilegesLoader {
	driver := GetDriver(db)
	return &SqlPrivilegesLoader{DB: db, Query: query, ParameterCount: parameterCount, NoSequence: noSequence, HandleDriver: handleDriver, Driver: driver}
}
func (l SqlPrivilegesLoader) Load(ctx context.Context, id string) ([]Privilege, error) {
	models := make([]Module, 0)
	p0 := make([]Privilege, 0)
	params := make([]interface{}, 0)
	params = append(params, id)
	if l.ParameterCount > 1 {
		for i := 2; i <= l.ParameterCount; i++ {
			params = append(params, id)
		}
	}
	driver := l.Driver
	if l.HandleDriver {
		if driver == DriverOracle || driver == DriverPostgres {
			var x string
			if driver == DriverOracle {
				x = ":val"
			} else {
				x = "$"
			}
			for i := 0; i < len(params); i++ {
				count := i + 1
				l.Query = strings.Replace(l.Query, "?", x+strconv.Itoa(count), 1)
			}
		}
	}
	rows, er1 := l.DB.Query(l.Query, params...)
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
	indexes, er3 := GetColumnIndexes(modelType, columns, driver)
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

func ScanType(rows *sql.Rows, modelTypes reflect.Type, indexes []int) (t []interface{}, err error) {
	for rows.Next() {
		initArray := reflect.New(modelTypes).Interface()
		if err = rows.Scan(StructScan(initArray, indexes)...); err == nil {
			t = append(t, initArray)
		}
	}
	return
}
func StructScan(s interface{}, indexColumns []int) (r []interface{}) {
	if s != nil {
		maps := reflect.Indirect(reflect.ValueOf(s))
		for _, index := range indexColumns {
			r = append(r, maps.Field(index).Addr().Interface())
		}
	}
	return
}
func GetColumnIndexes(modelType reflect.Type, columnsName []string, driver string) (indexes []int, err error) {
	if modelType.Kind() != reflect.Struct {
		return nil, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		column, ok := FindTag(ormTag, "column")
		if driver == DriverOracle {
			column = strings.ToUpper(column)
		}
		if ok {
			if contains(columnsName, column) {
				indexes = append(indexes, i)
			}
		}
	}
	return
}
func FindTag(tag string, key string) (string, bool) {
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
