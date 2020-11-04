package main

import (
	"os"
	"path/filepath"

	"lincast/webui/backend"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	logsFilename = "lincast.log"
)

func main() {
	devMode := os.Getenv("DEV_MODE") != ""

	setupLogger(logsFilename, devMode)

	if r := run(devMode); r != nil {
		log.Panicln("Error on run:", r.Error())
	}
}

func run(devMode bool) error {
	log.Infoln("Starting LinCast")

	// Arguments should be obtained from settings.
	err := backend.Run(8080, true, devMode, true)
	if err != nil {
		return err
	}

	return nil
}

func setupLogger(filename string, devMode bool) {
	log.SetReportCaller(true)

	if devMode {
		log.SetLevel(log.DebugLevel)

		return
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Panicln(err.Error())
	}

	log.SetOutput(&lumberjack.Logger{
		Filename:   filepath.Join(dir, filename),
		MaxBackups: 3,
		MaxSize:    50,
	})
}
