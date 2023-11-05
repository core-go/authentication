# Authentication
![Authentication](https://camo.githubusercontent.com/961908454560f4fdcd044a27e1741bc13d8440d794ad69d4fd6bf77023195701/68747470733a2f2f63646e2d696d616765732d312e6d656469756d2e636f6d2f6d61782f3830302f312a7652314a553030384e555234774b4567717766756f412e706e67)
- authenticator
- ldap authenticator
- 2 factor authentication
- oauth2
![oauth2](https://camo.githubusercontent.com/782b650c42e2a73f79e729e77176f3dbd5edf51b683e13ebdae0a6f5e4cdd7b2/68747470733a2f2f63646e2d696d616765732d312e6d656469756d2e636f6d2f6d61782f3830302f312a6153765054544461532d386c674f4164544d6e6335412e706e67)

## Installation
Please make sure to initialize a Go module before installing core-go/auth:

```shell
go get -u github.com/core-go/auth
```

Import:
```go
import "github.com/core-go/auth"
```

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

## OAuth2
### Models
- Configuration
- OAuth2Info
- User

### Services
- OAuth2Service
- Azure

### Repositories
- ConfigurationRepository
- OAuth2UserRepository
- UserRepository