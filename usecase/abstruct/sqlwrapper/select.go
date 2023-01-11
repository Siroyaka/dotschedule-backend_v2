package sqlwrapper

import "github.com/Siroyaka/dotschedule-backend_v2/utility"

type SelectRepository[X any] interface {
	Select(Scanable[X]) ([]X, utility.IError)
	SelectPrepare(Scanable[X], []interface{}) ([]X, utility.IError)
}

type IScan interface {
	Scan(...interface{}) error
}

type Scanable[X any] func(IScan) (X, utility.IError)
