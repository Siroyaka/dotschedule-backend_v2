package abstruct

import "github.com/Siroyaka/dotschedule-backend_v2/utility/utilerror"

type RepositoryRequest[P any, Res any] interface {
	Execute(requestParam P) (response Res, err utilerror.IError)
}
