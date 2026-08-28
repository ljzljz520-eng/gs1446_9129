package main

import (
	"freelanceledger/api"
	"freelanceledger/ingest"
	"freelanceledger/query"
	"freelanceledger/service"
	"freelanceledger/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	path := os.Getenv("LEDGER_DB")
	if path == "" {
		path = "ledger.db"
	}
	st, e := storage.Open(path)
	if e != nil {
		log.Fatal(e)
	}
	defer st.Close()
	srv := &api.Server{Importer: &ingest.Importer{Store: st}, Processor: &service.Processor{Store: st}, Query: &query.Service{Store: st}}
	log.Println(http.ListenAndServe(":8080", srv.Routes()))
}
