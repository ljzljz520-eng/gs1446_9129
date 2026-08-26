package model

import (
	"encoding/json"
	"fmt"
)

func Encode(v any) ([]byte, error) {
	b, e := json.Marshal(v)
	if e != nil {
		return nil, fmt.Errorf("encode: %w", e)
	}
	return b, nil
}
func Decode(data []byte, v any) error {
	if len(data) == 0 {
		return ErrNotFound
	}
	if e := json.Unmarshal(data, v); e != nil {
		return fmt.Errorf("decode: %w", e)
	}
	return nil
}
func CloneTags(in []string) []string { out := make([]string, len(in)); copy(out, in); return out }
