package protodate // import "go.pedge.io/proto/date"

import (
	"time"

	"go.pedge.io/googleapis/google/type"
)

var (
	// SystemDater is a Dater that uses the system UTC time.
	SystemDater = &systemDater{}
)

// Now returns the current date at UTC.
func Now() *google_type.Date {
	return TimeToDate(time.Now().UTC())
}

// TimeToDate converts a golang Time to a Date.
func TimeToDate(t time.Time) *google_type.Date {
	return &google_type.Date{
		Day:   int32(t.Day()),
		Month: int32(t.Month()),
		Year:  int32(t.Year()),
	}
}

// DateToTime converts a Date to a golang Time.
func DateToTime(d *google_type.Date) time.Time {
	if d == nil {
		return time.Unix(0, 0).UTC()
	}
	return time.Date(int(d.Year), time.Month(d.Month), int(d.Day), 0, 0, 0, 0, time.UTC)
}

// Dater provides the current date.
type Dater interface {
	Now() *google_type.Date
}

type systemDater struct{}

func (s *systemDater) Now() *google_type.Date {
	return Now()
}
