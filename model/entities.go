package model

import "time"

type Record struct {
	ID, ProfileID, Kind, Currency, Note string
	Amount                              int64
	OccurredAt                          time.Time
	Status                              string
	Tags                                []string
}
type Profile struct {
	ID, Name, Email, TaxRegion string
	Active                     bool
	CreatedAt                  time.Time
}
type Event struct {
	ID, RecordID, Type, Payload string
	At                          time.Time
}
type Audit struct {
	ID, Actor, Action, Target, Detail string
	At                                time.Time
}
type ArchiveEntry struct {
	RecordID, Reason string
	ArchivedAt       time.Time
}

func NewRecord(id, profile, kind string, amount int64, at time.Time) Record {
	return Record{ID: id, ProfileID: profile, Kind: kind, Amount: amount, OccurredAt: at, Currency: "USD", Status: "pending", Tags: []string{}}
}
func (r Record) Validate() error {
	if r.ID == "" {
		return ErrInvalidRecord
	}
	if r.ProfileID == "" {
		return ErrInvalidRecord
	}
	if r.Amount <= 0 {
		return ErrInvalidRecord
	}
	if r.Kind == "" {
		return ErrInvalidRecord
	}
	return nil
}
func (r Record) IsArchived() bool      { return r.Status == "archived" }
func (r Record) MarkProcessed() Record { r.Status = "processed"; return r }
func (r Record) MarkArchived() Record  { r.Status = "archived"; return r }
func NewProfile(id, name, email, region string) Profile {
	return Profile{ID: id, Name: name, Email: email, TaxRegion: region, Active: true, CreatedAt: time.Now().UTC()}
}
func (p Profile) Validate() error {
	if p.ID == "" || p.Name == "" {
		return ErrInvalidProfile
	}
	return nil
}
func (p Profile) Deactivate() Profile { p.Active = false; return p }
