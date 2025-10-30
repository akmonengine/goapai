package goapai

import "testing"

func TestEffectString_GetKey(t *testing.T) {
	effect := EffectString{Key: 42, Value: "test"}
	if got := effect.GetKey(); got != 42 {
		t.Errorf("GetKey() = %v, want 42", got)
	}
}

func TestEffectString_Check(t *testing.T) {
	tests := []struct {
		name       string
		stateVal   string
		effectKey  StateKey
		effectVal  string
		operator   arithmetic
		want       bool
	}{
		{
			name:       "SET match",
			stateVal:   "hello",
			effectKey:  1,
			effectVal:  "hello",
			operator:   SET,
			want:       true,
		},
		{
			name:       "SET no match",
			stateVal:   "hello",
			effectKey:  1,
			effectVal:  "world",
			operator:   SET,
			want:       false,
		},
		{
			name:       "key not found",
			stateVal:   "hello",
			effectKey:  99,
			effectVal:  "test",
			operator:   SET,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := CreateAgent(Goals{}, Actions{})
			if tt.effectKey == tt.effectKey { // Only set state if key matches
				SetState[string](&agent, 1, tt.stateVal)
			}

			effect := EffectString{Key: tt.effectKey, Value: tt.effectVal, Operator: tt.operator}
			got := effect.check(agent.w)

			if got != tt.want {
				t.Errorf("check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEffectString_Apply(t *testing.T) {
	tests := []struct {
		name       string
		setupState bool
		stateVal   string
		effectKey  StateKey
		effectVal  string
		operator   arithmetic
		wantVal    string
		wantErr    bool
	}{
		{
			name:       "SET new state",
			setupState: false,
			effectKey:  1,
			effectVal:  "hello",
			operator:   SET,
			wantVal:    "hello",
			wantErr:    false,
		},
		{
			name:       "SET existing state",
			setupState: true,
			stateVal:   "old",
			effectKey:  1,
			effectVal:  "new",
			operator:   SET,
			wantVal:    "new",
			wantErr:    false,
		},
		{
			name:       "ADD concatenation",
			setupState: true,
			stateVal:   "hello",
			effectKey:  1,
			effectVal:  " world",
			operator:   ADD,
			wantVal:    "hello world",
			wantErr:    false,
		},
		{
			name:       "ADD new state",
			setupState: false,
			effectKey:  1,
			effectVal:  "test",
			operator:   ADD,
			wantVal:    "test",
			wantErr:    false,
		},
		{
			name:       "SUBSTRACT not allowed",
			setupState: true,
			stateVal:   "hello",
			effectKey:  1,
			effectVal:  "test",
			operator:   SUBSTRACT,
			wantErr:    true,
		},
		{
			name:       "MULTIPLY not allowed",
			setupState: true,
			stateVal:   "hello",
			effectKey:  1,
			effectVal:  "test",
			operator:   MULTIPLY,
			wantErr:    true,
		},
		{
			name:       "DIVIDE not allowed",
			setupState: true,
			stateVal:   "hello",
			effectKey:  1,
			effectVal:  "test",
			operator:   DIVIDE,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := CreateAgent(Goals{}, Actions{})
			if tt.setupState {
				SetState[string](&agent, tt.effectKey, tt.stateVal)
			}

			effect := EffectString{Key: tt.effectKey, Value: tt.effectVal, Operator: tt.operator}
			err := effect.apply(&agent.w)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			idx := agent.w.states.GetIndex(tt.effectKey)
			if idx < 0 {
				t.Fatal("State not found after apply")
			}

			gotState := agent.w.states[idx].(State[string])
			if gotState.Value != tt.wantVal {
				t.Errorf("Value = %v, want %v", gotState.Value, tt.wantVal)
			}
		})
	}
}

func TestEffectBool_GetKey(t *testing.T) {
	effect := EffectBool{Key: 42, Value: true}
	if got := effect.GetKey(); got != 42 {
		t.Errorf("GetKey() = %v, want 42", got)
	}
}

func TestEffectBool_Apply_Errors(t *testing.T) {
	tests := []struct {
		name      string
		operator  arithmetic
		wantErr   bool
	}{
		{
			name:      "SET allowed",
			operator:  SET,
			wantErr:   false,
		},
		{
			name:      "ADD not allowed",
			operator:  ADD,
			wantErr:   true,
		},
		{
			name:      "SUBSTRACT not allowed",
			operator:  SUBSTRACT,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := CreateAgent(Goals{}, Actions{})
			SetState[bool](&agent, 1, true)

			effect := EffectBool{Key: 1, Value: false, Operator: tt.operator}
			err := effect.apply(&agent.w)

			if tt.wantErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestEffect_Apply_AllOperators(t *testing.T) {
	tests := []struct {
		name      string
		initial   int
		operator  arithmetic
		value     int
		wantVal   int
		wantErr   bool
	}{
		{
			name:     "SET",
			initial:  100,
			operator: SET,
			value:    50,
			wantVal:  50,
		},
		{
			name:     "ADD",
			initial:  100,
			operator: ADD,
			value:    50,
			wantVal:  150,
		},
		{
			name:     "SUBSTRACT",
			initial:  100,
			operator: SUBSTRACT,
			value:    30,
			wantVal:  70,
		},
		{
			name:     "MULTIPLY",
			initial:  10,
			operator: MULTIPLY,
			value:    5,
			wantVal:  50,
		},
		{
			name:     "DIVIDE",
			initial:  100,
			operator: DIVIDE,
			value:    4,
			wantVal:  25,
		},
		{
			name:     "ADD on non-existing key",
			initial:  0, // Not set
			operator: ADD,
			value:    50,
			wantVal:  50,
		},
		{
			name:     "SUBSTRACT on non-existing key",
			initial:  0, // Not set
			operator: SUBSTRACT,
			value:    50,
			wantVal:  -50,
		},
		{
			name:     "MULTIPLY on non-existing key error",
			initial:  0, // Not set
			operator: MULTIPLY,
			value:    50,
			wantErr:  true,
		},
		{
			name:     "DIVIDE on non-existing key error",
			initial:  0, // Not set
			operator: DIVIDE,
			value:    50,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := CreateAgent(Goals{}, Actions{})
			if tt.initial != 0 || (tt.operator != ADD && tt.operator != SUBSTRACT && tt.operator != MULTIPLY && tt.operator != DIVIDE) {
				SetState[int](&agent, 1, tt.initial)
			}

			effect := Effect[int]{Key: 1, Value: tt.value, Operator: tt.operator}
			err := effect.apply(&agent.w)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			idx := agent.w.states.GetIndex(1)
			if idx < 0 {
				t.Fatal("State not found after apply")
			}

			gotState := agent.w.states[idx].(State[int])
			if gotState.Value != tt.wantVal {
				t.Errorf("Value = %v, want %v", gotState.Value, tt.wantVal)
			}
		})
	}
}

func TestEffect_GetKey(t *testing.T) {
	effect := Effect[int]{Key: 42, Value: 100}
	if got := effect.GetKey(); got != 42 {
		t.Errorf("GetKey() = %v, want 42", got)
	}
}
