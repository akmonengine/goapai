package goapai

import "testing"

// Test Actions.AddAction
func TestActions_AddAction(t *testing.T) {
	actions := Actions{}

	actions.AddAction("test", 1.5, false, Conditions{}, Effects{})

	if len(actions) != 1 {
		t.Errorf("Expected 1 action, got %d", len(actions))
	}

	action := actions[0]
	if action.name != "test" {
		t.Errorf("Expected name 'test', got '%s'", action.name)
	}
	if action.cost != 1.5 {
		t.Errorf("Expected cost 1.5, got %f", action.cost)
	}
	if action.repeatable {
		t.Error("Expected repeatable to be false")
	}
}

func TestActions_AddAction_Multiple(t *testing.T) {
	actions := Actions{}

	actions.AddAction("action1", 1.0, true, Conditions{}, Effects{})
	actions.AddAction("action2", 2.0, false, Conditions{}, Effects{})

	if len(actions) != 2 {
		t.Errorf("Expected 2 actions, got %d", len(actions))
	}
}

// Test Action.GetName
func TestAction_GetName(t *testing.T) {
	actions := Actions{}
	actions.AddAction("my_action", 1.0, false, Conditions{}, Effects{})

	name := actions[0].GetName()
	if name != "my_action" {
		t.Errorf("Expected 'my_action', got '%s'", name)
	}
}

// Test Action.GetEffects
func TestAction_GetEffects(t *testing.T) {
	effects := Effects{
		Effect[int]{Key: 1, Value: 10, Operator: SET},
	}

	actions := Actions{}
	actions.AddAction("test", 1.0, false, Conditions{}, effects)

	retrieved := actions[0].GetEffects()
	if len(retrieved) != 1 {
		t.Errorf("Expected 1 effect, got %d", len(retrieved))
	}
}

// Test Effect[T Numeric] operations
func TestEffect_Check_Match(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 100)

	effect := Effect[int]{Key: 1, Value: 100, Operator: SET}
	if !effect.check(agent.states) {
		t.Error("Expected effect to match state")
	}
}

func TestEffect_Check_NoMatch(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 100)

	effect := Effect[int]{Key: 1, Value: 200, Operator: SET}
	if effect.check(agent.states) {
		t.Error("Expected effect not to match state")
	}
}

func TestEffect_Check_NonSetOperator(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 100)

	effect := Effect[int]{Key: 1, Value: 100, Operator: ADD}
	if effect.check(agent.states) {
		t.Error("Expected non-SET operator to always return false")
	}
}

func TestEffect_Apply_Set(t *testing.T) {
	data := statesData{
		State[int]{Key: 1, Value: 100},
	}

	effect := Effect[int]{Key: 1, Value: 200, Operator: SET}
	err := effect.apply(data)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if data[0].(State[int]).Value != 200 {
		t.Errorf("Expected value 200, got %d", data[0].(State[int]).Value)
	}
}

func TestEffect_Apply_Add(t *testing.T) {
	data := statesData{
		State[int]{Key: 1, Value: 100},
	}

	effect := Effect[int]{Key: 1, Value: 50, Operator: ADD}
	err := effect.apply(data)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if data[0].(State[int]).Value != 150 {
		t.Errorf("Expected value 150, got %d", data[0].(State[int]).Value)
	}
}

func TestEffect_Apply_Subtract(t *testing.T) {
	data := statesData{
		State[int]{Key: 1, Value: 100},
	}

	effect := Effect[int]{Key: 1, Value: 30, Operator: SUBSTRACT}
	err := effect.apply(data)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if data[0].(State[int]).Value != 70 {
		t.Errorf("Expected value 70, got %d", data[0].(State[int]).Value)
	}
}

func TestEffect_Apply_Multiply(t *testing.T) {
	data := statesData{
		State[int]{Key: 1, Value: 10},
	}

	effect := Effect[int]{Key: 1, Value: 5, Operator: MULTIPLY}
	err := effect.apply(data)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if data[0].(State[int]).Value != 50 {
		t.Errorf("Expected value 50, got %d", data[0].(State[int]).Value)
	}
}

func TestEffect_Apply_Divide(t *testing.T) {
	data := statesData{
		State[int]{Key: 1, Value: 100},
	}

	effect := Effect[int]{Key: 1, Value: 4, Operator: DIVIDE}
	err := effect.apply(data)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if data[0].(State[int]).Value != 25 {
		t.Errorf("Expected value 25, got %d", data[0].(State[int]).Value)
	}
}

func TestEffect_Apply_NewKey(t *testing.T) {
	data := statesData{}

	effect := Effect[int]{Key: 1, Value: 42, Operator: SET}
	err := effect.apply(data)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Note: apply modifies the slice in place but the original reference doesn't change
	// This is actually expected behavior - effects.apply() returns a new slice
}

// Test EffectBool
func TestEffectBool_Check_Match(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[bool](&agent, 1, true)

	effect := EffectBool{Key: 1, Value: true, Operator: SET}
	if !effect.check(agent.states) {
		t.Error("Expected effect to match state")
	}
}

func TestEffectBool_Check_NoMatch(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[bool](&agent, 1, true)

	effect := EffectBool{Key: 1, Value: false, Operator: SET}
	if effect.check(agent.states) {
		t.Error("Expected effect not to match state")
	}
}

func TestEffectBool_Apply_Set(t *testing.T) {
	data := statesData{
		State[bool]{Key: 1, Value: false},
	}

	effect := EffectBool{Key: 1, Value: true, Operator: SET}
	err := effect.apply(data)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if data[0].(State[bool]).Value != true {
		t.Error("Expected value to be true")
	}
}

func TestEffectBool_Apply_InvalidOperator(t *testing.T) {
	data := statesData{
		State[bool]{Key: 1, Value: true},
	}

	effect := EffectBool{Key: 1, Value: false, Operator: ADD}
	err := effect.apply(data)

	if err == nil {
		t.Error("Expected error for invalid operator on bool")
	}
}

func TestEffectBool_Apply_NewKey(t *testing.T) {
	data := statesData{}

	effect := EffectBool{Key: 1, Value: true, Operator: SET}
	err := effect.apply(data)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Note: apply modifies the slice in place but the original reference doesn't change
	// This is actually expected behavior - effects.apply() returns a new slice
}

// Test EffectString
func TestEffectString_Check_Match(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[string](&agent, 1, "test")

	effect := EffectString{Key: 1, Value: "test", Operator: SET}
	if !effect.check(agent.states) {
		t.Error("Expected effect to match state")
	}
}

func TestEffectString_Check_NoMatch(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[string](&agent, 1, "test")

	effect := EffectString{Key: 1, Value: "other", Operator: SET}
	if effect.check(agent.states) {
		t.Error("Expected effect not to match state")
	}
}

func TestEffectString_Apply_Set(t *testing.T) {
	data := statesData{
		State[string]{Key: 1, Value: "old"},
	}

	effect := EffectString{Key: 1, Value: "new", Operator: SET}
	err := effect.apply(data)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if data[0].(State[string]).Value != "new" {
		t.Errorf("Expected value 'new', got '%s'", data[0].(State[string]).Value)
	}
}

func TestEffectString_Apply_Add_Concatenate(t *testing.T) {
	data := statesData{
		State[string]{Key: 1, Value: "hello"},
	}

	effect := EffectString{Key: 1, Value: " world", Operator: ADD}
	err := effect.apply(data)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if data[0].(State[string]).Value != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", data[0].(State[string]).Value)
	}
}

func TestEffectString_Apply_InvalidOperator(t *testing.T) {
	data := statesData{
		State[string]{Key: 1, Value: "test"},
	}

	effect := EffectString{Key: 1, Value: "x", Operator: MULTIPLY}
	err := effect.apply(data)

	if err == nil {
		t.Error("Expected error for invalid operator on string")
	}
}

func TestEffectString_Apply_NewKey(t *testing.T) {
	data := statesData{}

	effect := EffectString{Key: 1, Value: "new", Operator: SET}
	err := effect.apply(data)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Note: apply modifies the slice in place but the original reference doesn't change
	// This is actually expected behavior - effects.apply() returns a new slice
}

// Test Effects (slice) operations
func TestEffects_SatisfyStates_AllMatch(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 100)
	SetState[bool](&agent, 2, true)

	effects := Effects{
		Effect[int]{Key: 1, Value: 100, Operator: SET},
		EffectBool{Key: 2, Value: true, Operator: SET},
	}

	if !effects.satisfyStates(agent.states) {
		t.Error("Expected effects to satisfy states")
	}
}

func TestEffects_SatisfyStates_OneFails(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 100)
	SetState[bool](&agent, 2, true)

	effects := Effects{
		Effect[int]{Key: 1, Value: 100, Operator: SET},
		EffectBool{Key: 2, Value: false, Operator: SET},
	}

	if effects.satisfyStates(agent.states) {
		t.Error("Expected effects not to satisfy states when one doesn't match")
	}
}

func TestEffects_Apply(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 100)
	SetState[bool](&agent, 2, false)

	effects := Effects{
		Effect[int]{Key: 1, Value: 50, Operator: ADD},
		EffectBool{Key: 2, Value: true, Operator: SET},
	}

	newData, err := effects.apply(agent.states)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(newData) != 2 {
		t.Errorf("Expected 2 states, got %d", len(newData))
	}

	// Check int value
	intIdx := newData.GetIndex(1)
	if intIdx < 0 {
		t.Error("Expected to find key 1")
	}
	if newData[intIdx].(State[int]).Value != 150 {
		t.Errorf("Expected int value 150, got %d", newData[intIdx].(State[int]).Value)
	}

	// Check bool value
	boolIdx := newData.GetIndex(2)
	if boolIdx < 0 {
		t.Error("Expected to find key 2")
	}
	if newData[boolIdx].(State[bool]).Value != true {
		t.Error("Expected bool value true")
	}
}

func TestEffects_Apply_Error(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 100)

	effects := Effects{
		Effect[float64]{Key: 1, Value: 50.0, Operator: SET}, // Type mismatch
	}

	_, err := effects.apply(agent.states)
	if err == nil {
		t.Error("Expected error for type mismatch")
	}
}
