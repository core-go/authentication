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

Please make sure to initialize a Go module before installing common-go/auth:

```shell
go get -u github.com/common-go/auth
```

Import:

```go
import "github.com/common-go/auth"
```

## Implementations of AuthenticationRepository
- [sql](https://github.com/common-go/auth-sql): requires [gorm](https://github.com/go-gorm/gorm)
- [mongo](https://github.com/common-go/auth-mongo)
- [dynamodb](https://github.com/common-go/auth-dynamodb)
- [firestore](https://github.com/common-go/auth-firestore)
- [elasticsearch](https://github.com/common-go/auth-elasticsearch)

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
