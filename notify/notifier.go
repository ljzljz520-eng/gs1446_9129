package notify

import (
	"context"
	"fmt"
	"freelanceledger/model"
	"sync"
)

type Message struct {
	To, Subject, Body string
	RecordID          string
}
type Sender interface {
	Send(context.Context, Message) error
}
type MemorySender struct {
	mu       sync.Mutex
	Messages []Message
}

func (m *MemorySender) Send(ctx context.Context, msg Message) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Messages = append(m.Messages, msg)
	return nil
}
func (m *MemorySender) Count() int { m.mu.Lock(); defer m.mu.Unlock(); return len(m.Messages) }
func BuildProcessed(p model.Profile, r model.Record) Message {
	return Message{To: p.Email, Subject: "Income processed", Body: fmt.Sprintf("Record %s amount %d", r.ID, r.Amount), RecordID: r.ID}
}
func BuildArchived(p model.Profile, r model.Record) Message {
	return Message{To: p.Email, Subject: "Income archived", Body: fmt.Sprintf("Record %s archived", r.ID), RecordID: r.ID}
}
func SendRecord(ctx context.Context, s Sender, p model.Profile, r model.Record) error {
	if p.Email == "" {
		return fmt.Errorf("missing recipient")
	}
	if r.IsArchived() {
		return s.Send(ctx, BuildArchived(p, r))
	}
	return s.Send(ctx, BuildProcessed(p, r))
}
