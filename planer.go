package goapai

import "fmt"

type Plan Actions

type GoalPriorityFn func(sensors Sensors) float32

// GetTotalCost returns the cost of Plan.
//
// This total cost is the sum of all actions's cost.
func (plan Plan) GetTotalCost() float32 {
	var cost float32

	for _, action := range plan {
		cost += action.cost
	}

	return cost
}

// GetPlan returns the current GoalName, and the best Plan to achieve this Goal.
//
// The maxDepth argument limits the number of actions required to match the goal.
// Plan can be empty if the number of actions required is upper than maxDepth, or if the goal is unreachable.
func GetPlan(agent Agent, maxDepth int) (GoalName, Plan) {
	goalName, err := agent.getPrioritizedGoalName()

	if err != nil {
		fmt.Println(err)

		return "", Plan{}
	}

	for _, state := range agent.w.states {
		state.Store(&agent.w)
	}

	return goalName, astar(agent.w, agent.goals[goalName], agent.actions, maxDepth)
}

func (agent *Agent) getPrioritizedGoalName() (GoalName, error) {
	var prioritizedGoalName GoalName
	var prioritizedValue float32

	for name, goal := range agent.goals {
		priority := goal.PriorityFn(agent.sensors)

		if priority > prioritizedValue {
			prioritizedGoalName = name
			prioritizedValue = priority
		}
	}

	if prioritizedValue > 0.0 {
		return prioritizedGoalName, nil
	} else {
		return prioritizedGoalName, fmt.Errorf("no goal available")
	}
}

// GetNextAction returns the first Action required to achieve the Plan.
//
// An error is returned if no action is available, meaning the Plan is empty.
func (plan Plan) GetNextAction() (Action, error) {
	if len(plan) > 0 {
		return *plan[0], nil
	}

	return Action{}, fmt.Errorf("no action available")
}
