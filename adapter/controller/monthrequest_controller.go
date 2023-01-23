package controller

import (
	"fmt"
	"net/http"

	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type MonthRequestController struct {
	common          utility.Common
	monthInteractor usecase.MonthInteractor
	contentType     string
}

func NewMonthRequestController(common utility.Common, MonthInteractor usecase.MonthInteractor, contentType string) MonthRequestController {
	return MonthRequestController{
		common:          common,
		monthInteractor: MonthInteractor,
		contentType:     contentType,
	}
}

func (c MonthRequestController) MonthRequestHandler() http.Handler {
	return http.HandlerFunc(c.monthRequest)
}

func (c MonthRequestController) monthRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", c.contentType)
	month := r.URL.Query().Get("month")
	if month == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, "Invalid Request")
		return
	}

	baseDate, err := c.common.CreateNewWrappedTimeFromLocalMonth(month)

	if err != nil {
		utility.LogError(err.WrapError())
		w.WriteHeader(400)
		fmt.Fprintf(w, "Invalid Month Format: %s", month)
		return
	}

	utility.LogInfo(fmt.Sprintf("Month Request: %s", baseDate.ToUTCFormatString()))

	list, err := c.monthInteractor.GetMonthData(baseDate)

	if err != nil {
		utility.LogError(err.WrapError())
		w.WriteHeader(500)
		fmt.Fprintf(w, "Data Fetch Error: %s", month)
		return
	}

	d := domain.NewAPIResponseData("ok", len(list), "", list)

	json, err := d.ToJson()

	if err != nil {
		utility.LogError(err.WrapError())
		//c.loggerInteractor.Fatal(fmt.Sprintf("Json convert error: %s", month), err, 0)
		w.WriteHeader(500)
		fmt.Fprintf(w, "Data Fetch Error: %s", month)
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, json)
}
