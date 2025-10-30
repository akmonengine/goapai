package goapai

import "testing"

func TestState_Distance(t *testing.T) {
	tests := []struct {
		name      string
		stateVal  interface{}
		condition ConditionInterface
		want      float32
	}{
		// Numeric types - EQUAL operator
		{
			name:      "int EQUAL satisfied",
			stateVal:  100,
			condition: &Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
			want:      0.0,
		},
		{
			name:      "int EQUAL below target",
			stateVal:  50,
			condition: &Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
			want:      50.0,
		},
		{
			name:      "int EQUAL above target",
			stateVal:  150,
			condition: &Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
			want:      50.0,
		},
		// UPPER_OR_EQUAL operator
		{
			name:      "int UPPER_OR_EQUAL satisfied",
			stateVal:  100,
			condition: &Condition[int]{Key: 1, Value: 80, Operator: UPPER_OR_EQUAL},
			want:      0.0,
		},
		{
			name:      "int UPPER_OR_EQUAL below target",
			stateVal:  50,
			condition: &Condition[int]{Key: 1, Value: 80, Operator: UPPER_OR_EQUAL},
			want:      30.0,
		},
		// UPPER operator
		{
			name:      "int UPPER satisfied",
			stateVal:  100,
			condition: &Condition[int]{Key: 1, Value: 80, Operator: UPPER},
			want:      0.0,
		},
		{
			name:      "int UPPER at target",
			stateVal:  80,
			condition: &Condition[int]{Key: 1, Value: 80, Operator: UPPER},
			want:      1.0,
		},
		// LOWER_OR_EQUAL operator
		{
			name:      "int LOWER_OR_EQUAL satisfied",
			stateVal:  50,
			condition: &Condition[int]{Key: 1, Value: 80, Operator: LOWER_OR_EQUAL},
			want:      0.0,
		},
		{
			name:      "int LOWER_OR_EQUAL above target",
			stateVal:  100,
			condition: &Condition[int]{Key: 1, Value: 80, Operator: LOWER_OR_EQUAL},
			want:      20.0,
		},
		// LOWER operator
		{
			name:      "int LOWER satisfied",
			stateVal:  50,
			condition: &Condition[int]{Key: 1, Value: 80, Operator: LOWER},
			want:      0.0,
		},
		{
			name:      "int LOWER at target",
			stateVal:  80,
			condition: &Condition[int]{Key: 1, Value: 80, Operator: LOWER},
			want:      1.0,
		},
		// NOT_EQUAL operator
		{
			name:      "int NOT_EQUAL satisfied",
			stateVal:  50,
			condition: &Condition[int]{Key: 1, Value: 100, Operator: NOT_EQUAL},
			want:      0.0,
		},
		{
			name:      "int NOT_EQUAL not satisfied",
			stateVal:  100,
			condition: &Condition[int]{Key: 1, Value: 100, Operator: NOT_EQUAL},
			want:      1.0,
		},
		// Float64 tests
		{
			name:      "float64 EQUAL",
			stateVal:  75.5,
			condition: &Condition[float64]{Key: 1, Value: 100.0, Operator: EQUAL},
			want:      24.5,
		},
		// uint64 tests
		{
			name:      "uint64 UPPER_OR_EQUAL",
			stateVal:  uint64(50),
			condition: &Condition[uint64]{Key: 1, Value: uint64(100), Operator: UPPER_OR_EQUAL},
			want:      50.0,
		},
		// int8 tests
		{
			name:      "int8 EQUAL",
			stateVal:  int8(50),
			condition: &Condition[int8]{Key: 1, Value: int8(100), Operator: EQUAL},
			want:      50.0,
		},
		// uint8 tests
		{
			name:      "uint8 LOWER_OR_EQUAL",
			stateVal:  uint8(100),
			condition: &Condition[uint8]{Key: 1, Value: uint8(80), Operator: LOWER_OR_EQUAL},
			want:      20.0,
		},
		// Bool tests
		{
			name:      "bool EQUAL satisfied",
			stateVal:  true,
			condition: &ConditionBool{Key: 1, Value: true, Operator: EQUAL},
			want:      0.0,
		},
		{
			name:      "bool EQUAL not satisfied",
			stateVal:  true,
			condition: &ConditionBool{Key: 1, Value: false, Operator: EQUAL},
			want:      1.0,
		},
		{
			name:      "bool NOT_EQUAL satisfied",
			stateVal:  true,
			condition: &ConditionBool{Key: 1, Value: false, Operator: NOT_EQUAL},
			want:      0.0,
		},
		{
			name:      "bool NOT_EQUAL not satisfied",
			stateVal:  true,
			condition: &ConditionBool{Key: 1, Value: true, Operator: NOT_EQUAL},
			want:      1.0,
		},
		// String tests
		{
			name:      "string EQUAL satisfied",
			stateVal:  "test",
			condition: &ConditionString{Key: 1, Value: "test", Operator: EQUAL},
			want:      0.0,
		},
		{
			name:      "string EQUAL not satisfied",
			stateVal:  "test",
			condition: &ConditionString{Key: 1, Value: "other", Operator: EQUAL},
			want:      1.0,
		},
		{
			name:      "string NOT_EQUAL satisfied",
			stateVal:  "test",
			condition: &ConditionString{Key: 1, Value: "other", Operator: NOT_EQUAL},
			want:      0.0,
		},
		{
			name:      "string NOT_EQUAL not satisfied",
			stateVal:  "test",
			condition: &ConditionString{Key: 1, Value: "test", Operator: NOT_EQUAL},
			want:      1.0,
		},
		// Key mismatch
		{
			name:      "key mismatch",
			stateVal:  100,
			condition: &Condition[int]{Key: 99, Value: 100, Operator: EQUAL},
			want:      0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := CreateAgent(Goals{}, Actions{})

			switch v := tt.stateVal.(type) {
			case int:
				SetState[int](&agent, 1, v)
			case int8:
				SetState[int8](&agent, 1, v)
			case uint8:
				SetState[uint8](&agent, 1, v)
			case uint64:
				SetState[uint64](&agent, 1, v)
			case float64:
				SetState[float64](&agent, 1, v)
			case bool:
				SetState[bool](&agent, 1, v)
			case string:
				SetState[string](&agent, 1, v)
			}

			if len(agent.w.states) == 0 {
				t.Fatal("Failed to set state")
			}

			state := agent.w.states[0]
			got := state.Distance(tt.condition)

			if got != tt.want {
				t.Errorf("Distance() = %v, want %v", got, tt.want)
			}
		})
	}
}
