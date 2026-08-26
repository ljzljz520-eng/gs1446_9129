package policy

import (
	"freelanceledger/model"
	"strings"
)

type Classification struct{ Kind, Confidence, Reason string }

func Classify(note string) Classification {
	n := strings.ToLower(note)
	switch {
	case strings.Contains(n, "client"):
		return Classification{"consulting", "high", "client keyword"}
	case strings.Contains(n, "flight") || strings.Contains(n, "hotel"):
		return Classification{"travel", "medium", "travel keyword"}
	case strings.Contains(n, "laptop") || strings.Contains(n, "software"):
		return Classification{"equipment", "medium", "equipment keyword"}
	default:
		return Classification{"other", "low", "no matching keyword"}
	}
}
func ApplyClassification(r model.Record) model.Record {
	if r.Kind == "" || r.Kind == "other" {
		r.Kind = Classify(r.Note).Kind
	}
	return r
}
