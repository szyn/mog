package logger

import (
	log "github.com/Sirupsen/logrus"
)

func Log(message string) {
	log.Debug(message)
}

func Info(message string) {
	log.Info(message)
}

func ErrorIf(err error) bool {
	if err != nil {
		log.Error(err.Error())
		return true
	}

	return false
}

func DieIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func PanicIf(err error) {
	if err != nil {
		log.Panic(err)
	}
}
