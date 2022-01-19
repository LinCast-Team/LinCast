package main

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"lincast/database"
	"lincast/models"
	"lincast/update"
	"lincast/webui"

	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm"
)

/* -------------------------------- Constants ------------------------------- */

// All this constants should be read from the configs file (see #94)
const (
	// Default filenames
	dbFilename   = "podcasts.sqlite"
	logsFilename = "lincast.log"

	// Default settings of the server
	serverPort  = 8080
	serverLocal = true
	serverLogs  = true

	// Default settings related with feeds' refresh
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

	manualFeedUpd := make(chan *models.Podcast)

	// Run the loop that updates the subscribed podcasts.
	go runUpdateQueue(db, updateFreq, manualFeedUpd)

	// Make a new instance of the server.
	sv := webui.New(serverPort, serverLocal, devMode, serverLogs, db, manualFeedUpd)

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

func runUpdateQueue(db *gorm.DB, updateInterval time.Duration, manualFeedUpd chan *models.Podcast) {
	log.WithField("updateInterval", updateInterval.String()).Debug("Starting feeds' update loop")

	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()
	qLength := runtime.NumCPU()

	updateQueue, err := update.NewUpdateQueue(db, qLength)
	if err != nil {
		log.WithField("error", errorx.Decorate(errorx.EnsureStackTrace(err), "error when creating update queue")).
			Panic("Cannot initialize the update queue")
	}

	log.Info("Updating feeds for first time since LinCast is running")
	err = updateAllPodcasts(db, updateQueue)
	if err != nil {
		log.WithField("error", errorx.Decorate(err, "Error when trying to update podcasts"))
	}

	for {
		select {
		case <-ticker.C:
			{
				log.Info("Updating podcasts' feeds")
				err := updateAllPodcasts(db, updateQueue)
				if err != nil {
					log.WithField("error", errorx.EnsureStackTrace(err)).Error("Error when trying to update podcasts' feeds")
				} else {
					log.Info("Podcasts' feeds updated correctly")
				}
			}
		case p := <-manualFeedUpd:
			{
				j := update.NewJob(p)

				log.WithFields(log.Fields{
					"podcastFeed":       p.FeedLink,
					"podcastID":         p.ID,
					"podcastSubscribed": p.Subscribed,
				}).Info("Sending podcast to the update queue (manual update)")

				updateQueue.Send(j)
			}
		}
	}
}

func updateAllPodcasts(db *gorm.DB, updateQueue *update.UpdateQueue) error {
	var subscribedPodcasts []models.Podcast
	if res := db.Where("subscribed", true).Find(&subscribedPodcasts); res.Error != nil {
		return errorx.InternalError.Wrap(res.Error, "error trying to get subscribed podcasts")
	}

	log.Debug("Starting loop to send podcasts to the update queue")
	for _, p := range subscribedPodcasts {
		j := update.NewJob(&p)

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
