package freelanceledger

import (
	"context"
	"freelanceledger/ingest"
	"freelanceledger/model"
	"freelanceledger/storage"
	"path/filepath"
	"testing"
	"time"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	p := filepath.Join(t.TempDir(), "ledger.db")
	s, e := storage.Open(p)
	if e != nil {
		t.Fatal(e)
	}
	r := model.NewRecord("persist-1", "p1", "consulting", 120, time.Now())
	if e = s.PutRecord(r); e != nil {
		t.Fatal(e)
	}
	s.Close()
	s, e = storage.Open(p)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	if _, e = s.GetRecord(r.ID); e != nil {
		t.Fatal(e)
	}
}
func TestStoreRoundTrip(t *testing.T) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	p := model.NewProfile("p", "Name", "n@x", "US")
	if e := s.PutProfile(p); e != nil {
		t.Fatal(e)
	}
	got, e := s.GetProfile("p")
	if e != nil || got.Name != "Name" {
		t.Fatal(got, e)
	}
}
func TestWorkflowOne(t *testing.T) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	i := &ingest.Importer{Store: s}
	r := model.NewRecord("w1", "p", "design", 10, time.Now())
	if _, e := i.Import(context.Background(), r); e != nil {
		t.Fatal(e)
	}
}
func TestWorkflowTwo(t *testing.T) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	r := model.NewRecord("w2", "p", "dev", 20, time.Now())
	if e := s.PutRecord(r); e != nil {
		t.Fatal(e)
	}
	if got, e := s.GetRecord("w2"); e != nil || got.Status != "pending" {
		t.Fatal(got, e)
	}
}
func TestWorkflowThree(t *testing.T) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	r := model.NewRecord("w3", "p", "writing", 30, time.Now())
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
}
func TestBusinessChain13(t *testing.T) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	i := &ingest.Importer{Store: s}
	r := model.NewRecord("dup", "p", "dev", 25, time.Now())
	if _, e := i.Import(context.Background(), r); e != nil {
		t.Fatal(e)
	}
	if _, e := i.ImportOrKeep(context.Background(), r); e == nil {
		t.Fatalf("duplicate import was accepted")
	}
}
