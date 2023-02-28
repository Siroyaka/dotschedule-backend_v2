package wrappedbasics

import (
	"fmt"
	"time"

	"github.com/Siroyaka/dotschedule-backend_v2/utility/config"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/logger"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"
)

const (
	utcLocation               = "UTC"
	wrappedTimeConfigKey      = "WRAPPED_TIME"
	localLocationConfigKey    = "LOCAL_LOCATION"
	datetimeFormatConfigKey   = "DATETIME_FORMAT"
	dateFormatConfigKey       = "DATE_FORMAT"
	monthFormatConfigKey      = "MONTH_FORMAT"
	otherFormatTypesConfigKey = "OTHER_FORMAT_TYPES"
	otherFormatsConfigKey     = "OTHER_FORMATS"
)

var (
	WrappedTimeProps = WrappedTimePropsStruct{
		localLocation:  "Asia/Tokyo",
		datetimeFormat: "2006-01-02 15:04:05",
		dateFormat:     "2006-01-02",
		monthFormat:    "2006-01",
		otherFormats:   make(map[string]WrappedTimeFormat),
	}
)

type WrappedTimePropsStruct struct {
	localLocation  string
	datetimeFormat WrappedTimeFormat
	dateFormat     WrappedTimeFormat
	monthFormat    WrappedTimeFormat
	otherFormats   map[string]WrappedTimeFormat
}

func InitializeWrappedTimeProps() {
	if !config.Has(wrappedTimeConfigKey) {
		logger.Info("Wrapped time config not found.")
		return
	}
	timeConfig := config.ReadChild("WRAPPED_TIME")

	if timeConfig.Has(localLocationConfigKey) {
		WrappedTimeProps.localLocation = timeConfig.Read(localLocationConfigKey)
	} else {
		logger.Info("local location use default.")
	}

	if timeConfig.Has(datetimeFormatConfigKey) {
		WrappedTimeProps.datetimeFormat = WrappedTimeFormat(timeConfig.Read(datetimeFormatConfigKey))
	} else {
		logger.Info("datetime format use default.")
	}

	if timeConfig.Has(dateFormatConfigKey) {
		WrappedTimeProps.dateFormat = WrappedTimeFormat(timeConfig.Read(dateFormatConfigKey))
	} else {
		logger.Info("date format use default.")
	}

	if timeConfig.Has(monthFormatConfigKey) {
		WrappedTimeProps.monthFormat = WrappedTimeFormat(timeConfig.Read(monthFormatConfigKey))
	} else {
		logger.Info("month format use default.")
	}

	if timeConfig.Has(otherFormatTypesConfigKey) && timeConfig.Has(otherFormatsConfigKey) {
		otherFormatsConfig := timeConfig.ReadChild(otherFormatsConfigKey)
		for _, key := range timeConfig.ReadStringList(otherFormatTypesConfigKey) {
			if otherFormatsConfig.Has(key) {
				WrappedTimeProps.otherFormats[key] = WrappedTimeFormat(otherFormatsConfig.Read(key))
			} else {
				logger.Info(fmt.Sprintf("original format %s is not found.", key))
			}
		}
	} else {
		logger.Info("other format nothing.")
	}
}

func (wtp WrappedTimePropsStruct) LocalLocation() string {
	return wtp.localLocation
}

func (wtp WrappedTimePropsStruct) DateTimeFormat() WrappedTimeFormat {
	return wtp.datetimeFormat
}

func (wtp WrappedTimePropsStruct) DateFormat() WrappedTimeFormat {
	return wtp.dateFormat
}

func (wtp WrappedTimePropsStruct) MonthFormat() WrappedTimeFormat {
	return wtp.monthFormat
}

func (wtp WrappedTimePropsStruct) OtherFormats(key string) WrappedTimeFormat {
	if v, ok := wtp.otherFormats[key]; ok {
		return v
	}
	logger.Info(fmt.Sprintf("Format %s is NotFound.", key))
	return time.RFC3339
}

type WrappedTimeFormat string

func (wtf WrappedTimeFormat) toString() string {
	return string(wtf)
}

type WrappedTime struct {
	time time.Time
}

type IWrappedTime interface {
	Before(c IWrappedTime) bool
	After(c IWrappedTime) bool
	Equal(c IWrappedTime) bool
	Time() time.Time
	ToLocalFormatString(format WrappedTimeFormat) string
	ToUTCFormatString(format WrappedTimeFormat) string
	Add(y int, m int, d int, hour int, min int, sec int) IWrappedTime
}

func Now() WrappedTime {
	now := time.Now()
	return WrappedTime{
		time: now.UTC(),
	}
}

func NewWrappedTimeFromUTC(t string, baseFormat WrappedTimeFormat) (IWrappedTime, utilerror.IError) {
	tz, err := time.LoadLocation(utcLocation)
	if err != nil {
		return nil, utilerror.New(err.Error(), utilerror.ERR_LOAD_TIMELOCATION, utcLocation)
	}
	localTime, err := time.ParseInLocation(baseFormat.toString(), t, tz)
	if err != nil {
		return nil, utilerror.New(err.Error(), utilerror.ERR_LOAD_TIMELOCATION, t, utcLocation, baseFormat.toString())
	}
	return WrappedTime{
		time: localTime.UTC(),
	}, nil
}

func NewWrappedTimeFromLocal(t string, baseFormat WrappedTimeFormat) (IWrappedTime, utilerror.IError) {
	tz, err := time.LoadLocation(WrappedTimeProps.localLocation)
	if err != nil {
		return nil, utilerror.New(err.Error(), utilerror.ERR_LOAD_TIMELOCATION, WrappedTimeProps.localLocation)
	}
	localTime, err := time.ParseInLocation(baseFormat.toString(), t, tz)
	if err != nil {
		return nil, utilerror.New(err.Error(), utilerror.ERR_LOAD_TIMELOCATION, t, WrappedTimeProps.localLocation, baseFormat.toString())
	}
	return WrappedTime{
		time: localTime.UTC(),
	}, nil
}

// a.Before(b)の場合はaのほうがbより前であるかを確認する
//
// つまり、a < b
func (t WrappedTime) Before(c IWrappedTime) bool {
	return t.Time().Before(c.Time())
}

// a.After(b)の場合はaのほうがbより後であるかを確認する
//
// つまり、a > b
func (t WrappedTime) After(c IWrappedTime) bool {
	return t.Time().After(c.Time())
}

func (t WrappedTime) Equal(c IWrappedTime) bool {
	return t.Time().Equal(c.Time())
}

func (t WrappedTime) Time() time.Time {
	return t.time
}

func (t WrappedTime) ToLocalFormatString(format WrappedTimeFormat) string {
	tz, err := time.LoadLocation(WrappedTimeProps.localLocation)
	if err != nil {
		return ""
	}
	return t.time.In(tz).Format(format.toString())
}

func (t WrappedTime) ToUTCFormatString(format WrappedTimeFormat) string {
	return t.time.UTC().Format(format.toString())
}

func (t WrappedTime) Add(y int, m int, d int, hour int, min int, sec int) IWrappedTime {
	n := t.time.AddDate(y, m, d).Add(time.Hour*time.Duration(hour) + time.Minute*time.Duration(min) + time.Second*time.Duration(sec))
	return WrappedTime{
		time: n,
	}
}
