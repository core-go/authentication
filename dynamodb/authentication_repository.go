package dynamodb

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/common-go/auth"
)

type AuthenticationRepository struct {
	Db                      *dynamodb.DynamoDB
	UserTableName           string
	PasswordTableName       string
	CheckTwoFactors         func(ctx context.Context, id string) (bool, error)
	ActivatedStatus         interface{}
	Status                  auth.UserStatusConfig
	UserName                string
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
	RolesName               string
	UserTypeName            string
	AccessDateFromName      string
	AccessDateToName        string
	AccessTimeFromName      string
	AccessTimeToName        string
	TwoFactorsName          string
}

func NewAuthenticationRepositoryByConfig(dynamoDB *dynamodb.DynamoDB, userTableName, passwordTableName string, activatedStatus interface{}, status auth.UserStatusConfig, c auth.SchemaConfig, options ...func(context.Context, string) (bool, error)) *AuthenticationRepository {
	return NewAuthenticationRepository(dynamoDB, userTableName, passwordTableName, activatedStatus, status, c.Username, c.SuccessTime, c.FailTime, c.FailCount, c.LockedUntilTime, c.Status, c.PasswordChangedTime, c.Password, c.Contact, c.Email, c.Phone, c.DisplayName, c.MaxPasswordAge, c.Roles, c.UserType, c.AccessDateFrom, c.AccessDateTo, c.AccessTimeFrom, c.AccessTimeTo, c.TwoFactors, options...)
}

func NewAuthenticationRepository(dynamoDB *dynamodb.DynamoDB, userTableName, passwordTableName string, activatedStatus interface{}, status auth.UserStatusConfig, userName, successTimeName, failTimeName, failCountName, lockedUntilTimeName, statusName, passwordChangedTimeName, passwordName, contactName, emailName, phoneName, displayNameName, maxPasswordAgeName, rolesName, userTypeName, accessDateFromName, accessDateToName, accessTimeFromName, accessTimeToName, twoFactors string, options ...func(context.Context, string) (bool, error)) *AuthenticationRepository {
	var checkTwoFactors func(context.Context, string) (bool, error)
	if len(options) > 0 && options[0] != nil {
		checkTwoFactors = options[0]
	}
	return &AuthenticationRepository{
		Db:                      dynamoDB,
		UserTableName:           userTableName,
		PasswordTableName:       passwordTableName,
		CheckTwoFactors:         checkTwoFactors,
		ActivatedStatus:         activatedStatus,
		Status:                  status,
		UserName:                userName,
		SuccessTimeName:         successTimeName,
		FailTimeName:            failTimeName,
		FailCountName:           failCountName,
		LockedUntilTimeName:     lockedUntilTimeName,
		StatusName:              statusName,
		PasswordChangedTimeName: passwordChangedTimeName,
		PasswordName:            passwordName,
		ContactName:             contactName,
		EmailName:               emailName,
		PhoneName:               phoneName,
		DisplayNameName:         displayNameName,
		MaxPasswordAgeName:      maxPasswordAgeName,
		RolesName:               rolesName,
		UserTypeName:            userTypeName,
		AccessDateFromName:      accessDateFromName,
		AccessDateToName:        accessDateToName,
		AccessTimeFromName:      accessTimeFromName,
		AccessTimeToName:        accessTimeToName,
		TwoFactorsName:          twoFactors,
	}
}

func (r *AuthenticationRepository) GetUserInfo(ctx context.Context, username string) (*auth.UserInfo, error) {
	userInfo := auth.UserInfo{}
	filter := expression.Equal(expression.Name("_id"), expression.Value(username))
	expr, _ := expression.NewBuilder().WithFilter(filter).Build()
	query := &dynamodb.ScanInput{
		TableName:                 aws.String(r.UserTableName),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}
	output, er1 := r.Db.ScanWithContext(ctx, query)
	if er1 != nil || len(output.Items) != 1 {
		return nil, er1
	}
	raw := make(map[string]interface{})
	er1 = dynamodbattribute.UnmarshalMap(output.Items[0], &raw)
	if er1 != nil {
		return nil, er1
	}
	if len(r.StatusName) > 0 {
		rawStatus := raw[r.StatusName]
		status, ok := rawStatus.(string)
		if !ok {
			iStatus, ok2 := rawStatus.(int32)
			if ok2 {
				status = strconv.Itoa(int(iStatus))
			} else {
				bStatus, ok3 := rawStatus.(bool)
				if ok3 {
					status = strconv.FormatBool(bStatus)
				}
			}
		}
		userInfo.Deactivated = status == r.Status.Deactivated
		userInfo.Suspended = status == r.Status.Suspended
		userInfo.Disable = status == r.Status.Disable
	}

	if len(r.ContactName) > 0 {
		if contact, ok := raw[r.ContactName].(string); ok {
			userInfo.Contact = contact
		}
	}
	if len(r.EmailName) > 0 {
		if email, ok := raw[r.EmailName].(string); ok {
			userInfo.Email = email
		}
	}
	if len(r.PhoneName) > 0 {
		if phone, ok := raw[r.PhoneName].(string); ok {
			userInfo.Phone = phone
		}
	}

	if len(r.DisplayNameName) > 0 {
		if displayName, ok := raw[r.DisplayNameName].(string); ok {
			userInfo.DisplayName = displayName
		}
	}

	if len(r.MaxPasswordAgeName) > 0 {
		if maxPasswordAgeName, ok := raw[r.MaxPasswordAgeName].(int64); ok {
			userInfo.MaxPasswordAge = maxPasswordAgeName
		}
	}

	if len(r.UserTypeName) > 0 {
		if userType, ok := raw[r.UserTypeName].(string); ok {
			userInfo.UserType = userType
		}
	}

	if len(r.AccessDateFromName) > 0 {
		if accessDateFrom, ok := raw[r.AccessDateFromName].(time.Time); ok {
			userInfo.AccessDateFrom = &accessDateFrom
		}
	}

	if len(r.AccessDateToName) > 0 {
		if accessDateTo, ok := raw[r.AccessDateToName].(time.Time); ok {
			userInfo.AccessDateTo = &accessDateTo
		}
	}

	if len(r.AccessTimeFromName) > 0 {
		if accessTimeFrom, ok := raw[r.AccessTimeFromName].(time.Time); ok {
			userInfo.AccessTimeFrom = &accessTimeFrom
		} else if accessTimeFrom, ok := raw[r.AccessTimeFromName].(string); ok {
			userInfo.AccessTimeFrom = getTime(accessTimeFrom)
		}
	}

	if len(r.AccessTimeToName) > 0 {
		if accessTimeTo, ok := raw[r.AccessTimeToName].(time.Time); ok {
			userInfo.AccessTimeTo = &accessTimeTo
		} else if accessTimeTo, ok := raw[r.AccessTimeToName].(string); ok {
			userInfo.AccessTimeTo = getTime(accessTimeTo)
		}
	}

	if r.CheckTwoFactors != nil {
		id := userInfo.Id
		if len(id) == 0 {
			id = username
		}
		ok, er2 := r.CheckTwoFactors(ctx, id)
		if er2 != nil {
			return &userInfo, er2
		}
		userInfo.TwoFactors = ok
	} else if len(r.TwoFactorsName) > 0 {
		if isTwoFactor, ok := raw[r.TwoFactorsName]; ok {
			if b, k := isTwoFactor.(bool); k {
				userInfo.TwoFactors = b
			}
		}
	}

	if r.UserTableName == r.PasswordTableName {
		return r.getPasswordInfo(ctx, &userInfo, raw), nil
	}

	queryPassword := &dynamodb.ScanInput{
		TableName:                 aws.String(r.PasswordTableName),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}
	outputPassword, er3 := r.Db.ScanWithContext(ctx, queryPassword)
	if er3 != nil {
		return nil, er3
	}
	rawPassword := make(map[string]interface{})
	er4 := dynamodbattribute.UnmarshalMap(outputPassword.Items[0], &rawPassword)
	if er4 != nil {
		return nil, er4
	}
	return r.getPasswordInfo(ctx, &userInfo, rawPassword), nil
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

func (r *AuthenticationRepository) getPasswordInfo(ctx context.Context, user *auth.UserInfo, raw map[string]interface{}) *auth.UserInfo {
	if len(r.PasswordName) > 0 {
		if pass, ok := raw[r.PasswordName].(string); ok {
			user.Password = pass
		}
	}

	if len(r.LockedUntilTimeName) > 0 {
		if lockedUntilTime, ok := raw[r.LockedUntilTimeName].(time.Time); ok {
			user.LockedUntilTime = &lockedUntilTime
		}
	}

	if len(r.SuccessTimeName) > 0 {
		if successTime, ok := raw[r.SuccessTimeName].(time.Time); ok {
			user.SuccessTime = &successTime
		}
	}

	if len(r.FailTimeName) > 0 {
		if failTime, ok := raw[r.FailTimeName].(time.Time); ok {
			user.FailTime = &failTime
		}
	}

	if len(r.FailCountName) > 0 {
		if failCountName, ok := raw[r.FailCountName].(int32); ok {
			user.FailCount = int(failCountName)
		}
	}

	if len(r.PasswordChangedTimeName) > 0 {
		if passwordChangedTime, ok := raw[r.PasswordChangedTimeName].(time.Time); ok {
			user.PasswordChangedTime = &passwordChangedTime
		}
	}
	return user
}

func (r *AuthenticationRepository) PassAuthentication(ctx context.Context, userId string) (int64, error) {
	return r.passAuthenticationAndActivate(ctx, userId, false)
}

func (r *AuthenticationRepository) PassAuthenticationAndActivate(ctx context.Context, userId string) (int64, error) {
	return r.passAuthenticationAndActivate(ctx, userId, true)
}

func (r *AuthenticationRepository) passAuthenticationAndActivate(ctx context.Context, userId string, updateStatus bool) (int64, error) {
	if len(r.SuccessTimeName) == 0 && len(r.FailCountName) == 0 && len(r.LockedUntilTimeName) == 0 {
		if !updateStatus || len(r.StatusName) == 0 {
			return 0, nil
		}
	}
	pass := make(map[string]interface{})
	pass["_id"] = userId
	if len(r.SuccessTimeName) > 0 {
		pass[r.SuccessTimeName] = time.Now()
	}
	if len(r.FailCountName) > 0 {
		pass[r.FailCountName] = 0
	}
	if len(r.LockedUntilTimeName) > 0 {
		pass[r.LockedUntilTimeName] = nil
	}
	if !updateStatus {
		return upsertOne(ctx, r.Db, r.PasswordTableName, pass)
	}
	if r.UserTableName == r.PasswordTableName {
		pass[r.StatusName] = r.ActivatedStatus
		return upsertOne(ctx, r.Db, r.PasswordTableName, pass)
	}
	k1, er1 := upsertOne(ctx, r.Db, r.PasswordTableName, pass)
	if er1 != nil {
		return k1, er1
	}
	user := make(map[string]interface{})
	user["_id"] = userId
	user[r.StatusName] = r.ActivatedStatus
	k2, er2 := upsertOne(ctx, r.Db, r.UserTableName, user)
	return k1 + k2, er2
}

func (r *AuthenticationRepository) WrongPassword(ctx context.Context, userId string, failCount int, lockedUntil *time.Time) error {
	if len(r.FailTimeName) == 0 && len(r.FailCountName) == 0 && len(r.LockedUntilTimeName) == 0 {
		return nil
	}
	pass := make(map[string]interface{})
	pass["_id"] = userId
	if len(r.FailTimeName) > 0 {
		pass[r.FailTimeName] = time.Now()
	}
	if len(r.FailCountName) > 0 {
		pass[r.FailCountName] = failCount
		if len(r.LockedUntilTimeName) > 0 {
			pass[r.LockedUntilTimeName] = lockedUntil
		}
	}
	_, err := patchOne(ctx, r.Db, r.PasswordTableName, []string{"_id"}, pass)
	return err
}

func upsertOne(ctx context.Context, db *dynamodb.DynamoDB, tableName string, model interface{}) (int64, error) {
	modelMap, err := dynamodbattribute.MarshalMap(model)
	if err != nil {
		return 0, err
	}
	params := &dynamodb.PutItemInput{
		TableName:              aws.String(tableName),
		Item:                   modelMap,
		ReturnConsumedCapacity: aws.String(dynamodb.ReturnConsumedCapacityTotal),
	}
	output, err := db.PutItemWithContext(ctx, params)
	if err != nil {
		return 0, err
	}
	return int64(aws.Float64Value(output.ConsumedCapacity.CapacityUnits)), nil
}

func patchOne(ctx context.Context, db *dynamodb.DynamoDB, tableName string, keys []string, model map[string]interface{}) (int64, error) {
	idMap := map[string]interface{}{}
	for i := range keys {
		idMap[keys[i]] = model[keys[i]]
		delete(model, keys[i])
	}
	keyMap, err := buildKeyMap(keys, idMap)
	if err != nil {
		return 0, err
	}
	updateBuilder := expression.UpdateBuilder{}
	for key, value := range model {
		updateBuilder = updateBuilder.Set(expression.Name(key), expression.Value(value))
	}
	var cond expression.ConditionBuilder
	for key, value := range idMap {
		if reflect.ValueOf(cond).IsZero() {
			cond = expression.Name(key).Equal(expression.Value(value))
		}
		cond = cond.And(expression.Name(key).Equal(expression.Value(value)))
	}
	expr, _ := expression.NewBuilder().WithUpdate(updateBuilder).WithCondition(cond).Build()
	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       keyMap,
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnConsumedCapacity:    aws.String(dynamodb.ReturnConsumedCapacityTotal),
	}
	output, err := db.UpdateItemWithContext(ctx, input)
	if err != nil {
		if strings.Index(err.Error(), "ConditionalCheckFailedException:") >= 0 {
			return 0, fmt.Errorf("object not found")
		}
		return 0, err
	}
	return int64(aws.Float64Value(output.ConsumedCapacity.CapacityUnits)), nil
}

func buildKeyMap(keys []string, value interface{}) (map[string]*dynamodb.AttributeValue, error) {
	idValue := reflect.ValueOf(value)
	idMap := map[string]interface{}{}
	switch idValue.Kind() {
	case reflect.Map:
		for _, key := range keys {
			if !idValue.MapIndex(reflect.ValueOf(key)).IsValid() {
				return nil, fmt.Errorf("wrong mapping key and value")
			}
			idMap[key] = idValue.MapIndex(reflect.ValueOf(key)).Interface()
		}
		if len(idMap) != idValue.Len() {
			return nil, fmt.Errorf("wrong mapping key and value")
		}
	case reflect.Slice, reflect.Array:
		if len(keys) != idValue.Len() {
			return nil, fmt.Errorf("wrong mapping key and value")
		}
		for idx := range keys {
			idMap[keys[idx]] = idValue.Index(idx).Interface()
		}
	default:
		idMap[keys[0]] = idValue.Interface()
	}
	keyMap := map[string]*dynamodb.AttributeValue{}
	for key, value := range idMap {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.String:
			keyMap[key] = &dynamodb.AttributeValue{S: aws.String(v.String())}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			keyMap[key] = &dynamodb.AttributeValue{N: aws.String(strconv.FormatInt(v.Int(), 10))}
		case reflect.Float32, reflect.Float64:
			keyMap[key] = &dynamodb.AttributeValue{N: aws.String(fmt.Sprintf("%g", v.Float()))}
		default:
			return keyMap, fmt.Errorf("data type not support")
		}
	}
	return keyMap, nil
}
