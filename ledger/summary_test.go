package ledger

import (
	"freelanceledger/model"
	"testing"
	"time"
)

func TestSummary(t *testing.T) {
	s := BuildSummary([]model.Record{{ProfileID: "p", Kind: "dev", Amount: 4, Status: "processed", OccurredAt: time.Now()}}, "p")
	if s.Gross != 4 || s.Count != 1 {
		t.Fatal(s)
	}
}
