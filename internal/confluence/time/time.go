// Package time provides functionality for working with dates and times in Confluence's expected formats.
package time

import (
	"fmt"
	gotime "time"

	"github.com/agilepathway/gauge-confluence/util"
)

// ConfluenceFormat is the format that Confluence uses to store times.
// The format is RFC3339 at millisecond precision.
//
// The format includes any trailing millisecond zeros, as Confluence stores any trailing zeros.
const confluenceFormat = "2006-01-02T15:04:05.000Z07:00"

// Time wraps Confluence time functionality around a time
type Time struct {
	time gotime.Time
}

// NewTime creates a new Time from the given time in Confluence format.
func NewTime(timeInConfluenceFormat string) Time {
	t, err := gotime.ParseInLocation(confluenceFormat, timeInConfluenceFormat, gotime.UTC)
	util.Fatal(fmt.Sprintf("Could not parse time: %s", timeInConfluenceFormat), err)

	return Time{t}
}

// CQLFormat formats a time so it can be used in Confluence CQL searches.
//
// See https://developer.atlassian.com/server/confluence/advanced-searching-using-cql/.
func (t Time) CQLFormat(cqlOffset int) string {
	cqlTime := t.cqlTime(cqlOffset)
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d",
		cqlTime.Year(), cqlTime.Month(), cqlTime.Day(), cqlTime.Hour(), cqlTime.Minute())
}

func (t Time) cqlTime(offset int) gotime.Time {
	d := gotime.Duration(abs(offset))
	u := t.time.UTC()

	switch {
	case offset == 0:
		return u
	case offset > 0:
		return u.Add(gotime.Hour * d)
	case offset < 0:
		return u.Add(-gotime.Hour * d)
	}

	panic("Unreachable code")
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func (t Time) String() string {
	return t.time.UTC().Format(confluenceFormat)
}

// Now returns the current UTC time (not the local time) in Confluence's expected time format.
func Now() Time {
	return Time{gotime.Now().UTC()}
}
