package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	a "github.com/core-go/auth"
	"google.golang.org/api/iterator"
	"strconv"
	"time"
)

type AuthenticationRepository struct {
	UserCollection          *firestore.CollectionRef
	PasswordCollection      *firestore.CollectionRef
	CheckTwoFactors         func(ctx context.Context, id string) (bool, error)
	ActivatedStatus         interface{}
	Status                  a.UserStatusConfig
	UserName                string
	SuccessTimeName         string
	FailTimeName            string
	FailCountName           string
	LockedUntilTimeName     string
	StatusName              string
	RoleName                string
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

func NewAuthenticationRepositoryByConfig(client *firestore.Client, userCollectionName, passwordCollectionName string, checkTwoFactors func(ctx context.Context, id string) (bool, error), activatedStatus interface{}, status a.UserStatusConfig, c a.SchemaConfig) *AuthenticationRepository {
	return NewAuthenticationRepository(client, userCollectionName, passwordCollectionName, checkTwoFactors, activatedStatus, status, c.Username, c.SuccessTime, c.FailTime, c.FailCount, c.LockedUntilTime, c.Status, c.Roles, c.PasswordChangedTime, c.Password, c.Contact, c.Email, c.Phone, c.DisplayName, c.MaxPasswordAge, c.UserType, c.AccessDateFrom, c.AccessDateTo, c.AccessTimeFrom, c.AccessTimeTo, c.TwoFactors)
}

func NewAuthenticationRepository(client *firestore.Client, userCollectionName, passwordCollectionName string, checkTwoFactors func(ctx context.Context, id string) (bool, error), activatedStatus interface{}, status a.UserStatusConfig, userName, successTimeName, failTimeName, failCountName, lockedUntilTimeName, statusName, roleName, passwordChangedTimeName, passwordName, contactName, emailName, phoneName, displayNameName, maxPasswordAgeName, userTypeName, accessDateFromName, accessDateToName, accessTimeFromName, accessTimeToName, twoFactorsName string) *AuthenticationRepository {
	passwordCollection := client.Collection(passwordCollectionName)
	userCollection := passwordCollection
	if passwordCollectionName != userCollectionName {
		userCollection = client.Collection(userCollectionName)
	}
	return &AuthenticationRepository{
		UserCollection:          userCollection,
		PasswordCollection:      passwordCollection,
		CheckTwoFactors:         checkTwoFactors,
		ActivatedStatus:         activatedStatus,
		Status:                  status,
		UserName:                userName,
		SuccessTimeName:         successTimeName,
		FailTimeName:            failTimeName,
		FailCountName:           failCountName,
		LockedUntilTimeName:     lockedUntilTimeName,
		StatusName:              statusName,
		RoleName:                roleName,
		PasswordChangedTimeName: passwordChangedTimeName,
		PasswordName:            passwordName,
		ContactName:             contactName,
		EmailName:               emailName,
		PhoneName:               phoneName,
		DisplayNameName:         displayNameName,
		MaxPasswordAgeName:      maxPasswordAgeName,
		UserTypeName:            userTypeName,
		AccessDateFromName:      accessDateFromName,
		AccessDateToName:        accessDateToName,
		AccessTimeFromName:      accessTimeFromName,
		AccessTimeToName:        accessTimeToName,
		TwoFactorsName:          twoFactorsName,
	}
}

func (r *AuthenticationRepository) GetUserInfo(ctx context.Context, auth a.AuthInfo) (*a.UserInfo, error) {
	userInfo := a.UserInfo{}
	//query := bson.M{"_id": id}
	iter := r.UserCollection.Where(r.UserName, "==", auth.Username).Documents(ctx)
	defer iter.Stop()
	result, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	rawStatus := result.Data()
	if rawStatus == nil {
		//should not happen
		return nil, errors.New("user info not found")
	}
	//raw, er1 := result.DecodeBytes()
	//if er1 != nil {
	//	return nil, er1
	//}

	if len(r.StatusName) > 0 {
		//rawStatus := raw.Lookup(r.StatusName)
		statusInfo, ok := rawStatus[r.StatusName]
		statusUserInfo := ""
		if ok {
			switch v := statusInfo.(type) {
			case int:
				statusUserInfo = strconv.Itoa(v)
			case int64:
				statusUserInfo = strconv.FormatInt(v, 10)
			case string:
				statusUserInfo = v
			case bool:
				statusUserInfo = strconv.FormatBool(v)
			default:
				return nil, fmt.Errorf(r.StatusName+": is of unsupported type %T", v)
			}
		}
		deactivated := statusUserInfo == r.Status.Deactivated
		userInfo.Deactivated = &deactivated
		userInfo.Suspended = statusUserInfo == r.Status.Suspended
		userInfo.Disable = statusUserInfo == r.Status.Disable
	}

	userInfo.Id = result.Ref.ID

	if len(r.UserName) > 0 {
		name, ok := rawStatus[r.UserName]
		if ok {
			if e, k := name.(string); k {
				userInfo.Username = e
			}
		}
	}
	if len(r.ContactName) > 0 {
		contact, ok := rawStatus[r.ContactName]
		if ok {
			if e, k := contact.(string); k {
				userInfo.Contact = &e
			}
		}
	}
	if len(r.EmailName) > 0 {
		email, ok := rawStatus[r.EmailName]
		if ok {
			if e, k := email.(string); k {
				userInfo.Email = &e
			}
		}
	}
	if len(r.PhoneName) > 0 {
		phone, ok := rawStatus[r.PhoneName]
		if ok {
			if e, k := phone.(string); k {
				userInfo.Phone = &e
			}
		}
	}

	if len(r.RoleName) > 0 {
		roles, ok := rawStatus[r.RoleName]
		if ok {
			if array, k := roles.([]interface{}); k {
				var tempRoles []string
				for _, i2 := range array {
					if e, k := i2.(string); k {
						tempRoles = append(tempRoles, e)
					}
				}
				userInfo.Roles = tempRoles
			}
		}
	}

	if len(r.DisplayNameName) > 0 {
		displayName, ok := rawStatus[r.DisplayNameName]
		if ok {
			if e, k := displayName.(string); k {
				userInfo.DisplayName = &e
			}
		}
	}

	if len(r.MaxPasswordAgeName) > 0 {
		maxPasswordAge, ok := rawStatus[r.MaxPasswordAgeName]
		if ok {
			if e, k := maxPasswordAge.(int32); k {
				userInfo.MaxPasswordAge = &e
			}
		}
	}

	if len(r.UserTypeName) > 0 {
		maxPasswordAge, ok := rawStatus[r.UserTypeName]
		if ok {
			if e, k := maxPasswordAge.(string); k {
				userInfo.UserType = &e
			}
		}
	}

	if len(r.AccessDateFromName) > 0 {
		accessDateFrom, ok := rawStatus[r.AccessDateFromName]
		if ok {
			if e, k := accessDateFrom.(time.Time); k {
				userInfo.AccessDateFrom = &e
			}
		}
	}
	if len(r.AccessDateToName) > 0 {
		accessDateTo, ok := rawStatus[r.AccessDateToName]
		if ok {
			if e, k := accessDateTo.(time.Time); k {
				userInfo.AccessDateTo = &e
			}
		}
	}

	if len(r.AccessTimeFromName) > 0 {
		accessTimeFrom, ok := rawStatus[r.AccessTimeFromName]
		if ok {
			if e, k := accessTimeFrom.(time.Time); k {
				userInfo.AccessTimeFrom = &e
			} else if s, k := accessTimeFrom.(string); k {
				userInfo.AccessTimeFrom = getTime(s)
			}
		}
	}

	if len(r.AccessTimeToName) > 0 {
		accessTimeTo, ok := rawStatus[r.AccessTimeToName]
		if ok {
			if e, k := accessTimeTo.(time.Time); k {
				userInfo.AccessTimeTo = &e
			} else if s, k := accessTimeTo.(string); k {
				userInfo.AccessTimeTo = getTime(s)
			}
		}
	}

	if r.CheckTwoFactors != nil {
		id := userInfo.Id
		if len(id) == 0 {
			id = auth.Username
		}
		ok, er2 := r.CheckTwoFactors(ctx, id)
		if er2 != nil {
			return &userInfo, er2
		}
		userInfo.TwoFactors = ok
	} else if len(r.TwoFactorsName) > 0 {
		isTwoFactor, ok := rawStatus[r.TwoFactorsName]
		if ok {
			if b, k := isTwoFactor.(bool); k {
				userInfo.TwoFactors = b
			}
		}
	}

	if r.UserCollection.ID == r.PasswordCollection.ID {
		return r.getPasswordInfo(ctx, &userInfo, rawStatus), nil
	}
	//temp, ok := rawStatus["_id"]
	//if !ok {
	//	return nil, fmt.Errorf("_id: does not exist")
	//}
	//id1, k := temp.(string)
	//if !k {
	//	return nil, fmt.Errorf("_id:is not a string")
	//}
	resultPass, err := r.PasswordCollection.Doc(result.Ref.ID).Get(ctx)
	if err != nil {
		return nil, err
	}

	rawPassStatus := resultPass.Data()
	if rawPassStatus == nil {
		//should not happen
		return nil, errors.New("pass info not found")
	}
	ctx = context.WithValue(ctx, "roles", userInfo.Roles)
	return r.getPasswordInfo(ctx, &userInfo, rawPassStatus), nil
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

func (r *AuthenticationRepository) getPasswordInfo(ctx context.Context, user *a.UserInfo, raw map[string]interface{}) *a.UserInfo {
	if len(r.PasswordName) > 0 {
		pass, ok := raw[r.PasswordName]
		if ok {
			if e, k := pass.(string); k {
				user.Password = e
			}
		}
	}

	if len(r.LockedUntilTimeName) > 0 {
		pass, ok := raw[r.LockedUntilTimeName]
		if ok {
			if e, k := pass.(time.Time); k {
				user.LockedUntilTime = &e
			}
		}
	}

	if len(r.SuccessTimeName) > 0 {
		pass, ok := raw[r.SuccessTimeName]
		if ok {
			if e, k := pass.(time.Time); k {
				user.SuccessTime = &e
			}
		}
	}

	if len(r.FailTimeName) > 0 {
		pass, ok := raw[r.FailTimeName]
		if ok {
			if e, k := pass.(time.Time); k {
				user.FailTime = &e
			}
		}
	}

	if len(r.FailCountName) > 0 {
		pass, ok := raw[r.FailCountName]
		if ok {
			if e, k := pass.(int64); k {
				failCount := int(e)
				user.FailCount = &failCount
			}
		}
	}

	if len(r.PasswordChangedTimeName) > 0 {
		pass, ok := raw[r.PasswordChangedTimeName]
		if ok {
			if e, k := pass.(time.Time); k {
				user.PasswordChangedTime = &e
			}
		}
	}
	return user
}

/*
func (r *AuthenticationRepository) GetPasswordInfo(ctx context.Context, id string) (*a.PasswordInfo, error) {
	authentication := a.PasswordInfo{}
	ok, err := f.FindOneAndDecode(ctx, r.collection, id, &authentication)
	if ok && err == nil {
		return &authentication, nil
	}
	return nil, err
}
//*/
//func (r *AuthenticationRepository) PassAuthenticationAndActivate(ctx context.Context, userId string) (int64, error) {
//	return r.PassAuthentication(ctx, userId)
//}

func (r *AuthenticationRepository) Pass(ctx context.Context, userId string, deactivated *bool) (int64, error) {
	return r.passAuthenticationAndActivate(ctx, userId, deactivated)
}

func (r *AuthenticationRepository) passAuthenticationAndActivate(ctx context.Context, userId string, updateStatus *bool) (int64, error) {
	if len(r.SuccessTimeName) == 0 && len(r.FailCountName) == 0 && len(r.LockedUntilTimeName) == 0 {
		if updateStatus != nil && !*updateStatus {
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
	if updateStatus != nil && !*updateStatus {
		return r.upsertWithMap(ctx, r.PasswordCollection, userId, pass)
	}
	if r.UserCollection.ID == r.PasswordCollection.ID {
		pass[r.StatusName] = r.ActivatedStatus
		return r.upsertWithMap(ctx, r.PasswordCollection, userId, pass)
	}
	k1, er1 := r.upsertWithMap(ctx, r.PasswordCollection, userId, pass)
	if er1 != nil {
		return k1, er1
	}
	user := make(map[string]interface{})
	user["_id"] = userId
	user[r.StatusName] = r.ActivatedStatus
	k2, er2 := r.upsertWithMap(ctx, r.UserCollection, userId, pass)
	return k1 + k2, er2
}

//func (r *AuthenticationRepository) PassAuthentication(ctx context.Context, userId string) (int64, error) {
//	pass := make(map[string]interface{})
//	pass["_id"] = userId
//	if len(r.SuccessTimeName) > 0 {
//		pass[r.successTimeName] = time.Now()
//	}
//	if len(r.failCountName) > 0 {
//		pass[r.failCountName] = 0
//	}
//	if len(r.lockedUntilTimeName) > 0 {
//		pass[r.lockedUntilTimeName] = nil
//	}
//	return r.upsertWithMap(ctx, r.collection, userId, pass)
//}

//func (r *AuthenticationRepository) WrongPassword(ctx context.Context, userId string, failCount int, lockedUntil *time.Time) error {
//	pass := make(map[string]interface{})
//	pass["_id"] = userId
//	if len(r.failTimeName) > 0 {
//		pass[r.failTimeName] = time.Now()
//	}
//	if len(r.failCountName) > 0 {
//		pass[r.failCountName] = failCount
//		if len(r.lockedUntilTimeName) > 0 {
//			pass[r.lockedUntilTimeName] = lockedUntil
//		}
//	}
//	_, err := r.upsertWithMap(ctx, r.collection, userId, pass)
//	return err
//}

func (r *AuthenticationRepository) Fail(ctx context.Context, userId string, failCount *int, lockedUntil *time.Time) error {
	if len(r.FailTimeName) == 0 && len(r.FailCountName) == 0 && len(r.LockedUntilTimeName) == 0 {
		return nil
	}
	pass := make(map[string]interface{})
	pass["_id"] = userId
	if len(r.FailTimeName) > 0 {
		pass[r.FailTimeName] = time.Now()
	}
	if failCount != nil && len(r.FailCountName) > 0 {
		pass[r.FailCountName] = *failCount + 1
		if len(r.LockedUntilTimeName) > 0 {
			pass[r.LockedUntilTimeName] = lockedUntil
		}
	}
	//query := bson.M{"_id": userId}
	_, err := r.upsertWithMap(ctx, r.PasswordCollection, userId, pass)
	return err
}

func (r *AuthenticationRepository) upsertWithMap(ctx context.Context, collection *firestore.CollectionRef, id string, data map[string]interface{}) (int64, error) {
	_, err := collection.Doc(id).Set(ctx, data, firestore.MergeAll)
	if err != nil {
		return 0, err
	}
	return 1, nil
}
