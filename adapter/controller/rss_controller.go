package controller

import (
	"fmt"

	"github.com/Siroyaka/dotschedule-backend_v2/usecase"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type RSSController struct {
	intr usecase.RSSInteractor
}

func NewRSSController(intr usecase.RSSInteractor) RSSController {
	return RSSController{
		intr: intr,
	}
}

func startExec() {
	utility.LogInfo("start rss batch")
}

func endExec() {
	utility.LogInfo("end rss batch")
}

func (rc RSSController) Exec() {
	startExec()
	defer endExec()
	err := rc.intr.GetMaster()
	if err != nil {
		utility.LogFatal(err.WrapError())
		return
	}
	totalInsert := 0
	totalUpdate := 0
	for rc.intr.Next() {
		list, err := rc.intr.GetRSSData()
		if err != nil {
			utility.LogFatal(err.WrapError())
			continue
		}
		insert, update, err := rc.intr.PushToDB(list)
		if err != nil {
			utility.LogFatal(err.WrapError())
			continue
		}
		if err = rc.intr.EndRow(); err != nil {
			utility.LogFatal(err.WrapError())
			continue
		}
		totalInsert += insert
		totalUpdate += update
	}
	utility.LogInfo(fmt.Sprintf("RSSFeed insert: %d, update: %d", totalInsert, totalUpdate))
}
