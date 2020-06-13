package auth

import "time"

func SetTokenExpiredTime(accessTimeFrom *time.Time, accessTimeTo *time.Time, expires int64) (time.Time, int64) {
	if accessTimeTo == nil || accessTimeFrom == nil {
		var tokenExpiredTime = time.Now().Add(time.Second * time.Duration(int(expires/1000)))
		return tokenExpiredTime, expires
	}
	if accessTimeTo.Before(*accessTimeFrom) || accessTimeTo.Equal(*accessTimeFrom) {
		tmp := accessTimeTo.Add(time.Hour * 24)
		accessTimeTo = &tmp
	}
	var tokenExpiredTime time.Time
	var jwtExpiredTime int64
	if time.Millisecond*time.Duration(expires) > accessTimeTo.Sub(time.Now()) {
		tokenExpiredTime = time.Now().Add(accessTimeTo.Sub(time.Now())).UTC()
		jwtExpiredTime = int64(accessTimeTo.Sub(time.Now()).Seconds() * 1000)
	} else {
		tokenExpiredTime = time.Now().Add(time.Millisecond * time.Duration(expires)).UTC()
		jwtExpiredTime = expires
	}
	return tokenExpiredTime, jwtExpiredTime
}


func IsAccessDateValid(fromDate, toDate *time.Time) bool {
	today := time.Now()
	if fromDate == nil && toDate == nil {
		return true
	} else if fromDate == nil {
		toDateStr := toDate.Add(time.Hour * 24)
		if toDateStr.After(today) {
			return true
		}
	} else if toDate == nil {
		if fromDate.Before(today) || fromDate.Equal(today) {
			return true
		}
	} else {
		toDateStr := toDate.Add(time.Hour * 24)
		if (fromDate.Before(today) || fromDate.Equal(today)) && toDateStr.After(today) {
			return true
		}
	}
	return false
}

func IsAccessTimeValid(fromTime, toTime *time.Time) bool {
	today := time.Now()
	location := time.Now().Location()
	if fromTime == nil && toTime == nil {
		return true
	} else if fromTime == nil {
		toTimeStr := toTime.In(location)
		if toTimeStr.After(today) || toTimeStr.Equal(today) {
			return true
		}
		return false
	} else if toTime == nil {
		fromTimeStr := fromTime.In(location)
		if fromTimeStr.Before(today) || fromTimeStr.Equal(today) {
			return true
		}
		return false
	}
	toTimeStr := toTime.In(location)
	fromTimeStr := fromTime.In(location)

	if toTimeStr.Before(fromTimeStr) || toTimeStr.Equal(fromTimeStr) {
		toTimeStr = toTimeStr.Add(time.Hour * 24)
	}
	if (fromTimeStr.Before(today) || fromTimeStr.Equal(today)) && (toTimeStr.After(today) || toTimeStr.Equal(today)) {
		return true
	}
	return false
}
