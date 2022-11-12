package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/core-go/auth/oauth2"
	"reflect"
)

type ConfigurationRepository struct {
	DB                     *sql.DB
	TableName              string
	OAuth2UserRepositories map[string]oauth2.OAuth2UserRepository
	Status                 string
	Active                 string
	Driver                 string
	BuildParam             func(i int) string
	configurationFields    map[string]int
}

func NewConfigurationRepository(db *sql.DB, tableName string, oAuth2PersonInfoServices map[string]oauth2.OAuth2UserRepository, status string, active string) (*ConfigurationRepository, error) {
	if len(status) == 0 {
		status = "status"
	}
	if len(active) == 0 {
		active = "A"
	}
	build := getBuild(db)
	driver := getDriver(db)
	var configuration oauth2.Configuration
	configurationType := reflect.TypeOf(configuration)
	configurationFields, err := getColumnIndexes(configurationType)
	if err != nil {
		return nil, err
	}
	return &ConfigurationRepository{DB: db, TableName: tableName, OAuth2UserRepositories: oAuth2PersonInfoServices, Status: status, Active: active, Driver: driver, BuildParam: build, configurationFields: configurationFields}, nil
}

func (s *ConfigurationRepository) GetConfiguration(ctx context.Context, id string) (*oauth2.Configuration, string, error) {
	var configurations []oauth2.Configuration
	limitRowsQL := "limit 1"
	driver := getDriver(s.DB)
	if driver == driverOracle {
		limitRowsQL = "and rownum = 1"
	}
	query := fmt.Sprintf(`select * from %s where %s = %s %s`, s.TableName, "id", s.BuildParam(0), limitRowsQL)
	err := queryWithMap(ctx, s.DB, s.configurationFields, &configurations, query, id)
	if err != nil {
		return nil, "", err
	}
	if len(configurations) == 0 {
		return nil, "", nil
	}
	model := configurations[0]
	clientId := model.ClientId
	clientId, err = s.OAuth2UserRepositories[id].GetRequestTokenOAuth(ctx, model.ClientId, model.ClientSecret)
	return &model, clientId, err
}

func (s *ConfigurationRepository) GetConfigurations(ctx context.Context) ([]oauth2.Configuration, error) {
	var configurations []oauth2.Configuration
	query := fmt.Sprintf(`select * from %s where %s = %s `, s.TableName, s.Status, s.BuildParam(1))
	err := queryWithMap(ctx, s.DB, s.configurationFields, &configurations, query, s.Active)
	if err != nil {
		return nil, err
	}
	return configurations, nil
}
