package util

import "time"

func TimeAddr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func TimeValue(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
