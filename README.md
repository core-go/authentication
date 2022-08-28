# Authentication
![Authentication](https://camo.githubusercontent.com/961908454560f4fdcd044a27e1741bc13d8440d794ad69d4fd6bf77023195701/68747470733a2f2f63646e2d696d616765732d312e6d656469756d2e636f6d2f6d61782f3830302f312a7652314a553030384e555234774b4567717766756f412e706e67)
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
