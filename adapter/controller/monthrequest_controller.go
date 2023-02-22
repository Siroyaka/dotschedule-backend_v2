package controller

import (
	"fmt"
	"net/http"

	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/interactor"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/logger"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

type MonthRequestController struct {
	monthInteractor     interactor.DaysParticipantsInteractor
	contentType         string
	localTimeDifference int
}

func NewMonthRequestController(
	MonthInteractor interactor.DaysParticipantsInteractor,
	contentType string,
	localTimeDifference int,
) MonthRequestController {
	return MonthRequestController{
		monthInteractor:     MonthInteractor,
		contentType:         contentType,
		localTimeDifference: localTimeDifference,
	}
}

func (c MonthRequestController) MonthRequestHandler() http.Handler {
	return http.HandlerFunc(c.monthRequest)
}

// Goの月の加算は、単純に月を加算し、もしその月に該当の日がなければその分次の月に繰り越して計算するようになっている
//
// なので、3/31に1ヶ月加算すると4/31 = 5/1となってしまう
//
// 加えて、UTCで計算されるようになっている。
//
// なので、JTCの4/1 0:00は3/31 15:00からの計算となる
//
// そのため、単純に4/1の1ヶ月後を計算するとJTCに直したときには5/2になってしまう
//
//	4/1 = 3/31 15:00 の1ヶ月後 = 5/1 15:00 = 5/2
//
// なので、時差を計算することで無理やりこれを解消する
func (c MonthRequestController) monthAdd(fromDate wrappedbasics.IWrappedTime, addRange int) wrappedbasics.IWrappedTime {
	return fromDate.Add(0, 0, 0, c.localTimeDifference, 0, 0).Add(0, addRange, 0, 0, 0, 0).Add(0, 0, 0, -1*c.localTimeDifference, 0, 0)
}

func (c MonthRequestController) monthRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", c.contentType)
	month := r.URL.Query().Get("month")
	if month == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, "Invalid Request")
		return
	}

	fromDate, err := wrappedbasics.NewWrappedTimeFromLocal(month, wrappedbasics.WrappedTimeProps.MonthFormat())

	if err != nil {
		logger.Error(err.WrapError())
		w.WriteHeader(400)
		fmt.Fprintf(w, "Invalid Month Format. %s", month)
		return
	}

	logger.Info(fmt.Sprintf("Month Request. month: %s", month))

	toDate := c.monthAdd(fromDate, 1)

	list, err := c.monthInteractor.GetMonthData(fromDate, toDate)

	if err != nil {
		logger.Error(err.WrapError())
		w.WriteHeader(500)
		fmt.Fprintf(w, "Data Fetch Error: %s", month)
		return
	}

	d := domain.NewAPIResponseData("ok", len(list), "", list)

	json, err := d.ToJson()

	if err != nil {
		logger.Error(err.WrapError())
		//c.loggerInteractor.Fatal(fmt.Sprintf("Json convert error: %s", month), err, 0)
		w.WriteHeader(500)
		fmt.Fprintf(w, "Data Fetch Error: %s", month)
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, json)
}
