package api

import (
	"net/http"
	"time"
)

func HealthHandler(start time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","uptime":"` + time.Since(start).String() + `"}`))
	}
}
