package middlewares

import (
	"log"
	"net/http"
	"os"
)

const FilePerms = 0666

var (
	LogFile, _ = os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, FilePerms)
	Logger     = log.New(LogFile, "", log.LstdFlags)
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
		Logger.Printf("REQUEST: %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(rw, r)
		if rw.statusCode >= 400 {
			Logger.Printf("ERROR: %s %s - Status Code: %d", r.Method, r.RequestURI, rw.statusCode)
		} else {
			Logger.Printf("RESPONSE: %s %s - Status Code: %d", r.Method, r.RequestURI, rw.statusCode)
		}
	})
}
