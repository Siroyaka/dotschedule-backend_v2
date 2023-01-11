package fullschedule

import "github.com/Siroyaka/dotschedule-backend_v2/utility"

type GetRepository[X any] interface {
	Get(func(utility.IScan) (X, error), ...any) ([]X, utility.IError)
}
