package main

import (
	"log"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/controller"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/firestorenews"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/fullschedule"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/streamermaster"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/streamingparticipants"
	"github.com/Siroyaka/dotschedule-backend_v2/infrastructure"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/config"
)

const (
	configPath  = "./config.json"
	projectName = "FIRESTORE_NEWS"

	config_public            = "PUBLIC"
	config_sql               = "SQLITE"
	config_query             = "QUERY"
	config_firestoreDocument = "DOCUMENT_NAME"

	config_sqlPath                 = "PATH"
	config_sqlReplaceTargetString  = "TARGETS_STRING"
	config_sqlReplacedChar         = "REPLACED_CHAR"
	config_sqlReplacedCharSplitter = "REPLACED_CHAR_SPLITTER"

	config_videoIDDocument      = "VIDEOID"
	config_videoStatusDocument  = "VIDEOSTATUS"
	config_participantsDocument = "PARTICIPANTS"
	config_updateAtDocument     = "UPDATEAT"

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

	config_platformIdGetQuery = "GET_STREAMERID_OF_PLATFORM"

	config_participantsGetQuery    = "GET_PARTICIPANTS"
	config_participantsInsertQuery = "INSERT_PARTICIPANTS"
	config_participantsDeleteQuery = "DELETE_PARTICIPANTS"
)

func main() {

	confLoader := repository.NewConfigLoader(infrastructure.NewFileReader[utility.ConfigData]())

	configValue, err := confLoader.ReadJsonConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	config.Setup(projectName, configValue)

	if err := utility.LoggerStart(); err != nil {
		log.Fatal(err)
	}

	publicConfig := config.ReadChild(config_public)
	sqlConfig := config.ReadChild(config_sql)

	rootConfig := config.ReadProjectConfig()
	queryConfig := rootConfig.ReadChild(config_query)
	firestoreDocumentConfig := rootConfig.ReadChild(config_firestoreDocument)

	common := utility.NewCommon(publicConfig)

	// infrastructure
	fs := infrastructure.NewFirestore(
		rootConfig.Read(config_credentialFilePath),
		rootConfig.Read(config_project_id),
	)
	defer fs.Close()

	sqlHandler := infrastructure.NewSqliteHandlerCGOLess(sqlConfig.Read(config_sqlPath))
	defer sqlHandler.Close()

	// repository
	fRepos := firestorenews.NewGetRepository(
		fs,
		rootConfig.Read(config_collectionName),
		rootConfig.Read(config_compareble),
		rootConfig.Read(config_operator),
		firestoreDocumentConfig.Read(config_videoIDDocument),
		firestoreDocumentConfig.Read(config_videoStatusDocument),
		firestoreDocumentConfig.Read(config_participantsDocument),
		firestoreDocumentConfig.Read(config_updateAtDocument),
	)

	cRepos := fullschedule.NewCountRepository(
		sqlHandler,
		queryConfig.Read(config_scheduleCountQuery),
	)

	ifRepos := fullschedule.NewInsertRepository(
		sqlHandler,
		queryConfig.Read(config_scheduleInsertQuery),
	)

	ufRepos := fullschedule.NewUpdateAnyColumnRepository(
		sqlHandler,
		queryConfig.Read(config_scheduleUpdateQuery),
	)

	gmRepos := streamermaster.NewGetPlatformIdRepository(
		sqlHandler,
		queryConfig.Read(config_platformIdGetQuery),
	)

	gpRepos := streamingparticipants.NewGetRepository(
		sqlHandler,
		queryConfig.Read(config_participantsGetQuery),
	)

	ipRepos := streamingparticipants.NewInsertRepository(
		sqlHandler,
		queryConfig.Read(config_participantsInsertQuery),
		sqlConfig.Read(config_sqlReplaceTargetString),
		sqlConfig.Read(config_sqlReplacedChar),
		sqlConfig.Read(config_sqlReplacedCharSplitter),
	)

	dpRepos := streamingparticipants.NewDeleteStreamingParticipants(
		sqlHandler,
		queryConfig.Read(config_participantsDeleteQuery),
		sqlConfig.Read(config_sqlReplaceTargetString),
		sqlConfig.Read(config_sqlReplacedChar),
		sqlConfig.Read(config_sqlReplacedCharSplitter),
	)

	// usecase
	firestoreNewsIntr := usecase.NewFirestoreNewsInteractor(
		fRepos,
		cRepos,
		ifRepos,
		ufRepos,
		gmRepos,
		gpRepos,
		ipRepos,
		dpRepos,
		common,
		rootConfig.ReadInteger(config_firestoreNewsTargetMin),
		rootConfig.Read(config_platform),
	)

	// controller
	controller := controller.NewFirestoreNewsController(firestoreNewsIntr)

	controller.Exec()

}
