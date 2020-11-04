package auth

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"strings"
)

type SqlPrivilegesLoader struct {
	DB             *sql.DB
	Query          string
	ParameterCount int
}

func NewSqlPrivilegesLoader(db *sql.DB, query string, parameterCount int) *SqlPrivilegesLoader {
	return &SqlPrivilegesLoader{DB: db, Query: query, ParameterCount: parameterCount}
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
func getColumnIndexes(modelType reflect.Type, columnsName []string) (indexes []int, err error) {
	if modelType.Kind() != reflect.Struct {
		return nil, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		column, ok := findTag(ormTag, "column")
		if ok {
			if contains(columnsName, column) {
				indexes = append(indexes, i)
			}
		}
	}
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
