package storage

import (
	"encoding/json"
	"fmt"
	"freelanceledger/model"
	"go.etcd.io/bbolt"
	"sync"
	"time"
)

var buckets = []string{"records", "profiles", "events", "audits", "archives"}

type Store struct {
	db   *bbolt.DB
	mu   sync.RWMutex
	path string
}

func Open(path string) (*Store, error) {
	db, e := bbolt.Open(path, 0600, nil)
	if e != nil {
		return nil, fmt.Errorf("open: %w", e)
	}
	s := &Store{db: db, path: path}
	e = db.Update(func(tx *bbolt.Tx) error {
		for _, n := range buckets {
			if _, x := tx.CreateBucketIfNotExists([]byte(n)); x != nil {
				return x
			}
		}
		return nil
	})
	if e != nil {
		db.Close()
		return nil, e
	}
	return s, nil
}
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	e := s.db.Close()
	s.db = nil
	return e
}
func (s *Store) put(bucket, key string, v any) error {
	b, e := json.Marshal(v)
	if e != nil {
		return fmt.Errorf("marshal: %w", e)
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte(bucket)).Put([]byte(key), b) })
}
func (s *Store) get(bucket, key string, v any) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket)).Get([]byte(key))
		if b == nil {
			return model.ErrNotFound
		}
		return json.Unmarshal(b, v)
	})
}
func (s *Store) Delete(bucket, key string) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte(bucket)).Delete([]byte(key)) })
}
func (s *Store) Path() string                   { return s.path }
func (s *Store) PutRecord(r model.Record) error { return s.put("records", r.ID, r) }
func (s *Store) GetRecord(id string) (model.Record, error) {
	var r model.Record
	e := s.get("records", id, &r)
	return r, e
}
func (s *Store) PutProfile(p model.Profile) error { return s.put("profiles", p.ID, p) }
func (s *Store) GetProfile(id string) (model.Profile, error) {
	var p model.Profile
	e := s.get("profiles", id, &p)
	return p, e
}
func (s *Store) PutEvent(v model.Event) error          { return s.put("events", v.ID, v) }
func (s *Store) PutAudit(v model.Audit) error          { return s.put("audits", v.ID, v) }
func (s *Store) PutArchive(v model.ArchiveEntry) error { return s.put("archives", v.RecordID, v) }
func (s *Store) ListRecords() ([]model.Record, error) {
	out := []model.Record{}
	e := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("records")).ForEach(func(_, v []byte) error {
			var r model.Record
			if e := json.Unmarshal(v, &r); e != nil {
				return e
			}
			out = append(out, r)
			return nil
		})
	})
	return out, e
}
func (s *Store) Count(bucket string) (int, error) {
	n := 0
	e := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return model.ErrNotFound
		}
		return b.ForEach(func(_, v []byte) error {
			if v != nil {
				n++
			}
			return nil
		})
	})
	return n, e
}
func NewEvent(id, rid, typ, payload string) model.Event {
	return model.Event{ID: id, RecordID: rid, Type: typ, Payload: payload, At: time.Now().UTC()}
}
