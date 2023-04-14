package main

import (
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/controller"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/fileio"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/otheroutbound"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository/sqlcontains"
	"github.com/Siroyaka/dotschedule-backend_v2/infrastructure"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/interactor"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/config"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/logger"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

const (
	configPath  = "./config.json"
	projectName = "REGISTRATION_REQUEST"

	config_public            = "PUBLIC"
	config_sql               = "SQLITE"
	config_query             = "QUERY"
	config_firestoreDocument = "DOCUMENT_NAME"

	config_sqlPath                 = "PATH"
	config_sqlReplaceTargetString  = "TARGETS_STRING"
	config_sqlReplacedChar         = "REPLACED_CHAR"
	config_sqlReplacedCharSplitter = "REPLACED_CHAR_SPLITTER"

	config_IDDocument           = "ID"
	config_platformDocument     = "PLATFORM"
	config_URLDocument          = "URL"
	config_updateAtDocument     = "UPDATE_AT"
	config_participantsDocument = "PARTICIPANTS"
	config_titleDocument        = "TITLE"
	config_streamerNameDocument = "STREAMER_NAME"
	config_streamerIDDocument   = "STREAMER_ID"
	config_startDateDocument    = "START_DATE"

	config_credentialFilePath     = "CREDENTIALS_FILE"
	config_project_id             = "GCP_PROJECT_ID"
	config_collectionName         = "COLLECTION_NAME"
	config_compareble             = "QUERY_COMPAREBLE"
	config_operator               = "QUERY_OPERATOR"
	config_firestoreNewsTargetMin = "NEWS_TARGET_MIN"
	config_platform               = "PLATFORM"

	config_scheduleCountQuery  = "COUNT_SCHEDULE"
	config_scheduleInsertQuery = "INSERT_SCHEDULE"
	config_scheduleUpdateQuery = "UPDATE_SCHEDULE"

	config_participantsGetQuery    = "GET_PARTICIPANTS"
	config_participantsInsertQuery = "INSERT_PARTICIPANTS"
	config_participantsDeleteQuery = "DELETE_PARTICIPANTS"
)

func loadConfig() config.IConfig {
	reader := fileio.NewReaderRepository[map[string]interface{}](infrastructure.NewFileReader[map[string]interface{}]())
	data, err := reader.ReadJson(configPath)
	if err != nil {
		panic(err)
	}
	return config.New(data)
}

func main() {
	config.Setup(projectName, loadConfig())

	logger.Start()
	wrappedbasics.InitializeWrappedTimeProps()

	sqlConfig := config.ReadChild(config_sql)

	rootConfig := config.ReadProjectConfig()
	queryConfig := rootConfig.ReadChild(config_query)
	firestoreDocumentConfig := rootConfig.ReadChild(config_firestoreDocument)

	// infrastructure
	fs := infrastructure.NewFirestore(
		rootConfig.Read(config_credentialFilePath),
		rootConfig.Read(config_project_id),
	)
	defer fs.Close()

	sqlHandler := infrastructure.NewSqliteHandlerCGOLess(sqlConfig.Read(config_sqlPath))
	defer sqlHandler.Close()

	// repository
	fRepos := otheroutbound.NewGetFirestoreRegistrationRequestRepository(
		fs,
		rootConfig.Read(config_collectionName),
		rootConfig.Read(config_compareble),
		rootConfig.Read(config_operator),
		firestoreDocumentConfig.Read(config_IDDocument),
		firestoreDocumentConfig.Read(config_URLDocument),
		firestoreDocumentConfig.Read(config_platformDocument),
		firestoreDocumentConfig.Read(config_participantsDocument),
		firestoreDocumentConfig.Read(config_updateAtDocument),
		firestoreDocumentConfig.Read(config_startDateDocument),
		firestoreDocumentConfig.Read(config_streamerIDDocument),
		firestoreDocumentConfig.Read(config_streamerNameDocument),
		firestoreDocumentConfig.Read(config_titleDocument),
	)

	cRepos := sqlcontains.NewContainsScheduleRepository(
		sqlHandler,
		queryConfig.Read(config_scheduleCountQuery),
	)

	ifRepos := sqlrepository.NewInsertFullScheduleRepository(
		sqlHandler,
		queryConfig.Read(config_scheduleInsertQuery),
	)

	ufRepos := sqlrepository.NewUpdateScheduleCompleteStatusRepository(
		sqlHandler,
		queryConfig.Read(config_scheduleUpdateQuery),
	)

	gpRepos := sqlrepository.NewSelectParticipantsRepository(
		sqlHandler,
		queryConfig.Read(config_participantsGetQuery),
	)

	ipRepos := sqlrepository.NewInsertParticipantsRepository2(
		sqlHandler,
		queryConfig.Read(config_participantsInsertQuery),
		sqlConfig.Read(config_sqlReplaceTargetString),
		sqlConfig.Read(config_sqlReplacedChar),
		sqlConfig.Read(config_sqlReplacedCharSplitter),
	)

	dpRepos := sqlrepository.NewDeleteParticipantsRepository(
		sqlHandler,
		queryConfig.Read(config_participantsDeleteQuery),
		sqlConfig.Read(config_sqlReplaceTargetString),
		sqlConfig.Read(config_sqlReplacedChar),
		sqlConfig.Read(config_sqlReplacedCharSplitter),
	)

	// usecase
	registrationRequestIntr := interactor.NewRegistrationRequestInteractor(
		fRepos,
		cRepos,
		ifRepos,
		ufRepos,
		gpRepos,
		ipRepos,
		dpRepos,
		rootConfig.ReadInteger(config_firestoreNewsTargetMin),
	)

	// controller
	controller := controller.NewFirestoreRegistraitonRequestController(registrationRequestIntr)

	controller.Exec()

}
