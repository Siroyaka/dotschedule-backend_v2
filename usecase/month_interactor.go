package usecase

import (
	"strings"

	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/abstruct/viewschedule"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type MonthInteractor struct {
	getRepository             viewschedule.GetRepository
	dataSplitter, arrSplitter string
	displayScheduleStatus     int
	common                    utility.Common
}

func NewMonthInteractor(
	getRepository viewschedule.GetRepository,
	dataSplitter, arrSplitter string,
	displayScheduleStatus int,
	common utility.Common,
) MonthInteractor {
	return MonthInteractor{
		getRepository:         getRepository,
		dataSplitter:          dataSplitter,
		arrSplitter:           arrSplitter,
		common:                common,
		displayScheduleStatus: displayScheduleStatus,
	}
}

func (intr MonthInteractor) dataAdapter(
	date string,
	memberData string,
	icons string,
) domain.MonthData {
	var memberList []domain.MonthMemberData
	if memberData != "" {
		for _, v := range strings.Split(memberData, intr.arrSplitter) {
			vSplit := strings.Split(v, intr.dataSplitter)
			mId := vSplit[0]
			mName := vSplit[1]
			memberList = append(memberList, domain.NewMonthMemberData(mId, mName))
		}
	}

	iconList := strings.Split(icons, intr.arrSplitter)

	return domain.NewMonthData(date, memberList, iconList)
}

func (intr MonthInteractor) GetMonthData(baseDate utility.WrappedTime) ([]domain.MonthData, utility.IError) {
	fromDate := baseDate
	// Month first day on utc is last day of month because it's jst. There is possibility that using add month then skip next month.
	// Therefore add date then add month then sub date.
	toDate := fromDate.Add(0, 0, 1, 0, 0, 0).Add(0, 1, 0, 0, 0, 0).Add(0, 0, -1, 0, 0, 0)
	return intr.getRepository.GetMonthData(fromDate, toDate, intr.displayScheduleStatus, intr.dataAdapter)
}
