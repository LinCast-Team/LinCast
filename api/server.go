package api

import (
	"net/http"
	"strconv"
	"time"

	"lincast/api/handlers"
	"lincast/models"

	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// New returns a new instance of the server. To execute it, the method `ListenAndServe` must be called.
func New(port uint, localServer bool, devMode bool, logRequests bool, db *gorm.DB, manualUpdate chan *models.Podcast) *http.Server {
	handlersManager := handlers.NewManager(db, manualUpdate)

	router := newRouter(devMode, logRequests, handlersManager)

	var addr string
	if localServer {
		addr = "127.0.0.1"
	}

	// TODO Check if this can be handled by some Chi middleware
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
func newRouter(_, logRequests bool, handlersManager *handlers.Manager) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	if logRequests {
		router.Use(loggingMiddleware)
	}

	// gzip.BestCompression = 9
	// gzip.BestSpeed = 1
	gzipMiddleware, _ := gziphandler.NewGzipLevelHandler(5) // Intermediate compression without letting aside the speed.

	// TODO Check if there is no Chi middleware to handle gzip compression
	router.Use(gzipMiddleware)

	router.Route("/api/v0", func(r chi.Router) {
		// TODO Implement user authentication
		r.Route("/auth", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("TODO!"))
			})
		})

		r.Route("/podcasts", func(r chi.Router) {
			r.Post("/subscribe", handlersManager.SubscribeToPodcastHandler)
			r.Put("/unsubscribe", handlersManager.UnsubscribeToPodcastHandler)
			r.Get("/{id:[0-9]+}", handlersManager.GetPodcastHandler)
			r.Get("/{id:[0-9]+}/episodes", handlersManager.GetEpisodesHandler)
			r.Get("/{id:[0-9]+}/episodes/{epID:[0-9]+}", handlersManager.EpisodeDetailsHandler)
			r.Get("/{id:[0-9]+}/episodes/{epID:[0-9]+}/progress", handlersManager.EpisodeProgressHandler)
			r.Put("/{id:[0-9]+}/episodes/{epID:[0-9]+}/progress", handlersManager.EpisodeProgressHandler)
			r.Put("/{id:[0-9]+}/episodes/{epID:[0-9]+}/status", handlersManager.SetEpisodeStatusHandler)
			r.Get("/podcasts/latest_eps", handlersManager.LatestEpisodesHandler)
		})

		r.Route("/user", func(r chi.Router) {
			r.Get("/subscriptions", handlersManager.GetUserPodcastsHandler)
		})

		r.Route("/player", func(r chi.Router) {
			r.Get("/playback_info", handlersManager.PlayerPlaybackInfoHandler)
			r.Put("/playback_info", handlersManager.PlayerPlaybackInfoHandler)

			r.Route("/queue", func(r chi.Router) {
				r.Get("/", handlersManager.QueueHandler)
				r.Put("/", handlersManager.QueueHandler)
				r.Delete("/", handlersManager.QueueHandler)
				// TODO Maybe these paths can be renamed
				r.Post("/add", handlersManager.AddToQueueHandler)
				r.Delete("/remove", handlersManager.DelFromQueueHandler)
			})
		})
	})

	return router
}
