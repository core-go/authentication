package auth

import "time"

func addSeconds(date time.Time, seconds int64) time.Time {
	return date.Add(time.Second * time.Duration(seconds))
}

func addDays(date time.Time, days int32) time.Time {
	return date.Add(time.Hour * time.Duration(days) * 24)
}

func compareDate(date1 time.Time, date2 time.Time) int {
	return int(date1.Sub(date2).Seconds())
}
