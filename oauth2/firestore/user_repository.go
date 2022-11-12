package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/core-go/auth"
	"strings"

	"github.com/core-go/auth/oauth2"
)

type UserRepository struct {
	Collection      *firestore.CollectionRef
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

func NewUserRepositoryByConfig(db *firestore.Client, collectionName, prefix string, activatedStatus string, services []string, c oauth2.OAuth2SchemaConfig, status *auth.UserStatusConfig, options ...oauth2.OAuth2GenderMapper) *UserRepository {
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
		updatedByName:   c.UpdatedBy,
	}
	return m
}

func NewUserRepository(db *firestore.Client, collectionName, prefix, activatedStatus string, services []string, pictureName, displayName, givenName, familyName, middleName, genderName string, status *auth.UserStatusConfig, options ...oauth2.OAuth2GenderMapper) *UserRepository {
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
	queries := []Query{
		{Key: r.UserName, Operator: "==", Value: email},
		{Key: r.EmailName, Operator: "==", Value: email},
		{Key: r.Prefix + r.OAuth2EmailName, Operator: "==", Value: email},
	}
	for _, sv := range r.Services {
		if sv != r.Prefix {
			queries = append(queries, Query{Key: sv + r.OAuth2EmailName, Operator: "==", Value: email})
		}
	}
	disable := false
	suspended := false
	snapShot, err := r.query(ctx, queries...)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return "", disable, suspended, nil
		}
		return "", disable, suspended, err
	}
	var userId, status string
	userId = snapShot.Ref.ID
	if r.Status != nil {
		data := snapShot.Data()
		if value, exist := data[r.StatusName]; exist {
			if s, exist := value.(string); exist {
				status = s
			}
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

func (r *UserRepository) query(ctx context.Context, queries ...Query) (*firestore.DocumentSnapshot, error) {
	for _, query := range queries {
		q := r.Collection.Where(query.Key, query.Operator, query.Value).Limit(1)
		iter := q.Documents(ctx)
		defer iter.Stop()
		doc, err := iter.Next()
		if err == nil {
			return doc, err
		}
	}
	return nil, fmt.Errorf("not found")
}

type Query struct {
	Key, Operator string
	Value         interface{}
}

func (r *UserRepository) Update(ctx context.Context, id, email, account string) (bool, error) {
	docSnap, err := r.Collection.Doc(id).Get(ctx)
	if err != nil || docSnap.Data() == nil {
		return false, err
	}

	updateValue := []firestore.Update{
		{Path: r.Prefix + r.OAuth2EmailName, Value: email},
		{Path: r.Prefix + r.AccountName, Value: account},
		{Path: r.Prefix + r.ActiveName, Value: true},
	}
	if len(r.updatedByName) > 0 {
		updateValue = append(updateValue, firestore.Update{Path: r.updatedByName, Value: id})
	}

	_, err = r.Collection.Doc(id).Update(ctx, updateValue)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *UserRepository) Insert(ctx context.Context, id string, user oauth2.User) (bool, error) {
	userMap := r.userToMap(ctx, id, user)
	_, err := r.Collection.Doc(id).Create(ctx, userMap)

	if err != nil {
		errMsg := err.Error()
		if strings.Index(errMsg, "Document already exists") >= 0 {
			return true, nil
		} else {
			return false, err
		}
	}
	return false, nil
}

func (r *UserRepository) userToMap(ctx context.Context, id string, user oauth2.User) map[string]interface{} {
	userMap := oauth2.UserToMap(ctx, id, user, r.GenderMapper, r.Schema)

	userMap[r.UserName] = user.Email
	userMap[r.StatusName] = r.ActivatedStatus

	userMap[r.Prefix+r.OAuth2EmailName] = user.Email
	userMap[r.Prefix+r.AccountName] = user.Account
	userMap[r.Prefix+r.ActiveName] = true
	return userMap
}
