package ledger

import (
	"freelanceledger/model"
	"time"
)

type Period struct {
	Start, End time.Time
	Records    []model.Record
}

func NewPeriod(year int, month time.Month) Period {
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	return Period{Start: start, End: start.AddDate(0, 1, 0)}
}
func (p Period) Include(r model.Record) bool {
	return !r.OccurredAt.Before(p.Start) && r.OccurredAt.Before(p.End)
}
func (p Period) Add(records []model.Record) Period {
	p.Records = nil
	for _, r := range records {
		if p.Include(r) {
			p.Records = append(p.Records, r)
		}
	}
	return p
}
func (p Period) Gross() int64 {
	var n int64
	for _, r := range p.Records {
		n += r.Amount
	}
	return n
}
func (p Period) Deductible() int64 {
	var n int64
	for _, r := range p.Records {
		n += Deduction(r.Kind, r.Amount)
	}
	return n
}
func (p Period) Net() int64 { return p.Gross() - p.Deductible() }
func (p Period) CountStatus(status string) int {
	n := 0
	for _, r := range p.Records {
		if r.Status == status {
			n++
		}
	}
	return n
}
