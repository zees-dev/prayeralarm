package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/zees-dev/prayeralarm/prayer"
)

type handler struct {
	Router    http.Handler
	prayerSvc *prayer.Service
}

func loggedRouter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"method=%s path=%s duration=%dms\n",
			r.Method,
			r.URL.EscapedPath(),
			time.Since(start).Milliseconds(),
		)
	})
}

func NewHandler(prayerSvc *prayer.Service) handler {
	h := handler{
		prayerSvc: prayerSvc,
	}

	r := mux.NewRouter()

	// Static file serve
	// r.HandleFunc("/", HomeHandler)

	r.HandleFunc("/api/health", h.healthHandler).Methods(http.MethodGet)

	r.HandleFunc("/api/timings", h.timingsHandler).Methods(http.MethodGet)

	r.HandleFunc("/api/timings/toggle/{index}", h.timingsUpdateHandler).Methods(http.MethodPost)

	r.HandleFunc("/api/timings/off", h.timingsTurnOffHandler).Methods(http.MethodPost)

	r.HandleFunc("/api/timings/on", h.timingsTurnOnHandler).Methods(http.MethodPost)

	h.Router = loggedRouter(r)

	return h
}

func (h handler) healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func (h handler) timingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "text/html" {
		w.Header().Set("Content-Type", "text/html")

		h.prayerSvc.DisplayPrayerTimings(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.prayerSvc.Prayers)
}

func (h handler) timingsUpdateHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	index, ok := params["index"]
	if !ok {
		http.Error(w, "missing index parameter", http.StatusBadRequest)
		return
	}

	intIndex, err := strconv.Atoi(index)
	if err != nil {
		http.Error(w, "integer value required for index parameter", http.StatusBadRequest)
		return
	}

	if intIndex > len(h.prayerSvc.Prayers) {
		http.Error(w, "invalid index value provided", http.StatusBadRequest)
		return
	}

	h.prayerSvc.ToggleAdhan(uint8(intIndex))
	w.WriteHeader(http.StatusNoContent)
}

func (h handler) timingsTurnOffHandler(w http.ResponseWriter, r *http.Request) {
	h.prayerSvc.TurnOffAllAdhan()
	w.WriteHeader(http.StatusNoContent)
}

func (h handler) timingsTurnOnHandler(w http.ResponseWriter, r *http.Request) {
	h.prayerSvc.TurnOnAllAdhan()
	w.WriteHeader(http.StatusNoContent)
}
