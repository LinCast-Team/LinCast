package main

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"lincast/podcasts"
	"lincast/webui/backend"

	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	dbFilename   = "podcasts.sqlite"
	logsFilename = "lincast.log"
)

func main() {
	devMode := os.Getenv("DEV_MODE") != ""

	setupLogger(logsFilename, devMode)

	if r := run(devMode); r != nil {
		log.Panicln("Error on run:", errorx.Decorate(r, "error on run"))
	}
}

func run(devMode bool) error {
	log.Info("Starting LinCast")

	log.Debug("Getting working directory")
	wd, err := os.Getwd()
	if err != nil {
		return errorx.InternalError.Wrap(err, "error when trying to get the working directory")
	}
	log.WithField("wd", wd).Debug("Working directory obtained")

	dbPath := filepath.Join(wd, "data/")

	log.WithField("dbPath", dbPath).Debug("Ensuring that the path of the database exists")
	err = os.MkdirAll(dbPath, os.ModePerm)
	if err != nil {
		return errorx.InternalError.Wrap(err, "error when trying to make the directory where the database"+
			" will be stored")
	}
	log.Info("Path of the database checked (or created) correctly")

	log.WithFields(log.Fields{"dbPath": dbPath, "dbFilename": dbFilename}).
		Debug("Creating a new instance of Database")
	db, err := podcasts.NewDB(dbPath, dbFilename)
	if err != nil {
		return errorx.InternalError.Wrap(errorx.EnsureStackTrace(err), "error when trying to initialize"+
			" the database in the path '%s'", filepath.Join(dbPath, dbFilename))
	}
	log.Info("Database instantiated correctly")

	// Run the loop that updates the subscribed podcasts.
	log.Debug("Running podcasts update loop")
	go runUpdateQueue(db, time.Minute*30)

	// Make a new instance of the server.
	log.Debug("Instantiating backend")
	sv := backend.New(8080, true, devMode, true, db)
	log.WithFields(log.Fields{
		"port":        0,
		"localServer": true,
		"devMode":     false,
		"logRequests": true,
	}).Info("Backend instantiated")

	log.Debug("Executing server's ListenAndServe method")
	err = sv.ListenAndServe()
	if err != nil {
		return errorx.InternalError.Wrap(err, "error on server ListenAndServe")
	}

	return nil
}

func runUpdateQueue(db *podcasts.Database, updateInterval time.Duration) {
	log.WithField("updateInterval", updateInterval.String()).Debug("Starting podcasts update loop")

	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()
	qLength := runtime.NumCPU()

	log.WithField("length", qLength).Debug("Instantiating a new UpdateQueue")
	updateQueue, err := podcasts.NewUpdateQueue(db, qLength)
	if err != nil {
		log.WithField("error", errorx.Decorate(errorx.EnsureStackTrace(err), "error when creating update queue")).
			Panic("Cannot initialize the update queue")
	}
	log.Debug("UpdateQueue initialized correctly")

	log.Info("Updating podcasts on boot")
	err = updatePodcasts(db, updateQueue)
	if err != nil {
		log.WithField("error", errorx.Decorate(err, "Error when trying to update podcasts"))
	}
	log.Info("Podcasts updated for first time correctly")

	for range ticker.C {
		log.Debug("Tick received, executing podcasts update")
		err := updatePodcasts(db, updateQueue)
		if err != nil {
			log.WithField("error", errorx.Decorate(err, "Error when trying to update podcasts"))
		} else {
			log.Info("Podcasts update executed correctly, waiting for next signal")
		}
	}
}

func updatePodcasts(db *podcasts.Database, updateQueue *podcasts.UpdateQueue) error {
	log.WithFields(log.Fields{
		"dbIsNil":          db == nil,
		"updateQueueIsNil": updateQueue == nil,
	}).Debug("Starting the update of podcasts...")

	log.Debug("Getting subscribed podcasts from the database")
	subscribedPodcasts, err := db.GetPodcastsBySubscribedStatus(true)
	if err != nil {
		return errorx.InternalError.Wrap(err, "error trying to get subscribed podcasts")
	}
	log.WithField("subscribedPodcastsN", len(*subscribedPodcasts)).Info("Subscribed podcasts obtained")

	log.Debug("Starting loop to send subscribed podcasts to UpdateQueue")
	for _, p := range *subscribedPodcasts {
		j := podcasts.NewJob(&p)

		log.WithFields(log.Fields{
			"jobIsNil":          j == nil,
			"podcastFeed":       p.FeedLink,
			"podcastID":         p.ID,
			"podcastSubscribed": p.Subscribed,
		}).Debug("Sending podcast to UpdateQueue as a new Job")

		updateQueue.Send(j)

		log.WithFields(log.Fields{
			"jobIsNil":          j == nil,
			"podcastFeed":       p.FeedLink,
			"podcastID":         p.ID,
			"podcastSubscribed": p.Subscribed,
		}).Debug("Podcast sent to UpdateQueue, worker in action")
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
