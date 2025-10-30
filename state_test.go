package goapai

import "testing"

func TestState_Operations(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "GetKey",
			testFunc: func(t *testing.T) {
				state := State[int]{Key: 42, Value: 100}
				if state.GetKey() != 42 {
					t.Errorf("Expected key 42, got %d", state.GetKey())
				}
			},
		},
		{
			name: "GetValue",
			testFunc: func(t *testing.T) {
				state := State[int]{Key: 1, Value: 100}
				if state.GetValue().(int) != 100 {
					t.Errorf("Expected value 100, got %v", state.GetValue())
				}
			},
		},
		{
			name: "Check match",
			testFunc: func(t *testing.T) {
				agent := CreateAgent(Goals{}, Actions{})
				SetState[int](&agent, 1, 100)

				state := State[int]{Key: 1, Value: 100}
				if !state.Check(agent.w, 1) {
					t.Error("Expected state to match")
				}
			},
		},
		{
			name: "Check no match",
			testFunc: func(t *testing.T) {
				agent := CreateAgent(Goals{}, Actions{})
				SetState[int](&agent, 1, 100)

				wrongState := State[int]{Key: 1, Value: 200}
				if wrongState.Check(agent.w, 1) {
					t.Error("Expected state not to match")
				}
			},
		},
		{
			name: "Check key not found",
			testFunc: func(t *testing.T) {
				agent := CreateAgent(Goals{}, Actions{})
				state := State[int]{Key: 99, Value: 100}

				if state.Check(agent.w, 99) {
					t.Error("Expected false for non-existent key")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func TestCondition_Operators(t *testing.T) {
	tests := []struct {
		name      string
		stateVal  int
		condVal   int
		operator  operator
		wantMatch bool
	}{
		{"equal match", 100, 100, EQUAL, true},
		{"equal no match", 100, 99, EQUAL, false},
		{"not equal match", 100, 99, NOT_EQUAL, true},
		{"not equal no match", 100, 100, NOT_EQUAL, false},
		{"lower true", 50, 100, LOWER, true},
		{"lower false equal", 50, 50, LOWER, false},
		{"lower false greater", 100, 50, LOWER, false},
		{"lower or equal true less", 50, 100, LOWER_OR_EQUAL, true},
		{"lower or equal true equal", 50, 50, LOWER_OR_EQUAL, true},
		{"lower or equal false", 100, 50, LOWER_OR_EQUAL, false},
		{"upper true", 100, 50, UPPER, true},
		{"upper false equal", 100, 100, UPPER, false},
		{"upper false less", 50, 100, UPPER, false},
		{"upper or equal true greater", 100, 50, UPPER_OR_EQUAL, true},
		{"upper or equal true equal", 100, 100, UPPER_OR_EQUAL, true},
		{"upper or equal false", 50, 100, UPPER_OR_EQUAL, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := CreateAgent(Goals{}, Actions{})
			SetState[int](&agent, 1, tt.stateVal)

			condition := Condition[int]{Key: 1, Value: tt.condVal, Operator: tt.operator}
			if got := condition.Check(agent.w); got != tt.wantMatch {
				t.Errorf("Check() = %v, want %v for state=%d, cond=%d, op=%v",
					got, tt.wantMatch, tt.stateVal, tt.condVal, tt.operator)
			}
		})
	}
}

func TestCondition_KeyNotFound(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	condition := Condition[int]{Key: 99, Value: 100, Operator: EQUAL}

	if condition.Check(agent.w) {
		t.Error("Expected condition to fail when key not found")
	}
}

func TestConditionBool(t *testing.T) {
	tests := []struct {
		name      string
		stateVal  bool
		condVal   bool
		operator  operator
		wantMatch bool
	}{
		{"equal true match", true, true, EQUAL, true},
		{"equal false match", false, false, EQUAL, true},
		{"equal no match", true, false, EQUAL, false},
		{"not equal match", true, false, NOT_EQUAL, true},
		{"not equal no match", true, true, NOT_EQUAL, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := CreateAgent(Goals{}, Actions{})
			SetState[bool](&agent, 1, tt.stateVal)

			condition := ConditionBool{Key: 1, Value: tt.condVal, Operator: tt.operator}
			if got := condition.Check(agent.w); got != tt.wantMatch {
				t.Errorf("Check() = %v, want %v", got, tt.wantMatch)
			}
		})
	}

	t.Run("key not found", func(t *testing.T) {
		agent := CreateAgent(Goals{}, Actions{})
		condition := ConditionBool{Key: 99, Value: true, Operator: EQUAL}

		if condition.Check(agent.w) {
			t.Error("Expected condition to fail when key not found")
		}
	})
}

func TestConditionString(t *testing.T) {
	tests := []struct {
		name      string
		stateVal  string
		condVal   string
		operator  operator
		wantMatch bool
	}{
		{"equal match", "test", "test", EQUAL, true},
		{"equal no match", "test", "other", EQUAL, false},
		{"not equal match", "test", "other", NOT_EQUAL, true},
		{"not equal no match", "test", "test", NOT_EQUAL, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := CreateAgent(Goals{}, Actions{})
			SetState[string](&agent, 1, tt.stateVal)

			condition := ConditionString{Key: 1, Value: tt.condVal, Operator: tt.operator}
			if got := condition.Check(agent.w); got != tt.wantMatch {
				t.Errorf("Check() = %v, want %v", got, tt.wantMatch)
			}
		})
	}
}

func TestConditionFn(t *testing.T) {
	tests := []struct {
		name       string
		sensorVal  int
		threshold  int
		wantResult bool
	}{
		{"above threshold", 100, 50, true},
		{"below threshold", 30, 50, false},
		{"at threshold", 50, 50, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := CreateAgent(Goals{}, Actions{})
			SetSensor(&agent, "value", tt.sensorVal)

			condition := &ConditionFn{
				Key: 1,
				CheckFn: func(sensors Sensors) bool {
					return sensors.GetSensor("value").(int) > tt.threshold
				},
			}

			if got := condition.Check(agent.w); got != tt.wantResult {
				t.Errorf("Check() = %v, want %v", got, tt.wantResult)
			}

			// Test caching
			if !condition.resolved {
				t.Error("Expected condition to be marked as resolved")
			}

			// Call again to test cache
			if got := condition.Check(agent.w); got != tt.wantResult {
				t.Error("Expected cached result to match")
			}
		})
	}
}

func TestConditions_Check(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*Agent)
		conditions Conditions
		wantMatch  bool
	}{
		{
			name: "all match",
			setup: func(a *Agent) {
				SetState[int](a, 1, 100)
				SetState[bool](a, 2, true)
			},
			conditions: Conditions{
				&Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
				&ConditionBool{Key: 2, Value: true, Operator: EQUAL},
			},
			wantMatch: true,
		},
		{
			name: "one fails",
			setup: func(a *Agent) {
				SetState[int](a, 1, 100)
				SetState[bool](a, 2, true)
			},
			conditions: Conditions{
				&Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
				&ConditionBool{Key: 2, Value: false, Operator: EQUAL},
			},
			wantMatch: false,
		},
		{
			name:       "empty conditions",
			setup:      func(a *Agent) {},
			conditions: Conditions{},
			wantMatch:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := CreateAgent(Goals{}, Actions{})
			tt.setup(&agent)

			if got := tt.conditions.Check(agent.w); got != tt.wantMatch {
				t.Errorf("Check() = %v, want %v", got, tt.wantMatch)
			}
		})
	}
}

func TestSensors_GetSensor(t *testing.T) {
	sensors := Sensors{
		"test": 42,
	}

	value := sensors.GetSensor("test")
	if value.(int) != 42 {
		t.Errorf("Expected 42, got %v", value)
	}
}

func TestState_Hash(t *testing.T) {
	tests := []struct {
		name     string
		state1   StateInterface
		state2   StateInterface
		wantSame bool
	}{
		{
			name:     "same int world",
			state1:   State[int]{Key: 1, Value: 100},
			state2:   State[int]{Key: 1, Value: 100},
			wantSame: true,
		},
		{
			name:     "different int values",
			state1:   State[int]{Key: 1, Value: 100},
			state2:   State[int]{Key: 1, Value: 200},
			wantSame: false,
		},
		{
			name:     "different keys",
			state1:   State[int]{Key: 1, Value: 100},
			state2:   State[int]{Key: 2, Value: 100},
			wantSame: false,
		},
		{
			name:     "same bool world",
			state1:   State[bool]{Key: 1, Value: true},
			state2:   State[bool]{Key: 1, Value: true},
			wantSame: true,
		},
		{
			name:     "different bool values",
			state1:   State[bool]{Key: 1, Value: true},
			state2:   State[bool]{Key: 1, Value: false},
			wantSame: false,
		},
		{
			name:     "same string world",
			state1:   State[string]{Key: 1, Value: "test"},
			state2:   State[string]{Key: 1, Value: "test"},
			wantSame: true,
		},
		{
			name:     "different string values",
			state1:   State[string]{Key: 1, Value: "test"},
			state2:   State[string]{Key: 1, Value: "other"},
			wantSame: false,
		},
		{
			name:     "same float64 world",
			state1:   State[float64]{Key: 1, Value: 3.14},
			state2:   State[float64]{Key: 1, Value: 3.14},
			wantSame: true,
		},
		{
			name:     "different float64 values",
			state1:   State[float64]{Key: 1, Value: 3.14},
			state2:   State[float64]{Key: 1, Value: 2.71},
			wantSame: false,
		},
		{
			name:     "same uint64 world",
			state1:   State[uint64]{Key: 1, Value: 12345},
			state2:   State[uint64]{Key: 1, Value: 12345},
			wantSame: true,
		},
		{
			name:     "different uint64 values",
			state1:   State[uint64]{Key: 1, Value: 12345},
			state2:   State[uint64]{Key: 1, Value: 54321},
			wantSame: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash1 := tt.state1.Hash()
			hash2 := tt.state2.Hash()

			if tt.wantSame && hash1 != hash2 {
				t.Errorf("Expected same hash, got %d and %d", hash1, hash2)
			}
			if !tt.wantSame && hash1 == hash2 {
				t.Errorf("Expected different hashes, both got %d", hash1)
			}
		})
	}
}
