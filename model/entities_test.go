package model

import (
	"testing"
	"time"
)

func TestRecordValidation(t *testing.T) {
	if NewRecord("", "p", "x", 1, time.Now()).Validate() == nil {
		t.Fatal("expected invalid")
	}
	if NewRecord("1", "p", "x", 1, time.Now()).Validate() != nil {
		t.Fatal("valid rejected")
	}
}
func TestProfileLifecycle(t *testing.T) {
	p := NewProfile("p", "A", "a@x", "US")
	if p.Deactivate().Active {
		t.Fatal("active")
	}
}
