package webui

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"time"

	"lincast/database"
	"lincast/psync"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/NYTimes/gziphandler"
)

const frontendPath = "frontend/dist"

//go:embed frontend/dist
var _embededFrontend embed.FS

func getFileSystem(devMode bool) http.FileSystem {
	if devMode {
		return http.FS(os.DirFS(frontendPath))
	}

	fsys, err := fs.Sub(_embededFrontend, frontendPath)
	if err != nil {
		log.WithError(err).Panic("Error when trying to get a subfs of the embedded frontend")
	}

	return http.FS(fsys)
}

// New returns a new instance of the server. To execute it, the method `ListenAndServe` must be called.
func New(port uint16, localServer bool, devMode bool, logRequests bool, podcastsDB *database.Database, playerSync *psync.PlayerSync) *http.Server {
	if podcastsDB == nil {
		log.Panic("'podcastsDB' is nil")
	}

	if playerSync == nil {
		log.Panic("'playerSync' is nil")
	}

	_podcastsDB = podcastsDB
	_playerSync = playerSync

	router := newRouter(devMode, logRequests)

	var addr string
	if localServer {
		addr = "127.0.0.1"
	}

	s := http.Server{
		Addr:           addr + ":" + strconv.Itoa(int(port)),
		Handler:        router,
		ReadTimeout:    time.Second * 15,
		WriteTimeout:   time.Second * 15,
		MaxHeaderBytes: 8000, // 8KB
	}

	log.WithFields(log.Fields{
		"address":        s.Addr,
		"readTimeout":    s.ReadTimeout.String(),
		"writeTimeout":   s.WriteTimeout.String(),
		"maxHeaderBytes": s.MaxHeaderBytes,
	}).Infoln("Starting server")

	return &s
}

// newRouter returns a new instance of the router with their paths already set.
func newRouter(devMode, logRequests bool) *mux.Router {
	router := mux.NewRouter()

	if logRequests {
		router.Use(loggingMiddleware)
	}

	router.Handle("/api/v0/podcasts/subscribe", gziphandler.GzipHandler(subscribeToPodcastHandler)).Methods("POST")
	router.Handle("/api/v0/podcasts/unsubscribe", gziphandler.GzipHandler(unsubscribeToPodcastHandler)).Methods("PUT")
	router.Handle("/api/v0/podcasts/user", gziphandler.GzipHandler(getUserPodcastsHandler)).Methods("GET")
	router.Handle("/api/v0/podcasts/{id:[0-9]+}/details", gziphandler.GzipHandler(getPodcastHandler)).Methods("GET")
	router.Handle("/api/v0/podcasts/{id:[0-9]+}/episodes", gziphandler.GzipHandler(getEpisodesHandler)).Methods("GET")
	router.Handle("/api/v0/player/progress", gziphandler.GzipHandler(playerProgressHandler)).Methods("GET", "PUT")
	router.Handle("/api/v0/player/queue", gziphandler.GzipHandler(queueHandler)).Methods("GET", "PUT", "DELETE")
	router.Handle("/api/v0/player/queue/add", gziphandler.GzipHandler(addToQueueHandler)).Methods("POST")
	router.Handle("/api/v0/player/queue/remove", gziphandler.GzipHandler(delFromQueueHandler)).Methods("DELETE")
	router.PathPrefix("/").Handler(gziphandler.GzipHandler(http.FileServer(getFileSystem(devMode))))

	return router
}
