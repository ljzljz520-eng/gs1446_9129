package ledger

import (
	"freelanceledger/model"
	"math"
	"sort"
	"time"
)

type Summary struct {
	ProfileID                           string
	Gross, Processed, Pending, Archived int64
	Count                               int
	Average                             float64
	ByKind                              map[string]int64
	ByMonth                             map[string]int64
}

func BuildSummary(records []model.Record, profile string) Summary {
	s := Summary{ProfileID: profile, ByKind: map[string]int64{}, ByMonth: map[string]int64{}}
	for _, r := range records {
		if profile != "" && r.ProfileID != profile {
			continue
		}
		s.Count++
		s.Gross += r.Amount
		s.ByKind[r.Kind] += r.Amount
		s.ByMonth[r.OccurredAt.UTC().Format("2006-01")] += r.Amount
		switch r.Status {
		case "processed":
			s.Processed += r.Amount
		case "pending":
			s.Pending += r.Amount
		case "archived":
			s.Archived += r.Amount
		}
	}
	if s.Count > 0 {
		s.Average = float64(s.Gross) / float64(s.Count)
	}
	return s
}
func NormalizeAmount(amount int64, currency string) int64 {
	rates := map[string]float64{"USD": 1, "EUR": 1.08, "GBP": 1.27, "JPY": 0.0067}
	r, ok := rates[currency]
	if !ok {
		r = 1
	}
	return int64(math.Round(float64(amount) * r))
}
func TopKinds(records []model.Record, n int) []string {
	totals := map[string]int64{}
	for _, r := range records {
		totals[r.Kind] += r.Amount
	}
	keys := make([]string, 0, len(totals))
	for k := range totals {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return totals[keys[i]] > totals[keys[j]] })
	if n > len(keys) {
		n = len(keys)
	}
	return keys[:n]
}
func MonthRange(records []model.Record) (time.Time, time.Time) {
	if len(records) == 0 {
		return time.Time{}, time.Time{}
	}
	min, max := records[0].OccurredAt, records[0].OccurredAt
	for _, r := range records {
		if r.OccurredAt.Before(min) {
			min = r.OccurredAt
		}
		if r.OccurredAt.After(max) {
			max = r.OccurredAt
		}
	}
	return min, max
}
func FilterSince(records []model.Record, since time.Time) []model.Record {
	out := []model.Record{}
	for _, r := range records {
		if !r.OccurredAt.Before(since) {
			out = append(out, r)
		}
	}
	return out
}
func FilterUntil(records []model.Record, until time.Time) []model.Record {
	out := []model.Record{}
	for _, r := range records {
		if r.OccurredAt.Before(until) {
			out = append(out, r)
		}
	}
	return out
}
func SortNewest(records []model.Record) []model.Record {
	out := append([]model.Record(nil), records...)
	sort.Slice(out, func(i, j int) bool { return out[i].OccurredAt.After(out[j].OccurredAt) })
	return out
}
