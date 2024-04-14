package sql

import (
	"context"
	"database/sql"
	"fmt"
	a "github.com/core-go/auth"
	"reflect"
	"strings"
	"time"
)

type SqlUserRepository struct {
	DB             *sql.DB
	Query          string
	Status         a.UserStatusConfig
	MaxPasswordAge int32
	NoTime         bool
	Driver         string
	userFields     map[string]int
	Param          func(int) string
	Conf           DBConfig
}
func NewUserAdapter(db *sql.DB, query string, conf DBConfig, status a.UserStatusConfig, options ...bool) (*SqlUserRepository, error) {
	return NewUserRepository(db, query, conf, status, options...)
}
func NewUserRepository(db *sql.DB, query string, conf DBConfig, status a.UserStatusConfig, options ...bool) (*SqlUserRepository, error) {
	var handleDriver bool
	if len(options) >= 1 {
		handleDriver = options[0]
	} else {
		handleDriver = true
	}
	if len(conf.Password) == 0 {
		conf.Password = conf.User
	}
	driver := getDriver(db)
	var param func(int) string
	if handleDriver {
		query = replaceQueryArgs(driver, query)
		param = GetBuildByDriver(driver)
	}
	var user a.UserInfo
	userType := reflect.TypeOf(user)
	userFields, err := getColumnIndexes(userType)
	if err != nil {
		return nil, err
	}
	return &SqlUserRepository{DB: db, Query: query, Driver: driver, userFields: userFields, Param: param, Conf: conf, Status:  status}, nil
}
func (l SqlUserRepository) GetUser(ctx context.Context, username string) (*a.UserInfo, error) {
	var models []a.UserInfo
	_, err := queryWithMap(ctx, l.DB, l.userFields, &models, l.Query, username)
	if err != nil {
		return nil, err
	}
	if len(models) > 0 {
		c := models[0]
		sts := c.Status
		st := l.Status
		if sts != nil && len(*sts) > 0 {
			if c.Status == &st.Deactivated {
				b := true
				c.Deactivated = &b
			}
			if c.Status == &st.Suspended {
				c.Suspended = true
			}
			if c.Status == &st.Disable {
				c.Disable = true
			}
		}
		mpa := l.MaxPasswordAge
		if mpa > 0 && c.MaxPasswordAge == nil {
			c.MaxPasswordAge = &mpa
		}
		return &c, nil
	}
	return nil, nil
}
func (l SqlUserRepository) Pass(ctx context.Context, id string, deactivated *bool) error {
	if len(l.Conf.User) == 0 && len(l.Conf.Password) == 0 {
		return nil
	}
	now := time.Now()
	i := 1
	cols := make([]string, 0)
	params := make([]interface{}, 0)
	if len(l.Conf.SuccessTime) > 0 {
		cols = append(cols, fmt.Sprintf("%s=%s", l.Conf.SuccessTime, l.Param(i)))
		params = append(params, now)
		i = i + 1
	}
	if len(l.Conf.FailTime) > 0 {
		cols = append(cols, fmt.Sprintf("%s=null", l.Conf.FailTime))
	}
	if len(l.Conf.FailCount) > 0 {
		cols = append(cols, fmt.Sprintf("%s=0", l.Conf.FailCount))
	}
	if len(l.Conf.LockedUntilTime) > 0 {
		cols = append(cols, fmt.Sprintf("%s=null", l.Conf.LockedUntilTime))
	}
	if len(l.Conf.Status) == 0 || len(l.Status.Activated) == 0 || (deactivated != nil && *deactivated == false) {
		if len(cols) == 0 {
			return nil
		}
		params = append(params, id)
		sql := fmt.Sprintf(`update %s set %s where %s = %s`, l.Conf.Password, strings.Join(cols, ","), l.Conf.Id, l.Param(i))
		_, err := l.DB.ExecContext(ctx, sql, params...)
		return err
	}
	if l.Conf.User == l.Conf.Password {
		cols = append(cols, fmt.Sprintf("%s=%s", l.Conf.Status, l.Param(i)))
		params = append(params, l.Status.Activated)
		i = i + 1

		params = append(params, id)
		sql := fmt.Sprintf(`update %s set %s where %s = %s`, l.Conf.Password, strings.Join(cols, ","), l.Conf.Id, l.Param(i))
		_, err := l.DB.ExecContext(ctx, sql, params...)
		return err
	}

	sqlU := fmt.Sprintf(`update %s set %s = %s where %s = %s`, l.Conf.User, l.Conf.Status, l.Param(1), l.Conf.Id, l.Param(2))
	tx, err := l.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, sqlU, l.Status.Activated, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(cols) > 0 {
		params = append(params, id)
		sql := fmt.Sprintf(`update %s set %s where %s = %s`, l.Conf.Password, strings.Join(cols, ","), l.Conf.Id, l.Param(i))
		_, err = tx.ExecContext(ctx, sql, params...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
func (l SqlUserRepository) Fail(ctx context.Context, id string, failCount *int, lockedUntilTime *time.Time) error {
	if len(l.Conf.User) == 0 && len(l.Conf.Password) == 0 {
		return nil
	}
	now := time.Now()
	i := 1
	cols := make([]string, 0)
	params := make([]interface{}, 0)
	if len(l.Conf.FailTime) > 0 {
		cols = append(cols, fmt.Sprintf("%s=%s", l.Conf.FailTime, l.Param(i)))
		params = append(params, now)
		i = i + 1
	}
	if len(l.Conf.FailCount) > 0 && failCount != nil {
		count := *failCount + 1
		cols = append(cols, fmt.Sprintf("%s=%d", l.Conf.FailCount, count))
	}
	if len(l.Conf.LockedUntilTime) > 0 && lockedUntilTime != nil {
		cols = append(cols, fmt.Sprintf("%s=%s", l.Conf.FailCount, l.Param(i)))
		params = append(params, lockedUntilTime)
		i = i + 1
	}
	if len(cols) == 0 {
		return nil
	}
	params = append(params, id)
	sql := fmt.Sprintf(`update %s set %s where %s = %s`, l.Conf.Password, strings.Join(cols, ","), l.Conf.Id, l.Param(i))
	_, err := l.DB.ExecContext(ctx, sql, params...)
	return err
}
