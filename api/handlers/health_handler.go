package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func ServiceHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{Message: "Health Ok"}); err != nil {
		panic(err)
	}
}

// PingHandler checks the database connection and returns "pong" if successful.
func PingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, "Database connection failed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}
}
