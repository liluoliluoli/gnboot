package time_util

import "time"

func FormatYYYYMMDD(t time.Time) string {
	return t.Format("20060102")
}

func FormatStrToYYYYMMDD(timeStr string) string {
	t, err := time.ParseInLocation("2006-01-02T15:04:05.0000000Z", timeStr, time.Local)
	if err != nil {
		return ""
	}
	return t.Format("20060101")
}

func ParseUtcTime(timeStr string) time.Time {
	parsedTime, _ := time.ParseInLocation("2006-01-02T15:04:05.0000000Z", timeStr, time.Local)
	t := parsedTime.Truncate(time.Second)
	return t
}
