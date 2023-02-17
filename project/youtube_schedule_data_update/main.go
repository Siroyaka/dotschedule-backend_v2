package main

import (
	"log"

	"github.com/Siroyaka/dotschedule-backend_v2/adapter/controller"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/discordpost"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/sqlrepository"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/repository/youtubedataapi"
	"github.com/Siroyaka/dotschedule-backend_v2/infrastructure"
	infraYoutube "github.com/Siroyaka/dotschedule-backend_v2/infrastructure/youtubedataapi"
	"github.com/Siroyaka/dotschedule-backend_v2/usecase/interactor"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/config"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/wrappedbasics"
)

const (
	configPath  = "./config.json"
	projectName = "YOUTUBE_SCHEDULE_DATA_UPDATE"

	config_public  = "PUBLIC"
	config_sql     = "SQLITE"
	config_parser  = "PARSER"
	config_query   = "QUERY"
	config_discord = "DISCORD"

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
	config_getParticipants     = "GET_PARTICIPANTS"
	config_insertParticipants  = "INSERT_PARTICIPANTS"

	// discord
	config_discordUrl                = "NORTIFICATION_URL"
	config_discordNortificationRange = "TARGET_TIME_RANGE"
)

func main() {
	confLoader := repository.NewConfigLoader(infrastructure.NewFileReader[utility.ConfigData]())

	configValue, err := confLoader.ReadJsonConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	config.Setup(projectName, configValue)

	utility.LoggerStart()

	wrappedbasics.InitializeWrappedTimeProps()

	sqlConfig := config.ReadChild(config_sql)
	rootConfig := config.ReadProjectConfig()
	queryConfig := rootConfig.ReadChild(config_query)
	parserConfig := rootConfig.ReadChild(config_parser)
	discordConfig := rootConfig.ReadChild(config_discord)

	// infrastructure
	sqlHandler := infrastructure.NewSqliteHandlerCGOLess(sqlConfig.Read(config_sqlPath))
	defer sqlHandler.Close()

	youtubeDataAPI := infraYoutube.NewYoutubeDataAPI(rootConfig.Read(config_developerKey))

	httpRequest := infrastructure.NewHTTPRequest()

	// repository
	youtubeVideoListRepos := youtubedataapi.NewGetSingleVideoDataRepository(youtubeDataAPI, rootConfig.ReadStringList(config_partList))

	getScheduleRepos := sqlrepository.NewSelectNormalizeSchedulesRepository(sqlHandler, queryConfig.Read(config_getSchedule))

	getStreamerMasterRepos := sqlrepository.NewSelectStreamerMasterWithPlatformMasterRepository(sqlHandler, queryConfig.Read(config_getStreamerMaster))

	updateScheduleRepos := sqlrepository.NewUpdateScheduleRepository(sqlHandler, queryConfig.Read(config_updateSchedule))

	updateScheduleTo100Repos := sqlrepository.NewUpdateScheduleStatusTo100Repository(sqlHandler, queryConfig.Read(config_updateScheduleTo100))

	getParticipantsRepos := sqlrepository.NewSelectAScheduleParticipantsRepository(sqlHandler, queryConfig.Read(config_getParticipants))

	insertParticipantsRepos := sqlrepository.NewInsertSingleParticipantsRepository(sqlHandler, queryConfig.Read(config_insertParticipants))

	discordPostRepos := discordpost.NewDiscordPostRepository(httpRequest, discordConfig.Read(config_discordUrl))

	// interactor
	intr := interactor.NewNormalizationYoutubeDataInteractor(
		getStreamerMasterRepos,
		getScheduleRepos,
		updateScheduleRepos,
		updateScheduleTo100Repos,
		getParticipantsRepos,
		insertParticipantsRepos,
		youtubeVideoListRepos,
		discordPostRepos,
		utility.NewYoutubeDurationParser(
			parserConfig.Read(config_parserHour),
			parserConfig.Read(config_parserMinite),
			parserConfig.Read(config_parserSecond),
		),
		rootConfig.Read(config_platform),
		rootConfig.Read(config_streamingUrlPrefic),
		discordConfig.ReadInteger(config_discordNortificationRange),
	)

	// controller
	controller := controller.NewNormalizationController(intr)

	controller.Execute()
}
