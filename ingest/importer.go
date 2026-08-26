package ingest

import (
	"context"
	"fmt"
	"freelanceledger/model"
	"freelanceledger/storage"
)

type Importer struct{ Store *storage.Store }

func (i *Importer) Import(ctx context.Context, r model.Record) (model.Record, error) {
	if err := ctx.Err(); err != nil {
		return r, fmt.Errorf("context: %w", err)
	}
	if err := r.Validate(); err != nil {
		return r, fmt.Errorf("validate: %w", err)
	}
	ev := storage.NewEvent("import-"+r.ID, r.ID, "imported", r.Note)
	if err := i.Store.SaveRecordWithEvent(r, ev); err != nil {
		return r, fmt.Errorf("import record: %w", err)
	}
	return r, nil
}
func (i *Importer) ImportBatch(ctx context.Context, records []model.Record) ([]model.Record, error) {
	out := []model.Record{}
	for _, r := range records {
		saved, e := i.Import(ctx, r)
		if e != nil {
			return out, e
		}
		out = append(out, saved)
	}
	return out, nil
}
func (i *Importer) EnsureProfile(p model.Profile) error {
	if e := p.Validate(); e != nil {
		return e
	}
	return i.Store.PutProfile(p)
}
func (i *Importer) RetryImport(ctx context.Context, r model.Record, attempts int) error {
	var last error
	for n := 0; n < attempts; n++ {
		if _, e := i.Import(ctx, r); e == nil {
			return nil
		} else {
			last = e
		}
	}
	return last
}
func (i *Importer) ImportOrKeep(ctx context.Context, r model.Record) (model.Record, error) {
	saved, e := i.Import(ctx, r)
	if e == nil {
		return saved, nil
	}
	if e == model.ErrDuplicate {
		return r, model.ErrDuplicate
	}
	if err := i.Store.PutRecord(r); err != nil {
		return r, err
	}
	return r, nil
}
