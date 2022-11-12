package oauth2

import (
	"context"
	"time"
)

func UserToMap(ctx context.Context, id string, user User, genderMapper OAuth2GenderMapper, c *OAuth2SchemaConfig) map[string]interface{} {
	userMap := make(map[string]interface{})
	if c == nil {
		return userMap
	}
	if len(c.Picture) > 0 && len(user.Picture) > 0 {
		userMap[c.Picture] = user.Picture
	}
	if len(c.DisplayName) > 0 && len(user.DisplayName) > 0 {
		userMap[c.DisplayName] = user.DisplayName
	}
	if len(c.GivenName) > 0 && len(user.GivenName) > 0 {
		userMap[c.GivenName] = user.GivenName
	}
	if len(c.FamilyName) > 0 && len(user.FamilyName) > 0 {
		userMap[c.FamilyName] = user.FamilyName
	}
	if len(c.MiddleName) > 0 && len(user.MiddleName) > 0 {
		userMap[c.MiddleName] = user.MiddleName
	}
	if len(c.Gender) > 0 && user.Gender != nil {
		if genderMapper != nil {
			userMap[c.Gender] = genderMapper.Map(ctx, *user.Gender)
		} else  {
			userMap[c.Gender] = user.Gender
		}
	}

	now := time.Now()
	if len(c.CreatedTime) > 0 {
		userMap[c.CreatedTime] = now
	}
	if len(c.UpdatedTime) > 0 {
		userMap[c.UpdatedTime] = now
	}
	if len(c.CreatedBy) > 0 {
		userMap[c.CreatedBy] = id
	}
	if len(c.UpdatedBy) > 0 {
		userMap[c.UpdatedBy] = id
	}
	if len(c.Version) > 0 {
		userMap[c.Version] = 1
	}
	return userMap
}
