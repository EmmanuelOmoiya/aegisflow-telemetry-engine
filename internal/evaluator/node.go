package evaluator

import (
	"strings"
)

type Node interface {
	Evaluate(context map[string]interface{}) (interface{}, error)
}

type DataType string

const (
	TypeString  DataType = "STRING"
	TypeNumber  DataType = "NUMBER"
	TypeBoolean DataType = "BOOLEAN"
	TypeUnknown DataType = "UNKNOWN"
)

type OperationalNode interface {
	Node
	OpCode() string
}

type LiteralNode struct {
	Value interface{}
	Kind  DataType
}

func (n *LiteralNode) Evaluate(context map[string]interface{}) (interface{}, error) {
	return n.Value, nil
}

type IdentifierNode struct {
	Name string
}

func (n *IdentifierNode) Evaluate(context map[string]interface{}) (interface{}, error) {
	if !strings.Contains(n.Name, ".") {
		val, exists := context[n.Name]
		if !exists {
			return nil, nil // Gracefully return nil if the property doesn't exist
		}
		return val, nil
	}

	// Handle dot-notation path traversal for nested JSON attributes
	parts := strings.Split(n.Name, ".")
	var current interface{} = context

	for _, part := range parts {
		currentMap, ok := current.(map[string]interface{})
		if !ok {
			return nil, nil // Path breaks prematurely, return nil
		}

		val, exists := currentMap[part]
		if !exists {
			return nil, nil
		}
		current = val
	}

	return current, nil
}