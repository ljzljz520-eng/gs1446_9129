package workflow

import (
	"freelanceledger/model"
	"testing"
)

func TestStages(t *testing.T) {
	if e := ValidateStages([]string{"a", "b", "c", "d"}, 4); e != nil {
		t.Fatal(e)
	}
	if model.ErrNotFound == nil {
		t.Fatal("sentinel")
	}
}
