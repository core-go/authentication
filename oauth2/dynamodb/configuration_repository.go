package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/core-go/authentication/oauth2"
)

type ConfigurationRepository struct {
	DB                     *dynamodb.DynamoDB
	ConfigurationTableName string
	OAuth2UserRepositories map[string]oauth2.OAuth2UserRepository
	Status                 string
	Active                 string
}

func NewConfigurationRepository(db *dynamodb.DynamoDB, configurationTableName string, oauth2UserRepositories map[string]oauth2.OAuth2UserRepository, status string, active string) *ConfigurationRepository {
	if len(status) == 0 {
		status = "status"
	}
	if len(active) == 0 {
		active = "A"
	}
	return &ConfigurationRepository{DB: db, ConfigurationTableName: configurationTableName, OAuth2UserRepositories: oauth2UserRepositories, Status: status, Active: active}
}

func (s *ConfigurationRepository) GetConfiguration(ctx context.Context, id string) (*oauth2.Configuration, string, error) {
	var model oauth2.Configuration

	filter := expression.Equal(expression.Name("id"), expression.Value(id))
	expr, _ := expression.NewBuilder().WithFilter(filter).Build()
	query := &dynamodb.ScanInput{
		TableName:                 aws.String(s.ConfigurationTableName),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}
	output, err := s.DB.ScanWithContext(ctx, query)
	k := &model
	err = dynamodbattribute.UnmarshalMap(output.Items[0], k)
	if err != nil {
		return nil, "", err
	}
	clientId := model.ClientId
	k.ClientId, err = s.OAuth2UserRepositories[id].GetRequestTokenOAuth(ctx, model.ClientId, model.ClientSecret)
	return k, clientId, err
}
func (s *ConfigurationRepository) GetConfigurations(ctx context.Context) ([]oauth2.Configuration, error) {
	var models []oauth2.Configuration
	var model oauth2.Configuration
	filter := expression.Equal(expression.Name(s.Status), expression.Value(s.Active))
	expr, _ := expression.NewBuilder().WithFilter(filter).Build()
	query := &dynamodb.ScanInput{
		TableName:                 aws.String(s.ConfigurationTableName),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}
	output, err := s.DB.ScanWithContext(ctx, query)
	for _, ele := range output.Items {
		_ = dynamodbattribute.UnmarshalMap(ele, model)
		models = append(models, model)
	}
	err = dynamodbattribute.UnmarshalMap(output.Items[0], models)
	if err != nil {
		return nil, err
	}
	return models, nil
}
