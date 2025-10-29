package goapai

import "testing"

func TestCreateAgent(t *testing.T) {
	goals := Goals{
		"test_goal": {
			Conditions: Conditions{
				&ConditionBool{Key: 1, Value: true},
			},
			PriorityFn: func(sensors Sensors) float32 {
				return 1.0
			},
		},
	}

	actions := Actions{}
	actions.AddAction("test_action", 1.0, false, Conditions{}, Effects{})

	agent := CreateAgent(goals, actions)

	if len(agent.actions) != 1 {
		t.Errorf("Expected 1 action, got %d", len(agent.actions))
	}

	if len(agent.goals) != 1 {
		t.Errorf("Expected 1 goal, got %d", len(agent.goals))
	}

	if agent.sensors == nil {
		t.Error("Expected sensors to be initialized")
	}

	if agent.states.Agent == nil {
		t.Error("Expected states.Agent to be non-nil")
	}
}

func TestSetState(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(*Agent)
		checkFunc func(*testing.T, Agent)
	}{
		{
			name: "int state",
			setupFunc: func(a *Agent) {
				SetState[int](a, 1, 42)
			},
			checkFunc: func(t *testing.T, a Agent) {
				if len(a.states.data) != 1 {
					t.Errorf("Expected 1 state, got %d", len(a.states.data))
				}
				state := a.states.data[0].(State[int])
				if state.Key != 1 || state.Value != 42 {
					t.Errorf("Expected key=1, value=42, got key=%d, value=%d", state.Key, state.Value)
				}
			},
		},
		{
			name: "bool state",
			setupFunc: func(a *Agent) {
				SetState[bool](a, 2, true)
			},
			checkFunc: func(t *testing.T, a Agent) {
				if len(a.states.data) != 1 {
					t.Errorf("Expected 1 state, got %d", len(a.states.data))
				}
				state := a.states.data[0].(State[bool])
				if state.Key != 2 || state.Value != true {
					t.Errorf("Expected key=2, value=true, got key=%d, value=%v", state.Key, state.Value)
				}
			},
		},
		{
			name: "string state",
			setupFunc: func(a *Agent) {
				SetState[string](a, 3, "test")
			},
			checkFunc: func(t *testing.T, a Agent) {
				if len(a.states.data) != 1 {
					t.Errorf("Expected 1 state, got %d", len(a.states.data))
				}
				state := a.states.data[0].(State[string])
				if state.Key != 3 || state.Value != "test" {
					t.Errorf("Expected key=3, value='test', got key=%d, value='%s'", state.Key, state.Value)
				}
			},
		},
		{
			name: "multiple states",
			setupFunc: func(a *Agent) {
				SetState[int](a, 1, 100)
				SetState[bool](a, 2, false)
				SetState[string](a, 3, "hello")
			},
			checkFunc: func(t *testing.T, a Agent) {
				if len(a.states.data) != 3 {
					t.Errorf("Expected 3 states, got %d", len(a.states.data))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := CreateAgent(Goals{}, Actions{})
			tt.setupFunc(&agent)
			tt.checkFunc(t, agent)
		})
	}
}

func TestSetSensor(t *testing.T) {
	type TestEntity struct {
		health int
	}

	tests := []struct {
		name      string
		setupFunc func(*Agent)
		checkFunc func(*testing.T, Agent)
	}{
		{
			name: "single sensor",
			setupFunc: func(a *Agent) {
				entity := &TestEntity{health: 100}
				SetSensor(a, "entity", entity)
			},
			checkFunc: func(t *testing.T, a Agent) {
				if len(a.sensors) != 1 {
					t.Errorf("Expected 1 sensor, got %d", len(a.sensors))
				}
				retrieved := a.sensors.GetSensor("entity").(*TestEntity)
				if retrieved.health != 100 {
					t.Errorf("Expected health 100, got %d", retrieved.health)
				}
			},
		},
		{
			name: "multiple sensors",
			setupFunc: func(a *Agent) {
				SetSensor(a, "sensor1", "value1")
				SetSensor(a, "sensor2", 42)
				SetSensor(a, "sensor3", true)
			},
			checkFunc: func(t *testing.T, a Agent) {
				if len(a.sensors) != 3 {
					t.Errorf("Expected 3 sensors, got %d", len(a.sensors))
				}
				if a.sensors.GetSensor("sensor1").(string) != "value1" {
					t.Error("sensor1 value mismatch")
				}
				if a.sensors.GetSensor("sensor2").(int) != 42 {
					t.Error("sensor2 value mismatch")
				}
				if a.sensors.GetSensor("sensor3").(bool) != true {
					t.Error("sensor3 value mismatch")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := CreateAgent(Goals{}, Actions{})
			tt.setupFunc(&agent)
			tt.checkFunc(t, agent)
		})
	}
}
