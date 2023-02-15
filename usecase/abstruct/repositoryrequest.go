package abstruct

import "github.com/Siroyaka/dotschedule-backend_v2/utility"

type RepositoryRequest[P any, Res any] interface {
	Execute(requestParam P) (response Res, err utility.IError)
}
