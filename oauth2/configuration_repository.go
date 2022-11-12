package oauth2

import "context"

type ConfigurationRepository interface {
	GetConfiguration(ctx context.Context, id string) (*Configuration, string, error)
	GetConfigurations(ctx context.Context) ([]Configuration, error)
}
