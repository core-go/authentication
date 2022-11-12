package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/core-go/auth"
	dyn "github.com/core-go/dynamodb"
	"strings"
	"time"

	"github.com/core-go/auth/oauth2"
)

type UserRepository struct {
	DB              *dynamodb.DynamoDB
	UserTableName   string
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

func NewUserRepositoryByConfig(db *dynamodb.DynamoDB, userTableName, prefix string, activatedStatus string, services []string, c oauth2.OAuth2SchemaConfig, status *auth.UserStatusConfig, options ...oauth2.OAuth2GenderMapper) *UserRepository {
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

	m := &UserRepository{
		DB:              db,
		UserTableName:   userTableName,
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

func NewUserRepository(db *dynamodb.DynamoDB, userTableName, prefix, activatedStatus string, services []string, pictureName, displayName, givenName, familyName, middleName, genderName string, status *auth.UserStatusConfig, options ...oauth2.OAuth2GenderMapper) *UserRepository {
	var genderMapper oauth2.OAuth2GenderMapper
	if len(options) >= 1 {
		genderMapper = options[0]
	}

	m := &UserRepository{
		DB:              db,
		UserTableName:   userTableName,
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

	projection := expression.NamesList(expression.Name("id"), expression.Name(r.StatusName))
	filter1 := expression.Equal(expression.Name(r.UserName), expression.Value(email))
	filter2 := expression.Equal(expression.Name(r.EmailName), expression.Value(email))
	filter3 := expression.Equal(expression.Name(r.Prefix+r.OAuth2EmailName), expression.Value(email))
	var sliceFilter []expression.ConditionBuilder
	sliceFilter = append(sliceFilter, filter3)

	for _, sv := range r.Services {
		if sv != r.Prefix {
			sliceFilter = append(sliceFilter, expression.Equal(expression.Name(sv+r.OAuth2EmailName), expression.Value(email)))
		}
	}

	filter := expression.Or(filter1, filter2, sliceFilter...)

	expr, _ := expression.NewBuilder().WithProjection(projection).WithFilter(filter).Build()
	query := &dynamodb.ScanInput{
		TableName:                 aws.String(r.UserTableName),
		ProjectionExpression:      expr.Projection(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}
	output, err := r.DB.ScanWithContext(ctx, query)
	disable := false
	suspended := false
	if err != nil {
		return "", disable, suspended, err
	}
	if len(output.Items) != 1 {
		return "", disable, suspended, err
	}
	var result map[string]string
	err = dynamodbattribute.UnmarshalMap(output.Items[0], &result)
	if err != nil {
		return "", disable, suspended, err
	}

	userId := result["id"]
	if r.Status != nil {
		status := result[r.StatusName]
		if status == r.Status.Disable {
			disable = true
		}
		if status == r.Status.Suspended {
			suspended = true
		}
	}

	return userId, disable, suspended, err

}

func (r *UserRepository) Update(ctx context.Context, id, email, account string) (bool, error) {

	user := make(map[string]interface{})

	user["id"] = id

	user[r.Prefix+r.OAuth2EmailName] = email
	user[r.Prefix+r.AccountName] = account
	user[r.Prefix+r.ActiveName] = true

	if len(r.updatedTimeName) > 0 {
		user[r.updatedTimeName] = time.Now()
	}

	if len(r.updatedByName) > 0 {
		user[r.updatedByName] = id
	}

	result, err:=  dyn.PatchOne(ctx,r.DB,r.UserTableName,[]string{"id"},user)
	return result > 0 , err

}

func (r *UserRepository) Insert(ctx context.Context, id string, user oauth2.User) (bool, error) {
	userMap := r.userToMap(ctx, id, user)

	_, err := dyn.InsertOne(ctx, r.DB, r.UserTableName, []string{"id"}, userMap)
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

	userMap["id"] = id
	userMap[r.UserName] = user.Email
	userMap[r.StatusName] = r.ActivatedStatus

	userMap[r.Prefix+r.OAuth2EmailName] = user.Email
	userMap[r.Prefix+r.AccountName] = user.Account
	userMap[r.Prefix+r.ActiveName] = true
	return userMap
}
