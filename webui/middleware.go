package webui

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"URI":           r.RequestURI,
			"method":        r.Method,
			"remoteAddr":    r.RemoteAddr,
			"tls":           r.TLS != nil,
			"contentLength": r.ContentLength,
		}).Debug("New request")
		next.ServeHTTP(w, r)
	})
}
