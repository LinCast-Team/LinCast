package backend

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"lincast/podcasts"

	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	log "github.com/sirupsen/logrus"
)

const (
	// Frontend path ("/" is the root of the project).
	frontendPath = "/webui/frontend/dist"
)

type spaHandler struct {
	staticPath string
	devMode    bool
}

func (s spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(s.staticPath, r.URL.Path)

	if s.devMode {
		log.Debugln("Dev mode enabled, adding cache prevention headers")

		// Avoid cache if we are on development mode.
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
	}

	// Check if the file exists.
	_, err := pkger.Stat(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		return
	}

	// If there are no errors, serve the requested file from pkger.
	http.FileServer(pkger.Dir(s.staticPath)).ServeHTTP(w, r)
}

// New returns a new instance of the server. To execute it, the method `ListenAndServe` must be called.
func New(port uint16, localServer bool, devMode bool, logRequests bool, podcastsDB *podcasts.Database) *http.Server {
	if podcastsDB == nil {
		log.Panic("'podcastsDB' is nil")
	}

	_podcastsDB = podcastsDB

	// Include the frontend inside the binary.
	_ = pkger.Include(frontendPath)
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

func newRouter(devMode, logRequests bool) *mux.Router {
	router := mux.NewRouter()

	if logRequests {
		router.Use(loggingMiddleware)
	}

	spa := spaHandler{
		staticPath: frontendPath,
		devMode:    devMode,
	}

	router.HandleFunc("/api/v0/feeds/subscribe", subscribeToPodcastHandler).Methods("POST")
	router.HandleFunc("/api/v0/feeds/unsubscribe", unsubscribeToPodcastHandler).Methods("PUT")
	router.PathPrefix("/").Handler(spa)

	return router
}
