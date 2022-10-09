package cronworker

import (
	"time"

	"github.com/gorhill/cronexpr"
)

// parseAtTime with input format using crontab standard, for example https://crontab.guru
func parseAtTime(t string) (duration, nextDuration time.Duration, err error) {
	now := time.Now()

	cronExpr := cronexpr.MustParse(t)

	// Next time from now
	next := cronExpr.Next(time.Now())

	// Subtract next time with now
	duration = next.Sub(now)

	// Set nextDuration
	nextDuration = cronExpr.Next(next).Sub(next)

	return
}
