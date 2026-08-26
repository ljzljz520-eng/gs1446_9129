package service

import (
	"errors"
	"freelanceledger/model"
	"freelanceledger/storage"
)

func CheckDuplicate(s *storage.Store, r model.Record) (bool, error) {
	old, e := s.GetRecord(r.ID)
	if e != nil {
		if errors.Is(e, model.ErrNotFound) {
			return false, nil
		}
		return false, e
	}
	if old.ID == r.ID {
		return true, nil
	}
	return false, nil
}
func PreserveExisting(s *storage.Store, r model.Record) error {
	dup, e := CheckDuplicate(s, r)
	if e != nil {
		return e
	}
	if dup {
		return model.ErrDuplicate
	}
	return s.PutRecord(r)
}
