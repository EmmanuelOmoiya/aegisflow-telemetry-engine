package tests

import (
	"testing"

	"github.com/EmmanuelOmoiya/aegisflow-telemetry-engine/internal/evaluator"
)

func TestRuleASTCompilationAndEvaluation(t *testing.T) {
	// Setup test matrix with varied complex rule scenarios
	tests := []struct {
		name       string
		ruleStr    string
		context    map[string]interface{}
		wantResult bool
		expectErr  bool
	}{
		{
			name:    "Simple String Equality Match",
			ruleStr: `device_id == "sensor-alpha"`,
			context: map[string]interface{}{
				"device_id": "sensor-alpha",
			},
			wantResult: true,
			expectErr:  false,
		},
		{
			name:    "Numeric Relational Comparison",
			ruleStr: `payload.temperature > 85.5`,
			context: map[string]interface{}{
				"payload": map[string]interface{}{
					"temperature": 92.4,
				},
			},
			wantResult: true,
			expectErr:  false,
		},
		{
			name:    "Type Normalization (Integer matching Float representation)",
			ruleStr: `payload.status_code == 200.0`,
			context: map[string]interface{}{
				"payload": map[string]interface{}{
					"status_code": 200, // passed as native int
				},
			},
			wantResult: true,
			expectErr:  false,
		},
		{
			name:    "Compound AND Condition with Short-Circuiting",
			ruleStr: `source == "gateway_01" AND payload.humidity <= 45.0`,
			context: map[string]interface{}{
				"source": "gateway_01",
				"payload": map[string]interface{}{
					"humidity": 40.2,
				},
			},
			wantResult: true,
			expectErr:  false,
		},
		{
			name:    "Compound OR Condition Short-Circuit (Left Side True)",
			ruleStr: `device_id == "critical-node" OR payload.vibration > 10.0`,
			context: map[string]interface{}{
				"device_id": "critical-node",
				// payload.vibration is completely missing; short-circuit prevents field lookup failure
			},
			wantResult: true,
			expectErr:  false,
		},
		{
			name:    "Graceful Missing Field Resolution",
			ruleStr: `payload.optional_metric == "present"`,
			context: map[string]interface{}{
				"payload": map[string]interface{}{
					// metric is missing entirely
				},
			},
			wantResult: false,
			expectErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := evaluator.NewLexer(tt.ruleStr)
			p := evaluator.NewParser(l)
			
			rootNode, err := p.ParseRule()
			if err != nil {
				if !tt.expectErr {
					t.Fatalf("Unexpected compilation failure: %v", err)
				}
				return
			}

			res, err := rootNode.Evaluate(tt.context)
			if err != nil {
				if !tt.expectErr {
					t.Fatalf("Unexpected execution evaluation failure: %v", err)
				}
				return
			}

			gotResult, ok := res.(bool)
			if !ok {
				t.Fatalf("Expected boolean evaluation result, got: %v", res)
			}

			if gotResult != tt.wantResult {
				t.Errorf("Evaluation mismatch! Rule: (%s) -> Got %t, Want %t", tt.ruleStr, gotResult, tt.wantResult)
			}
		})
	}
}