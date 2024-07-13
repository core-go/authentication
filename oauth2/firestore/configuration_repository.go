package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"strings"

	"github.com/core-go/authentication/oauth2"
)

type ConfigurationRepository struct {
	Collection             *firestore.CollectionRef
	OAuth2UserRepositories map[string]oauth2.OAuth2UserRepository
	Status                 string
	Active                 string
}

func NewConfigurationRepository(db *firestore.Client, collectionName string, oAuth2PersonInfoServices map[string]oauth2.OAuth2UserRepository, status string, active string) *ConfigurationRepository {
	collection := db.Collection(collectionName)
	return &ConfigurationRepository{Collection: collection, OAuth2UserRepositories: oAuth2PersonInfoServices, Status: status, Active: active}
}

func (s *ConfigurationRepository) GetConfiguration(ctx context.Context, sourceType string) (*oauth2.Configuration, string, error) {
	var model oauth2.Configuration
	doc, err := s.Collection.Doc(sourceType).Get(ctx)
	if !doc.Exists() {
		return nil, "", err
	}
	if err != nil {
		if strings.Index(err.Error(), "no more items in iterator") >= 0 {
			return nil, "", nil
		}
		return nil, "", err
	}

	k := &model
	err = doc.DataTo(k)
	if err != nil {
		return nil, "", err
	}

	clientId := model.ClientId
	k.ClientId, err = s.OAuth2UserRepositories[sourceType].GetRequestTokenOAuth(ctx, model.ClientId, model.ClientSecret)
	return k, clientId, err
}
func (s *ConfigurationRepository) GetConfigurations(ctx context.Context) ([]oauth2.Configuration, error) {
	arr := make([]oauth2.Configuration, 0)
	q := s.Collection.Where(s.Status, "=", s.Active)
	iter := q.Documents(ctx)
	for {
		doc, er1 := iter.Next()
		if er1 == iterator.Done {
			break
		}
		if er1 != nil {
			return nil, er1
		}
		var configuration oauth2.Configuration
		er2 := doc.DataTo(&configuration)
		if er2 != nil {
			return nil, er2
		}
		arr = append(arr, configuration)
	}
	return arr, nil
}
