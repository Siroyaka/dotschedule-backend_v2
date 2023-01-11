package utility

import "log"

func LogDebug(msg string) {
	log.Printf("DEBUG\t%s\n", msg)
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
