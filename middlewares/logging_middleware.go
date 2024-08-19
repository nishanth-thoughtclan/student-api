package middlewares

import (
	"log"
	"net/http"
	"os"
)

const filePerms = 0666

var (
	logFile, _ = os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePerms)
	logger     = log.New(logFile, "", log.LstdFlags)
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		logger.Printf("REQUEST: %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(rw, r)
		if rw.statusCode >= 400 {
			logger.Printf("ERROR: %s %s - Status Code: %d", r.Method, r.RequestURI, rw.statusCode)
		} else {
			logger.Printf("RESPONSE: %s %s - Status Code: %d", r.Method, r.RequestURI, rw.statusCode)
		}
	})
}
