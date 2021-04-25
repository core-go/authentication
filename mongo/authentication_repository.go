package mongo

import (
	"context"
	"fmt"
	"github.com/common-go/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type AuthenticationRepository struct {
	UserCollection          *mongo.Collection
	PasswordCollection      *mongo.Collection
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
	EmailName               string
	PhoneName               string
	ContactName             string
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

func NewAuthenticationRepositoryByConfig(db *mongo.Database, userCollectionName, passwordCollectionName string, activatedStatus interface{}, status auth.UserStatusConfig, c auth.SchemaConfig, options ...func(context.Context, string) (bool, error)) *AuthenticationRepository {
	return NewAuthenticationRepository(db, userCollectionName, passwordCollectionName, activatedStatus, status, c.Username, c.SuccessTime, c.FailTime, c.FailCount, c.LockedUntilTime, c.Status, c.PasswordChangedTime, c.Password, c.Contact, c.Email, c.Phone, c.DisplayName, c.MaxPasswordAge, c.Roles, c.UserType, c.AccessDateFrom, c.AccessDateTo, c.AccessTimeFrom, c.AccessTimeTo, c.TwoFactors, options...)
}

func NewAuthenticationRepository(db *mongo.Database, userCollectionName, passwordCollectionName string, activatedStatus interface{}, status auth.UserStatusConfig, userName, successTimeName, failTimeName, failCountName, lockedUntilTimeName, statusName, passwordChangedTimeName, passwordName, contactName, emailName, phoneName, displayNameName, maxPasswordAgeName, rolesName, userTypeName, accessDateFromName, accessDateToName, accessTimeFromName, accessTimeToName, twoFactorsName string, options ...func(context.Context, string) (bool, error)) *AuthenticationRepository {
	passwordCollection := db.Collection(passwordCollectionName)
	userCollection := passwordCollection
	if passwordCollectionName != userCollectionName {
		userCollection = db.Collection(userCollectionName)
	}
	var checkTwoFactors func(context.Context, string) (bool, error)
	if len(options) >= 1 {
		checkTwoFactors = options[0]
	}
	return &AuthenticationRepository{UserCollection: userCollection, PasswordCollection: passwordCollection, CheckTwoFactors: checkTwoFactors, ActivatedStatus: activatedStatus, Status: status, UserName: userName, SuccessTimeName: successTimeName, FailTimeName: failTimeName, FailCountName: failCountName, LockedUntilTimeName: lockedUntilTimeName, StatusName: statusName, PasswordChangedTimeName: passwordChangedTimeName, PasswordName: passwordName, ContactName: contactName, EmailName: emailName, PhoneName: phoneName, DisplayNameName: displayNameName, MaxPasswordAgeName: maxPasswordAgeName, RolesName: rolesName, UserTypeName: userTypeName, AccessDateFromName: accessDateFromName, AccessDateToName: accessDateToName, AccessTimeFromName: accessTimeFromName, AccessTimeToName: accessTimeToName, TwoFactorsName: twoFactorsName}
}

func (r *AuthenticationRepository) GetUserInfo(ctx context.Context, username string) (*auth.UserInfo, error) {
	userInfo := auth.UserInfo{}
	query := bson.M{r.UserName: username}
	result := r.UserCollection.FindOne(ctx, query)
	if result.Err() != nil {
		if fmt.Sprint(result.Err()) == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, result.Err()
	}

	raw, er1 := result.DecodeBytes()
	if er1 != nil {
		return nil, er1
	}

	if id, ok := raw.Lookup("_id").StringValueOK(); ok {
		userInfo.Id = id
	}

	if len(r.StatusName) > 0 {
		rawStatus := raw.Lookup(r.StatusName)
		status, ok := rawStatus.StringValueOK()
		if !ok {
			iStatus, ok2 := rawStatus.Int32OK()
			if ok2 {
				status = strconv.Itoa(int(iStatus))
			} else {
				bStatus, ok3 := rawStatus.BooleanOK()
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
		if contact, ok := raw.Lookup(r.ContactName).StringValueOK(); ok {
			userInfo.Contact = contact
		}
	}
	if len(r.EmailName) > 0 {
		if email, ok := raw.Lookup(r.EmailName).StringValueOK(); ok {
			userInfo.Email = email
		}
	}
	if len(r.PhoneName) > 0 {
		if phone, ok := raw.Lookup(r.PhoneName).StringValueOK(); ok {
			userInfo.Phone = phone
		}
	}

	if len(r.DisplayNameName) > 0 {
		if displayName, ok := raw.Lookup(r.DisplayNameName).StringValueOK(); ok {
			userInfo.DisplayName = displayName
		}
	}

	if len(r.MaxPasswordAgeName) > 0 {
		if raw.Lookup(r.MaxPasswordAgeName).IsNumber() == true {
			userInfo.MaxPasswordAge = raw.Lookup(r.MaxPasswordAgeName).Int64()
		}
	}

	if len(r.UserTypeName) > 0 {
		if userType, ok := raw.Lookup(r.UserTypeName).StringValueOK(); ok {
			userInfo.UserType = userType
		}
	}

	if len(r.AccessDateFromName) > 0 {
		if accessDateFrom, ok := raw.Lookup(r.AccessDateFromName).TimeOK(); ok {
			userInfo.AccessDateFrom = &accessDateFrom
		}
	}

	if len(r.AccessDateToName) > 0 {
		if accessDateTo, ok := raw.Lookup(r.AccessDateToName).TimeOK(); ok {
			userInfo.AccessDateTo = &accessDateTo
		}
	}

	if len(r.AccessTimeFromName) > 0 {
		if accessTimeFrom, ok := raw.Lookup(r.AccessTimeFromName).TimeOK(); ok {
			userInfo.AccessTimeFrom = &accessTimeFrom
		} else if accessTimeFrom, ok := raw.Lookup(r.AccessTimeFromName).StringValueOK(); ok {
			userInfo.AccessTimeFrom = getTime(accessTimeFrom)
		}
	}

	if len(r.AccessTimeToName) > 0 {
		if accessTimeTo, ok := raw.Lookup(r.AccessTimeToName).TimeOK(); ok {
			userInfo.AccessTimeTo = &accessTimeTo
		} else if accessTimeTo, ok := raw.Lookup(r.AccessTimeToName).StringValueOK(); ok {
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
		if isTwoFactor, ok := raw.Lookup(r.TwoFactorsName).BooleanOK(); ok {
			userInfo.TwoFactors = isTwoFactor
		}
	}

	if r.UserCollection.Name() == r.PasswordCollection.Name() {
		return r.getPasswordInfo(ctx, &userInfo, raw), nil
	}
	id1 := raw.Lookup("_id").StringValue()
	query2 := bson.M{"_id": id1}
	resultPass := r.PasswordCollection.FindOne(ctx, query2)
	if resultPass.Err() != nil {
		return nil, resultPass.Err()
	}
	rawPassword, er3 := resultPass.DecodeBytes()
	if er3 != nil {
		return nil, er3
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

func (r *AuthenticationRepository) getPasswordInfo(ctx context.Context, user *auth.UserInfo, raw bson.Raw) *auth.UserInfo {
	if len(r.PasswordName) > 0 {
		if pass, ok := raw.Lookup(r.PasswordName).StringValueOK(); ok {
			user.Password = pass
		}
	}

	if len(r.LockedUntilTimeName) > 0 {
		if lockedUntilTime, ok := raw.Lookup(r.LockedUntilTimeName).TimeOK(); ok {
			user.LockedUntilTime = &lockedUntilTime
		}
	}

	if len(r.SuccessTimeName) > 0 {
		if successTime, ok := raw.Lookup(r.SuccessTimeName).TimeOK(); ok {
			user.SuccessTime = &successTime
		}
	}

	if len(r.FailTimeName) > 0 {
		if failTime, ok := raw.Lookup(r.FailTimeName).TimeOK(); ok {
			user.FailTime = &failTime
		}
	}

	if len(r.FailCountName) > 0 {
		if raw.Lookup(r.FailCountName).IsNumber() == true {
			user.FailCount = int(raw.Lookup(r.FailCountName).Int32())
		}
	}

	if len(r.PasswordChangedTimeName) > 0 {
		if passwordChangedTime, ok := raw.Lookup(r.PasswordChangedTimeName).TimeOK(); ok {
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
		if !updateStatus {
			return 0, nil
		} else if len(r.StatusName) == 0 {
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
	query := bson.M{"_id": userId}
	if !updateStatus {
		return upsertOne(ctx, r.PasswordCollection, query, pass)
	}
	if r.UserCollection.Name() == r.PasswordCollection.Name() {
		pass[r.StatusName] = r.ActivatedStatus
		return upsertOne(ctx, r.PasswordCollection, query, pass)
	}
	k1, er1 := upsertOne(ctx, r.PasswordCollection, query, pass)
	if er1 != nil {
		return k1, er1
	}
	user := make(map[string]interface{})
	user["_id"] = userId
	user[r.StatusName] = r.ActivatedStatus
	k2, er2 := upsertOne(ctx, r.UserCollection, user, query)
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
	query := bson.M{"_id": userId}
	_, err := upsertOne(ctx, r.PasswordCollection, query, pass)
	return err
}

func exist(ctx context.Context, collection *mongo.Collection, id interface{}, objectId bool) (bool, error) {
	query := bson.M{"_id": id}
	if objectId {
		objId, err := primitive.ObjectIDFromHex(id.(string))
		if err != nil {
			return false, err
		}
		query = bson.M{"_id": objId}
	}
	x := collection.FindOne(ctx, query)
	if x.Err() != nil {
		if fmt.Sprint(x.Err()) == "mongo: no documents in result" {
			return false, nil
		} else {
			return false, x.Err()
		}
	}
	return true, nil
}

func upsertOne(ctx context.Context, collection *mongo.Collection, filter bson.M, model interface{}) (int64, error) {
	defaultObjID, _ := primitive.ObjectIDFromHex("000000000000")

	if idValue := filter["_id"]; idValue == "" || idValue == 0 || idValue == defaultObjID {
		return insertOne(ctx, collection, model)
	} else {
		isExisted, err := exist(ctx, collection, idValue, false)
		if err != nil {
			return 0, err
		}
		if isExisted {
			update := bson.M{
				"$set": model,
			}
			result := collection.FindOneAndUpdate(ctx, filter, update)
			if result.Err() != nil {
				if fmt.Sprint(result.Err()) == "mongo: no documents in result" {
					return 0, nil
				} else {
					return 0, result.Err()
				}
			}
			return 1, result.Err()
		} else {
			return insertOne(ctx, collection, model)
		}
	}
}

func insertOne(ctx context.Context, collection *mongo.Collection, model interface{}) (int64, error) {
	result, err := collection.InsertOne(ctx, model)
	if err != nil {
		errMsg := err.Error()
		if strings.Index(errMsg, "duplicate key error collection:") >= 0 {
			if strings.Index(errMsg, "dup key: { _id: ") >= 0 {
				return -1, nil
			} else {
				return -2, nil
			}
		} else {
			return 0, err
		}
	} else {
		if idValue, ok := result.InsertedID.(primitive.ObjectID); ok {
			valueOfModel := reflect.Indirect(reflect.ValueOf(model))
			typeOfModel := valueOfModel.Type()
			idIndex, _ := findIdField(typeOfModel)
			if idIndex != -1 {
				mapObjectIdToModel(idValue, valueOfModel, idIndex)
			}
		}
		return 1, err
	}
}

func findIdField(modelType reflect.Type) (int, string) {
	return findField(modelType, "_id")
}

func findField(modelType reflect.Type, bsonName string) (int, string) {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		bsonTag := field.Tag.Get("bson")
		tags := strings.Split(bsonTag, ",")
		for _, tag := range tags {
			if strings.TrimSpace(tag) == bsonName {
				return i, field.Name
			}
		}
	}
	return -1, ""
}

func mapObjectIdToModel(id primitive.ObjectID, valueOfModel reflect.Value, idIndex int) {
	switch reflect.Indirect(valueOfModel).Field(idIndex).Kind() {
	case reflect.String:
		if _, err := setValue(valueOfModel, idIndex, id.Hex()); err != nil {
			log.Println("Err: ", err)
		}
		break
	default:
		if _, err := setValue(valueOfModel, idIndex, id); err != nil {
			log.Println("Err: ", err)
		}
		break
	}
}

func setValue(model interface{}, index int, value interface{}) (interface{}, error) {
	valueObject := reflect.Indirect(reflect.ValueOf(model))
	switch reflect.ValueOf(model).Kind() {
	case reflect.Ptr:
		{
			valueObject.Field(index).Set(reflect.ValueOf(value))
			return model, nil
		}
	default:
		if modelWithTypeValue, ok := model.(reflect.Value); ok {
			_, err := setValueWithTypeValue(modelWithTypeValue, index, value)
			return modelWithTypeValue.Interface(), err
		}
	}
	return model, nil
}

func setValueWithTypeValue(model reflect.Value, index int, value interface{}) (reflect.Value, error) {
	trueValue := reflect.Indirect(model)
	switch trueValue.Kind() {
	case reflect.Struct:
		{
			val := reflect.Indirect(reflect.ValueOf(value))
			if trueValue.Field(index).Kind() == val.Kind() {
				trueValue.Field(index).Set(reflect.ValueOf(value))
				return trueValue, nil
			} else {
				return trueValue, fmt.Errorf("value's kind must same as field's kind")
			}
		}
	default:
		return trueValue, nil
	}
}
