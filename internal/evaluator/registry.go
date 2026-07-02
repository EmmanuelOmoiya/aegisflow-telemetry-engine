package evaluator

import (
	"sync"
)

type RuleRegistry struct {
	mu    sync.RWMutex
	rules map[string]Node
}

var GlobalRegistry = NewRuleRegistry()

func NewRuleRegistry() *RuleRegistry {
	return &RuleRegistry{
		rules: make(map[string]Node),
	}
}

func (r *RuleRegistry) Register(ruleID string, ruleStr string) error {
	l := NewLexer(ruleStr)
	p := NewParser(l)
	
	astRoot, err := p.ParseRule()
	if err != nil {
		return err
	}

	r.mu.Lock()
	r.rules[ruleID] = astRoot
	r.mu.Unlock()
	return nil
}

func (r *RuleRegistry) Lookup(ruleID string) (Node, bool) {
	r.mu.RLock()
	node, exists := r.rules[ruleID]
	r.mu.RUnlock()
	return node, exists
}

func (r *RuleRegistry) Clear() {
	r.mu.Lock()
	r.rules = make(map[string]Node)
	r.mu.Unlock()
}