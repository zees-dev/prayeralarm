package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/zees-dev/prayeralarm/prayer"
)

type server struct {
	router    *mux.Router
	prayerSvc *prayer.Service
}

func NewServer(prayerSvc *prayer.Service) *server {
	s := &server{
		router:    mux.NewRouter(),
		prayerSvc: prayerSvc,
	}
	s.initializeRoutes()
	return s
}

func (s *server) Run(port uint) {
	httpServer := &http.Server{
		Handler:      handlers.CORS()(loggedHandler(s.router)),
		Addr:         fmt.Sprintf("127.0.0.1:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Running prayeralarm server on port %d...", port)
	log.Fatal(httpServer.ListenAndServe())
}

func loggedHandler(next http.Handler) http.Handler {
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

func (s *server) initializeRoutes() {
	s.router.HandleFunc("/api/health", s.healthHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/api/timings", s.timingsHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/api/timings/toggle/{index}", s.timingsUpdateHandler).Methods(http.MethodPost)
	s.router.HandleFunc("/api/timings/off", s.timingsTurnOffHandler).Methods(http.MethodPost)
	s.router.HandleFunc("/api/timings/on", s.timingsTurnOnHandler).Methods(http.MethodPost)
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("client/public")))
}

func (s *server) healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func (s *server) timingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "text/html" {
		w.Header().Set("Content-Type", "text/html")

		s.prayerSvc.DisplayPrayerTimings(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.prayerSvc.DailyPrayerTimings)
}

func (s *server) timingsUpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	prayer, err := s.prayerSvc.ToggleAdhan(uint8(intIndex))
	if err != nil {
		http.Error(w, fmt.Sprintf("error otggling adhan; err=%s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prayer)
}

func (s *server) timingsTurnOffHandler(w http.ResponseWriter, r *http.Request) {
	s.prayerSvc.TurnOffAllAdhan()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.prayerSvc.DailyPrayerTimings)
}

func (s *server) timingsTurnOnHandler(w http.ResponseWriter, r *http.Request) {
	s.prayerSvc.TurnOnAllAdhan()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.prayerSvc.DailyPrayerTimings)
}
