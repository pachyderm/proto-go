package protodate // import "go.pedge.io/proto/date"

import (
	"sync"
	"time"

	"go.pedge.io/googleapis/google/type"
)

var (
	// SystemDater is a Dater that uses the system UTC time.
	SystemDater = &systemDater{}
)

// NewDate is a convienence function to create a new Date.
func NewDate(month int32, day int32, year int32) *google_type.Date {
	return &google_type.Date{
		Month: month,
		Day:   day,
		Year:  year,
	}
}

// Now returns the current date at UTC.
func Now() *google_type.Date {
	return TimeToDate(time.Now().UTC())
}

// TimeToDate converts a golang Time to a Date.
func TimeToDate(t time.Time) *google_type.Date {
	return NewDate(int32(t.Month()), int32(t.Day()), int32(t.Year()))
}

// DateToTime converts a Date to a golang Time.
func DateToTime(d *google_type.Date) time.Time {
	if d == nil {
		return time.Unix(0, 0).UTC()
	}
	return time.Date(int(d.Year), time.Month(d.Month), int(d.Day), 0, 0, 0, 0, time.UTC)
}

// DateLess returns true if i is before j.
func DateLess(i *google_type.Date, j *google_type.Date) bool {
	if j == nil {
		return false
	}
	if i == nil {
		return true
	}
	if i.Year < j.Year {
		return true
	}
	if i.Year > j.Year {
		return false
	}
	if i.Month < j.Month {
		return true
	}
	if i.Month > j.Month {
		return false
	}
	return i.Day < j.Day
}

// DateInRange returns whether d is within start to end, inclusive.
// The given date is expected to not be nil.
// If start is nil, it checks whether d is less than or equal to end.
// If end is nil it checks whether d is greater than or equal to end.
// If start and end are nil, it returns true.
func DateInRange(d *google_type.Date, start *google_type.Date, end *google_type.Date) bool {
	if start == nil && end == nil {
		return true
	}
	if start == nil {
		return DateLess(d, end) || DateEqual(d, end)
	}
	if end == nil {
		return DateLess(start, d) || DateEqual(start, d)
	}
	return DateEqual(d, start) || DateEqual(d, end) || (DateLess(start, d) && DateLess(d, end))
}

// DateEqual returns true if i equals j.
func DateEqual(i *google_type.Date, j *google_type.Date) bool {
	return ((i == nil) == (j == nil)) && ((i == nil) || (*i == *j))
}

// Dater provides the current date.
type Dater interface {
	Now() *google_type.Date
}

// FakeDater is a Dater for testing.
type FakeDater interface {
	Dater
	Set(month int32, day int32, year int32)
}

// NewFakeDater returns a new FakeDater with the initial date.
func NewFakeDater(month int32, day int32, year int32) FakeDater {
	return newFakeDater(month, day, year)
}

type systemDater struct{}

func (s *systemDater) Now() *google_type.Date {
	return Now()
}

type fakeDater struct {
	curDate *google_type.Date
	lock    *sync.RWMutex
}

func newFakeDater(month int32, day int32, year int32) *fakeDater {
	return &fakeDater{NewDate(month, day, year), &sync.RWMutex{}}
}

func (f *fakeDater) Now() *google_type.Date {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return copyDate(f.curDate)
}

func (f *fakeDater) Set(month int32, day int32, year int32) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.curDate = NewDate(month, day, year)
}

func copyDate(date *google_type.Date) *google_type.Date {
	c := *date
	return &c
}
