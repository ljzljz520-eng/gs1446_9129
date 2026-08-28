package policy

import (
	"fmt"
	"freelanceledger/model"
	"strings"
)

type Rule struct {
	Name    string
	Limit   int64
	Kinds   []string
	Enabled bool
}
type Engine struct{ Rules []Rule }

func (e *Engine) Check(r model.Record) error {
	for _, rule := range e.Rules {
		if !rule.Enabled {
			continue
		}
		if rule.Limit > 0 && r.Amount > rule.Limit {
			return fmt.Errorf("%s: amount exceeds limit", rule.Name)
		}
		if len(rule.Kinds) > 0 && !contains(rule.Kinds, r.Kind) {
			return fmt.Errorf("%s: kind rejected", rule.Name)
		}
	}
	return nil
}
func contains(xs []string, v string) bool {
	for _, x := range xs {
		if strings.EqualFold(x, v) {
			return true
		}
	}
	return false
}
func Default() Engine {
	return Engine{Rules: []Rule{{Name: "positive", Limit: 0, Enabled: true}, {Name: "review", Limit: 100000, Enabled: true}}}
}
func (e *Engine) Add(rule Rule) { e.Rules = append(e.Rules, rule) }
func (e *Engine) Remove(name string) {
	out := e.Rules[:0]
	for _, r := range e.Rules {
		if r.Name != name {
			out = append(out, r)
		}
	}
	e.Rules = out
}
func (e *Engine) Explain(r model.Record) []string {
	out := []string{}
	for _, rule := range e.Rules {
		if rule.Enabled {
			if e.CheckSingle(rule, r) == nil {
				out = append(out, rule.Name+":pass")
			} else {
				out = append(out, rule.Name+":fail")
			}
		}
	}
	return out
}
func (e *Engine) CheckSingle(rule Rule, r model.Record) error {
	if rule.Limit > 0 && r.Amount > rule.Limit {
		return fmt.Errorf("limit")
	}
	if len(rule.Kinds) > 0 && !contains(rule.Kinds, r.Kind) {
		return fmt.Errorf("kind")
	}
	return nil
}
func ValidateProfile(p model.Profile) error {
	if !p.Active {
		return fmt.Errorf("profile inactive")
	}
	if p.Email == "" {
		return model.ErrInvalidProfile
	}
	return nil
}
