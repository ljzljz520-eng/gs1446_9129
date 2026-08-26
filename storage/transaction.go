package storage

import (
	"fmt"
	"freelanceledger/model"
	"go.etcd.io/bbolt"
)

func (s *Store) SaveRecordWithEvent(r model.Record, e model.Event) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		rb, er := model.Encode(r)
		if er != nil {
			return fmt.Errorf("record: %w", er)
		}
		if old := tx.Bucket([]byte("records")).Get([]byte(r.ID)); old != nil {
			return fmt.Errorf("duplicate: %w", model.ErrDuplicate)
		}
		if er = tx.Bucket([]byte("records")).Put([]byte(r.ID), rb); er != nil {
			return er
		}
		eb, er := model.Encode(e)
		if er != nil {
			return er
		}
		return tx.Bucket([]byte("events")).Put([]byte(e.ID), eb)
	})
}
func (s *Store) UpdateRecord(r model.Record) error        { return s.PutRecord(r) }
func (s *Store) ArchiveRecord(a model.ArchiveEntry) error { return s.PutArchive(a) }
