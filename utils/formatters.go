package utils

import "time"

func FormatTime(t time.Time) string {
	return t.Format("15:04:05")
}
