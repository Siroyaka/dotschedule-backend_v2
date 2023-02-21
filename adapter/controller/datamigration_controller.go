package controller

import (
	"fmt"

	"github.com/Siroyaka/dotschedule-backend_v2/usecase/interactor"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type DataMigrationController struct {
	dataMigrationInteractor interactor.DataMigrationInteractor
}

func NewDataMigrationController(dataMigrationInteractor interactor.DataMigrationInteractor) DataMigrationController {
	return DataMigrationController{
		dataMigrationInteractor: dataMigrationInteractor,
	}
}

func (dc DataMigrationController) Migration() {
	utility.LogInfo(fmt.Sprintf("start migration: data file count %d", dc.dataMigrationInteractor.Len()))
	var allTotal, allInsert, allRegisterd, allErr int
	fileTotal := dc.dataMigrationInteractor.Len()
	finishedFileCount := 0
	for dc.dataMigrationInteractor.Next() {
		total, insert, registered, errCount, err := dc.dataMigrationInteractor.Migration()
		if err != nil {
			utility.LogFatal(err.WrapError())
			continue
		}
		finishedFileCount++
		allTotal += total
		allInsert += insert
		allRegisterd += registered
		allErr += errCount
	}
	utility.LogInfo(fmt.Sprintf("end migration fileTotal: %d, fileFinished: %d, total: %d, insert: %d, registered: %d, error: %d", fileTotal, finishedFileCount, allTotal, allInsert, allRegisterd, allErr))
}
