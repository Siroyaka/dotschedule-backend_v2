package fullschedule

import "github.com/Siroyaka/dotschedule-backend_v2/utility"

type CountRepository interface {
	Count(...any) (int, utility.IError)
}
