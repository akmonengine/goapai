package goapai

import (
	"slices"
	"testing"
)

// Test getImpactingActions
func TestGetImpactingActions(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 100)
	SetState[bool](&agent, 2, false)

	actions := Actions{}
	// Action that changes state - should be included
	actions.AddAction("change", 1.0, false, Conditions{}, Effects{
		Effect[int]{Key: 1, Value: 200, Operator: SET},
	})
	// Action with effects matching current state - should be excluded
	actions.AddAction("no_change", 1.0, false, Conditions{}, Effects{
		Effect[int]{Key: 1, Value: 100, Operator: SET},
		EffectBool{Key: 2, Value: false, Operator: SET},
	})

	impacting := getImpactingActions(agent.states, actions)

	if len(impacting) != 1 {
		t.Errorf("Expected 1 impacting action, got %d", len(impacting))
	}

	if impacting[0].name != "change" {
		t.Errorf("Expected 'change' action, got '%s'", impacting[0].name)
	}
}

func TestGetImpactingActions_AllImpacting(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 100)

	actions := Actions{}
	actions.AddAction("action1", 1.0, false, Conditions{}, Effects{
		Effect[int]{Key: 1, Value: 200, Operator: SET},
	})
	actions.AddAction("action2", 1.0, false, Conditions{}, Effects{
		Effect[int]{Key: 1, Value: 50, Operator: SET},
	})

	impacting := getImpactingActions(agent.states, actions)

	if len(impacting) != 2 {
		t.Errorf("Expected 2 impacting actions, got %d", len(impacting))
	}
}

// Test getLessCostlyNodeKey
func TestGetLessCostlyNodeKey(t *testing.T) {
	nodes := []*node{
		{totalCost: 10.0},
		{totalCost: 5.0},
		{totalCost: 15.0},
	}

	key := getLessCostlyNodeKey(nodes)
	if key != 1 {
		t.Errorf("Expected key 1 (lowest cost), got %d", key)
	}
}

func TestGetLessCostlyNodeKey_Empty(t *testing.T) {
	nodes := []*node{}

	key := getLessCostlyNodeKey(nodes)
	if key != -1 {
		t.Errorf("Expected -1 for empty list, got %d", key)
	}
}

// Test fetchNode
func TestFetchNode_Found(t *testing.T) {
	agent1 := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent1, 1, 100)
	agent1.states.data.sort()
	agent1.states.hash = agent1.states.data.hashStates()

	agent2 := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent2, 1, 100)
	agent2.states.data.sort()
	agent2.states.hash = agent2.states.data.hashStates()

	nodes := []*node{
		{states: agent1.states},
		{states: agent2.states},
	}

	key, found := fetchNode(nodes, agent1.states)
	if !found {
		t.Error("Expected to find node")
	}
	if key != 0 {
		t.Errorf("Expected key 0, got %d", key)
	}
}

func TestFetchNode_NotFound(t *testing.T) {
	agent1 := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent1, 1, 100)
	agent1.states.data.sort()
	agent1.states.hash = agent1.states.data.hashStates()

	agent2 := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent2, 1, 200)
	agent2.states.data.sort()
	agent2.states.hash = agent2.states.data.hashStates()

	nodes := []*node{
		{states: agent1.states},
	}

	_, found := fetchNode(nodes, agent2.states)
	if found {
		t.Error("Expected not to find node")
	}
}

// Test buildPlanFromNode
func TestBuildPlanFromNode(t *testing.T) {
	action1 := &Action{name: "action1", cost: 1.0}
	action2 := &Action{name: "action2", cost: 1.0}
	action3 := &Action{name: "action3", cost: 1.0}

	// Build a chain: nil -> node1 -> node2 -> node3
	node1 := &node{Action: action1, parentNode: nil, depth: 1}
	node2 := &node{Action: action2, parentNode: node1, depth: 2}
	node3 := &node{Action: action3, parentNode: node2, depth: 3}

	plan := buildPlanFromNode(node3)

	if len(plan) != 3 {
		t.Errorf("Expected plan length 3, got %d", len(plan))
	}

	if plan[0].name != "action1" {
		t.Errorf("Expected first action to be 'action1', got '%s'", plan[0].name)
	}
	if plan[1].name != "action2" {
		t.Errorf("Expected second action to be 'action2', got '%s'", plan[1].name)
	}
	if plan[2].name != "action3" {
		t.Errorf("Expected third action to be 'action3', got '%s'", plan[2].name)
	}
}

func TestBuildPlanFromNode_SingleNode(t *testing.T) {
	action1 := &Action{name: "action1", cost: 1.0}
	node1 := &node{Action: action1, parentNode: nil, depth: 1}

	plan := buildPlanFromNode(node1)

	if len(plan) != 1 {
		t.Errorf("Expected plan length 1, got %d", len(plan))
	}
}

// Test simulateActionState
func TestSimulateActionState(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 100)

	action := &Action{
		effects: Effects{
			Effect[int]{Key: 1, Value: 50, Operator: ADD},
		},
	}

	newStates, ok := simulateActionState(action, agent.states)
	if !ok {
		t.Error("Expected simulation to succeed")
	}

	idx := newStates.data.GetIndex(1)
	if idx < 0 {
		t.Error("Expected to find key 1 in new states")
	}

	if newStates.data[idx].(State[int]).Value != 150 {
		t.Errorf("Expected value 150, got %d", newStates.data[idx].(State[int]).Value)
	}
}

func TestSimulateActionState_NoChange(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 100)

	action := &Action{
		effects: Effects{
			Effect[int]{Key: 1, Value: 100, Operator: SET},
		},
	}

	_, ok := simulateActionState(action, agent.states)
	if ok {
		t.Error("Expected simulation to fail when effects match current state")
	}
}

// Test allowedRepetition
func TestAllowedRepetition_Repeatable(t *testing.T) {
	action := &Action{name: "test", repeatable: true}
	parentNode := &node{Action: action}

	if !allowedRepetition(action, parentNode) {
		t.Error("Expected repeatable action to be allowed")
	}
}

func TestAllowedRepetition_NonRepeatable_NotUsed(t *testing.T) {
	action1 := &Action{name: "action1", repeatable: false}
	action2 := &Action{name: "action2", repeatable: false}

	node1 := &node{Action: action2, parentNode: nil}

	if !allowedRepetition(action1, node1) {
		t.Error("Expected non-repeated action to be allowed")
	}
}

func TestAllowedRepetition_NonRepeatable_AlreadyUsed(t *testing.T) {
	action := &Action{name: "test", repeatable: false}

	node1 := &node{Action: action, parentNode: nil}
	node2 := &node{Action: &Action{name: "other"}, parentNode: node1}

	if allowedRepetition(action, node2) {
		t.Error("Expected repeated non-repeatable action to be disallowed")
	}
}

// Test countMissingGoal
func TestCountMissingGoal_AllMet(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 100)
	SetState[bool](&agent, 2, true)

	goal := goalInterface{
		Conditions: Conditions{
			&Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
			&ConditionBool{Key: 2, Value: true, Operator: EQUAL},
		},
	}

	count := countMissingGoal(goal, agent.states)
	if count != 0 {
		t.Errorf("Expected 0 missing goals, got %d", count)
	}
}

func TestCountMissingGoal_OneMissing(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 100)
	SetState[bool](&agent, 2, false)

	goal := goalInterface{
		Conditions: Conditions{
			&Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
			&ConditionBool{Key: 2, Value: true, Operator: EQUAL},
		},
	}

	count := countMissingGoal(goal, agent.states)
	if count != 1 {
		t.Errorf("Expected 1 missing goal, got %d", count)
	}
}

func TestCountMissingGoal_AllMissing(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 0)
	SetState[bool](&agent, 2, false)

	goal := goalInterface{
		Conditions: Conditions{
			&Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
			&ConditionBool{Key: 2, Value: true, Operator: EQUAL},
		},
	}

	count := countMissingGoal(goal, agent.states)
	if count != 2 {
		t.Errorf("Expected 2 missing goals, got %d", count)
	}
}

// Test computeHeuristic
func TestComputeHeuristic(t *testing.T) {
	fromAgent := CreateAgent(Goals{}, Actions{})
	SetState[int](&fromAgent, 1, 0)

	currentAgent := CreateAgent(Goals{}, Actions{})
	SetState[int](&currentAgent, 1, 50)

	goal := goalInterface{
		Conditions: Conditions{
			&Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
		},
	}

	heuristic := computeHeuristic(fromAgent.states, goal, currentAgent.states)
	if heuristic <= 0 {
		t.Error("Expected positive heuristic for unmet goal")
	}
}

func TestComputeHeuristic_GoalMet(t *testing.T) {
	fromAgent := CreateAgent(Goals{}, Actions{})
	SetState[int](&fromAgent, 1, 0)

	currentAgent := CreateAgent(Goals{}, Actions{})
	SetState[int](&currentAgent, 1, 100)

	goal := goalInterface{
		Conditions: Conditions{
			&Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
		},
	}

	heuristic := computeHeuristic(fromAgent.states, goal, currentAgent.states)
	if heuristic != 0 {
		t.Errorf("Expected 0 heuristic for met goal, got %f", heuristic)
	}
}

// Test astar integration
func TestAstar_SimpleGoal(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 0)

	actions := Actions{}
	actions.AddAction("increment", 1.0, true, Conditions{}, Effects{
		Effect[int]{Key: 1, Value: 10, Operator: ADD},
	})

	goal := goalInterface{
		Conditions: Conditions{
			&Condition[int]{Key: 1, Value: 30, Operator: EQUAL},
		},
	}

	plan := astar(agent.states, goal, actions, 10)

	// Plan includes the root node with empty action
	if len(plan) != 4 {
		t.Errorf("Expected plan with 4 actions (including root), got %d", len(plan))
	}
}

func TestAstar_UnreachableGoal(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 0)

	actions := Actions{}
	// No actions available

	goal := goalInterface{
		Conditions: Conditions{
			&Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
		},
	}

	plan := astar(agent.states, goal, actions, 10)

	if len(plan) != 0 {
		t.Errorf("Expected empty plan for unreachable goal, got %d actions", len(plan))
	}
}

func TestAstar_MaxDepth(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 0)

	actions := Actions{}
	actions.AddAction("increment", 1.0, true, Conditions{}, Effects{
		Effect[int]{Key: 1, Value: 1, Operator: ADD},
	})

	goal := goalInterface{
		Conditions: Conditions{
			&Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
		},
	}

	// Max depth of 5 should prevent reaching goal of 100
	plan := astar(agent.states, goal, actions, 5)

	if len(plan) > 5 {
		t.Errorf("Expected plan to respect maxDepth of 5, got %d actions", len(plan))
	}
}

func TestAstar_AlreadyAtGoal(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 100)

	actions := Actions{}

	goal := goalInterface{
		Conditions: Conditions{
			&Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
		},
	}

	plan := astar(agent.states, goal, actions, 10)

	// Plan includes root node with empty action when already at goal
	if len(plan) != 1 {
		t.Errorf("Expected plan with 1 action (root) when already at goal, got %d actions", len(plan))
	}
}

func TestAstar_PreferLowerCost(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 0)

	actions := Actions{}
	actions.AddAction("expensive", 10.0, false, Conditions{}, Effects{
		Effect[int]{Key: 1, Value: 100, Operator: SET},
	})
	actions.AddAction("cheap", 1.0, false, Conditions{}, Effects{
		Effect[int]{Key: 1, Value: 100, Operator: SET},
	})

	goal := goalInterface{
		Conditions: Conditions{
			&Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
		},
	}

	plan := astar(agent.states, goal, actions, 10)

	// Plan includes root + action
	if len(plan) != 2 {
		t.Errorf("Expected plan with 2 actions (root + 1), got %d", len(plan))
	}

	if plan[1].name != "cheap" {
		t.Errorf("Expected cheaper action to be chosen, got '%s'", plan[1].name)
	}
}

func TestAstar_RespectConditions(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 0)
	SetState[bool](&agent, 2, false)

	actions := Actions{}
	// This action requires key 2 to be true
	actions.AddAction("conditional", 1.0, false, Conditions{
		&ConditionBool{Key: 2, Value: true, Operator: EQUAL},
	}, Effects{
		Effect[int]{Key: 1, Value: 100, Operator: SET},
	})
	// This action enables the conditional action
	actions.AddAction("enabler", 1.0, false, Conditions{}, Effects{
		EffectBool{Key: 2, Value: true, Operator: SET},
	})

	goal := goalInterface{
		Conditions: Conditions{
			&Condition[int]{Key: 1, Value: 100, Operator: EQUAL},
		},
	}

	plan := astar(agent.states, goal, actions, 10)

	// Plan includes root + 2 actions
	if len(plan) != 3 {
		t.Errorf("Expected plan with 3 actions (root + 2), got %d", len(plan))
	}

	// Second action (index 1) should be the enabler
	if plan[1].name != "enabler" {
		t.Errorf("Expected second action to be 'enabler', got '%s'", plan[1].name)
	}
	// Third action (index 2) should be the conditional one
	if plan[2].name != "conditional" {
		t.Errorf("Expected third action to be 'conditional', got '%s'", plan[2].name)
	}
}

func TestAstar_NonRepeatableActions(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 0)

	actions := Actions{}
	actions.AddAction("increment", 1.0, false, Conditions{}, Effects{
		Effect[int]{Key: 1, Value: 10, Operator: ADD},
	})

	goal := goalInterface{
		Conditions: Conditions{
			&Condition[int]{Key: 1, Value: 30, Operator: EQUAL},
		},
	}

	plan := astar(agent.states, goal, actions, 10)

	// Non-repeatable action can only be used once, so goal is unreachable
	if len(plan) != 0 {
		t.Errorf("Expected empty plan for non-repeatable action, got %d actions", len(plan))
	}
}

func TestAstar_DataCloning(t *testing.T) {
	agent := CreateAgent(Goals{}, Actions{})
	SetState[int](&agent, 1, 100)

	originalData := slices.Clone(agent.states.data)

	actions := Actions{}
	actions.AddAction("modify", 1.0, false, Conditions{}, Effects{
		Effect[int]{Key: 1, Value: 200, Operator: SET},
	})

	goal := goalInterface{
		Conditions: Conditions{
			&Condition[int]{Key: 1, Value: 200, Operator: EQUAL},
		},
	}

	_ = astar(agent.states, goal, actions, 10)

	// Original state should not be modified
	if len(agent.states.data) != len(originalData) {
		t.Error("Original state was modified")
	}
	if agent.states.data[0].(State[int]).Value != originalData[0].(State[int]).Value {
		t.Error("Original state values were modified")
	}
}
