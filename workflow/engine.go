package workflow

import (
	"context"
	"fmt"
	"freelanceledger/archive"
	"freelanceledger/audit"
	"freelanceledger/ingest"
	"freelanceledger/model"
	"freelanceledger/service"
)

type Engine struct {
	Importer  *ingest.Importer
	Processor *service.Processor
	Archive   *archive.Service
	Audit     *audit.Logger
}
type Result struct {
	Record model.Record
	Stages []string
}

func (e *Engine) Register(ctx context.Context, r model.Record) (Result, error) {
	out := Result{Stages: []string{"received"}}
	if _, x := e.Importer.Import(ctx, r); x != nil {
		return out, x
	}
	out.Record = r
	out.Stages = append(out.Stages, "validated", "saved")
	if e.Audit != nil {
		if x := e.Audit.RecordImport(r); x != nil {
			return out, fmt.Errorf("audit: %w", x)
		}
	}
	out.Stages = append(out.Stages, "displayed")
	return out, nil
}
func (e *Engine) Process(ctx context.Context, id string) (Result, error) {
	out := Result{Stages: []string{"registered"}}
	r, x := e.Processor.Process(ctx, id)
	if x != nil {
		return out, x
	}
	out.Record = r
	out.Stages = append(out.Stages, "processed", "archived-ready", "queried")
	return out, nil
}
func (e *Engine) Close(ctx context.Context, id, reason string) (Result, error) {
	out := Result{Stages: []string{"submitted"}}
	if x := e.Archive.Archive(ctx, id, reason); x != nil {
		return out, x
	}
	r, x := e.Processor.Process(ctx, id)
	if x != nil {
		return out, x
	}
	out.Record = r.MarkArchived()
	out.Stages = append(out.Stages, "reviewed", "notified", "tracked")
	if e.Audit != nil {
		_ = e.Audit.RecordArchive(id, reason)
	}
	return out, nil
}
func ValidateStages(stages []string, min int) error {
	if len(stages) < min {
		return fmt.Errorf("workflow has %d stages", len(stages))
	}
	return nil
}
