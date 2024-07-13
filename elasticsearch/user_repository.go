package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	auth "github.com/core-go/authentication"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type AuthenticationRepository struct {
	Client                  *elasticsearch.Client
	UserIndexName           string
	PasswordIndexName       string
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

func NewAuthenticationRepositoryByConfig(client *elasticsearch.Client, userIndexName, passwordIndexName string, activatedStatus interface{}, status auth.UserStatusConfig, c auth.SchemaConfig, options ...func(context.Context, string) (bool, error)) *AuthenticationRepository {
	return NewAuthenticationRepository(client, userIndexName, passwordIndexName, activatedStatus, status, c.Username, c.SuccessTime, c.FailTime, c.FailCount, c.LockedUntilTime, c.Status, c.PasswordChangedTime, c.Password, c.Contact, c.Email, c.Phone, c.DisplayName, c.MaxPasswordAge, c.Roles, c.UserType, c.AccessDateFrom, c.AccessDateTo, c.AccessTimeFrom, c.AccessTimeTo, c.TwoFactors, options...)
}

func NewAuthenticationRepository(client *elasticsearch.Client, userIndexName, passwordIndexName string, activatedStatus interface{}, status auth.UserStatusConfig, userName, successTimeName, failTimeName, failCountName, lockedUntilTimeName, statusName, passwordChangedTimeName, passwordName, contactName, emailName, phoneName, displayNameName, maxPasswordAgeName, rolesName, userTypeName, accessDateFromName, accessDateToName, accessTimeFromName, accessTimeToName string, twoFactorsName string, options ...func(context.Context, string) (bool, error)) *AuthenticationRepository {
	var checkTwoFactors func(context.Context, string) (bool, error)
	if len(options) >= 1 {
		checkTwoFactors = options[0]
	}
	return &AuthenticationRepository{Client: client, UserIndexName: userIndexName, PasswordIndexName: passwordIndexName, CheckTwoFactors: checkTwoFactors, ActivatedStatus: activatedStatus, Status: status, UserName: userName, SuccessTimeName: successTimeName, FailTimeName: failTimeName, FailCountName: failCountName, LockedUntilTimeName: lockedUntilTimeName, StatusName: statusName, PasswordChangedTimeName: passwordChangedTimeName, PasswordName: passwordName, ContactName: contactName, EmailName: emailName, PhoneName: phoneName, DisplayNameName: displayNameName, MaxPasswordAgeName: maxPasswordAgeName, RolesName: rolesName, UserTypeName: userTypeName, AccessDateFromName: accessDateFromName, AccessDateToName: accessDateToName, AccessTimeFromName: accessTimeFromName, AccessTimeToName: accessTimeToName, TwoFactorsName: twoFactorsName}
}

func (r *AuthenticationRepository) GetUserInfo(ctx context.Context, username string) (*auth.UserInfo, error) {
	userInfo := auth.UserInfo{}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"match": map[string]interface{}{
					"_id": username,
				},
			},
		},
	}
	raw := make(map[string]interface{})
	ok, err := findOneAndDecode(ctx, r.Client, []string{r.UserIndexName}, query, &raw)
	if !ok || err != nil {
		return nil, err
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
		if maxPasswordAgeName, ok := raw[r.MaxPasswordAgeName].(int32); ok {
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
		if isTwoFactor, ok := raw[r.AccessTimeToName].(bool); ok {
			userInfo.TwoFactors = isTwoFactor
		}
	}

	if r.UserIndexName == r.PasswordIndexName {
		return r.getPasswordInfo(ctx, &userInfo, raw), nil
	}

	rawPassword := make(map[string]interface{})
	ok1, er1 := findOneAndDecode(ctx, r.Client, []string{r.UserIndexName}, query, &rawPassword)
	if !ok1 || er1 != nil {
		return nil, er1
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

func (r *AuthenticationRepository) Pass(ctx context.Context, userId string) (int64, error) {
	return r.passAuthenticationAndActivate(ctx, userId, false)
}

func (r *AuthenticationRepository) PassAndActivate(ctx context.Context, userId string) (int64, error) {
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
		return upsertOne(ctx, r.Client, r.PasswordIndexName, userId, pass)
	}
	if r.UserIndexName == r.PasswordIndexName {
		pass[r.StatusName] = r.ActivatedStatus
		return upsertOne(ctx, r.Client, r.PasswordIndexName, userId, pass)
	}
	k1, er1 := upsertOne(ctx, r.Client, r.PasswordIndexName, userId, pass)
	if er1 != nil {
		return k1, er1
	}
	user := make(map[string]interface{})
	user["_id"] = userId
	user[r.StatusName] = r.ActivatedStatus
	k2, er2 := upsertOne(ctx, r.Client, r.UserIndexName, userId, user)
	return k1 + k2, er2
}

func (r *AuthenticationRepository) Fail(ctx context.Context, userId string, failCount int, lockedUntil *time.Time) error {
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
	_, err := upsertOne(ctx, r.Client, r.PasswordIndexName, userId, pass)
	return err
}

func findOneAndDecode(ctx context.Context, es *elasticsearch.Client, index []string, query map[string]interface{}, result interface{}) (bool, error) {
	req := esapi.SearchRequest{
		Index:          index,
		Body:           esutil.NewJSONReader(query),
		TrackTotalHits: true,
		Pretty:         true,
	}
	res, err := req.Do(ctx, es)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return false, errors.New("response error")
	} else {
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return false, err
		} else {
			hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
			total := int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
			if total >= 1 {
				if err := json.NewDecoder(esutil.NewJSONReader(hits[0])).Decode(&result); err != nil {
					return false, err
				}
				return true, nil
			}
			return false, nil
		}
	}
}

func upsertOne(ctx context.Context, es *elasticsearch.Client, indexName string, id string, model interface{}) (int64, error) {
	body := buildQueryWithoutIdFromObject(model)
	req := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: id,
		Body:       esutil.NewJSONReader(body),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, es)
	if err != nil {
		return -1, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return -1, errors.New("document ID not exists in the index")
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return -1, err
	}
	successful := int64(r["_shards"].(map[string]interface{})["successful"].(float64))
	return successful, nil
}

func buildQueryWithoutIdFromObject(object interface{}) map[string]interface{} {
	valueOf := reflect.Indirect(reflect.ValueOf(object))
	idIndex, _ := findIdField(valueOf.Type())
	result := map[string]interface{}{}
	for i := 0; i < valueOf.NumField(); i++ {
		if i != idIndex {
			_, jsonName := findFieldByIndex(valueOf.Type(), i)
			result[jsonName] = valueOf.Field(i).Interface()
		}
	}
	return result
}

func findIdField(modelType reflect.Type) (int, string) {
	return findFieldByJson(modelType, "_id")
}

func findFieldByJson(modelType reflect.Type, jsonTagName string) (index int, fieldName string) {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		tag1, ok1 := field.Tag.Lookup("json")
		if ok1 && strings.Split(tag1, ",")[0] == jsonTagName {
			return i, field.Name
		}
	}
	return -1, jsonTagName
}

func findFieldByIndex(modelType reflect.Type, fieldIndex int) (fieldName, jsonTagName string) {
	if fieldIndex < modelType.NumField() {
		field := modelType.Field(fieldIndex)
		jsonTagName := ""
		if jsonTag, ok := field.Tag.Lookup("json"); ok {
			jsonTagName = strings.Split(jsonTag, ",")[0]
		}
		return field.Name, jsonTagName
	}
	return "", ""
}
