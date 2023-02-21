package controller

import (
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/interactor"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
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
		utility.LogFatal(err.WrapError())
		return
	}

	controller.firestoreNewsIntr.UpdateDB(firestoreData)
}
