package repo

import (
	"context"
	"database/sql"
	"fmt"
	auth "github.com/core-go/authentication"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type AuthenticationRepository struct {
	db                      *sql.DB
	BuildParam              func(i int) string
	userTableName           string
	passwordTableName       string
	CheckTwoFactors         func(ctx context.Context, id string) (bool, error)
	activatedStatus         interface{}
	Status                  auth.UserStatusConfig
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

func NewAuthenticationRepositoryByConfig(db *sql.DB, buildParam func(i int) string, userTableName, passwordTableName string, activatedStatus string, status auth.UserStatusConfig, c auth.SchemaConfig, options ...func(context.Context, string) (bool, error)) *AuthenticationRepository {
	return NewAuthenticationRepository(db, buildParam, userTableName, passwordTableName, activatedStatus, status, c.Id, c.Username, c.UserId, c.SuccessTime, c.FailTime, c.FailCount, c.LockedUntilTime, c.Status, c.PasswordChangedTime, c.Password, c.Contact, c.Email, c.Phone, c.DisplayName, c.MaxPasswordAge, c.UserType, c.AccessDateFrom, c.AccessDateTo, c.AccessTimeFrom, c.AccessTimeTo, c.TwoFactors, options...)
}

func NewAuthenticationRepository(db *sql.DB, buildParam func(i int) string, userTableName, passwordTableName string, activatedStatus string, status auth.UserStatusConfig, idName, userName, userID, successTimeName, failTimeName, failCountName, lockedUntilTimeName, statusName, passwordChangedTimeName, passwordName, contactName, emailName, phoneName, displayNameName, maxPasswordAgeName, userTypeName, accessDateFromName, accessDateToName, accessTimeFromName, accessTimeToName, twoFactorsName string, options ...func(context.Context, string) (bool, error)) *AuthenticationRepository {
	var checkTwoFactors func(context.Context, string) (bool, error)
	if len(options) >= 1 {
		checkTwoFactors = options[0]
	}
	var b = buildParam
	if b == nil {
		b = getBuild(db)
	}
	return &AuthenticationRepository{
		db:                      db,
		BuildParam:              b,
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

func (r *AuthenticationRepository) GetUserInfo(ctx context.Context, userid string) (*auth.UserInfo, error) {
	userInfo := auth.UserInfo{}
	strSQL := ""
	if len(r.StatusName) > 0 {
		strSQL += r.StatusName + ", "
	}
	if len(r.UserId) > 0 {
		strSQL += r.UserId + ", "
	}
	if len(r.IdName) > 0 {
		strSQL += r.IdName + ", "
	}
	if len(r.UserName) > 0 {
		strSQL += "userid as " + r.UserName + ", "
	}
	if len(r.ContactName) > 0 {
		strSQL += r.ContactName + ", "
	}
	if len(r.EmailName) > 0 {
		strSQL += r.EmailName + ", "
	}
	if len(r.PhoneName) > 0 {
		strSQL += r.PhoneName + ", "
	}
	if len(r.DisplayNameName) > 0 {
		strSQL += "CONCAT(firstname, ' ', lastname) as " + r.DisplayNameName + ", "
	}

	if len(r.MaxPasswordAgeName) > 0 {
		strSQL += r.MaxPasswordAgeName + ", "
	}

	if len(r.UserTypeName) > 0 {
		strSQL += "roletype as " + r.UserTypeName + ", "
	}

	if len(r.AccessDateFromName) > 0 {
		strSQL += "datefrom as " + r.AccessDateFromName + ", "
	}
	if len(r.AccessDateToName) > 0 {
		strSQL += "dateto as " + r.AccessDateToName + ", "
	}

	if len(r.AccessTimeFromName) > 0 {
		strSQL += `CONVERT(concat(DATE_FORMAT(NOW(), "%Y-%m-%d"), ' ', timefrom), datetime) as ` + r.AccessTimeFromName + ", "
	}

	if len(r.AccessTimeToName) > 0 {
		strSQL += `CONVERT(concat(DATE_FORMAT(NOW(), "%Y-%m-%d"), ' ', timeto), datetime) as ` + r.AccessTimeToName + ", "
	}

	if len(r.PasswordName) > 0 {
		strSQL += r.PasswordName + ", "
	}

	if len(r.LockedUntilTimeName) > 0 {
		strSQL += r.LockedUntilTimeName + ", "
	}

	if len(r.SuccessTimeName) > 0 {
		strSQL += r.SuccessTimeName + ", "
	}

	if len(r.FailTimeName) > 0 {
		strSQL += r.FailTimeName + ", "
	}

	if len(r.FailCountName) > 0 {
		strSQL += r.FailCountName + ", "
	}

	if len(r.PasswordChangedTimeName) > 0 {
		strSQL += r.PasswordChangedTimeName + ", "
	}
	strSQL = strings.TrimRight(strSQL, ", ")
	if r.userTableName == r.passwordTableName {
		query := `SELECT ` + strSQL +
			` FROM ` + r.userTableName +
			` WHERE userid = ` + r.BuildParam(1) +
			` LIMIT 1`
		rows, err := r.db.QueryContext(ctx, query, userid)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			sqlScanStruct(rows, &userInfo)
		}
	} else {
		query := `SELECT ` + strSQL +
			` FROM ` + r.userTableName +
			` INNER JOIN ` + r.passwordTableName +
			` ON ` + r.passwordTableName + `.` + r.UserId + " = " + r.userTableName + "." + r.UserId +
			` WHERE ` + r.userTableName + `.` + `userid = ` + r.BuildParam(1)
		rows, err := r.db.QueryContext(ctx, query, userid)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			sqlScanStruct(rows, &userInfo)
		}
	}
	return &userInfo, nil
}

func sqlScanStruct(rows *sql.Rows, outputStruct interface{}) error {
	v := reflect.ValueOf(outputStruct).Elem()
	if v.Kind() != reflect.Struct {
		return nil // bail if it's not a struct
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	countColumn := len(cols)
	values := make([]interface{}, countColumn)
	valuePtrs := make([]interface{}, countColumn)

	for i, _ := range valuePtrs {
		valuePtrs[i] = &values[i]
	}

	if err := rows.Scan(valuePtrs...); err != nil {
		return err
	}

	valueMap := make(map[string]interface{})
	for id, colName := range cols {
		val := values[id]
		if val != nil {
			if b, ok := val.([]byte); ok {
				valueMap[colName] = string(b)
			} else {
				valueMap[colName] = val
			}
		} else {
			valueMap[colName] = nil
		}
	}

	n := v.NumField() // number of fields in struct

	for i := 0; i < n; i = i + 1 {
		if !v.Field(i).CanSet() {
			continue
		}

		var fieldValue interface{}

		if fV, ok := valueMap[v.Type().Field(i).Tag.Get("sql")]; ok {
			fieldValue = fV
		} else if fV, ok := valueMap[string(v.Type().Field(i).Tag)]; ok {
			fieldValue = fV
		} else if fV, ok := valueMap[v.Type().Field(i).Name]; ok {
			fieldValue = fV
		} else {
			continue
		}

		if fieldValue == nil {
			continue
		}

		f := v.Field(i)
		switch f.Kind() {
		case reflect.String:
			v.Field(i).SetString(fmt.Sprintf("%v", fieldValue))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intValue, _ := strconv.ParseInt(fmt.Sprintf("%v", fieldValue), 10, 64)
			v.Field(i).SetInt(intValue)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			uintValue, _ := strconv.ParseUint(fmt.Sprintf("%v", fieldValue), 10, 64)
			v.Field(i).SetUint(uintValue)
		case reflect.Float64, reflect.Float32:
			floatValue, _ := strconv.ParseFloat(fmt.Sprintf("%v", fieldValue), 64)
			v.Field(i).SetFloat(floatValue)
		case reflect.Bool:
			boolValue, _ := strconv.ParseBool(fmt.Sprintf("%v", fieldValue))
			v.Field(i).SetBool(boolValue)
		case reflect.Ptr:
			switch f.Type() {
			case reflect.PtrTo(reflect.TypeOf(string(""))):
				svalue := fmt.Sprintf("%v", fieldValue)
				v.Field(i).Set(reflect.ValueOf(&svalue))
			case reflect.PtrTo(reflect.TypeOf(int(0))):
				int64Value, _ := strconv.ParseInt(fmt.Sprintf("%v", fieldValue), 10, 64)
				intValue := int(int64Value)
				v.Field(i).Set(reflect.ValueOf(&intValue))
			case reflect.PtrTo(reflect.TypeOf(int8(0))):
				int64Value, _ := strconv.ParseInt(fmt.Sprintf("%v", fieldValue), 10, 64)
				int8Value := int8(int64Value)
				v.Field(i).Set(reflect.ValueOf(&int8Value))
			case reflect.PtrTo(reflect.TypeOf(int16(0))):
				int64Value, _ := strconv.ParseInt(fmt.Sprintf("%v", fieldValue), 10, 64)
				int16Value := int16(int64Value)
				v.Field(i).Set(reflect.ValueOf(&int16Value))
			case reflect.PtrTo(reflect.TypeOf(int32(0))):
				int64Value, _ := strconv.ParseInt(fmt.Sprintf("%v", fieldValue), 10, 64)
				int32Value := int32(int64Value)
				v.Field(i).Set(reflect.ValueOf(&int32Value))
			case reflect.PtrTo(reflect.TypeOf(int64(0))):
				int64Value, _ := strconv.ParseInt(fmt.Sprintf("%v", fieldValue), 10, 64)
				v.Field(i).Set(reflect.ValueOf(&int64Value))
			case reflect.PtrTo(reflect.TypeOf(uint(0))):
				uint64Value, _ := strconv.ParseUint(fmt.Sprintf("%v", fieldValue), 10, 64)
				uintValue := uint(uint64Value)
				v.Field(i).Set(reflect.ValueOf(&uintValue))
			case reflect.PtrTo(reflect.TypeOf(uint8(0))):
				uint64Value, _ := strconv.ParseUint(fmt.Sprintf("%v", fieldValue), 10, 64)
				uint8Value := uint8(uint64Value)
				v.Field(i).Set(reflect.ValueOf(&uint8Value))
			case reflect.PtrTo(reflect.TypeOf(uint16(0))):
				uint64Value, _ := strconv.ParseUint(fmt.Sprintf("%v", fieldValue), 10, 64)
				uint16Value := uint16(uint64Value)
				v.Field(i).Set(reflect.ValueOf(&uint16Value))
			case reflect.PtrTo(reflect.TypeOf(uint32(0))):
				uint64Value, _ := strconv.ParseUint(fmt.Sprintf("%v", fieldValue), 10, 64)
				uint32Value := uint32(uint64Value)
				v.Field(i).Set(reflect.ValueOf(&uint32Value))
			case reflect.PtrTo(reflect.TypeOf(uint64(0))):
				uint64Value, _ := strconv.ParseUint(fmt.Sprintf("%v", fieldValue), 10, 64)
				v.Field(i).Set(reflect.ValueOf(&uint64Value))
			case reflect.PtrTo(reflect.TypeOf(float32(0))):
				float64Value, _ := strconv.ParseFloat(fmt.Sprintf("%v", fieldValue), 64)
				float32Value := float32(float64Value)
				v.Field(i).Set(reflect.ValueOf(&float32Value))
			case reflect.PtrTo(reflect.TypeOf(float64(0))):
				float64Value, _ := strconv.ParseFloat(fmt.Sprintf("%v", fieldValue), 64)
				v.Field(i).Set(reflect.ValueOf(&float64Value))
			case reflect.PtrTo(reflect.TypeOf(bool(false))):
				boolValue, _ := strconv.ParseBool(fmt.Sprintf("%v", fieldValue))
				v.Field(i).Set(reflect.ValueOf(&boolValue))
			}
		default:
		}
	}

	return nil
}

func (r *AuthenticationRepository) PassAuthentication(ctx context.Context, userId string) (int64, error) {
	return r.passAuthenticationAndActivate(ctx, userId, false)
}
func (r *AuthenticationRepository) PassAuthenticationAndActivate(ctx context.Context, userId string) (int64, error) {
	return r.passAuthenticationAndActivate(ctx, userId, true)
}

func (r *AuthenticationRepository) passAuthenticationAndActivate(ctx context.Context, userId string, updateStatus bool) (int64, error) {
	if len(r.SuccessTimeName) == 0 && len(r.FailCountName) == 0 && len(r.LockedUntilTimeName) == 0 {
		if !updateStatus {
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
	if !updateStatus {
		return patch(ctx, r.db, r.passwordTableName, pass, query)
	}

	if r.userTableName == r.passwordTableName {
		pass[r.StatusName] = r.activatedStatus
		return patch(ctx, r.db, r.passwordTableName, pass, query)
	}

	k1, err := patch(ctx, r.db, r.passwordTableName, pass, query)
	if err != nil {
		return k1, err
	}

	user := make(map[string]interface{})
	user[r.IdName] = userId
	user[r.StatusName] = r.activatedStatus
	k2, err1 := patch(ctx, r.db, r.userTableName, user, query)
	return k1 + k2, err1
}

func (r *AuthenticationRepository) WrongPassword(ctx context.Context, userId string, failCount int, lockedUntil *time.Time) error {
	if len(r.FailTimeName) == 0 && len(r.FailCountName) == 0 && len(r.LockedUntilTimeName) == 0 {
		return nil
	}
	pass := make(map[string]interface{})
	pass[r.IdName] = userId
	if len(r.FailTimeName) > 0 {
		pass[r.FailTimeName] = time.Now()
	}
	if len(r.FailCountName) > 0 {
		pass[r.FailCountName] = failCount
		if len(r.LockedUntilTimeName) > 0 {
			pass[r.LockedUntilTimeName] = lockedUntil
		}
	}
	query := map[string]interface{}{
		r.IdName: userId,
	}
	_, err := patch(ctx, r.db, r.passwordTableName, pass, query)
	return err
}

func getTime(accessTime string) *time.Time {
	const LAYOUT = "2006-01-02T15:04"
	if len(accessTime) > 0 {
		today := time.Now()
		location := time.Now().Location()
		x := today.Format("2006-01-02") + "T" + accessTime
		t, e := time.ParseInLocation(LAYOUT, x, location)
		if e == nil {
			return &t
		}
	}
	return nil
}

/*
	func patch(db *gorm.DB, table string, model map[string]interface{}, query map[string]interface{}) (int64, error) {
		result := db.Table(table).Where(query).Updates(model)
		if err := result.Error; err != nil {
			return result.RowsAffected, err
		}
		return result.RowsAffected, nil
	}
*/
func patch(ctx context.Context, db *sql.DB, table string, model map[string]interface{}, query map[string]interface{}) (int64, error) {
	objectUpdate := ""
	objectUpdateValue := ""
	keyUpdate := ""
	keyValue := ""
	for k, v := range model {
		objectUpdate = k
		objectUpdateValue = fmt.Sprintf("%v", v)
	}
	for k, v := range query {
		keyUpdate = k
		keyValue = fmt.Sprintf("%v", v)
	}
	strSql := `UPDATE ` + table + `
 SET ` + objectUpdate + " = " + objectUpdateValue +
		` WHERE ` + keyUpdate + " = " + keyValue
	result, err := db.ExecContext(ctx, strSql)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
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
