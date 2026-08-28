package service

import (
	"context"
	"errors"
	"fmt"
	"freelanceledger/model"
	"freelanceledger/storage"
)

type Processor struct{ Store *storage.Store }

func (p *Processor) Process(ctx context.Context, id string) (model.Record, error) {
	if err := ctx.Err(); err != nil {
		return model.Record{}, err
	}
	r, e := p.Store.GetRecord(id)
	if e != nil {
		return r, fmt.Errorf("load: %w", e)
	}
	if r.IsArchived() {
		return r, model.ErrArchived
	}
	r = r.MarkProcessed()
	if e = p.Store.UpdateRecord(r); e != nil {
		return r, fmt.Errorf("update: %w", e)
	}
	return r, nil
}
func (p *Processor) ImportOrKeep(ctx context.Context, r model.Record) (model.Record, error) {
	e := p.Store.PutRecord(r)
	if e != nil {
		return r, e
	}
	return r, nil
}
func IsDuplicate(err error) bool { return errors.Is(err, model.ErrDuplicate) }
func (p *Processor) Reconcile(ctx context.Context, ids []string) (int, error) {
	count := 0
	for _, id := range ids {
		if _, e := p.Process(ctx, id); e != nil {
			if errors.Is(e, model.ErrArchived) {
				continue
			}
			return count, e
		}
		count++
	}
	return count, nil
}
func (p *Processor) MarkStatus(id, status string) error {
	r, e := p.Store.GetRecord(id)
	if e != nil {
		return e
	}
	if status == "archived" {
		r = r.MarkArchived()
	} else if status == "processed" {
		r = r.MarkProcessed()
	} else {
		return fmt.Errorf("status %s", status)
	}
	return p.Store.UpdateRecord(r)
}
