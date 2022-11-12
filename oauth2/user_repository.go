package oauth2

import "context"

type UserRepository interface {
	GetUser(ctx context.Context, email string) (string, bool, bool, error)
	Update(ctx context.Context, id, email, account string) (bool, error)
	Insert(ctx context.Context, id string, user *User) (bool, error)
}
