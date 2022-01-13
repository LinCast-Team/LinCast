package webui

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"userAgent":     r.UserAgent(),
			"remoteAddr":    r.RemoteAddr,
			"tls":           r.TLS != nil,
			"contentLength": r.ContentLength,
		}).Debug("Request received")

		next.ServeHTTP(w, r)
	})
}
