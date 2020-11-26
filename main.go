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
	log.Infoln("Starting LinCast")

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

	db, err := podcasts.NewDB(dbPath, dbFilename)
	if err != nil {
		return errorx.InternalError.Wrap(err, "error when trying to initialize the database in the path"+
			" '%s'", filepath.Join(dbPath, dbFilename))
	}

	// Run the loop that updates the subscribed podcasts.
	go runUpdateQueue(db, time.Minute*30)

	// Make a new instance of the server.
	sv := backend.New(8080, true, devMode, true)
	err = sv.ListenAndServe()
	if err != nil {
		return errorx.InternalError.Wrap(err, "error on server ListenAndServe")
	}

	return nil
}

func runUpdateQueue(db *podcasts.Database, updateInterval time.Duration) {
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	updateQueue, err := podcasts.NewUpdateQueue(db, runtime.NumCPU())
	if err != nil {
		log.WithField("error", errorx.Decorate(err, "error when creating update queue")).
			Panic("Cannot initialize the update queue")
	}

	log.Info("Updating podcasts on boot")
	err = updatePodcasts(db, updateQueue)
	if err != nil {
		log.WithField("error", errorx.Decorate(err, "Error when trying to update podcasts"))
	}

	for range ticker.C {
		err := updatePodcasts(db, updateQueue)
		if err != nil {
			log.WithField("error", errorx.Decorate(err, "Error when trying to update podcasts"))
		}
	}
}

func updatePodcasts(db *podcasts.Database, updateQueue *podcasts.UpdateQueue) error {
	log.WithFields(log.Fields{
		"dbIsNil":          db == nil,
		"updateQueueIsNil": updateQueue == nil,
	}).Info("Starting the update of podcasts...")

	log.Info("Getting subscribed podcasts from the database")
	subscribedPodcasts, err := db.GetPodcastsBySubscribedStatus(true)
	if err != nil {
		return errorx.InternalError.Wrap(err, "error trying to get subscribed podcasts")
	}

	log.WithField("subscribedPodcastsN", len(*subscribedPodcasts)).Info("Subscribed podcasts obtained")

	log.Info("Starting loop to send subscribed podcasts to UpdateQueue")
	for _, p := range *subscribedPodcasts {
		j := podcasts.NewJob(&p)

		log.WithFields(log.Fields{
			"jobIsNil":          j == nil,
			"podcastFeed":       p.FeedLink,
			"podcastID":         p.ID,
			"podcastSubscribed": p.Subscribed,
		}).Info("Sending podcast to UpdateQueue as a new Job")

		updateQueue.Send(j)

		log.WithFields(log.Fields{
			"jobIsNil":          j == nil,
			"podcastFeed":       p.FeedLink,
			"podcastID":         p.ID,
			"podcastSubscribed": p.Subscribed,
		}).Info("Podcast sent to UpdateQueue, worker in action")
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
