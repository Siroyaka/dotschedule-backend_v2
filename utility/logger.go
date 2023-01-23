package utility

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Siroyaka/dotschedule-backend_v2/mode"
	"github.com/Siroyaka/dotschedule-backend_v2/utility/config"
)

const (
	logConfigParent          = "LOG"
	logConfigDir             = "DIR"
	logConfigName            = "NAME"
	logConfigTimeFormat      = "TIME_FORMAT"
	logConfigOSOutput        = "OS_OUTPUT"
	logConfigSwitcher        = "SWITCHER"
	defaultLogDir            = "./"
	defaultLogName           = "default"
	defaultLogFileTimeFormat = "20060102"
	defaultLogOSOutput       = true
)

var (
	logger_Time     = ""
	logger_switcher = "none"
	logDir          = ""
	logName         = ""
	logOSOutput     = true
	logTimeFormat   = ""
)

func nowTime(format string) string {
	now := time.Now()
	return now.UTC().Format(format)
}

func loadLoggerConfig() config.IConfig {
	if projectConfig := config.ReadProjectConfig(); projectConfig != nil && projectConfig.Has(logConfigParent) {
		return projectConfig.ReadChild(logConfigParent)
	}
	return nil
}

func loggerConfigSetup() (string, bool, IError) {
	logDir = defaultLogDir
	logName = defaultLogName
	logTimeFormat = defaultLogFileTimeFormat
	logOSOutput = defaultLogOSOutput

	// read config
	logConfig := loadLoggerConfig()
	if logConfig != nil {
		if logConfig.Has(logConfigDir) {
			logDir = logConfig.Read(logConfigDir)
		}
		if logConfig.Has(logConfigName) {
			logName = logConfig.Read(logConfigName)
		}
		if logConfig.Has(logConfigTimeFormat) {
			logTimeFormat = logConfig.Read(logConfigTimeFormat)
		}
		if logConfig.Has(logConfigOSOutput) {
			logOSOutput = logConfig.ReadBoolean(logConfigOSOutput)
		}
		if logConfig.Has(logConfigSwitcher) {
			switch sw := logConfig.Read(logConfigSwitcher); sw {
			case "TIME":
				logger_switcher = sw
			default:
				logger_switcher = "none"
			}
		}
	}

	if logDir == "" {
		return "", logOSOutput, nil
	}

	if f, err := os.Stat(logDir); os.IsNotExist(err) || !f.IsDir() {
		return "", logOSOutput, NewError(err.Error(), "")
	}
	t := nowTime(logTimeFormat)

	logFileName := fmt.Sprintf("%s_%s.log", logName, t)
	return filepath.Join(logDir, logFileName), logOSOutput, nil
}

func loggerSwitch() {
	switch logger_switcher {
	case "none":
		return
	default:
		panic(fmt.Sprintf("logger switcher case error. %s", logger_switcher))
	}
}

func loggerSetup(logFilePath string, isOsOutput bool) IError {
	logfile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return NewError(err.Error(), "")
	}
	if isOsOutput {
		log.SetOutput(io.MultiWriter(os.Stdout, logfile))
	} else {
		log.SetOutput(logfile)
	}
	return nil
}

func LoggerStart() IError {
	logFilePath, isOsOutput, ierr := loggerConfigSetup()
	if ierr != nil {
		return ierr.WrapError()
	}
	if logFilePath == "" {
		return nil
	}

	err := loggerSetup(logFilePath, isOsOutput)
	if err != nil {
		return err.WrapError()
	}
	return nil
}

func LogDebug(msg string) {
	if mode.DEBUG {
		log.Printf("DEBUG\t%s\n", msg)
	}
}

func LogInfo(msg string) {
	log.Printf("INFO\t%s\n", msg)
}

func LogError(err error) {
	log.Printf("ERROR\t%s\n", err)
}

func LogFatal(err error) {
	log.Printf("FATAL\t%s\n", err)
}
