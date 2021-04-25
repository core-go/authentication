package sql

import (
	"context"
	"database/sql"
	"reflect"
	"time"

	a "github.com/common-go/auth"
)

type SqlConfig struct {
	Query           string `mapstructure:"query" json:"query,omitempty" gorm:"column:query" bson:"query,omitempty" dynamodbav:"query,omitempty" firestore:"query,omitempty"`
	SqlPass         string `mapstructure:"pass" json:"pass,omitempty" gorm:"column:pass" bson:"pass,omitempty" dynamodbav:"pass,omitempty" firestore:"pass,omitempty"`
	SqlFail         string `mapstructure:"fail" json:"fail,omitempty" gorm:"column:fail" bson:"fail,omitempty" dynamodbav:"fail,omitempty" firestore:"fail,omitempty"`
	DisableStatus   string `mapstructure:"disable" json:"disable,omitempty" gorm:"column:disable" bson:"disable,omitempty" dynamodbav:"disable,omitempty" firestore:"disable,omitempty"`
	SuspendedStatus string `mapstructure:"suspended" json:"suspended,omitempty" gorm:"column:suspended" bson:"suspended,omitempty" dynamodbav:"suspended,omitempty" firestore:"suspended,omitempty"`
	NoTime          bool   `mapstructure:"no_time" json:"noTime,omitempty" gorm:"column:notime" bson:"noTime,omitempty" dynamodbav:"noTime,omitempty" firestore:"noTime,omitempty"`
}
type SqlUserInfoService struct {
	DB              *sql.DB
	Query           string
	SqlPass         string
	SqlFail         string
	DisableStatus   string
	SuspendedStatus string
	NoTime          bool
	Driver          string
}

func NewSqlUserInfoService(db *sql.DB, query, sqlPass, sqlFail, disableStatus string, suspendedStatus string, noTime bool, options...bool) *SqlUserInfoService {
	var handleDriver bool
	if len(options) >= 1 {
		handleDriver = options[0]
	} else {
		handleDriver = true
	}
	driver := getDriver(db)
	if handleDriver {
		query = replaceQueryArgs(driver, query)
		sqlPass = replaceQueryArgs(driver, sqlPass)
	}
	return &SqlUserInfoService{DB: db, Query: query, SqlPass: sqlPass, SqlFail: sqlFail, DisableStatus: disableStatus, SuspendedStatus: suspendedStatus, NoTime: noTime, Driver: driver}
}

func NewSqlUserInfoByConfig(db *sql.DB, c SqlConfig, options...bool) *SqlUserInfoService {
	return NewSqlUserInfoService(db, c.Query, c.SqlPass, c.SqlFail, c.DisableStatus, c.SuspendedStatus, c.NoTime, options...)
}
func (l SqlUserInfoService) GetUserInfo(ctx context.Context, auth a.AuthInfo) (*a.UserInfo, error) {
	models := make([]a.UserInfo, 0)
	rows, er1 := l.DB.QueryContext(ctx, l.Query, auth.Username)
	if er1 != nil {
		switch er1 {
		case sql.ErrNoRows:
			return nil, nil
		default:
			return nil, er1
		}
	}
	defer rows.Close()
	modelTypes := reflect.TypeOf(models).Elem()
	modelType := reflect.TypeOf(a.UserInfo{})
	columns, er2 := rows.Columns()
	if er2 != nil {
		return nil, er2
	}
	// get list indexes column
	indexes, er3 := getColumnIndexes(modelType, columns, l.Driver)
	if er3 != nil {
		return nil, er3
	}
	tb, er4 := scanType(rows, modelTypes, indexes)
	if er4 != nil {
		return nil, er4
	}
	if len(tb) > 0 {
		if c, ok := tb[0].(*a.UserInfo); ok {
			if len(c.Status) > 0 {
				if c.Status == l.SuspendedStatus {
					c.Suspended = true
				}
				if c.Status == l.DisableStatus {
					c.Disable = true
				}
			}
			return c, nil
		}
	}
	return nil, nil
}
func (l SqlUserInfoService) Pass(ctx context.Context, user a.UserInfo) error {
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
func (l SqlUserInfoService) Fail(ctx context.Context, user a.UserInfo) error {
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
