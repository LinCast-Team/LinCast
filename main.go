package main

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"lincast/database"
	"lincast/psync"
	"lincast/queue"
	"lincast/webui"

	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

/* -------------------------------- Constants ------------------------------- */

// Default filenames (should be read of the configurations in the future).
const (
	dbFilename   = "podcasts.sqlite"
	logsFilename = "lincast.log"
)

// Default settings of the server (should be read of the configurations in the future).
const (
	serverPort  = 8080
	serverLocal = true
	serverLogs  = true
)

// Default settings related with feeds' refresh (should be read of the configurations in the future).
const (
	updateFreq = time.Minute * 30
)

/* -------------------------------------------------------------------------- */

func main() {
	devMode := os.Getenv("DEV_MODE") != ""

	setupLogger(logsFilename, devMode)

	log.Info("Starting LinCast")

	if r := run(devMode); r != nil {
		log.Panicln("Error on run:", errorx.Decorate(r, "error on run"))
	}
}

func run(devMode bool) error {
	wd, err := os.Getwd()
	if err != nil {
		return errorx.InternalError.Wrap(err, "error when trying to get the working directory")
	}

	dbPath := filepath.Join(wd, "data/")
	err = os.MkdirAll(dbPath, os.ModePerm)
	if err != nil {
		return errorx.InternalError.Wrap(err, "error when trying to make the directory where the database"+
			" will be stored")
	}

	db, err := database.New(dbPath, dbFilename)
	if err != nil {
		return errorx.InternalError.Wrap(errorx.EnsureStackTrace(err), "error when trying to initialize"+
			" the database in the path '%s'", filepath.Join(dbPath, dbFilename))
	}

	playerSync, err := psync.New(db)
	if err != nil {
		return errorx.InternalError.Wrap(errorx.EnsureStackTrace(err), "error when trying to instantiate the"+
			" synchronizer")
	}

	// Run the loop that updates the subscribed podcasts.
	go runUpdateQueue(db, updateFreq)

	// Make a new instance of the server.
	sv := webui.New(serverPort, serverLocal, devMode, serverLogs, db, playerSync)

	log.WithFields(log.Fields{
		"port":        serverPort,
		"localServer": serverLocal,
		"devMode":     devMode,
		"logRequests": serverLogs,
	}).Info("Starting server")

	err = sv.ListenAndServe()
	if err != nil {
		return errorx.InternalError.Wrap(err, "error on server ListenAndServe")
	}

	return nil
}

func runUpdateQueue(db *database.Database, updateInterval time.Duration) {
	log.WithField("updateInterval", updateInterval.String()).Debug("Starting feeds' update loop")

	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()
	qLength := runtime.NumCPU()

	updateQueue, err := queue.NewUpdateQueue(db, qLength)
	if err != nil {
		log.WithField("error", errorx.Decorate(errorx.EnsureStackTrace(err), "error when creating update queue")).
			Panic("Cannot initialize the update queue")
	}

	log.Info("Updating feeds for first time since LinCast is running")
	err = updatePodcasts(db, updateQueue)
	if err != nil {
		log.WithField("error", errorx.Decorate(err, "Error when trying to update podcasts"))
	}

	for range ticker.C {
		log.Info("Updating podcasts' feeds")
		err := updatePodcasts(db, updateQueue)
		if err != nil {
			log.WithField("error", errorx.EnsureStackTrace(err)).Error("Error when trying to update podcasts' feeds")
		} else {
			log.Info("Podcasts' feeds updated correctly")
		}
	}
}

func updatePodcasts(db *database.Database, updateQueue *queue.UpdateQueue) error {
	subscribedPodcasts, err := db.GetPodcastsBySubscribedStatus(true)
	if err != nil {
		return errorx.InternalError.Wrap(err, "error trying to get subscribed podcasts")
	}

	log.Debug("Starting loop to send podcasts to the update queue")
	for _, p := range *subscribedPodcasts {
		j := queue.NewJob(&p)

		log.WithFields(log.Fields{
			"podcastFeed":       p.FeedLink,
			"podcastID":         p.ID,
			"podcastSubscribed": p.Subscribed,
		}).Info("Sending podcast to the update queue")

		updateQueue.Send(j)
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
