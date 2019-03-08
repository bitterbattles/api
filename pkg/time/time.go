package time

import "time"

// NowUnix returns the current Unix timestamp
func NowUnix() int64 {
	return time.Now().Unix()
}
