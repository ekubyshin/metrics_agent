package utils

import "time"

func IntToDuration(val int) time.Duration {
	return time.Duration(val) * time.Second
}
