package sql

import (
	"context"
	"database/sql"
	a "github.com/core-go/auth"
	"reflect"
	"time"
)

type SqlConfig struct {
	Query           string `yaml:"query" mapstructure:"query" json:"query,omitempty" gorm:"column:query" bson:"query,omitempty" dynamodbav:"query,omitempty" firestore:"query,omitempty"`
	SqlPass         string `yaml:"pass" mapstructure:"pass" json:"pass,omitempty" gorm:"column:pass" bson:"pass,omitempty" dynamodbav:"pass,omitempty" firestore:"pass,omitempty"`
	SqlFail         string `yaml:"fail" mapstructure:"fail" json:"fail,omitempty" gorm:"column:fail" bson:"fail,omitempty" dynamodbav:"fail,omitempty" firestore:"fail,omitempty"`
	DisableStatus   string `yaml:"disable" mapstructure:"disable" json:"disable,omitempty" gorm:"column:disable" bson:"disable,omitempty" dynamodbav:"disable,omitempty" firestore:"disable,omitempty"`
	SuspendedStatus string `yaml:"suspended" mapstructure:"suspended" json:"suspended,omitempty" gorm:"column:suspended" bson:"suspended,omitempty" dynamodbav:"suspended,omitempty" firestore:"suspended,omitempty"`
	NoTime          bool   `yaml:"no_time" mapstructure:"no_time" json:"noTime,omitempty" gorm:"column:notime" bson:"noTime,omitempty" dynamodbav:"noTime,omitempty" firestore:"noTime,omitempty"`
}
type SqlUserRepository struct {
	DB              *sql.DB
	Query           string
	SqlPass         string
	SqlFail         string
	Status          a.UserStatusConfig
	MaxPasswordAge  int32
	NoTime          bool
	Driver          string
	userFields      map[string]int
	Param           func(int) string
	Id              string
	Password        string
	FailTime        string
	FailCount       string
	LockedUntilTime string
}

func NewSqlUserRepository(db *sql.DB, query, sqlPass, sqlFail, disableStatus string, suspendedStatus string, noTime bool, maxPasswordAge int32, options ...bool) (*SqlUserRepository, error) {
	var handleDriver bool
	if len(options) >= 1 {
		handleDriver = options[0]
	} else {
		handleDriver = true
	}
	driver := getDriver(db)
	var param func(int) string
	if handleDriver {
		query = replaceQueryArgs(driver, query)
		sqlPass = replaceQueryArgs(driver, sqlPass)
		param = GetBuildByDriver(driver)
	}
	var user a.UserInfo
	userType := reflect.TypeOf(user)
	userFields, err := getColumnIndexes(userType)
	if err != nil {
		return nil, err
	}
	return &SqlUserRepository{DB: db, Query: query, SqlPass: sqlPass, SqlFail: sqlFail, DisableStatus: disableStatus, SuspendedStatus: suspendedStatus, NoTime: noTime, MaxPasswordAge: maxPasswordAge, Driver: driver, userFields: userFields, Param: param}, nil
}

func NewSqlUserInfoByConfig(db *sql.DB, c SqlConfig, options ...bool) (*SqlUserRepository, error) {
	return NewSqlUserRepository(db, c.Query, c.SqlPass, c.SqlFail, c.DisableStatus, c.SuspendedStatus, c.NoTime, 0, options...)
}
func (l SqlUserRepository) GetUser(ctx context.Context, auth a.AuthInfo) (*a.UserInfo, error) {
	var models []a.UserInfo
	_, err := queryWithMap(ctx, l.DB, l.userFields, &models, l.Query, auth.Username)
	if err != nil {
		return nil, err
	}
	if len(models) > 0 {
		c := models[0]
		if len(c.Status) > 0 {
			if c.Status == l.Status.Deactivated {
				b := true
				c.Deactivated = &b
			}
			if c.Status == l.Status.Suspended {
				c.Suspended = true
			}
			if c.Status == l.Status.Disable {
				c.Disable = true
			}
		}
		return &c, nil
	}
	return nil, nil
}
func (l SqlUserRepository) Pass(ctx context.Context, user a.UserInfo) error {
	if len(l.SqlPass) == 0 {
		return nil
	}
	if l.NoTime {
		_, err := l.DB.ExecContext(ctx, l.SqlPass, user.Id)
		return err
	} else {
		_, err := l.DB.ExecContext(ctx, l.SqlPass, time.Now(), user.Id)
		return err
	}
}
func (l SqlUserRepository) Fail(ctx context.Context, id, failCount *int, lockedUntilTime *time.Time) error {
	if failCount == nil && lockedUntilTime == nil {
		return nil
	}
	if len(l.SqlFail) == 0 {
		return nil
	}
	if l.NoTime {
		_, err := l.DB.ExecContext(ctx, l.SqlFail, user.Id)
		return err
	} else {
		_, err := l.DB.ExecContext(ctx, l.SqlFail, time.Now(), user.Id)
		return err
	}
}
