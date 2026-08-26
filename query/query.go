package query

import (
	"freelanceledger/model"
	"freelanceledger/storage"
	"sort"
	"strings"
)

type Service struct{ Store *storage.Store }

func (s *Service) All() ([]model.Record, error) {
	r, e := s.Store.ListRecords()
	sort.Slice(r, func(i, j int) bool { return r[i].OccurredAt.Before(r[j].OccurredAt) })
	return r, e
}
func (s *Service) ByProfile(profile string) ([]model.Record, error) {
	all, e := s.All()
	if e != nil {
		return nil, e
	}
	out := []model.Record{}
	for _, r := range all {
		if r.ProfileID == profile {
			out = append(out, r)
		}
	}
	return out, nil
}
func (s *Service) Search(term string) ([]model.Record, error) {
	all, e := s.All()
	if e != nil {
		return nil, e
	}
	term = strings.ToLower(term)
	out := []model.Record{}
	for _, r := range all {
		if strings.Contains(strings.ToLower(r.Note), term) || strings.Contains(strings.ToLower(r.Kind), term) {
			out = append(out, r)
		}
	}
	return out, nil
}
func (s *Service) Total(profile string) int64 {
	rs, _ := s.ByProfile(profile)
	var n int64
	for _, r := range rs {
		if r.Status != "archived" {
			n += r.Amount
		}
	}
	return n
}
func (s *Service) ByStatus(status string) ([]model.Record, error) {
	all, e := s.All()
	if e != nil {
		return nil, e
	}
	out := []model.Record{}
	for _, r := range all {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out, nil
}
