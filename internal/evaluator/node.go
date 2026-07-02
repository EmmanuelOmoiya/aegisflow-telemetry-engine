package evaluator

import (
	"fmt"
	"math"
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

type ComparisonNode struct {
	Left     Node
	Operator string
	Right    Node
}

func (n *ComparisonNode) OpCode() string { return n.Operator }

func (n *ComparisonNode) Evaluate(context map[string]interface{}) (interface{}, error) {
	leftVal, err := n.Left.Evaluate(context)
	if err != nil {
		return false, err
	}
	rightVal, err := n.Right.Evaluate(context)
	if err != nil {
		return false, err
	}

	// Handle nil comparisons gracefully (e.g., checking if a missing field is equal to something)
	if leftVal == nil || rightVal == nil {
		switch n.Operator {
		case "==":
			return leftVal == rightVal, nil
		case "!=":
			return leftVal != rightVal, nil
		default:
			return false, nil // Relational operators like >, < return false against nil values
		}
	}

	// Normalize numeric types to float64 to enable clean, direct comparisons
	leftNum, isLeftNum := convertToFloat(leftVal)
	rightNum, isRightNum := convertToFloat(rightVal)

	if isLeftNum && isRightNum {
		switch n.Operator {
		case "==":
			return math.Abs(leftNum-rightNum) < 1e-9, nil
		case "!=":
			return math.Abs(leftNum-rightNum) >= 1e-9, nil
		case ">":
			return leftNum > rightNum, nil
		case "<":
			return leftNum < rightNum, nil
		case ">=":
			return leftNum >= rightNum, nil
		case "<=":
			return leftNum <= rightNum, nil
		default:
			return false, fmt.Errorf("unsupported numerical evaluation operator: %s", n.Operator)
		}
	}

	// Handle string comparison fallback blocks
	leftStr, isLeftStr := leftVal.(string)
	rightStr, isRightStr := rightVal.(string)

	if isLeftStr && isRightStr {
		switch n.Operator {
		case "==":
			return leftStr == rightStr, nil
		case "!=":
			return leftStr != rightStr, nil
		default:
			return false, fmt.Errorf("unsupported alphanumeric evaluation operator: %s", n.Operator)
		}
	}

	// Fallback equality check for matching boolean primitives safely
	if n.Operator == "==" {
		return leftVal == rightVal, nil
	} else if n.Operator == "!=" {
		return leftVal != rightVal, nil
	}

	return false, fmt.Errorf("type mismatch: cannot evaluate comparison %v %s %v", leftVal, n.Operator, rightVal)
}

func convertToFloat(v interface{}) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case float32:
		return float64(t), true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case int32:
		return float64(t), true
	default:
		return 0, false
	}
}

type LogicalNode struct {
	Left     Node
	Operator string
	Right    Node
}

func (n *LogicalNode) OpCode() string { return n.Operator }

func (n *LogicalNode) Evaluate(context map[string]interface{}) (interface{}, error) {
	leftVal, err := n.Left.Evaluate(context)
	if err != nil {
		return false, err
	}

	leftBool, ok := leftVal.(bool)
	if !ok {
		return false, fmt.Errorf("logical operator %s requires boolean operands, got: %v", n.Operator, leftVal)
	}

	// Apply short-circuit optimization patterns early
	switch n.Operator {
	case "AND":
		if !leftBool {
			return false, nil // Short-circuit: false AND anything is always false
		}
	case "OR":
		if leftBool {
			return true, nil // Short-circuit: true OR anything is always true
		}
	}

	// Evaluate the right branch only if the short-circuit conditions above were not met
	rightVal, err := n.Right.Evaluate(context)
	if err != nil {
		return false, err
	}

	rightBool, ok := rightVal.(bool)
	if !ok {
		return false, fmt.Errorf("logical operator %s requires boolean operands, got: %v", n.Operator, rightVal)
	}

	switch n.Operator {
	case "AND":
		return leftBool && rightBool, nil
	case "OR":
		return leftBool || rightBool, nil
	default:
		return false, fmt.Errorf("unsupported logical condition operator: %s", n.Operator)
	}
}
