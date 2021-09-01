package helper

import (
	"fmt"
	"time"
)

const (
	minute = 60
	hour   = 60 * 60
	day    = 60 * 60 * 24
	month  = 60 * 60 * 24 * 30
	year   = 60 * 60 * 24 * 365
)

// FormatDur returns a formatted duration
func FormatDur(duration time.Duration) string {
	seconds := uint64(duration.Seconds())
	if seconds < minute {
		return fmt.Sprintf("%d second%s", uint64(duration.Seconds()), mult(seconds > 1))
	} else if seconds < hour {
		return fmt.Sprintf("%d minute%s,%d second%s", seconds/minute, mult(seconds/minute > 1), seconds-(seconds/minute*minute), mult(seconds-(seconds/minute*minute) > 1))
	} else if seconds < day {
		return fmt.Sprintf("%d hour%s,%d minute%s", seconds/hour, mult(seconds/hour > 1), (seconds-(seconds/hour*hour))/minute, mult((seconds-(seconds/hour*hour))/minute > 1))
	} else if seconds < month {
		return fmt.Sprintf("%d day%s,%d hour%s", seconds/day, mult(seconds/day > 1), (seconds-(seconds/day*day))/hour, mult((seconds-(seconds/day*day))/hour > 1))
	} else if seconds < year {
		return fmt.Sprintf("%d month%s,%d day%s", seconds/month, mult(seconds/month > 1), (seconds-(seconds/month*month))/day, mult((seconds-(seconds/month*month))/day > 1))
	} else {
		return fmt.Sprintf("%d year%s,%d month%s,%d day%s", seconds/year, mult(seconds/year > 1), (seconds-(seconds/year*year))/month, mult((seconds-(seconds/year*year))/month > 1), (seconds-(seconds/month*month))/day, mult((seconds-(seconds/month*month))/day > 1))
	}
	return "unknown"
}

func mult(isMultiple bool) string {
	if isMultiple {
		return "s"
	}
	return ""
}
