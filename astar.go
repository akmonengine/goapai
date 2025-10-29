package goapai

import (
	"slices"
	"sync"
)

type node struct {
	*Action
	states states

	parentNode *node
	cost       float32
	totalCost  float32
	heuristic  float32
	depth      uint16
}

var nodesPool = sync.Pool{
	New: func() any {
		return make([]*node, 0, 32)
	},
}

func astar(from states, goal goalInterface, actions Actions, maxDepth int) Plan {
	availableActions := getImpactingActions(from, actions)
	openNodes := nodesPool.Get().([]*node)
	closedNodes := nodesPool.Get().([]*node)

	defer func() {
		nodesPool.Put(openNodes[:0])
		nodesPool.Put(closedNodes[:0])
	}()

	data := slices.Clone(from.data)
	openNodes = append(openNodes, &node{
		Action: &Action{},
		states: states{
			Agent: from.Agent,
			data:  data,
			hash:  data.hashStates(),
		},
		parentNode: nil,
		cost:       0,
		totalCost:  0,
		heuristic:  0,
		depth:      0,
	})

	for openNodeKey := 0; openNodeKey != -1; openNodeKey = getLessCostlyNodeKey(openNodes) {
		parentNode := openNodes[openNodeKey]
		if parentNode.depth > uint16(maxDepth) {
			openNodes = append(openNodes[:openNodeKey], openNodes[openNodeKey+1:]...)
			closedNodes = append(closedNodes, parentNode)
			continue
		}

		// Simulate world state, and check if we are at current state
		if countMissingGoal(goal, parentNode.states) == 0 {
			return buildPlanFromNode(parentNode)
		}

		for _, action := range availableActions {
			if !allowedRepetition(action, parentNode) {
				continue
			}

			if !action.conditions.Check(parentNode.states) {
				continue
			}

			simulatedStates, ok := simulateActionState(action, parentNode.states)
			if !ok {
				continue
			}

			if nodeKey, found := fetchNode(openNodes, simulatedStates); found {
				node := openNodes[nodeKey]
				if (parentNode.cost + action.cost) < node.cost {
					node.Action = action
					node.states = simulatedStates
					node.parentNode = parentNode
					node.cost = parentNode.cost + action.cost
					node.totalCost = parentNode.cost + action.cost + node.heuristic
					node.depth = parentNode.depth + 1

					openNodes[nodeKey] = node
				}
			} else if nodeKey, found := fetchNode(closedNodes, simulatedStates); found {
				node := closedNodes[nodeKey]
				if (parentNode.cost + action.cost) < node.cost {
					node.Action = action
					node.states = simulatedStates
					node.parentNode = parentNode
					node.cost = parentNode.cost + action.cost
					node.totalCost = parentNode.cost + action.cost + node.heuristic
					node.depth = parentNode.depth + 1

					openNodes[openNodeKey] = node
					closedNodes = append(closedNodes[:nodeKey], closedNodes[nodeKey+1:]...)
				}
			} else {
				heuristic := computeHeuristic(from, goal, simulatedStates)
				openNodes = append(openNodes, &node{
					Action:     action,
					states:     simulatedStates,
					parentNode: parentNode,
					cost:       parentNode.cost + action.cost,
					totalCost:  parentNode.cost + action.cost + heuristic,
					heuristic:  heuristic,
					depth:      parentNode.depth + 1,
				})
			}
		}

		openNodes = append(openNodes[:openNodeKey], openNodes[openNodeKey+1:]...)
		closedNodes = append(closedNodes, parentNode)
	}

	return Plan{}
}

// All the actions similar to initial states are useless:
// we consider they are not going towards the goal and are dead end
func getImpactingActions(from states, actions Actions) Actions {
	var availableActions Actions

	for _, action := range actions {
		if !action.effects.satisfyStates(from) {
			availableActions = append(availableActions, action)
		}
	}

	return availableActions
}

func getLessCostlyNodeKey(openNodes []*node) int {
	lowestKey := -1

	for key, node := range openNodes {
		if lowestKey < 0 || node.totalCost < openNodes[lowestKey].totalCost {
			lowestKey = key
		}
	}

	return lowestKey
}

func fetchNode(nodes []*node, states states) (int, bool) {
	for k, n := range nodes {
		if n.states.Check(states) {
			return k, true
		}
	}

	return 0, false
}

func buildPlanFromNode(node *node) Plan {
	plan := make(Plan, 0, node.depth)

	for node != nil {
		plan = append(plan, node.Action)
		node = node.parentNode
	}

	slices.Reverse(plan)

	return plan
}

func simulateActionState(action *Action, nodeStates states) (states, bool) {
	/* If action effects implies no changes to current worldState,
	then avoid generating huge chunks of memory */
	if action.effects.satisfyStates(nodeStates) {
		return states{}, false
	}

	data, err := action.effects.apply(nodeStates)
	if err != nil {
		return states{}, false
	}

	// Calculate hash incrementally by tracking changes
	newHash := nodeStates.hash

	// For each effect, we need to XOR out the old state and XOR in the new state
	for _, effect := range action.effects {
		// Find old state if it exists
		oldIndex := nodeStates.data.GetIndex(effect.GetKey())
		if oldIndex >= 0 {
			newHash ^= nodeStates.data[oldIndex].Hash() // Remove old
		}

		// Find new state in modified data
		newIndex := data.GetIndex(effect.GetKey())
		if newIndex >= 0 {
			newHash ^= data[newIndex].Hash() // Add new
		}
	}

	return states{
		Agent: nodeStates.Agent,
		data:  data,
		hash:  newHash,
	}, true
}

func allowedRepetition(action *Action, parentNode *node) bool {
	if action.repeatable {
		return true
	}

	node := parentNode
	for node != nil {
		if node.Action.name == action.name {
			return false
		}

		node = node.parentNode
	}

	return true
}

func countMissingGoal(goal goalInterface, states states) int {
	count := 0
	for _, condition := range goal.Conditions {
		if !condition.Check(states) {
			count++
		}
	}

	return count
}

/*
A very simple (empiristic) model for h using:
  - how much required states are met

We try to be conservative and reduce the number of steps
*/
func computeHeuristic(fromStates states, goal goalInterface, states states) float32 {
	missingGoalsCount := float32(countMissingGoal(goal, states))

	h := missingGoalsCount

	return h
}
