package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
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

var shutdownSignal = make(chan os.Signal)

func main() {
	handleCmdArgs()

	devMode := os.Getenv("DEV_MODE") != ""

	setupLogger(logsFilename, devMode)

	log.Info("Starting LinCast")

	run(devMode)
}

func run(devMode bool) {
	// Subscribe to signals related with the stop of the program
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	wd, err := os.Getwd()
	if err != nil {
		log.WithError(err).Panicln("Error when trying to get the working directory")
	}

	dbPath := filepath.Join(wd, "data/")
	err = os.MkdirAll(dbPath, os.ModePerm)
	if err != nil {
		log.WithError(err).Panicln("Error when trying to make the directory where the database will be stored")
	}

	db, err := database.New(dbPath, dbFilename)
	if err != nil {
		log.WithError(errorx.EnsureStackTrace(err)).WithField(
			"path", filepath.Join(dbPath, dbFilename),
		).Panicln("Error when trying to initialize the database")
	}

	manualFeedUpd := make(chan *models.Podcast)

	// Run the loop that updates the subscribed podcasts.
	go runUpdateQueue(db, updateFreq, manualFeedUpd)

	go func() {
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
			log.WithError(
				errorx.InternalError.Wrap(err, "error on server ListenAndServe"),
			).Panicln("")
		}
	}()

	<-shutdownSignal
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

func handleCmdArgs() {
	serviceCmd := flag.String("service", "", "Manage the service of LinCast. Commands: 'install', " +
		"'uninstall', 'start', 'stop', 'restart' and 'status'.")
	flag.Parse()

	if *serviceCmd == "" {
		fmt.Printf("The flag to manage the service can't be used without the action to perform. " +
			"\nSee --help for better understanding.\n")
		os.Exit(1)
	}

	manageService(*serviceCmd)

	os.Exit(0)
}
