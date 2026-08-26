package service

import (
	"context"
	"freelanceledger/model"
	"freelanceledger/storage"
	"path/filepath"
	"testing"
	"time"
)

func TestProcess(t *testing.T) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	r := model.NewRecord("p", "u", "dev", 1, time.Now())
	s.PutRecord(r)
	p := &Processor{Store: s}
	got, e := p.Process(context.Background(), r.ID)
	if e != nil || got.Status != "processed" {
		t.Fatal(got, e)
	}
}
