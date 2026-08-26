package audit

import (
	"fmt"
	"freelanceledger/model"
	"freelanceledger/storage"
	"time"
)

type Logger struct{ Store *storage.Store }

func (l *Logger) Log(actor, action, target, detail string) error {
	a := model.Audit{ID: fmt.Sprintf("%d-%s", time.Now().UnixNano(), target), Actor: actor, Action: action, Target: target, Detail: detail, At: time.Now().UTC()}
	return l.Store.PutAudit(a)
}
func (l *Logger) RecordImport(r model.Record) error { return l.Log("system", "import", r.ID, r.Kind) }
func (l *Logger) RecordArchive(id, reason string) error {
	return l.Log("system", "archive", id, reason)
}
func (l *Logger) Count() int {
	n, e := l.Store.Count("audits")
	if e != nil {
		return 0
	}
	return n
}
