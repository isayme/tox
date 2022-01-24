package util

import "time"

func NowInMills() int64 {
	now := time.Now()
	return now.UnixMilli()
}
