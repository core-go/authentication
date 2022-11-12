package oauth2

import "context"

type OAuth2GenderMapper interface {
	Map(ctx context.Context, gender string) interface{}
}
