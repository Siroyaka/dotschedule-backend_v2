package main

import (
	"log"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/controller"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlwrapper"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/youtubedataapi"
	"github.com/Siroyaka/dotschedule-backend_v2/domain"
	"github.com/Siroyaka/dotschedule-backend_v2/infrastructure"
	infraYoutube "github.com/Siroyaka/dotschedule-backend_v2/infrastructure/youtubedataapi"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/interactor"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

const (
	configPath = "./config.json"

	config_public = "PUBLIC"
	config_sql    = "SQLITE"
	config_root   = "YOUTUBE_SCHEDULE_DATA_UPDATE"
	config_parser = "PARSER"
	config_query  = "QUERY"

	config_sqlPath = "PATH"

	// root
	config_platform           = "PLATFORM"
	config_status100Deferment = "STATUS100_DEFERMENT"
	config_developerKey       = "DEVELOPER_KEY"
	config_streamingUrlPrefic = "STREAMING_URL_PREFIX"
	config_partList           = "PART_LIST"

	// parser
	config_parserHour   = "HOUR"
	config_parserMinite = "MINITE"
	config_parserSecond = "SECOND"

	// query
	config_getSchedule         = "GET_SCHEDULE"
	config_getStreamerMaster   = "GET_STREAMER_MASTER"
	config_updateSchedule      = "UPDATE_SCHEDULE"
	config_updateScheduleTo100 = "UPDATE_SCHEDULE_TO_STATUS100"
	config_countParticipants   = "COUNT_PARTICIPANTS"
	config_insertParticipants  = "INSERT_PARTICIPANTS"
)

func main() {
	confLoader := repository.NewConfigLoader(infrastructure.NewFileReader[utility.ConfigData]())

	config, err := confLoader.ReadJsonConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	publicConfig := config.ReadChild(config_public)
	sqlConfig := config.ReadChild(config_sql)
	rootConfig := config.ReadChild(config_root)
	queryConfig := rootConfig.ReadChild(config_query)
	parserConfig := rootConfig.ReadChild(config_parser)

	// import
	common := utility.NewCommon(publicConfig)

	// infrastructure
	sqlHandler := infrastructure.NewSqliteHandler(sqlConfig.Read(config_sqlPath))
	defer sqlHandler.Close()

	youtubeDataAPI := infraYoutube.NewYoutubeDataAPI(rootConfig.Read(config_developerKey))

	// repository
	youtubeVideoListRepos := youtubedataapi.NewVideoListRepository(youtubeDataAPI)

	getScheduleRepos := sqlwrapper.NewSelectRepository[domain.FullScheduleData](sqlHandler, queryConfig.Read(config_getSchedule))

	getStreamerMasterRepos := sqlwrapper.NewSelectRepository[domain.StreamerMasterWithPlatformData](sqlHandler, queryConfig.Read(config_getStreamerMaster))

	updateScheduleRepos := sqlwrapper.NewUpdateRepository(sqlHandler, queryConfig.Read(config_updateSchedule))

	updateScheduleTo100Repos := sqlwrapper.NewUpdateRepository(sqlHandler, queryConfig.Read(config_updateScheduleTo100))

	countParticipantsRepos := sqlwrapper.NewSelectRepository[int](sqlHandler, queryConfig.Read(config_countParticipants))

	insertParticipantsRepos := sqlwrapper.NewUpdateRepository(sqlHandler, queryConfig.Read(config_insertParticipants))

	// interactor
	intr := interactor.NewNormalizationYoutubeDataInteractor(
		getStreamerMasterRepos,
		getScheduleRepos,
		updateScheduleRepos,
		updateScheduleTo100Repos,
		countParticipantsRepos,
		insertParticipantsRepos,
		youtubeVideoListRepos,
		common,
		utility.NewYoutubeDurationParser(
			parserConfig.Read(config_parserHour),
			parserConfig.Read(config_parserMinite),
			parserConfig.Read(config_parserSecond),
		),
		rootConfig.Read(config_platform),
		rootConfig.Read(config_streamingUrlPrefic),
		rootConfig.ReadStringList(config_partList),
	)

	// controller
	controller := controller.NewNormalizationController(intr)

	controller.Execute()
}
