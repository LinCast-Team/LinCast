package webui

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"time"

	"lincast/models"
	"lincast/webui/handlers"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
func New(port uint16, localServer bool, devMode bool, logRequests bool, db *gorm.DB, manualUpdate chan *models.Podcast) *http.Server {
	if db == nil {
		log.Panic("'podcastsDB' is nil")
	}

	handlersManager := handlers.NewManager(db, manualUpdate)

	router := newRouter(devMode, logRequests, handlersManager)

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
func newRouter(devMode, logRequests bool, handlersManager *handlers.Manager) *mux.Router {
	router := mux.NewRouter()

	if logRequests {
		router.Use(loggingMiddleware)
	}

	// gzip.BestCompression = 9
	// gzip.BestSpeed = 1
	gzipMiddleware, _ := gziphandler.NewGzipLevelHandler(5) // Intermediate compression without letting aside the speed.

	router.Use(gzipMiddleware)

	router.HandleFunc("/api/v0/podcasts/subscribe", handlersManager.SubscribeToPodcastHandler).Methods("POST")
	router.HandleFunc("/api/v0/podcasts/unsubscribe", handlersManager.UnsubscribeToPodcastHandler).Methods("PUT")
	router.HandleFunc("/api/v0/user/subscriptions", handlersManager.GetUserPodcastsHandler).Methods("GET")
	router.HandleFunc("/api/v0/podcasts/{id:[0-9]+}/details", handlersManager.GetPodcastHandler).Methods("GET")
	router.HandleFunc("/api/v0/podcasts/{id:[0-9]+}/episodes", handlersManager.GetEpisodesHandler).Methods("GET")
	router.HandleFunc("/api/v0/podcasts/{pID:[0-9]+}/episodes/{epID:[0-9]+}/progress", handlersManager.EpisodeProgressHandler).Methods("GET", "PUT")
	router.HandleFunc("/api/v0/player/playback_info", handlersManager.PlayerPlaybackInfoHandler).Methods("GET", "PUT")
	router.HandleFunc("/api/v0/player/queue", handlersManager.QueueHandler).Methods("GET", "PUT", "DELETE")
	router.HandleFunc("/api/v0/player/queue/add", handlersManager.AddToQueueHandler).Methods("POST")
	router.HandleFunc("/api/v0/player/queue/remove", handlersManager.DelFromQueueHandler).Methods("DELETE")
	router.PathPrefix("/").Handler(http.FileServer(getFileSystem(devMode)))

	return router
}
