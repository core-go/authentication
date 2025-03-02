# Authentication
![Authentication](https://cdn-images-1.medium.com/max/800/1*vR1JU008NUR4wKEgqwfuoA.png)
- authenticator
- ldap authenticator
- 2 factor authentication
- oauth2

![oauth2](https://cdn-images-1.medium.com/max/800/1*aSvPTTDaS-8lgOAdTMnc5A.png)

## Installation
Please make sure to initialize a Go module before installing core-go/auth:

```shell
go get -u github.com/core-go/authentication
```

Import:
```go
import "github.com/core-go/authentication"
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