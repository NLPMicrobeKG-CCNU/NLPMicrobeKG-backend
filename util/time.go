package util

import "time"

func GetCurrentTime() *time.Time {
	t := time.Now().UTC().Add(8 * time.Hour)
	return &t
}

func GetCurrentDate() string {
	date, _ := FormatTime(*GetCurrentTime())
	return date
}

func FormatTime(t time.Time) (string, string) {
	return t.Format("2006-01-02"), t.Format("15:04:05")
}
