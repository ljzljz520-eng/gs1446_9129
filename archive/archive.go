package archive

import (
	"context"
	"fmt"
	"freelanceledger/model"
	"freelanceledger/storage"
	"time"
)

type Service struct{ Store *storage.Store }

func (a *Service) Archive(ctx context.Context, id, reason string) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	r, e := a.Store.GetRecord(id)
	if e != nil {
		return e
	}
	if r.IsArchived() {
		return model.ErrArchived
	}
	r = r.MarkArchived()
	if e = a.Store.UpdateRecord(r); e != nil {
		return fmt.Errorf("archive update: %w", e)
	}
	return a.Store.ArchiveRecord(model.ArchiveEntry{RecordID: id, Reason: reason, ArchivedAt: time.Now().UTC()})
}
func (a *Service) Restore(id string) error {
	r, e := a.Store.GetRecord(id)
	if e != nil {
		return e
	}
	if !r.IsArchived() {
		return nil
	}
	r.Status = "processed"
	return a.Store.UpdateRecord(r)
}
func (a *Service) Purge(id string) error {
	if e := a.Store.Delete("records", id); e != nil {
		return e
	}
	return a.Store.Delete("archives", id)
}
func (a *Service) Archived() ([]model.Record, error) {
	out, e := a.Store.ListRecords()
	if e != nil {
		return nil, e
	}
	r := []model.Record{}
	for _, x := range out {
		if x.IsArchived() {
			r = append(r, x)
		}
	}
	return r, nil
}
