package ledger

import (
	"fmt"
	"freelanceledger/model"
	"time"
)

func ValidateRecord(r model.Record) error {
	if e := r.Validate(); e != nil {
		return e
	}
	if r.Amount > 100000000 {
		return fmt.Errorf("amount exceeds operational maximum")
	}
	if r.OccurredAt.After(time.Now().Add(24 * time.Hour)) {
		return fmt.Errorf("future record")
	}
	return nil
}
func ValidatePeriod(p Period) error {
	if p.Start.IsZero() || p.End.IsZero() {
		return fmt.Errorf("period bounds required")
	}
	if !p.End.After(p.Start) {
		return fmt.Errorf("period order invalid")
	}
	return nil
}
func Reconcile(records []model.Record) error {
	seen := map[string]bool{}
	for _, r := range records {
		if seen[r.ID] {
			return model.ErrDuplicate
		}
		seen[r.ID] = true
		if e := ValidateRecord(r); e != nil {
			return e
		}
	}
	return nil
}
func StatusSummary(records []model.Record) map[string]int {
	out := map[string]int{}
	for _, r := range records {
		out[r.Status]++
	}
	return out
}
func AmountByCurrency(records []model.Record) map[string]int64 {
	out := map[string]int64{}
	for _, r := range records {
		out[r.Currency] += r.Amount
	}
	return out
}
func AmountByKind(records []model.Record) map[string]int64 {
	out := map[string]int64{}
	for _, r := range records {
		out[r.Kind] += r.Amount
	}
	return out
}
func Earliest(records []model.Record) model.Record {
	var out model.Record
	for i, r := range records {
		if i == 0 || r.OccurredAt.Before(out.OccurredAt) {
			out = r
		}
	}
	return out
}
func Latest(records []model.Record) model.Record {
	var out model.Record
	for i, r := range records {
		if i == 0 || r.OccurredAt.After(out.OccurredAt) {
			out = r
		}
	}
	return out
}
func ActiveRecords(records []model.Record) []model.Record {
	out := []model.Record{}
	for _, r := range records {
		if !r.IsArchived() {
			out = append(out, r)
		}
	}
	return out
}
func ArchiveCandidates(records []model.Record, before time.Time) []model.Record {
	out := []model.Record{}
	for _, r := range records {
		if r.OccurredAt.Before(before) && !r.IsArchived() {
			out = append(out, r)
		}
	}
	return out
}
func WithTag(records []model.Record, tag string) []model.Record {
	out := []model.Record{}
	for _, r := range records {
		for _, t := range r.Tags {
			if t == tag {
				out = append(out, r)
				break
			}
		}
	}
	return out
}
func AddTag(r model.Record, tag string) model.Record {
	for _, t := range r.Tags {
		if t == tag {
			return r
		}
	}
	r.Tags = append(r.Tags, tag)
	return r
}
func RemoveTag(r model.Record, tag string) model.Record {
	out := r.Tags[:0]
	for _, t := range r.Tags {
		if t != tag {
			out = append(out, t)
		}
	}
	r.Tags = out
	return r
}
func CurrencySupported(code string) bool {
	switch code {
	case "USD", "EUR", "GBP", "JPY":
		return true
	default:
		return false
	}
}
func NormalizeCurrency(r model.Record) model.Record {
	if r.Currency == "" {
		r.Currency = "USD"
	}
	return r
}
