package common

import (
	"time"
)

// RealClock - Realization of ClockInterface that returns real time
type RealClock struct{}

// Now returns current time
func (r *RealClock) Now() time.Time {
	return time.Now()
}
