# Authentication
## Models
- AuthInfo
- AuthResult
- UserAccount
- Privilege
- UserInfo
- StoredUser

## Services
- Authenticator
- PrivilegesLoader
- UserInfoService

## Token
- TokenConfig
- TokenGenerator

## Installation
Please make sure to initialize a Go module before installing core-go/auth:

```shell
go get -u github.com/core-go/auth
```

Import:
```go
import "github.com/core-go/auth"
```

## Details:
#### authenticator.go
```go
type Authenticator interface {
	Authenticate(ctx context.Context, user AuthInfo) (AuthResult, error)
}
```

#### privileges_loader.go
```go
type PrivilegesLoader interface {
	Load(ctx context.Context, id string) ([]Privilege, error)
}
```

#### user_info_service.go
```go
type UserInfoService interface {
	GetUserInfo(ctx context.Context, auth AuthInfo) (*UserInfo, error)
	Pass(ctx context.Context, user UserInfo) (bool, error)
	Fail(ctx context.Context, user UserInfo) (bool, error)
}
```
