package api

import (
	"context"
	"encoding/json"
	"freelanceledger/ingest"
	"freelanceledger/model"
	"freelanceledger/query"
	"freelanceledger/service"
	"net/http"
)

type Server struct {
	Importer  *ingest.Importer
	Processor *service.Processor
	Query     *query.Service
}

func (s *Server) Routes() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/records", s.records)
	m.HandleFunc("/records/process", s.process)
	return m
}
func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var rec model.Record
		if e := json.NewDecoder(r.Body).Decode(&rec); e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		if _, e := s.Importer.Import(r.Context(), rec); e != nil {
			http.Error(w, e.Error(), 409)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}
	all, e := s.Query.All()
	if e != nil {
		http.Error(w, e.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(all)
}
func (s *Server) process(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	rec, e := s.Processor.Process(context.Background(), id)
	if e != nil {
		http.Error(w, e.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(rec)
}
