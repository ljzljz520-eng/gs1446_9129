package query

import (
	"freelanceledger/model"
	"freelanceledger/storage"
	"path/filepath"
	"testing"
	"time"
)

func TestQueries(t *testing.T) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	s.PutRecord(model.NewRecord("q", "p", "dev", 9, time.Now()))
	q := &Service{Store: s}
	rows, e := q.ByProfile("p")
	if e != nil || len(rows) != 1 {
		t.Fatal("missing")
	}
}
