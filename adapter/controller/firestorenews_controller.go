package controller

import (
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/interactor"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/logger"
)

type FirestoreNewsController struct {
	firestoreNewsIntr interactor.FirestoreNewsInteractor
}

func NewFirestoreNewsController(firestoreNewsIntr interactor.FirestoreNewsInteractor) FirestoreNewsController {
	return FirestoreNewsController{
		firestoreNewsIntr: firestoreNewsIntr,
	}
}

func (controller FirestoreNewsController) Exec() {
	firestoreData, err := controller.firestoreNewsIntr.DataFetchFromFirestore()
	if err != nil {
		logger.Fatal(err.WrapError())
		return
	}

	if len(firestoreData) == 0 {
		logger.Debug("firestoreNews no data")
		return
	}

	controller.firestoreNewsIntr.UpdateDB(firestoreData)
}
