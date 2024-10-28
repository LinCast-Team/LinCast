package server

import (
	"lincast/internal/app/service"
	"lincast/internal/domain/entities"
	"lincast/internal/infra/http/handlers"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
)

// NewServer creates a new instance of http.Server with the specified configurations.
// It takes the following parameters:
// - port: The port number on which the server will listen.
// - localServer: A boolean indicating whether the server should only listen on the local loopback interface (127.0.0.1).
// - devMode: A boolean indicating whether the server is running in development mode.
// - logRequests: A boolean indicating whether to log incoming requests.
// - db: A pointer to a gorm.DB instance representing the database connection.
// - manualUpdate: A channel used for manual updates of podcast data.
//
// It returns a pointer to the created http.Server instance.
func NewServer(
	port uint,
	localServer bool,
	devMode bool,
	logRequests bool,
	manualUpdate chan *entities.Podcast,
	userApplicationService *service.UserApplicationService,
	tokenAuth *jwtauth.JWTAuth,
) *http.Server {
	handlersManager := handlers.NewManager(
		manualUpdate,
	)

	authHandlers := handlers.NewAuthHandlers(userApplicationService)

	router := createRouter(handlersManager, authHandlers, tokenAuth)

	var addr string
	if localServer {
		addr = "127.0.0.1"
	}

	s := createServer(addr, int(port), router)

	log.WithFields(log.Fields{
		"address":        s.Addr,
		"readTimeout":    s.ReadTimeout.String(),
		"writeTimeout":   s.WriteTimeout.String(),
		"maxHeaderBytes": s.MaxHeaderBytes,
	}).Debugln("Starting server")

	return s
}

// createServer creates a new HTTP server with the given address, port, and router.
// The server is configured with a read timeout and write timeout of 15 seconds each,
// and a maximum header size of 8KB.
//
// Parameters:
// - addr: The address to bind the server to.
// - port: The port number to listen on.
// - router: The router to handle incoming HTTP requests.
//
// Returns:
// - A pointer to the created http.Server instance.
func createServer(addr string, port int, router http.Handler) *http.Server {
	return &http.Server{
		Addr:           addr + ":" + strconv.Itoa(port),
		Handler:        router,
		ReadTimeout:    time.Second * 15,
		WriteTimeout:   time.Second * 15,
		MaxHeaderBytes: 8000, // 8KB
	}
}

// createRouter returns a new instance of the router with their paths already set.
func createRouter(handlersManager *handlers.Manager, authHandlers *handlers.AuthHandlers, tokenAuth *jwtauth.JWTAuth) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.Compress(5))
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Route("/api/v0", func(r chi.Router) {
		// TODO Implement user authentication
		r.Route("/auth", func(r chi.Router) {
			r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("TODO!"))
			})
			r.Post("/signup", authHandlers.SignUpHandler)
		})

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))

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
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("TODO! /api/v0/user"))
				})
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
	})

	return router
}
