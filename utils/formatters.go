package utils

import (
	"fmt"
	"strings"
	"time"
)

func FormatTime(t time.Time) string {
	return t.Format("15:04:05")
}

func SliceToString(slice []string) string {
	result := strings.Builder{}
	for _, v := range slice {
		result.WriteString(fmt.Sprintf(" %s", v))
	}
	return result.String()
}
