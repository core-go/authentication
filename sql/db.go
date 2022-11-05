package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strconv"
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

func getColumnIndexes(modelType reflect.Type) (map[string]int, error) {
	ma := make(map[string]int, 0)
	if modelType.Kind() != reflect.Struct {
		return ma, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		ormTag := field.Tag.Get("gorm")
		column, ok := findTag(ormTag, "column")
		column = strings.ToLower(column)
		if ok {
			ma[column] = i
		}
	}
	return ma, nil
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
func appendToArray(arr interface{}, item interface{}) interface{} {
	arrValue := reflect.ValueOf(arr)
	elemValue := reflect.Indirect(arrValue)

	itemValue := reflect.ValueOf(item)
	if itemValue.Kind() == reflect.Ptr {
		itemValue = reflect.Indirect(itemValue)
	}
	elemValue.Set(reflect.Append(elemValue, itemValue))
	return arr
}
func queryWithMap(ctx context.Context, db *sql.DB, fieldsIndex map[string]int, results interface{}, sql string, values ...interface{}) ([]string, error) {
	return queryWithMapAndArray(ctx, db, fieldsIndex, results, nil, sql, values...)
}
func queryWithMapAndArray(ctx context.Context, db *sql.DB, fieldsIndex map[string]int, results interface{}, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, sql string, values ...interface{}) ([]string, error) {
	rows, er1 := db.QueryContext(ctx, sql, values...)
	if er1 != nil {
		return nil, er1
	}
	defer rows.Close()
	columns, er2 := rows.Columns()
	if er2 != nil {
		return columns, er2
	}
	modelType := reflect.TypeOf(results).Elem().Elem()
	tb, er3 := scan(rows, modelType, fieldsIndex, toArray)
	if er3 != nil {
		return columns, er3
	}
	for _, element := range tb {
		appendToArray(results, element)
	}
	er4 := rows.Close()
	if er4 != nil {
		return columns, er4
	}
	// Rows.Err will report the last error encountered by Rows.Scan.
	if er5 := rows.Err(); er5 != nil {
		return columns, er5
	}
	return columns, nil
}
func getColumns(cols []string, err error) ([]string, error) {
	if cols == nil || err != nil {
		return cols, err
	}
	c2 := make([]string, 0)
	for _, c := range cols {
		s := strings.ToLower(c)
		c2 = append(c2, s)
	}
	return c2, nil
}
func scan(rows *sql.Rows, modelType reflect.Type, fieldsIndex map[string]int, options ...func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) (t []interface{}, err error) {
	if fieldsIndex == nil {
		fieldsIndex, err = getColumnIndexes(modelType)
		if err != nil {
			return
		}
	}
	var toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
	if len(options) > 0 {
		toArray = options[0]
	}
	columns, er0 := getColumns(rows.Columns())
	if er0 != nil {
		return nil, er0
	}
	for rows.Next() {
		initModel := reflect.New(modelType).Interface()
		r, swapValues := structScan(initModel, columns, fieldsIndex, toArray)
		if err = rows.Scan(r...); err == nil {
			swapValuesToBool(initModel, &swapValues)
			t = append(t, initModel)
		}
	}
	return
}
func structScan(s interface{}, columns []string, fieldsIndex map[string]int, options ...func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) (r []interface{}, swapValues map[int]interface{}) {
	var toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
	if len(options) > 0 {
		toArray = options[0]
	}
	return structScanAndIgnore(s, columns, fieldsIndex, toArray, -1)
}
func structScanAndIgnore(s interface{}, columns []string, fieldsIndex map[string]int, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}, indexIgnore int) (r []interface{}, swapValues map[int]interface{}) {
	if s != nil {
		modelType := reflect.TypeOf(s).Elem()
		swapValues = make(map[int]interface{}, 0)
		maps := reflect.Indirect(reflect.ValueOf(s))

		if columns == nil {
			for i := 0; i < maps.NumField(); i++ {
				tagBool := modelType.Field(i).Tag.Get("true")
				if tagBool == "" {
					r = append(r, maps.Field(i).Addr().Interface())
				} else {
					var str string
					swapValues[i] = reflect.New(reflect.TypeOf(str)).Elem().Addr().Interface()
					r = append(r, swapValues[i])
				}
			}
			return
		}

		for i, columnsName := range columns {
			if i == indexIgnore {
				continue
			}
			var index int
			var ok bool
			var modelField reflect.StructField
			var valueField reflect.Value
			if fieldsIndex == nil {
				if modelField, ok = modelType.FieldByName(columnsName); !ok {
					var t interface{}
					r = append(r, &t)
					continue
				}
				valueField = maps.FieldByName(columnsName)
			} else {
				if index, ok = fieldsIndex[columnsName]; !ok {
					var t interface{}
					r = append(r, &t)
					continue
				}
				modelField = modelType.Field(index)
				valueField = maps.Field(index)
			}
			x := valueField.Addr().Interface()
			tagBool := modelField.Tag.Get("true")
			if tagBool == "" {
				if toArray != nil && valueField.Kind() == reflect.Slice {
					x = toArray(x)
				}
				r = append(r, x)
			} else {
				var str string
				y := reflect.New(reflect.TypeOf(str))
				swapValues[index] = y.Elem().Addr().Interface()
				r = append(r, swapValues[index])
			}
		}
	}
	return
}
func swapValuesToBool(s interface{}, swap *map[int]interface{}) {
	if s != nil {
		modelType := reflect.TypeOf(s).Elem()
		maps := reflect.Indirect(reflect.ValueOf(s))
		for index, element := range *swap {
			dbValue2, ok2 := element.(*bool)
			if ok2 {
				if maps.Field(index).Kind() == reflect.Ptr {
					maps.Field(index).Set(reflect.ValueOf(dbValue2))
				} else {
					maps.Field(index).SetBool(*dbValue2)
				}
			} else {
				dbValue, ok := element.(*string)
				if ok {
					var isBool bool
					if *dbValue == "true" {
						isBool = true
					} else if *dbValue == "false" {
						isBool = false
					} else {
						boolStr := modelType.Field(index).Tag.Get("true")
						isBool = *dbValue == boolStr
					}
					if maps.Field(index).Kind() == reflect.Ptr {
						maps.Field(index).Set(reflect.ValueOf(&isBool))
					} else {
						maps.Field(index).SetBool(isBool)
					}
				}
			}
		}
	}
}
func GetBuildByDriver(driver string) func(i int) string {
	switch driver {
	case driverPostgres:
		return BuildDollarParam
	case driverOracle:
		return BuildOracleParam
	case driverMssql:
		return BuildMsSqlParam
	default:
		return BuildParam
	}
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
func BuildParam(i int) string {
	return "?"
}
func BuildOracleParam(i int) string {
	return ":" + strconv.Itoa(i)
}
func BuildMsSqlParam(i int) string {
	return "@p" + strconv.Itoa(i)
}
func BuildDollarParam(i int) string {
	return "$" + strconv.Itoa(i)
}
