package rlock

import "time"

type FixedIntervalRetry struct {
	Interval time.Duration
	Cnt      uint8
	Max      uint8
}

func (f *FixedIntervalRetry) Next() (time.Duration, bool) {
	f.Cnt++
	return f.Interval, f.Cnt <= f.Max
}

func (f *FixedIntervalRetry) Reset() {
	f.Cnt = 0
}
