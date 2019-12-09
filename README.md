# Authentication
## Models
- AuthInfo
- AuthResult
- UserAccount
- Privilege
- UserInfo
- StoredUser

## Services
- PrivilegeService
- TokenGenerator

## Config Model
- TokenConfig

## Installation

Please make sure to initialize a Go module before installing common-go/auth:

```shell
go get -u github.com/common-go/auth
```

Import:

```go
import "github.com/common-go/auth"
```

## Details:
#### auth_result.go
```go
type AuthResult struct {
	Status  AuthStatus   `json:"status,omitempty" bson:"status,omitempty" gorm:"column:status"`
	User    *UserAccount `json:"user,omitempty" bson:"user,omitempty" gorm:"column:user"`
	Message string       `json:"message,omitempty" bson:"message,omitempty" gorm:"column:message"`
}
```

#### user_account.go
```go
type UserAccount struct {
	UserId              string       `json:"userId,omitempty" bson:"userId,omitempty" gorm:"column:userid"`
	UserName            string       `json:"userName,omitempty" bson:"userName,omitempty" gorm:"column:username"`
	Email               string       `json:"email,omitempty" bson:"email,omitempty" gorm:"column:email"`
	DisplayName         string       `json:"displayName,omitempty" bson:"displayName,omitempty" gorm:"column:displayname"`
	PasswordExpiredTime *time.Time   `json:"passwordExpiredTime,omitempty" bson:"passwordExpiredTime,omitempty" gorm:"column:passwordexpiredtime"`
	Token               string       `json:"token,omitempty" bson:"token,omitempty" gorm:"column:token"`
	TokenExpiredTime    *time.Time   `json:"tokenExpiredTime,omitempty" bson:"tokenExpiredTime,omitempty" gorm:"column:tokenexpiredtime"`
	NewUser             bool         `json:"newUser,omitempty" bson:"newUser,omitempty" gorm:"column:newuser"`
	UserType            string       `json:"userType,omitempty" bson:"userType,omitempty" gorm:"column:usertype"`
	Roles               *[]string    `json:"roles,omitempty" bson:"roles,omitempty" gorm:"column:roles"`
	Privileges          *[]Privilege `json:"privileges,omitempty" bson:"privileges,omitempty" gorm:"column:privileges"`
}
```

#### privilege.go
```go
type Privilege struct {
	Id          string       `json:"id,omitempty" bson:"_id,omitempty" gorm:"column:id"`
	Name        string       `json:"name,omitempty" bson:"name,omitempty" gorm:"column:name"`
	ResourceKey string       `json:"resourceKey,omitempty" bson:"resourceKey,omitempty" gorm:"column:resourcekey"`
	Path        string       `json:"path,omitempty" bson:"path,omitempty" gorm:"column:path"`
	Icon        string       `json:"icon,omitempty" bson:"icon,omitempty" gorm:"column:icon"`
	Sequence    int          `json:"sequence,omitempty" bson:"sequence,omitempty" gorm:"column:sequence"`
	Children    *[]Privilege `json:"children,omitempty" bson:"children,omitempty" gorm:"column:children"`
}
```

#### user_info.go
```go
type UserInfo struct {
	UserId               string     `json:"userId,omitempty" bson:"userId,omitempty" gorm:"column:userid"`
	UserName             string     `json:"userName,omitempty" bson:"userName,omitempty" gorm:"column:username"`
	Email                string     `json:"email,omitempty" bson:"email,omitempty" gorm:"column:email"`
	DisplayName          string     `json:"displayName,omitempty" bson:"displayName,omitempty" gorm:"column:displayname"`
	Password             string     `json:"password,omitempty" bson:"password,omitempty" gorm:"column:password"`
	Disable              bool       `json:"disable,omitempty" bson:"disable,omitempty" gorm:"column:disable"`
	Deactivated          bool       `json:"deactivated,omitempty" bson:"deactivated,omitempty" gorm:"column:deactivated"`
	Suspended            bool       `json:"blocked,omitempty" bson:"blocked,omitempty" gorm:"column:suspended"`
	LockedUntilTime      *time.Time `json:"lockedUntilTime,omitempty" bson:"lockedUntilTime,omitempty" gorm:"column:lockeduntiltime"`
	SuccessTime          *time.Time `json:"successTime,omitempty" bson:"successTime,omitempty" gorm:"column:successtime"`
	FailTime             *time.Time `json:"failTime,omitempty" bson:"failTime,omitempty" gorm:"column:failtime"`
	FailCount            int        `json:"failCount,omitempty" bson:"failCount,omitempty" gorm:"column:failcount"`
	PasswordModifiedTime *time.Time `json:"passwordModifiedTime,omitempty" bson:"passwordModifiedTime,omitempty" gorm:"column:passwordmodifiedtime"`
	MaxPasswordAge       int        `json:"maxPasswordAge,omitempty" bson:"maxPasswordAge,omitempty" gorm:"column:maxpasswordage"`
	UserType             string     `json:"userType,omitempty" bson:"userType,omitempty" gorm:"column:usertype"`
	Roles                *[]string  `son:"roles,omitempty" bson:"roles,omitempty" gorm:"column:roles"`
	AccessDateFrom       *time.Time `json:"accessDateFrom,omitempty" bson:"accessDateFrom,omitempty" gorm:"type:date;column:accessdatefrom"`
	AccessDateTo         *time.Time `json:"accessDateTo,omitempty" bson:"accessDateTo,omitempty" gorm:"column:accessDateTo"`
	AccessTimeFrom       *time.Time `json:"accessTimeFrom,omitempty" bson:"accessTimeFrom,omitempty" gorm:"type:date;column:accesstimefrom"`
	AccessTimeTo         *time.Time `json:"accessTimeTo,omitempty" bson:"accessTimeTo,omitempty" gorm:"column:accesstimeto"`
}
```