# Authentication
![Authentication](https://camo.githubusercontent.com/a394ea3c13f690ecb9cf4a1747973cce1bdc8558e659040995003e96e486f88a/68747470733a2f2f63646e2d696d616765732d312e6d656469756d2e636f6d2f6d61782f3830302f312a56504f343261596a736c6d524937424c6369796a77412e706e67)
- authenticator
- ldap authenticator
- 2 factor authentication

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
