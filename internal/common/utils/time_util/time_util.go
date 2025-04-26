package time_util

import "time"

func FormatYYYYMMDD(t time.Time) string {
	return t.Format("20060102")
}
