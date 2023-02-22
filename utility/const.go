package utility

import "time"

const (
	DEFAULT_TIME_FORMAT = time.RFC3339
	EmptyString         = ""
)

type IScan interface {
	Scan(...interface{}) error
}
