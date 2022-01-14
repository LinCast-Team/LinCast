package webui

import (
	"net/http"

	"lincast/utils/safe"

	log "github.com/sirupsen/logrus"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"url":           r.URL.String(),
			"userAgent":     safe.Sanitize(r.UserAgent()),
			"remoteAddr":    r.RemoteAddr,
			"tls":           r.TLS != nil,
			"contentLength": r.ContentLength,
		}).Debug("Request received")

		next.ServeHTTP(w, r)
	})
}
