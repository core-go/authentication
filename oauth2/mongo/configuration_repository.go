package mongo

import (
	"context"
	"fmt"
	"github.com/core-go/authentication/oauth2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type ConfigurationRepository struct {
	Collection             *mongo.Collection
	OAuth2UserRepositories map[string]oauth2.OAuth2UserRepository
	Status                 string
	Active                 string
}

func NewConfigurationRepository(db *mongo.Database, collectionName string, oauth2UserRepositories map[string]oauth2.OAuth2UserRepository, status string, active string) *ConfigurationRepository {
	if len(status) == 0 {
		status = "status"
	}
	if len(active) == 0 {
		active = "A"
	}
	collection := db.Collection(collectionName)
	return &ConfigurationRepository{Collection: collection, OAuth2UserRepositories: oauth2UserRepositories, Status: status, Active: active}
}

func (s *ConfigurationRepository) GetConfiguration(ctx context.Context, id string) (*oauth2.Configuration, string, error) {
	var model oauth2.Configuration
	query := bson.M{"_id": id}
	x := s.Collection.FindOne(ctx, query)
	if x.Err() != nil {
		if strings.Compare(fmt.Sprint(x.Err()), "mongo: no documents in result") == 0 {
			return nil, "", nil
		}
		return nil, "", x.Err()
	}
	k := &model
	err := x.Decode(k)
	if err != nil {
		return nil, "", err
	}

	clientId := model.ClientId
	k.ClientId, err = s.OAuth2UserRepositories[id].GetRequestTokenOAuth(ctx, model.ClientId, model.ClientSecret)
	return k, clientId, err
}
func (s *ConfigurationRepository) GetConfigurations(ctx context.Context) ([]oauth2.Configuration, error) {
	var configurations []oauth2.Configuration
	query := bson.M{}
	cursor, er1 := s.Collection.Find(ctx, query)
	if er1 != nil {
		return nil, er1
	}
	er2 := cursor.All(ctx, &configurations)
	if er2 != nil {
		return nil, er2
	}
	return configurations, nil
}
