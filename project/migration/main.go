package main

import (
	"log"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/controller"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/fileio"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/fullschedule"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/migration"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/streamingparticipants"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/infrastructure"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/config"
)

const (
	configPath  = "./config.json"
	projectName = "MIGRATION"

	config_public = "PUBLIC"
	config_sql    = "SQLITE"
	config_query  = "QUERY"

	config_sqlPath                 = "PATH"
	config_sqlReplaceTargetString  = "TARGETS_STRING"
	config_sqlReplacedChar         = "REPLACED_CHAR"
	config_sqlReplacedCharSplitter = "REPLACED_CHAR_SPLITTER"

	config_dataFileDirectoryPath = "DATA_FILE_DIRECTORY_PATH"
	config_platformName          = "PLATFORM"

	// query
	config_getSchedule        = "GET_SCHEDULE"
	config_getStreamerMaster  = "GET_STREAMER_MASTER"
	config_insertSchedule     = "INSERT_SCHEDULE"
	config_insertParticipants = "INSERT_PARTICIPANTS"
)

func main() {
	confLoader := repository.NewConfigLoader(infrastructure.NewFileReader[utility.ConfigData]())

	configValue, err := confLoader.ReadJsonConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	config.Setup(projectName, configValue)

	if err := utility.LoggerSetup(); err != nil {
		log.Fatal(err)
	}

	publicConfig := config.ReadChild(config_public)
	sqlConfig := config.ReadChild(config_sql)
	rootConfig := config.ReadProjectConfig()
	queryConfig := rootConfig.ReadChild(config_query)

	// import
	common := utility.NewCommon(publicConfig)

	// infrastructure
	sqlHandler := infrastructure.NewSqliteHandlerCGOLess(sqlConfig.Read(config_sqlPath))
	defer sqlHandler.Close()

	fileReader := infrastructure.NewFileReader[domain.JsonScheduleList]()

	// repository
	fileReaderRepos := fileio.NewReaderRepository[domain.JsonScheduleList](fileReader)

	getDataForMigrationRepos := migration.NewGetRepository(
		sqlHandler,
		queryConfig.Read(config_getSchedule),
		sqlConfig.Read(config_sqlReplacedChar),
		sqlConfig.Read(config_sqlReplacedCharSplitter),
		sqlConfig.Read(config_sqlReplaceTargetString),
		queryConfig.Read(config_getStreamerMaster),
	)

	insertParticipantsRepos := streamingparticipants.NewInsertRepository(
		sqlHandler,
		queryConfig.Read(config_insertParticipants),
		sqlConfig.Read(config_sqlReplaceTargetString),
		sqlConfig.Read(config_sqlReplacedChar),
		sqlConfig.Read(config_sqlReplacedCharSplitter),
	)

	insertFullScheduleRepos := fullschedule.NewInsertRepository(
		sqlHandler,
		queryConfig.Read(config_insertSchedule),
	)

	// interactor
	dataMigrationIntr := usecase.NewDataMigrationInteractor(
		common,
		fileReaderRepos,
		rootConfig.Read(config_dataFileDirectoryPath),
		rootConfig.Read(config_platformName),
		getDataForMigrationRepos,
		insertFullScheduleRepos,
		insertParticipantsRepos,
	)

	// controller
	dataMigrationController := controller.NewDataMigrationController(dataMigrationIntr)

	dataMigrationController.Migration()
}
