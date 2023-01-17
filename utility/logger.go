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
	defaultLogDir            = "./"
	defaultLogName           = "default"
	defaultLogFileTimeFormat = "20060102"
	defaultLogOSOutput       = true
)

func loggerConfigSetup() (string, bool, IError) {
	logDir := defaultLogDir
	logName := defaultLogName
	logTimeFormat := defaultLogFileTimeFormat
	logOSOutput := defaultLogOSOutput

	// read config
	if projectConfig := config.ReadProjectConfig(); projectConfig != nil && projectConfig.Has(logConfigParent) {
		logConfig := projectConfig.ReadChild(logConfigParent)
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
	}

	if logDir == "" {
		return "", logOSOutput, nil
	}

	if f, err := os.Stat(logDir); os.IsNotExist(err) || !f.IsDir() {
		return "", logOSOutput, NewError(err.Error(), "")
	}
	now := time.Now()
	t := now.UTC().Format(logTimeFormat)

	logFileName := fmt.Sprintf("%s_%s.log", logName, t)
	return filepath.Join(logDir, logFileName), logOSOutput, nil
}

func LoggerSetup() IError {
	logFilePath, isOsOutput, ierr := loggerConfigSetup()
	if ierr != nil {
		return ierr.WrapError()
	}
	if logFilePath == "" {
		return nil
	}

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
