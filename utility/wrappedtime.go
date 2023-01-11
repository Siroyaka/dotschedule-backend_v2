package utility

import (
	"time"
)

const (
	utcLocation = "UTC"
)

type WrappedTime struct {
	time           time.Time
	localLocations string
	baseFormat     string
}

func Now(localLocations string, baseFormat string) (WrappedTime, IError) {
	_, err := time.LoadLocation(localLocations)
	if err != nil {
		return WrappedTime{}, NewError(err.Error(), ERR_LOAD_TIMELOCATION, localLocations)
	}
	now := time.Now()
	return WrappedTime{
		time:           now.UTC(),
		localLocations: localLocations,
		baseFormat:     baseFormat,
	}, nil
}

func NewWrappedTimeFromUTC(t string, localLocations string, baseFormat string) (WrappedTime, IError) {
	tz, err := time.LoadLocation(utcLocation)
	if err != nil {
		return WrappedTime{}, NewError(err.Error(), ERR_LOAD_TIMELOCATION, utcLocation)
	}
	localTime, err := time.ParseInLocation(baseFormat, t, tz)
	if err != nil {
		return WrappedTime{}, NewError(err.Error(), ERR_LOAD_TIMELOCATION, t, utcLocation, baseFormat)
	}
	return WrappedTime{
		time:           localTime.UTC(),
		localLocations: localLocations,
		baseFormat:     baseFormat,
	}, nil
}

func NewWrappedTimeFromLocal(t string, localLocations string, baseFormat string) (WrappedTime, IError) {
	tz, err := time.LoadLocation(localLocations)
	if err != nil {
		return WrappedTime{}, NewError(err.Error(), ERR_LOAD_TIMELOCATION, localLocations)
	}
	localTime, err := time.ParseInLocation(baseFormat, t, tz)
	if err != nil {
		return WrappedTime{}, NewError(err.Error(), ERR_LOAD_TIMELOCATION, t, localLocations, baseFormat)
	}
	return WrappedTime{
		time:           localTime.UTC(),
		localLocations: localLocations,
		baseFormat:     baseFormat,
	}, nil
}

func (t *WrappedTime) SetTimeFormat(timeFormat string) {
	t.baseFormat = timeFormat
}

func (t *WrappedTime) SetLocation(localLocation string) IError {
	_, err := time.LoadLocation(localLocation)
	if err != nil {
		return NewError(err.Error(), ERR_LOAD_TIMELOCATION, localLocation)
	}
	t.localLocations = localLocation
	return nil
}

func (t WrappedTime) Before(c WrappedTime) bool {
	return t.time.Before(c.time)
}

func (t WrappedTime) After(c WrappedTime) bool {
	return t.time.After(c.time)
}

func (t WrappedTime) Equal(c WrappedTime) bool {
	return t.time.Equal(c.time)
}

func (t *WrappedTime) ToLocalFormatString() string {
	tz, err := time.LoadLocation(t.localLocations)
	if err != nil {
		return ""
	}
	return t.time.In(tz).Format(t.baseFormat)
}

func (t *WrappedTime) ToUTCFormatString() string {
	return t.time.UTC().Format(t.baseFormat)
}

func (t WrappedTime) Add(y int, m int, d int, hour int, min int, sec int) WrappedTime {
	n := t.time.AddDate(y, m, d).Add(time.Hour*time.Duration(hour) + time.Minute*time.Duration(min) + time.Second*time.Duration(sec))
	return WrappedTime{
		time:           n,
		localLocations: t.localLocations,
		baseFormat:     t.baseFormat,
	}
}
