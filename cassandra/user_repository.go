package cassandra

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	a "github.com/core-go/auth"
	"github.com/gocql/gocql"
)

type UserRepository struct {
	Session                 *gocql.Session
	userTableName           string
	passwordTableName       string
	CheckTwoFactors         func(ctx context.Context, id string) (bool, error)
	activatedStatus         interface{}
	Status                  a.UserStatusConfig
	IdName                  string
	UserName                string
	UserId                  string
	SuccessTimeName         string
	FailTimeName            string
	FailCountName           string
	LockedUntilTimeName     string
	StatusName              string
	PasswordChangedTimeName string
	PasswordName            string
	ContactName             string
	EmailName               string
	PhoneName               string
	DisplayNameName         string
	MaxPasswordAgeName      string
	UserTypeName            string
	AccessDateFromName      string
	AccessDateToName        string
	AccessTimeFromName      string
	AccessTimeToName        string
	TwoFactorsName          string
}

func NewUserRepositoryByConfig(session *gocql.Session, userTableName, passwordTableName string, activatedStatus string, status a.UserStatusConfig, c a.SchemaConfig, options ...func(context.Context, string) (bool, error)) *UserRepository {
	return NewUserRepository(session, userTableName, passwordTableName, activatedStatus, status, c.Id, c.Username, c.UserId, c.SuccessTime, c.FailTime, c.FailCount, c.LockedUntilTime, c.Status, c.PasswordChangedTime, c.Password, c.Contact, c.Email, c.Phone, c.DisplayName, c.MaxPasswordAge, c.UserType, c.AccessDateFrom, c.AccessDateTo, c.AccessTimeFrom, c.AccessTimeTo, c.TwoFactors, options...)
}

func NewUserRepository(session *gocql.Session, userTableName, passwordTableName string, activatedStatus string, status a.UserStatusConfig, idName, userName, userID, successTimeName, failTimeName, failCountName, lockedUntilTimeName, statusName, passwordChangedTimeName, passwordName, contactName, emailName, phoneName, displayNameName, maxPasswordAgeName, userTypeName, accessDateFromName, accessDateToName, accessTimeFromName, accessTimeToName, twoFactorsName string, options ...func(context.Context, string) (bool, error)) *UserRepository {
	var checkTwoFactors func(context.Context, string) (bool, error)
	if len(options) >= 1 {
		checkTwoFactors = options[0]
	}
	return &UserRepository{
		Session:                 session,
		userTableName:           strings.ToLower(userTableName),
		passwordTableName:       strings.ToLower(passwordTableName),
		CheckTwoFactors:         checkTwoFactors,
		activatedStatus:         strings.ToLower(activatedStatus),
		Status:                  status,
		IdName:                  strings.ToLower(idName),
		UserName:                strings.ToLower(userName),
		UserId:                  strings.ToLower(userID),
		SuccessTimeName:         strings.ToLower(successTimeName),
		FailTimeName:            strings.ToLower(failTimeName),
		FailCountName:           strings.ToLower(failCountName),
		LockedUntilTimeName:     strings.ToLower(lockedUntilTimeName),
		StatusName:              strings.ToLower(statusName),
		PasswordChangedTimeName: strings.ToLower(passwordChangedTimeName),
		PasswordName:            strings.ToLower(passwordName),
		ContactName:             strings.ToLower(contactName),
		EmailName:               strings.ToLower(emailName),
		PhoneName:               strings.ToLower(phoneName),
		DisplayNameName:         strings.ToLower(displayNameName),
		MaxPasswordAgeName:      strings.ToLower(maxPasswordAgeName),
		UserTypeName:            strings.ToLower(userTypeName),
		AccessDateFromName:      strings.ToLower(accessDateFromName),
		AccessDateToName:        strings.ToLower(accessDateToName),
		AccessTimeFromName:      strings.ToLower(accessTimeFromName),
		AccessTimeToName:        strings.ToLower(accessTimeToName),
		TwoFactorsName:          strings.ToLower(twoFactorsName),
	}
}

func (r *UserRepository) GetUser(ctx context.Context, username string) (*a.UserInfo, error) {
	session := r.Session
	userInfo := a.UserInfo{}
	query := "SELECT * FROM " + r.userTableName + " WHERE " + r.UserName + " = ? ALLOW FILTERING"
	raws := session.Query(query, username).Iter()
	userInfo.Username = username
	for {
		// New map each iteration
		row := make(map[string]interface{})
		if !raws.MapScan(row) {
			break
		}
		// Do things with row
		if id, ok := row[r.UserId]; ok {
			userInfo.Id = id.(string)
		}
		if len(r.StatusName) > 0 {
			if status, ok := row[r.StatusName]; ok {
				r.StatusName = status.(string)
			}
		}
		if len(r.ContactName) > 0 {
			if contact, ok := row[r.ContactName]; ok {
				s := contact.(string)
				userInfo.Contact = &s
			}
		}
		if len(r.EmailName) > 0 {
			if email, ok := row[r.EmailName]; ok {
				s := email.(string)
				userInfo.Email = &s
			}
		}
		if len(r.PhoneName) > 0 {
			if phone, ok := row[r.PhoneName]; ok {
				s := phone.(string)
				userInfo.Phone = &s
			}
		}
		if len(r.DisplayNameName) > 0 {
			if displayName, ok := row[r.DisplayNameName]; ok {
				s := displayName.(string)
				userInfo.DisplayName = &s
			}
		}
		if len(r.MaxPasswordAgeName) > 0 {
			if maxPasswordAgeName, ok := row[r.MaxPasswordAgeName]; ok {
				i := int32(maxPasswordAgeName.(int))
				userInfo.MaxPasswordAge = &i
			}
		}
		if len(r.UserTypeName) > 0 {
			if userType, ok := row[r.UserTypeName]; ok {
				s := userType.(string)
				userInfo.UserType = &s
			}
		}
		if len(r.AccessDateFromName) > 0 {
			if accessDateFrom, ok := row[r.AccessDateFromName]; ok {
				userInfo.AccessDateFrom = accessDateFrom.(*time.Time)
			}
		}
		if len(r.AccessDateToName) > 0 {
			if accessDateTo, ok := row[r.AccessDateToName]; ok {
				userInfo.AccessDateTo = accessDateTo.(*time.Time)
			}
		}
		if len(r.AccessTimeFromName) > 0 {
			if accessTimeFrom, ok := row[r.AccessTimeFromName]; ok {
				userInfo.AccessTimeFrom = accessTimeFrom.(*time.Time)
			}
		}
		if len(r.AccessTimeToName) > 0 {
			if accessTimeTo, ok := row[r.AccessTimeToName]; ok {
				userInfo.AccessTimeTo = accessTimeTo.(*time.Time)
			}
		}
	}
	queryPasswordTable := "Select * From " + r.passwordTableName + " WHERE userid = ? ALLOW FILTERING"
	rawPassword := session.Query(queryPasswordTable, userInfo.Id).Iter()
	for {
		row := make(map[string]interface{})
		if !rawPassword.MapScan(row) {
			break
		}
		if len(r.PasswordName) > 0 {
			if pass, ok := row[r.PasswordName]; ok {
				userInfo.Password = pass.(string)
			}
		}
		if len(r.LockedUntilTimeName) > 0 {
			if lockedUntilTime, ok := row[r.LockedUntilTimeName]; ok {
				a1 := lockedUntilTime.(time.Time)
				userInfo.LockedUntilTime = &a1
			}
		}
		if len(r.SuccessTimeName) > 0 {
			if successTime, ok := row[r.SuccessTimeName]; ok {
				a2 := successTime.(time.Time)
				userInfo.SuccessTime = &a2
			}
		}
		if len(r.FailTimeName) > 0 {
			if failTime, ok := row[r.FailTimeName]; ok {
				a3 := failTime.(time.Time)
				userInfo.FailTime = &a3
			}
		}

		if len(r.FailCountName) > 0 {
			if failCountName, ok := row[r.FailCountName]; ok {
				i := failCountName.(int)
				userInfo.FailCount = &i
			}
		}

		if len(r.PasswordChangedTimeName) > 0 {
			if passwordChangedTime, ok := row[r.PasswordChangedTimeName]; ok {
				a4 := passwordChangedTime.(time.Time)
				userInfo.PasswordChangedTime = &a4
			}
		}
	}
	return &userInfo, nil
}

func (r *UserRepository) Pass(ctx context.Context, userId string, deactivated *bool) error {
	_, err := r.passAuthenticationAndActivate(ctx, userId, deactivated)
	return err
}
func (r *UserRepository) passAuthenticationAndActivate(ctx context.Context, userId string, updateStatus *bool) (int64, error) {
	session := r.Session
	if len(r.SuccessTimeName) == 0 && len(r.FailCountName) == 0 && len(r.LockedUntilTimeName) == 0 {
		if updateStatus != nil && !*updateStatus {
			return 0, nil
		} else if len(r.StatusName) == 0 {
			return 0, nil
		}
	}
	pass := make(map[string]interface{})
	if len(r.SuccessTimeName) > 0 {
		pass[r.SuccessTimeName] = time.Now()
	}
	if len(r.FailCountName) > 0 {
		pass[r.FailCountName] = 0
	}
	if len(r.LockedUntilTimeName) > 0 {
		pass[r.LockedUntilTimeName] = nil
	}
	query := map[string]interface{}{
		r.IdName: userId,
	}
	if updateStatus != nil && !*updateStatus {
		return patch(ctx, session, r.passwordTableName, pass, query)
	}

	if r.userTableName == r.passwordTableName {
		pass[r.StatusName] = r.activatedStatus
		return patch(ctx, session, r.passwordTableName, pass, query)
	}

	k1, err := patch(ctx, session, r.passwordTableName, pass, query)
	if err != nil {
		return k1, err
	}

	user := make(map[string]interface{})
	user[r.IdName] = userId
	user[r.StatusName] = r.activatedStatus
	k2, err1 := patch(ctx, session, r.userTableName, user, query)
	return k1 + k2, err1
}

func (r *UserRepository) Fail(ctx context.Context, userId string, failCount *int, lockedUntil *time.Time) error {
	if len(r.FailTimeName) == 0 && len(r.FailCountName) == 0 && len(r.LockedUntilTimeName) == 0 {
		return nil
	}
	pass := make(map[string]interface{})
	pass[r.IdName] = userId
	if len(r.FailTimeName) > 0 {
		pass[r.FailTimeName] = time.Now()
	}
	if failCount != nil && len(r.FailCountName) > 0 {
		pass[r.FailCountName] = *failCount + 1
		if len(r.LockedUntilTimeName) > 0 {
			pass[r.LockedUntilTimeName] = lockedUntil
		}
	}
	query := map[string]interface{}{
		r.IdName: userId,
	}
	_, err := patch(ctx, r.Session, r.passwordTableName, pass, query)
	return err
}

func patch(ctx context.Context, session *gocql.Session, table string, model map[string]interface{}, query map[string]interface{}) (int64, error) {
	keyUpdate := ""
	keyValue := ""
	for k, v := range query {
		keyUpdate = k
		keyValue = fmt.Sprintf("%v", v)
	}
	str := "SELECT * FROM " + table + " WHERE " + keyUpdate + " = ? ALLOW FILTERING"
	rows := session.Query(str, keyValue).Iter()
	for k, _ := range model {
		flag := false
		for row := range rows.Columns() {
			if rows.Columns()[row].Name == k {
				flag = true
			}
		}
		if !flag {
			if k == "failtime" || k == "lockeduntiltime" || k == "successtime" {
				queryAddCol := "ALTER TABLE " + table + " ADD " + k + " timestamp"
				er0 := session.Query(queryAddCol).Exec()
				if er0 != nil {
					return 0, er0
				}
			} else {
				queryAddCol := "ALTER TABLE " + table + " ADD " + k + " int"
				er0 := session.Query(queryAddCol).Exec()
				if er0 != nil {
					return 0, er0
				}
			}
		}
	}
	objectUpdate := make([]string, 0)
	objectUpdateValue := make([]interface{}, 0)
	for k, v := range model {
		objectUpdate = append(objectUpdate, fmt.Sprintf("%s = ? ", k))
		objectUpdateValue = append(objectUpdateValue, v)
	}
	for k, v := range query {
		keyUpdate = k
		keyValue = fmt.Sprintf("'%v'", v)
	}
	strSql := `UPDATE ` + table + ` SET ` + strings.Join(objectUpdate, ",") + ` WHERE ` + keyUpdate + " = " + keyValue
	result := session.Query(strSql, objectUpdateValue...)
	if result.Exec() != nil {
		log.Println(result.Exec())
		return 0, result.Exec()
	}
	return 1, nil
}
