package ingest

import (
	"freelanceledger/model"
	"strings"
	"time"
)

func Normalize(r model.Record) model.Record {
	r.ID = strings.TrimSpace(r.ID)
	r.ProfileID = strings.TrimSpace(r.ProfileID)
	r.Kind = strings.ToLower(strings.TrimSpace(r.Kind))
	r.Currency = strings.ToUpper(strings.TrimSpace(r.Currency))
	if r.Currency == "" {
		r.Currency = "USD"
	}
	if r.Status == "" {
		r.Status = "pending"
	}
	if r.OccurredAt.IsZero() {
		r.OccurredAt = time.Now().UTC()
	}
	r.Tags = model.CloneTags(r.Tags)
	return r
}
func NormalizeBatch(rs []model.Record) []model.Record {
	out := make([]model.Record, 0, len(rs))
	for _, r := range rs {
		out = append(out, Normalize(r))
	}
	return out
}
func ValidateBatch(rs []model.Record) error {
	seen := map[string]bool{}
	for _, r := range rs {
		r = Normalize(r)
		if e := r.Validate(); e != nil {
			return e
		}
		if seen[r.ID] {
			return model.ErrDuplicate
		}
		seen[r.ID] = true
	}
	return nil
}
