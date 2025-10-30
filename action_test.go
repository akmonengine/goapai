package goapai

import "testing"

func TestActions_AddAction(t *testing.T) {
	tests := []struct {
		name           string
		actionsToAdd   []struct{ name string; cost float32; repeatable bool }
		wantCount      int
		checkFirst     bool
		wantFirstName  string
		wantFirstCost  float32
		wantRepeatble  bool
	}{
		{
			name: "single action",
			actionsToAdd: []struct{ name string; cost float32; repeatable bool }{
				{name: "test", cost: 1.5, repeatable: false},
			},
			wantCount:     1,
			checkFirst:    true,
			wantFirstName: "test",
			wantFirstCost: 1.5,
			wantRepeatble: false,
		},
		{
			name: "multiple actions",
			actionsToAdd: []struct{ name string; cost float32; repeatable bool }{
				{name: "action1", cost: 1.0, repeatable: true},
				{name: "action2", cost: 2.0, repeatable: false},
			},
			wantCount:  2,
			checkFirst: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actions := Actions{}

			for _, a := range tt.actionsToAdd {
				actions.AddAction(a.name, a.cost, a.repeatable, Conditions{}, Effects{})
			}

			if len(actions) != tt.wantCount {
				t.Errorf("Expected %d actions, got %d", tt.wantCount, len(actions))
			}

			if tt.checkFirst {
				action := actions[0]
				if action.name != tt.wantFirstName {
					t.Errorf("Expected name '%s', got '%s'", tt.wantFirstName, action.name)
				}
				if action.cost != tt.wantFirstCost {
					t.Errorf("Expected cost %f, got %f", tt.wantFirstCost, action.cost)
				}
				if action.repeatable != tt.wantRepeatble {
					t.Errorf("Expected repeatable to be %v", tt.wantRepeatble)
				}
			}
		})
	}
}

func TestAction_GetName(t *testing.T) {
	tests := []struct {
		name       string
		actionName string
		want       string
	}{
		{
			name:       "basic name",
			actionName: "my_action",
			want:       "my_action",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actions := Actions{}
			actions.AddAction(tt.actionName, 1.0, false, Conditions{}, Effects{})

			got := actions[0].GetName()
			if got != tt.want {
				t.Errorf("Expected '%s', got '%s'", tt.want, got)
			}
		})
	}
}

func TestAction_GetEffects(t *testing.T) {
	tests := []struct {
		name        string
		effects     Effects
		wantCount   int
	}{
		{
			name: "single effect",
			effects: Effects{
				Effect[int]{Key: 1, Value: 10, Operator: SET},
			},
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actions := Actions{}
			actions.AddAction("test", 1.0, false, Conditions{}, tt.effects)

			retrieved := actions[0].GetEffects()
			if len(retrieved) != tt.wantCount {
				t.Errorf("Expected %d effect(s), got %d", tt.wantCount, len(retrieved))
			}
		})
	}
}

func TestEffect_Check(t *testing.T) {
	tests := []struct {
		name       string
		stateValue int
		effectKey  StateKey
		effectVal  int
		operator   arithmetic
		want       bool
	}{
		{
			name:       "match with SET operator",
			stateValue: 100,
			effectKey:  1,
			effectVal:  100,
			operator:   SET,
			want:       true,
		},
		{
			name:       "no match with SET operator",
			stateValue: 100,
			effectKey:  1,
			effectVal:  200,
			operator:   SET,
			want:       false,
		},
		{
			name:       "non-SET operator always false",
			stateValue: 100,
			effectKey:  1,
			effectVal:  100,
			operator:   ADD,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := CreateAgent(Goals{}, Actions{})
			SetState[int](&agent, tt.effectKey, tt.stateValue)

			effect := Effect[int]{Key: tt.effectKey, Value: tt.effectVal, Operator: tt.operator}
			got := effect.check(agent.w)

			if got != tt.want {
				t.Errorf("Expected %v, got %v", tt.want, got)
			}
		})
	}
}
