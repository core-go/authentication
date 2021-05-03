package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func query(ctx context.Context, db *sql.DB, results interface{}, sql string, values ...interface{}) ([]string, error) {
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

	fieldsIndex, er3 := getColumnIndexes(modelType)
	if er3 != nil {
		return columns, er3
	}

	tb, er4 := scans(rows, modelType, fieldsIndex)
	if er4 != nil {
		return columns, er4
	}
	for _, element := range tb {
		appendToArray(results, element)
	}
	er5 := rows.Close()
	if er5 != nil {
		return columns, er5
	}
	// Rows.Err will report the last error encountered by Rows.Scan.
	if er6 := rows.Err(); er6 != nil {
		return columns, er6
	}
	return columns, nil
}
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
func scans(rows *sql.Rows, modelType reflect.Type, fieldsIndex map[string]int) (t []interface{}, err error) {
	columns, er0 := getColumns(rows.Columns())
	if er0 != nil {
		return nil, er0
	}
	for rows.Next() {
		initModel := reflect.New(modelType).Interface()
		r, swapValues := structScan(initModel, columns, fieldsIndex, -1)
		if err = rows.Scan(r...); err == nil {
			swapValuesToBool(initModel, &swapValues)
			t = append(t, initModel)
		}
	}
	return
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
func structScan(s interface{}, columns []string, fieldsIndex map[string]int, indexIgnore int) (r []interface{}, swapValues map[int]interface{}) {
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
			tagBool := modelField.Tag.Get("true")
			if tagBool == "" {
				r = append(r, valueField.Addr().Interface())
			} else {
				var str string
				swapValues[index] = reflect.New(reflect.TypeOf(str)).Elem().Addr().Interface()
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
		for index, element := range (*swap) {
			var isBool bool
			boolStr := modelType.Field(index).Tag.Get("true")
			var dbValue = element.(*string)
			isBool = *dbValue == boolStr
			if maps.Field(index).Kind() == reflect.Ptr {
				maps.Field(index).Set(reflect.ValueOf(&isBool))
			} else {
				maps.Field(index).SetBool(isBool)
			}
		}
	}
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
