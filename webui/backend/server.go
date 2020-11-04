package backend

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

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
		log.WithField("requestedPath", r.RequestURI).Warnln("Unrecognized path requested, redirecting to root")
		// If not, then redirect the request to the root path.
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)

		return
	}

	// If there are no errors, serve the requested file from pkger.
	http.FileServer(pkger.Dir(s.staticPath)).ServeHTTP(w, r)
}

func setup(port uint16, localServer bool, devMode bool, logRequests bool) error {
	router := mux.NewRouter()

	if logRequests {
		router.Use(loggingMiddleware)
	}

	spa := spaHandler{
		staticPath: frontendPath,
		devMode:    devMode,
	}

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

	router.PathPrefix("/").Handler(spa)
	// APIs here

	log.WithFields(log.Fields{
		"address":        s.Addr,
		"readTimeout":    s.ReadTimeout.String(),
		"writeTimeout":   s.WriteTimeout.String(),
		"maxHeaderBytes": s.MaxHeaderBytes,
	}).Infoln("Starting server")

	return s.ListenAndServe()
}

func Run(port uint16, localServer bool, devMode bool, logRequests bool) error {
	// Include the frontend inside the binary.
	_ = pkger.Include(frontendPath)

	return setup(port, localServer, devMode, logRequests)
}
