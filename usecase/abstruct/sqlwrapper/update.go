package sqlwrapper

import "github.com/Siroyaka/dotschedule-backend_v2/utility"

type UpdateRepository interface {
	Update() (count int64, id int64, err utility.IError)
	UpdatePrepare([]interface{}) (count int64, id int64, err utility.IError)
}
