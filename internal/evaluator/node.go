package evaluator

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