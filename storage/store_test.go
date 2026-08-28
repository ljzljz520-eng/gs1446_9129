package storage

import (
	"freelanceledger/model"
	"path/filepath"
	"testing"
	"time"
)

func TestStoreDuplicate(t *testing.T) {
	s, _ := Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	r := model.NewRecord("x", "p", "k", 1, now())
	e := s.SaveRecordWithEvent(r, NewEvent("e", "x", "import", ""))
	if e != nil {
		t.Fatal(e)
	}
	if e = s.SaveRecordWithEvent(r, NewEvent("e2", "x", "import", "")); e == nil {
		t.Fatal("expected duplicate")
	}
}
func now() time.Time { return time.Now() }
