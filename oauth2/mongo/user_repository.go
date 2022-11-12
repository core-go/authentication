package mongo

import (
	"context"
	"github.com/core-go/auth"
	"github.com/core-go/auth/oauth2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"strings"
	"time"
)

type UserRepository struct {
	Collection      *mongo.Collection
	Prefix          string
	ActivatedStatus string
	Services        []string
	StatusName      string
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
}

func NewUserRepositoryByConfig(db *mongo.Database, collectionName, prefix string, activatedStatus string, services []string, c oauth2.OAuth2SchemaConfig, status *auth.UserStatusConfig, options ...oauth2.OAuth2GenderMapper) *UserRepository {
	var genderMapper oauth2.OAuth2GenderMapper
	if len(options) >= 1 {
		genderMapper = options[0]
	}
	if len(c.UserName) == 0 {
		c.UserName = "userName"
	}
	if len(c.Email) == 0 {
		c.Email = "email"
	}
	if len(c.Status) == 0 {
		c.Status = "status"
	}
	if len(c.OAuth2Email) == 0 {
		c.OAuth2Email = "Email"
	}
	if len(c.Account) == 0 {
		c.Account = "Account"
	}
	if len(c.Active) == 0 {
		c.Active = "Active"
	}
	collection := db.Collection(collectionName)
	m := &UserRepository{
		Collection:      collection,
		Prefix:          prefix,
		ActivatedStatus: activatedStatus,
		Services:        services,
		GenderMapper:    genderMapper,
		Status:          status,
		Schema:          &c,
		updatedTimeName: c.UpdatedTime,
		updatedByName:   c.UpdatedBy,

		StatusName:      c.Status,
		UserName:        c.UserName,
		EmailName:       c.Email,
		OAuth2EmailName: c.Email,
		AccountName:     c.Account,
		ActiveName:      c.Active,
	}
	return m
}

func NewUserRepository(db *mongo.Database, collectionName, prefix, activatedStatus string, services []string, pictureName, displayName, givenName, familyName, middleName, genderName string, status *auth.UserStatusConfig, options ...oauth2.OAuth2GenderMapper) *UserRepository {
	var genderMapper oauth2.OAuth2GenderMapper
	if len(options) >= 1 {
		genderMapper = options[0]
	}
	collection := db.Collection(collectionName)

	m := &UserRepository{
		Collection:      collection,
		Prefix:          prefix,
		ActivatedStatus: activatedStatus,
		StatusName:      "status",
		Services:        services,
		UserName:        "userName",
		EmailName:       "email",
		OAuth2EmailName: "Email",
		AccountName:     "Account",
		ActiveName:      "Active",
		GenderMapper:    genderMapper,
		Status:          status,
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

func (r *UserRepository) GetUser(ctx context.Context, email string) (string, bool, bool, error) {
	// query := bson.M{"$or": []bson.M{{"userName": email}, {"email": email}, {"linkedinEmail": email}, {"facebookEmail": email}, {"googleEmail": email}}}
	queries := []bson.M{{r.UserName: email}, {r.EmailName: email}, {r.Prefix + r.OAuth2EmailName: email}}
	for _, sv := range r.Services {
		if sv != r.Prefix {
			v := bson.M{sv + r.OAuth2EmailName: email}
			queries = append(queries, v)
		}
	}
	query := bson.M{"$or": queries}
	x := r.Collection.FindOne(ctx, query)
	k, er3 := x.DecodeBytes()
	disable := false
	suspended := false
	if er3 != nil {
		if strings.Contains(er3.Error(), "mongo: no documents in result") {
			return "", disable, suspended, nil
		}
		return "", disable, suspended, er3
	}
	userId := k.Lookup("_id").StringValue()
	if r.Status != nil {
		f := k.Lookup(r.StatusName)
		var status string
		if f.IsNumber() {
			cInt := f.Int32()
			status = strconv.Itoa(int(cInt))
		} else {
			status = k.Lookup(r.StatusName).StringValue()
		}
		if status == r.Status.Disable {
			disable = true
		}
		if status == r.Status.Suspended {
			suspended = true
		}
	}
	return userId, disable, suspended, nil
}

func (r *UserRepository) Update(ctx context.Context, id, email, account string) (bool, error) {
	user := make(map[string]interface{})

	user[r.Prefix+r.OAuth2EmailName] = email
	user[r.Prefix+r.AccountName] = account
	user[r.Prefix+r.ActiveName] = true

	if len(r.updatedTimeName) > 0 {
		user[r.updatedTimeName] = time.Now()
	}
	if len(r.updatedByName) > 0 {
		user[r.updatedByName] = id
	}

	updateQuery := bson.M{
		"$set": user,
	}

	result, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id}, updateQuery)

	return result.ModifiedCount+result.UpsertedCount+result.MatchedCount > 0, err
}

func (r *UserRepository) Insert(ctx context.Context, id string, user oauth2.User) (bool, error) {
	userMap := r.userToMap(ctx, id, user)
	_, err := r.Collection.InsertOne(ctx, userMap)
	if err != nil {
		errMsg := err.Error()
		if strings.Index(errMsg, "duplicate key error collection:") >= 0 {
			return true, nil
		} else {
			return false, err
		}
	}
	return false, nil
}

func (r *UserRepository) userToMap(ctx context.Context, id string, user oauth2.User) map[string]interface{} {
	userMap := oauth2.UserToMap(ctx, id, user, r.GenderMapper, r.Schema)

	userMap["_id"] = id
	userMap[r.UserName] = user.Email
	userMap[r.StatusName] = r.ActivatedStatus

	userMap[r.Prefix+r.OAuth2EmailName] = user.Email
	userMap[r.Prefix+r.AccountName] = user.Account
	userMap[r.Prefix+r.ActiveName] = true
	return userMap
}
