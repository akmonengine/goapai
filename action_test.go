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
	if !effect.check(agent.w) {
		t.Error("Expected effect to match state")
	}
}

func TestEffect_Check_NoMatch(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 100)

	effect := Effect[int]{Key: 1, Value: 200, Operator: SET}
	if effect.check(agent.w) {
		t.Error("Expected effect not to match state")
	}
}

func TestEffect_Check_NonSetOperator(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 100)

	effect := Effect[int]{Key: 1, Value: 100, Operator: ADD}
	if effect.check(agent.w) {
		t.Error("Expected non-SET operator to always return false")
	}
}
