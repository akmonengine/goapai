package goapai

import "testing"

func TestPlan_GetTotalCost(t *testing.T) {
	tests := []struct {
		name string
		plan Plan
		want float32
	}{
		{"plan 1", Plan{}, 0.0},
		{"plan 2", Plan{{
			name: "action 1",
			cost: 1.0,
		}}, 1.0},
		{"plan 3", Plan{{
			name: "action 1",
			cost: 1.0,
		}, {
			name: "action 2",
			cost: 2.0,
		}}, 3.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.plan.GetTotalCost(); got != tt.want {
				t.Errorf("GetTotalCost() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test GetPlan
func TestGetPlan_SimpleGoal(t *testing.T) {
	actions := Actions{}
	actions.AddAction("increment", 1.0, true, Conditions{}, Effects{
		Effect[int]{Key: 1, Value: 10, Operator: ADD},
	})

	goals := Goals{
		"reach_30": {
			Conditions: Conditions{
				&Condition[int]{Key: 1, Value: 30, Operator: EQUAL},
			},
			PriorityFn: func(sensors Sensors) float32 {
				return 1.0
			},
		},
	}

	agent := CreateAgent(goals, actions)
	SetState[int](&agent, 1, 0)

	goalName, plan := GetPlan(agent, 10)

	if goalName != "reach_30" {
		t.Errorf("Expected goal 'reach_30', got '%s'", goalName)
	}

	// Plan includes root node
	if len(plan) != 4 {
		t.Errorf("Expected plan with 4 actions (including root), got %d", len(plan))
	}
}

func TestGetPlan_AlreadyAtGoal(t *testing.T) {
	actions := Actions{}

	goals := Goals{
		"be_at_100": {
			Conditions: Conditions{
				&Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
			},
			PriorityFn: func(sensors Sensors) float32 {
				return 1.0
			},
		},
	}

	agent := CreateAgent(goals, actions)
	SetState[int](&agent, 1, 100)

	goalName, plan := GetPlan(agent, 10)

	if goalName != "be_at_100" {
		t.Errorf("Expected goal 'be_at_100', got '%s'", goalName)
	}

	// Plan includes root node even when already at goal
	if len(plan) != 1 {
		t.Errorf("Expected plan with 1 action (root) when already at goal, got %d actions", len(plan))
	}
}

func TestGetPlan_NoGoalsAvailable(t *testing.T) {
	actions := Actions{}

	goals := Goals{
		"impossible": {
			Conditions: Conditions{
				&Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
			},
			PriorityFn: func(sensors Sensors) float32 {
				return 0.0 // Zero priority means goal is not active
			},
		},
	}

	agent := CreateAgent(goals, actions)
	SetState[int](&agent, 1, 0)

	goalName, plan := GetPlan(agent, 10)

	if goalName != "" {
		t.Errorf("Expected empty goal name, got '%s'", goalName)
	}

	if len(plan) != 0 {
		t.Errorf("Expected empty plan when no goals available, got %d actions", len(plan))
	}
}

func TestGetPlan_UnreachableGoal(t *testing.T) {
	actions := Actions{}
	// No actions to reach the goal

	goals := Goals{
		"unreachable": {
			Conditions: Conditions{
				&Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
			},
			PriorityFn: func(sensors Sensors) float32 {
				return 1.0
			},
		},
	}

	agent := CreateAgent(goals, actions)
	SetState[int](&agent, 1, 0)

	goalName, plan := GetPlan(agent, 10)

	if goalName != "unreachable" {
		t.Errorf("Expected goal 'unreachable', got '%s'", goalName)
	}

	if len(plan) != 0 {
		t.Errorf("Expected empty plan for unreachable goal, got %d actions", len(plan))
	}
}

func TestGetPlan_MaxDepthExceeded(t *testing.T) {
	actions := Actions{}
	actions.AddAction("increment", 1.0, true, Conditions{}, Effects{
		Effect[int]{Key: 1, Value: 1, Operator: ADD},
	})

	goals := Goals{
		"reach_100": {
			Conditions: Conditions{
				&Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
			},
			PriorityFn: func(sensors Sensors) float32 {
				return 1.0
			},
		},
	}

	agent := CreateAgent(goals, actions)
	SetState[int](&agent, 1, 0)

	// Max depth of 5 should prevent reaching 100
	_, plan := GetPlan(agent, 5)

	if len(plan) > 5 {
		t.Errorf("Expected plan to respect maxDepth of 5, got %d actions", len(plan))
	}
}

// Test getPrioritizedGoalName
func TestGetPrioritizedGoalName_SingleGoal(t *testing.T) {
	goals := Goals{
		"goal1": {
			Conditions: Conditions{},
			PriorityFn: func(sensors Sensors) float32 {
				return 1.0
			},
		},
	}

	agent := CreateAgent(goals, Actions{})

	goalName, err := agent.getPrioritizedGoalName()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if goalName != "goal1" {
		t.Errorf("Expected 'goal1', got '%s'", goalName)
	}
}

func TestGetPrioritizedGoalName_MultipleGoals(t *testing.T) {
	goals := Goals{
		"low_priority": {
			Conditions: Conditions{},
			PriorityFn: func(sensors Sensors) float32 {
				return 0.5
			},
		},
		"high_priority": {
			Conditions: Conditions{},
			PriorityFn: func(sensors Sensors) float32 {
				return 2.0
			},
		},
		"medium_priority": {
			Conditions: Conditions{},
			PriorityFn: func(sensors Sensors) float32 {
				return 1.0
			},
		},
	}

	agent := CreateAgent(goals, Actions{})

	goalName, err := agent.getPrioritizedGoalName()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if goalName != "high_priority" {
		t.Errorf("Expected 'high_priority', got '%s'", goalName)
	}
}

func TestGetPrioritizedGoalName_ZeroPriority(t *testing.T) {
	goals := Goals{
		"inactive": {
			Conditions: Conditions{},
			PriorityFn: func(sensors Sensors) float32 {
				return 0.0
			},
		},
	}

	agent := CreateAgent(goals, Actions{})

	_, err := agent.getPrioritizedGoalName()

	if err == nil {
		t.Error("Expected error when all goals have zero priority")
	}
}

func TestGetPrioritizedGoalName_UsingSensors(t *testing.T) {
	type Entity struct {
		health int
	}

	goals := Goals{
		"heal": {
			Conditions: Conditions{},
			PriorityFn: func(sensors Sensors) float32 {
				entity := sensors.GetSensor("entity").(*Entity)
				if entity.health < 50 {
					return 2.0
				}
				return 0.1
			},
		},
		"explore": {
			Conditions: Conditions{},
			PriorityFn: func(sensors Sensors) float32 {
				return 1.0
			},
		},
	}

	agent := CreateAgent(goals, Actions{})
	entity := &Entity{health: 20}
	SetSensor(&agent, "entity", entity)

	goalName, err := agent.getPrioritizedGoalName()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if goalName != "heal" {
		t.Errorf("Expected 'heal' for low health, got '%s'", goalName)
	}
}

// Test GetNextAction
func TestGetNextAction_WithActions(t *testing.T) {
	plan := Plan{
		{name: "action1", cost: 1.0},
		{name: "action2", cost: 2.0},
	}

	action, err := plan.GetNextAction()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if action.name != "action1" {
		t.Errorf("Expected 'action1', got '%s'", action.name)
	}
}

func TestGetNextAction_EmptyPlan(t *testing.T) {
	plan := Plan{}

	_, err := plan.GetNextAction()

	if err == nil {
		t.Error("Expected error for empty plan")
	}
}

func TestGetNextAction_SingleAction(t *testing.T) {
	plan := Plan{
		{name: "only_action", cost: 1.0},
	}

	action, err := plan.GetNextAction()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if action.name != "only_action" {
		t.Errorf("Expected 'only_action', got '%s'", action.name)
	}
}

// Integration test with complex scenario
func TestGetPlan_ComplexScenario(t *testing.T) {
	actions := Actions{}

	// Need wood to make fire
	actions.AddAction("get_wood", 2.0, false, Conditions{}, Effects{
		EffectBool{Key: 1, Value: true, Operator: SET}, // has_wood
	})

	// Need matches to make fire
	actions.AddAction("get_matches", 1.0, false, Conditions{}, Effects{
		EffectBool{Key: 2, Value: true, Operator: SET}, // has_matches
	})

	// Make fire requires both wood and matches
	actions.AddAction("make_fire", 1.0, false, Conditions{
		&ConditionBool{Key: 1, Value: true, Operator: EQUAL}, // has_wood
		&ConditionBool{Key: 2, Value: true, Operator: EQUAL}, // has_matches
	}, Effects{
		EffectBool{Key: 3, Value: true, Operator: SET}, // has_fire
	})

	goals := Goals{
		"stay_warm": {
			Conditions: Conditions{
				&ConditionBool{Key: 3, Value: true, Operator: EQUAL}, // has_fire
			},
			PriorityFn: func(sensors Sensors) float32 {
				return 1.0
			},
		},
	}

	agent := CreateAgent(goals, actions)
	SetState[bool](&agent, 1, false) // has_wood
	SetState[bool](&agent, 2, false) // has_matches
	SetState[bool](&agent, 3, false) // has_fire

	goalName, plan := GetPlan(agent, 10)

	if goalName != "stay_warm" {
		t.Errorf("Expected goal 'stay_warm', got '%s'", goalName)
	}

	// Plan includes root + 3 actions
	if len(plan) != 4 {
		t.Errorf("Expected plan with 4 actions (root + 3), got %d", len(plan))
	}

	// Verify the plan makes sense (should get matches first as it's cheaper)
	// Index 1 should be a resource gathering action
	if plan[1].name != "get_matches" && plan[1].name != "get_wood" {
		t.Errorf("Expected second action to gather resources, got '%s'", plan[1].name)
	}

	// Last action should be make_fire
	if plan[3].name != "make_fire" {
		t.Errorf("Expected last action to be 'make_fire', got '%s'", plan[3].name)
	}

	// Verify total cost (root has cost 0)
	totalCost := plan.GetTotalCost()
	if totalCost != 4.0 {
		t.Errorf("Expected total cost 4.0, got %f", totalCost)
	}
}
