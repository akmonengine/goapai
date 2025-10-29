package goapai

import "testing"

func TestState_Operations(t *testing.T) {
	tests := []struct {
		name      string
		testFunc  func(*testing.T)
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
				if !state.Check(agent.states, 1) {
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
				if wrongState.Check(agent.states, 1) {
					t.Error("Expected state not to match")
				}
			},
		},
		{
			name: "Check key not found",
			testFunc: func(t *testing.T) {
				agent := CreateAgent(Goals{}, Actions{})
				state := State[int]{Key: 99, Value: 100}

				if state.Check(agent.states, 99) {
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

func TestStates_Check(t *testing.T) {
	tests := []struct {
		name      string
		setup1    func(*Agent)
		setup2    func(*Agent)
		wantMatch bool
	}{
		{
			name: "matching states",
			setup1: func(a *Agent) {
				SetState[int](a, 1, 100)
			},
			setup2: func(a *Agent) {
				SetState[int](a, 1, 100)
			},
			wantMatch: true,
		},
		{
			name: "different states",
			setup1: func(a *Agent) {
				SetState[int](a, 1, 100)
			},
			setup2: func(a *Agent) {
				SetState[int](a, 1, 100)
				SetState[int](a, 2, 200)
			},
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent1 := CreateAgent(Goals{}, Actions{})
			tt.setup1(&agent1)
			agent1.states.data.sort()
			agent1.states.hash = agent1.states.data.hashStates()

			agent2 := CreateAgent(Goals{}, Actions{})
			tt.setup2(&agent2)
			agent2.states.data.sort()
			agent2.states.hash = agent2.states.data.hashStates()

			if got := agent1.states.Check(agent2.states); got != tt.wantMatch {
				t.Errorf("Check() = %v, want %v", got, tt.wantMatch)
			}
		})
	}
}

func TestStatesData_Operations(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "GetIndex found",
			testFunc: func(t *testing.T) {
				data := statesData{
					State[int]{Key: 1, Value: 100},
					State[int]{Key: 2, Value: 200},
					State[bool]{Key: 3, Value: true},
				}

				idx := data.GetIndex(2)
				if idx != 1 {
					t.Errorf("Expected index 1, got %d", idx)
				}
			},
		},
		{
			name: "GetIndex not found",
			testFunc: func(t *testing.T) {
				data := statesData{
					State[int]{Key: 1, Value: 100},
				}

				idx := data.GetIndex(99)
				if idx != -1 {
					t.Errorf("Expected index -1 for missing key, got %d", idx)
				}
			},
		},
		{
			name: "sort",
			testFunc: func(t *testing.T) {
				data := statesData{
					State[int]{Key: 3, Value: 300},
					State[int]{Key: 1, Value: 100},
					State[int]{Key: 2, Value: 200},
				}

				data.sort()

				keys := []StateKey{1, 2, 3}
				for i, expected := range keys {
					if data[i].GetKey() != expected {
						t.Errorf("Expected key %d at position %d, got %d", expected, i, data[i].GetKey())
					}
				}
			},
		},
		{
			name: "hashStates same data",
			testFunc: func(t *testing.T) {
				data1 := statesData{
					State[int]{Key: 1, Value: 100},
					State[int]{Key: 2, Value: 200},
				}

				data2 := statesData{
					State[int]{Key: 1, Value: 100},
					State[int]{Key: 2, Value: 200},
				}

				hash1 := data1.hashStates()
				hash2 := data2.hashStates()

				if hash1 != hash2 {
					t.Error("Expected identical data to produce same hash")
				}
			},
		},
		{
			name: "hashStates different data",
			testFunc: func(t *testing.T) {
				data1 := statesData{
					State[int]{Key: 1, Value: 100},
					State[int]{Key: 2, Value: 200},
				}

				data2 := statesData{
					State[int]{Key: 1, Value: 100},
					State[int]{Key: 2, Value: 999},
				}

				hash1 := data1.hashStates()
				hash2 := data2.hashStates()

				if hash1 == hash2 {
					t.Error("Expected different data to produce different hash")
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
			if got := condition.Check(agent.states); got != tt.wantMatch {
				t.Errorf("Check() = %v, want %v for state=%d, cond=%d, op=%v",
					got, tt.wantMatch, tt.stateVal, tt.condVal, tt.operator)
			}
		})
	}
}

func TestCondition_KeyNotFound(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	condition := Condition[int]{Key: 99, Value: 100, Operator: EQUAL}

	if condition.Check(agent.states) {
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
			if got := condition.Check(agent.states); got != tt.wantMatch {
				t.Errorf("Check() = %v, want %v", got, tt.wantMatch)
			}
		})
	}

	t.Run("key not found", func(t *testing.T) {
		agent := CreateAgent(Goals{}, Actions{})
		condition := ConditionBool{Key: 99, Value: true, Operator: EQUAL}

		if condition.Check(agent.states) {
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
			if got := condition.Check(agent.states); got != tt.wantMatch {
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

			if got := condition.Check(agent.states); got != tt.wantResult {
				t.Errorf("Check() = %v, want %v", got, tt.wantResult)
			}

			// Test caching
			if !condition.resolved {
				t.Error("Expected condition to be marked as resolved")
			}

			// Call again to test cache
			if got := condition.Check(agent.states); got != tt.wantResult {
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

			if got := tt.conditions.Check(agent.states); got != tt.wantMatch {
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
