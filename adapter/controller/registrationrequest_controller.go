package controller

import (
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/interactor"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/logger"
)

type FirestoreRegistrationRequestController struct {
	intr interactor.RegistrationRequestInteractor
}

func NewFirestoreRegistraitonRequestController(intr interactor.RegistrationRequestInteractor) FirestoreRegistrationRequestController {
	return FirestoreRegistrationRequestController{
		intr: intr,
	}
}

func (controller FirestoreRegistrationRequestController) Exec() {
	firestoreData, err := controller.intr.DataFetchFromFirestore()
	if err != nil {
		logger.Fatal(err.WrapError())
		return
	}

	if len(firestoreData) == 0 {
		logger.Debug("registration request no data")
		return
	}

	controller.intr.UpdateDB(firestoreData)
}
