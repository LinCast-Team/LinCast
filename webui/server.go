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
)

const frontendPath = "frontend/dist"

//go:embed frontend/dist
var embededFrontend embed.FS

func getFileSystem(devMode bool) http.FileSystem {
	if devMode {
		return http.FS(os.DirFS(frontendPath))
	}

	// // Esto es necesario (creo)
	fsys, err := fs.Sub(embededFrontend, frontendPath)
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

// New returns a new instance of the server. To execute it, the method `ListenAndServe` must be called.
func New(port uint16, localServer bool, devMode bool, logRequests bool, podcastsDB *database.Database, playerSynchronizer *psync.Synchronizer) *http.Server {
	if podcastsDB == nil {
		log.Panic("'podcastsDB' is nil")
	}

	if playerSynchronizer == nil {
		log.Panic("'playerSynchronizer' is nil")
	}

	_podcastsDB = podcastsDB
	_pSynchronizer = playerSynchronizer

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

	router.HandleFunc("/api/v0/podcasts/subscribe", subscribeToPodcastHandler).Methods("POST")
	router.HandleFunc("/api/v0/podcasts/unsubscribe", unsubscribeToPodcastHandler).Methods("PUT")
	router.HandleFunc("/api/v0/podcasts/user", getUserPodcastsHandler).Methods("GET")
	router.HandleFunc("/api/v0/podcasts/{id:[0-9]+}/details", getPodcastHandler).Methods("GET")
	router.HandleFunc("/api/v0/podcasts/{id:[0-9]+}/episodes", getEpisodesHandler).Methods("GET")
	router.HandleFunc("/api/v0/player/progress", playerProgressHandler).Methods("GET", "PUT")
	router.PathPrefix("/").Handler(http.FileServer(getFileSystem(devMode)))

	return router
}
