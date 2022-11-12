package cassandra

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/core-go/auth"
	"github.com/gocql/gocql"

	"github.com/core-go/auth/oauth2"
)

type UserRepository struct {
	Session         *gocql.Session
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
	Status          *auth.UserStatusConfig
	GenderMapper    oauth2.OAuth2GenderMapper
	Schema          *oauth2.OAuth2SchemaConfig
	BuildParam      func(i int) string
}

func NewUserRepositoryByConfig(session *gocql.Session, tableName, prefix string, activatedStatus string, services []string, c oauth2.OAuth2SchemaConfig, status *auth.UserStatusConfig, options ...oauth2.OAuth2GenderMapper) *UserRepository {
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
	m := &UserRepository{
		Session:         session,
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

func NewUserRepository(session *gocql.Session, tableName, prefix, activatedStatus string, services []string, pictureName, displayName, givenName, familyName, middleName, genderName string, status *auth.UserStatusConfig, options ...oauth2.OAuth2GenderMapper) *UserRepository {
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

	m := &UserRepository{
		Session:         session,
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
	userId := ""
	statusUser := ""
	queryString := (`SELECT %s, %s FROM %s WHERE %s = ? ALLOW FILTERING`)
	queryUserName := fmt.Sprintf(queryString, s.Schema.UserId, s.Schema.Status, s.TableName, s.Schema.UserName)

	session := s.Session
	resultUserName := session.Query(queryUserName, email)
	for _, _ = range resultUserName.Iter().Columns() {
		// New map each iteration
		row := make(map[string]interface{})
		if !resultUserName.Iter().MapScan(row) {
			break
		}
		// Do things with row
		if userIdRow, ok := row[s.Schema.UserId]; ok {
			userId = userIdRow.(string)
		}
		if statusDb, ok := row[s.Schema.Status]; ok {
			statusUser = statusDb.(string)
		}
	}
	if len(userId) <= 0 {
		queryEmail := fmt.Sprintf(queryString, s.Schema.UserId, s.Schema.Status, s.TableName, s.Schema.Email)
		resultEmail := session.Query(queryEmail, email)

		for _, _ = range resultUserName.Iter().Columns() {
			// New map each iteration
			row := make(map[string]interface{})
			if !resultEmail.Iter().MapScan(row) {
				break
			}
			// Do things with row
			if userIdRow, ok := row[s.Schema.UserId]; ok {
				userId = userIdRow.(string)
			}
			if statusDb, ok := row[s.Schema.Status]; ok {
				statusUser = statusDb.(string)
			}
		}
	}
	if len(userId) <= 0 {
		queryOAuth2Email := fmt.Sprintf(queryString, s.Schema.UserId, s.Schema.Status, s.TableName, s.Prefix+s.Schema.Email)
		resultqueryOAuth2Email := session.Query(queryOAuth2Email, email)
		for _, _ = range resultUserName.Iter().Columns() {
			// New map each iteration
			row := make(map[string]interface{})
			if !resultqueryOAuth2Email.Iter().MapScan(row) {
				break
			}
			// Do things with row
			if userIdRow, ok := row[s.Schema.UserId]; ok {
				userId = userIdRow.(string)
			}
			if statusDb, ok := row[s.Schema.Status]; ok {
				statusUser = statusDb.(string)
			}
		}
	}
	disable := false
	suspended := false
	if s.Status != nil {
		status := statusUser
		if status == s.Status.Disable {
			disable = true
		}
		if status == s.Status.Suspended {
			suspended = true
		}
	}
	return userId, disable, suspended, nil
}

func (s *UserRepository) Update(ctx context.Context, id, email, account string) (bool, error) {
	session := s.Session
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
	query, values := BuildUpdate(s.TableName, user, s.Schema.UserId, id, "?")
	result := session.Query(query, values...)
	if result.Exec() != nil {
		return false, result.Exec()
	}
	r := result.Attempts()
	// if err2 != nil {
	// 	return false, err2
	// }
	return r > 0, nil
}

func (s *UserRepository) Insert(ctx context.Context, id string, personInfo oauth2.User) (bool, error) {
	session := s.Session
	user := s.userToMap(ctx, id, personInfo)
	query, values := BuildQuery(s.TableName, user)
	result := session.Query(query, values...)
	if result.Exec() != nil {
		return false, result.Exec()
	}
	return false, nil
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

func BuildQuery(tableName string, user map[string]interface{}) (string, []interface{}) {
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
		arrValue = append(arrValue, "?")
	}
	value := fmt.Sprintf("(%v)", strings.Join(arrValue, ","))
	return fmt.Sprintf("INSERT INTO %v %v VALUES %v", tableName, column, value), values
}

func BuildUpdate(table string, model map[string]interface{}, idname string, id interface{}, buildParam string) (string, []interface{}) {
	var values []interface{}
	querySet := make([]string, 0)
	for colName, v2 := range model {
		values = append(values, v2)
		querySet = append(querySet, fmt.Sprintf("%v="+buildParam, colName))
	}
	values = append(values, id)
	queryWhere := fmt.Sprintf(" %s = %s",
		idname,
		buildParam,
	)
	query := fmt.Sprintf("update %v set %v where %v", table, strings.Join(querySet, ","), queryWhere)
	return query, values
}
