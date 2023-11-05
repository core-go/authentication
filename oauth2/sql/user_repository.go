package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/core-go/auth"
	"github.com/core-go/auth/oauth2"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	driverPostgres   = "postgres"
	driverMysql      = "mysql"
	driverMssql      = "mssql"
	driverOracle     = "oracle"
	driverSqlite3    = "sqlite3"
	driverNotSupport = "no support"
)

type UserRepository struct {
	DB              *sql.DB
	Driver          string
	TableName       string
	Prefix          string
	ActivatedStatus string
	Services        []string
	StatusName      string
	UserIdName      string
	UserName        string
	EmailName       string
	OAuth2EmailName string
	AccountName     string
	ActiveName      string

	updatedTimeName string
	updatedByName   string
	UseId           bool
	Status          *auth.UserStatusConfig
	GenderMapper    oauth2.OAuth2GenderMapper
	Schema          *oauth2.OAuth2SchemaConfig
	BuildParam      func(int) string
}

func NewUserRepositoryByConfig(db *sql.DB, tableName, prefix string, activatedStatus string, services []string, c oauth2.OAuth2SchemaConfig, status *auth.UserStatusConfig, options ...oauth2.OAuth2GenderMapper) *UserRepository {
	var genderMapper oauth2.OAuth2GenderMapper
	if len(options) >= 1 {
		genderMapper = options[0]
	}
	c.UserId = strings.ToLower(c.UserId)
	c.UserName = strings.ToLower(c.UserName)
	c.Email = strings.ToLower(c.Email)
	c.Status = strings.ToLower(c.Status)
	c.OAuth2Email = strings.ToLower(c.OAuth2Email)
	c.Account = strings.ToLower(c.Account)
	c.Active = strings.ToLower(c.Active)
	c.DisplayName = strings.ToLower(c.DisplayName)
	c.Picture = strings.ToLower(c.Picture)
	c.Locale = strings.ToLower(c.Locale)
	c.Gender = strings.ToLower(c.Gender)
	c.DateOfBirth = strings.ToLower(c.DateOfBirth)
	c.GivenName = strings.ToLower(c.GivenName)
	c.MiddleName = strings.ToLower(c.MiddleName)
	c.FamilyName = strings.ToLower(c.FamilyName)
	c.CreatedTime = strings.ToLower(c.CreatedTime)
	c.CreatedBy = strings.ToLower(c.CreatedBy)
	c.UpdatedTime = strings.ToLower(c.UpdatedTime)
	c.UpdatedBy = strings.ToLower(c.UpdatedBy)
	c.Version = strings.ToLower(c.Version)
	s := make([]string, 0)
	for _, sv := range services {
		s = append(s, strings.ToLower(sv))
	}

	if len(c.UserName) == 0 {
		c.UserName = "username"
	}
	if len(c.Email) == 0 {
		c.Email = "email"
	}
	if len(c.Status) == 0 {
		c.Status = "status"
	}
	if len(c.OAuth2Email) == 0 {
		c.OAuth2Email = "email"
	}
	if len(c.Account) == 0 {
		c.Account = "account"
	}
	if len(c.Active) == 0 {
		c.Active = "active"
	}
	build := getBuild(db)
	driver := getDriver(db)
	m := &UserRepository{
		DB:              db,
		BuildParam:      build,
		Driver:          driver,
		TableName:       tableName,
		Prefix:          prefix,
		ActivatedStatus: activatedStatus,
		Services:        s,
		GenderMapper:    genderMapper,
		Schema:          &c,
		updatedTimeName: c.UpdatedTime,
		updatedByName:   c.UpdatedBy,
		Status:          status,
	}
	return m
}

func NewUserRepository(db *sql.DB, tableName, prefix, activatedStatus string, services []string, pictureName, displayName, givenName, familyName, middleName, genderName string, status *auth.UserStatusConfig, options ...oauth2.OAuth2GenderMapper) *UserRepository {
	var genderMapper oauth2.OAuth2GenderMapper
	if len(options) >= 1 {
		genderMapper = options[0]
	}

	pictureName = strings.ToLower(pictureName)
	displayName = strings.ToLower(displayName)
	givenName = strings.ToLower(givenName)
	familyName = strings.ToLower(familyName)
	middleName = strings.ToLower(middleName)
	genderName = strings.ToLower(genderName)

	build := getBuild(db)
	driver := getDriver(db)
	m := &UserRepository{
		DB:              db,
		BuildParam:      build,
		Driver:          driver,
		TableName:       tableName,
		Prefix:          prefix,
		ActivatedStatus: activatedStatus,
		StatusName:      "status",
		Services:        services,
		UserName:        "username",
		EmailName:       "email",
		OAuth2EmailName: "email",
		AccountName:     "account",
		ActiveName:      "active",
		Status:          status,
		GenderMapper:    genderMapper,
	}
	if len(pictureName) > 0 || len(displayName) > 0 || len(givenName) > 0 || len(middleName) > 0 || len(familyName) > 0 || len(genderName) > 0 {
		c := &oauth2.OAuth2SchemaConfig{}
		c.Picture = pictureName
		c.DisplayName = displayName
		c.GivenName = givenName
		c.MiddleName = middleName
		c.FamilyName = familyName
		c.Gender = genderName
		m.Schema = c
	}
	return m
}

func (s *UserRepository) GetUser(ctx context.Context, email string) (string, bool, bool, error) {
	arr := make(map[string]interface{})
	columns := make([]interface{}, 0)
	values := make([]interface{}, 0)
	i := 0
	columns = append(columns, s.Schema.UserId, s.Schema.Status, s.TableName,
		s.Schema.UserName, s.BuildParam(i),
		s.Schema.Email, s.BuildParam(i+1), s.Prefix+s.Schema.OAuth2Email, s.BuildParam(i+2))
	var sel strings.Builder
	sel.WriteString(`SELECT %s, %s FROM %s WHERE `)
	var where strings.Builder
	if s.UseId {
		values = append(values, email)
		where.WriteString(`%s = %s`)
	} else {
		values = append(values, email, email, email)
		where.WriteString(`%s = %s OR %s = %s OR %s = %s`)
		i = 3
		for _, sv := range s.Services {
			if sv != s.Prefix {
				where.WriteString(` OR %s = `)
				where.WriteString(s.BuildParam(i))
				i++
				columns = append(columns, sv+s.Schema.OAuth2Email)
				values = append(values, email)
			}
		}
	}
	sel.WriteString(where.String())
	query := fmt.Sprintf(sel.String(), columns...)
	rows, err := s.DB.Query(query, values...)
	disable := false
	suspended := false
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return "", disable, suspended, nil
		}
		return "", disable, suspended, err
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		if err1 := rows.Scan(columnPointers...); err1 != nil {
			return "", disable, suspended, err1
		}

		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			arr[colName] = *val
		}
	}
	err2 := rows.Err()
	if err2 != nil {
		return "", disable, suspended, err2
	}

	if len(arr) == 0 {
		return "", disable, suspended, nil
	}
	if s.Status != nil {
		status := string(arr[s.Schema.Status].([]byte))
		if status == s.Status.Disable {
			disable = true
		}
		if status == s.Status.Suspended {
			suspended = true
		}
	}
	return string(arr[s.Schema.UserId].([]byte)), disable, suspended, nil
}

func (s *UserRepository) Update(ctx context.Context, id, email, account string) (bool, error) {
	user := make(map[string]interface{})

	user[s.Prefix+s.Schema.OAuth2Email] = email
	user[s.Prefix+s.Schema.Account] = account
	user[s.Prefix+s.Schema.Active] = true

	if len(s.updatedTimeName) > 0 {
		user[s.updatedTimeName] = time.Now()
	}
	if len(s.updatedByName) > 0 {
		user[s.updatedByName] = id
	}

	query, values := BuildUpdate(s.TableName, user, s.Schema.UserId, id, s.BuildParam)
	result, err1 := s.DB.ExecContext(ctx, query, values...)
	if err1 != nil {
		return false, err1
	}
	r, err2 := result.RowsAffected()
	if err2 != nil {
		return false, err2
	}
	return r > 0, err2
}

func (s *UserRepository) Insert(ctx context.Context, id string, personInfo oauth2.User) (bool, error) {
	user := s.userToMap(ctx, id, personInfo)
	query, values := BuildQuery(s.TableName, user, s.BuildParam)
	_, err := s.DB.ExecContext(ctx, query, values...)
	if err != nil {
		return handleDuplicate(s.Driver, err)
	}
	return false, err
}

func handleDuplicate(driver string, err error) (bool, error) {
	switch driver {
	case driverPostgres:
		if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
			return true, nil
		}
		return false, err
	case driverMysql:
		if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
			return true, nil
		}
		return false, err
	case driverMssql:
		if strings.Contains(err.Error(), "Violation of PRIMARY KEY constraint") {
			return true, nil
		}
		return false, err
	case driverOracle:
		if strings.Contains(err.Error(), "ORA-00001: unique constraint") {
			return true, nil
		}
		return false, err
	case driverSqlite3:
		if strings.Contains(err.Error(), "UNIQUE constraint failed:") {
			return true, nil
		}
		return false, err
	default:
		return false, err
	}
}

func (s *UserRepository) userToMap(ctx context.Context, id string, user oauth2.User) map[string]interface{} {
	userMap := oauth2.UserToMap(ctx, id, user, s.GenderMapper, s.Schema)
	//userMap := User{}
	userMap[s.Schema.UserId] = id
	userMap[s.Schema.UserName] = user.Email
	userMap[s.Schema.Status] = s.ActivatedStatus

	userMap[s.Prefix+s.Schema.OAuth2Email] = user.Email
	userMap[s.Prefix+s.Schema.Account] = user.Account
	userMap[s.Prefix+s.Schema.Active] = true
	return userMap
}

func BuildQuery(tableName string, user map[string]interface{}, buildParam func(i int) string) (string, []interface{}) {
	var cols []string
	var values []interface{}
	for col, v := range user {
		cols = append(cols, col)
		values = append(values, v)
	}
	column := fmt.Sprintf("(%v)", strings.Join(cols, ","))
	numCol := len(cols)
	var arrValue []string
	for i := 0; i < numCol; i++ {
		arrValue = append(arrValue, buildParam(i))
	}
	value := fmt.Sprintf("(%v)", strings.Join(arrValue, ","))
	return fmt.Sprintf("INSERT INTO %v %v VALUES %v", tableName, column, value), values
}

func BuildUpdate(table string, model map[string]interface{}, idname string, id interface{}, buildParam func(i int) string) (string, []interface{}) {
	colNumber := 0
	var values []interface{}
	querySet := make([]string, 0)
	for colName, v2 := range model {
		values = append(values, v2)
		querySet = append(querySet, fmt.Sprintf("%v="+buildParam(colNumber), colName))
		colNumber++
	}
	values = append(values, id)
	queryWhere := fmt.Sprintf(" %s = %s",
		idname,
		buildParam(colNumber),
	)
	query := fmt.Sprintf("update %v set %v where %v", table, strings.Join(querySet, ","), queryWhere)
	return query, values
}

func buildParam(i int) string {
	return "?"
}
func buildOracleParam(i int) string {
	return ":val" + strconv.Itoa(i)
}
func buildMsSqlParam(i int) string {
	return "@p" + strconv.Itoa(i)
}
func buildDollarParam(i int) string {
	return "$" + strconv.Itoa(i)
}
func getBuild(db *sql.DB) func(i int) string {
	driver := reflect.TypeOf(db.Driver()).String()
	switch driver {
	case "*pq.Driver":
		return buildDollarParam
	case "*godror.drv":
		return buildOracleParam
	case "*mssql.Driver":
		return buildMsSqlParam
	default:
		return buildParam
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
