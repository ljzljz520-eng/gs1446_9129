package archive

import (
	"context"
	"freelanceledger/model"
	"freelanceledger/storage"
	"path/filepath"
	"testing"
	"time"
)

func TestArchive(t *testing.T) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	r := model.NewRecord("a", "p", "dev", 9, time.Now())
	s.PutRecord(r)
	a := &Service{Store: s}
	if e := a.Archive(context.Background(), r.ID, "year-end"); e != nil {
		t.Fatal(e)
	}
}
