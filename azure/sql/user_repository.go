package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/core-go/authentication/azure"
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
	ActivatedStatus string
	Schema          *SchemaConfig
	BuildParam      func(int) string
}

func NewUserRepositoryByConfig(db *sql.DB, tableName, activatedStatus string, c SchemaConfig) *UserRepository {
	c.Id = strings.ToLower(c.Id)
	c.Username = strings.ToLower(c.Username)
	c.PrincipalName = strings.ToLower(c.PrincipalName)
	c.Status = strings.ToLower(c.Status)
	c.DisplayName = strings.ToLower(c.DisplayName)
	c.GivenName = strings.ToLower(c.GivenName)
	c.Surname = strings.ToLower(c.Surname)
	c.CreatedTime = strings.ToLower(c.CreatedTime)
	c.CreatedBy = strings.ToLower(c.CreatedBy)
	c.UpdatedTime = strings.ToLower(c.UpdatedTime)
	c.UpdatedBy = strings.ToLower(c.UpdatedBy)
	c.Version = strings.ToLower(c.Version)

	if len(c.Username) == 0 {
		c.Username = "username"
	}
	if len(c.PrincipalName) == 0 {
		c.PrincipalName = "email"
	}
	if len(c.Status) == 0 {
		c.Status = "status"
	}
	build := getBuild(db)
	driver := getDriver(db)
	m := &UserRepository{
		DB:              db,
		BuildParam:      build,
		Driver:          driver,
		TableName:       tableName,
		ActivatedStatus: activatedStatus,
		Schema:          &c,
	}
	return m
}

func NewUserRepository(db *sql.DB, tableName, activatedStatus string, displayName, givenName, surname string) *UserRepository {
	displayName = strings.ToLower(displayName)
	givenName = strings.ToLower(givenName)
	surname = strings.ToLower(surname)

	build := getBuild(db)
	driver := getDriver(db)
	m := &UserRepository{
		DB:              db,
		BuildParam:      build,
		Driver:          driver,
		TableName:       tableName,
		ActivatedStatus: activatedStatus,
	}
	if len(displayName) > 0 || len(givenName) > 0 || len(surname) > 0 {
		c := &SchemaConfig{}
		c.DisplayName = displayName
		c.GivenName = givenName
		c.Surname = surname
		c.Status = "status"
		c.Username = "username"
		c.PrincipalName = "email"
		m.Schema = c
	}
	return m
}

func (s *UserRepository) Exist(ctx context.Context, id string) (bool, string, error) {
	var displayName string
	query := fmt.Sprintf(`select %s, %s from %s where %s = %s`, s.Schema.Id, s.Schema.DisplayName, s.TableName, s.Schema.Id, s.BuildParam(1))
	rows, err := s.DB.QueryContext(ctx, query, id)
	if err != nil {
		return false, displayName, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &displayName)
		if err != nil {
			return false, displayName, err
		}
		return true, displayName, nil
	}
	return false, displayName, nil
}

func (s *UserRepository) Insert(ctx context.Context, id string, personInfo *azure.AzureUser) (bool, error) {
	user := s.userToMap(id, personInfo)
	query, values := BuildQuery(s.TableName, user, s.BuildParam)
	_, err := s.DB.ExecContext(ctx, query, values...)
	if err != nil {
		return handleDuplicate(s.Driver, err)
	}
	return true, err
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

func (s *UserRepository) userToMap(id string, user *azure.AzureUser) map[string]interface{} {
	userMap := UserToMap(id, user, s.Schema)
	userMap[s.Schema.Id] = id
	if len(s.Schema.Username) > 0 {
		userMap[s.Schema.Username] = user.UserPrincipalName
	}
	if len(s.Schema.PrincipalName) > 0 {
		userMap[s.Schema.PrincipalName] = user.UserPrincipalName
	}
	userMap[s.Schema.Status] = s.ActivatedStatus
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
		arrValue = append(arrValue, buildParam(i+1))
	}
	value := fmt.Sprintf("(%v)", strings.Join(arrValue, ","))
	return fmt.Sprintf("insert into %v %v values %v", tableName, column, value), values
}
func UserToMap(id string, user *azure.AzureUser, c *SchemaConfig) map[string]interface{} {
	userMap := make(map[string]interface{})
	if c == nil {
		return userMap
	}

	if len(c.DisplayName) > 0 && len(user.DisplayName) > 0 {
		userMap[c.DisplayName] = user.DisplayName
	}
	if len(c.GivenName) > 0 && len(user.GivenName) > 0 {
		userMap[c.GivenName] = user.GivenName
	}
	if len(c.Surname) > 0 && len(user.Surname) > 0 {
		userMap[c.Surname] = user.Surname
	}
	if len(c.JobTitle) > 0 && len(user.JobTitle) > 0 {
		userMap[c.JobTitle] = user.JobTitle
	}
	if len(c.Language) > 0 && len(user.PreferredLanguage) > 0 {
		userMap[c.Language] = user.PreferredLanguage
	}

	now := time.Now()
	if len(c.CreatedTime) > 0 {
		userMap[c.CreatedTime] = now
	}
	if len(c.UpdatedTime) > 0 {
		userMap[c.UpdatedTime] = now
	}
	if len(c.CreatedBy) > 0 {
		userMap[c.CreatedBy] = id
	}
	if len(c.UpdatedBy) > 0 {
		userMap[c.UpdatedBy] = id
	}
	if len(c.Version) > 0 {
		userMap[c.Version] = 1
	}
	return userMap
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
