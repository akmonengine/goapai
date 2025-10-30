package goapai

import (
	"container/heap"
	"slices"
)

type node struct {
	*Action
	world world

	parentNode *node
	cost       float32
	totalCost  float32
	heuristic  float32
	depth      uint16
	heapIndex  int  // Index in the heap, needed for heap.Fix
	closed     bool // true = closed node, false = open node
}

func astar(from world, goal goalInterface, actions Actions, maxDepth int) Plan {
	availableActions := getImpactingActions(from, actions)

	startNode := &node{
		Action: &Action{},
		world: world{
			Agent:  from.Agent,
			states: slices.Clone(from.states),
			hash:   from.hash,
		},
		parentNode: nil,
		heapIndex:  -1,
		closed:     false,
	}

	nodesHeap := nodeHeap{}
	heap.Init(&nodesHeap)
	heap.Push(&nodesHeap, startNode)

	for nodesHeap.Len() > 0 {
		parentNode := heap.Pop(&nodesHeap).(*node)

		if parentNode.depth > uint16(maxDepth) {
			parentNode.closed = true
			heap.Fix(&nodesHeap, parentNode.heapIndex)
			continue
		}

		// Simulate world state, and check if we are at current state
		if countMissingGoal(goal, parentNode.world) == 0 {
			return buildPlanFromNode(parentNode)
		}

		for _, action := range availableActions {
			if !allowedRepetition(action, parentNode) {
				continue
			}

			if !action.conditions.Check(parentNode.world) {
				continue
			}

			simulatedStates, ok := simulateActionState(action, parentNode.world)
			if !ok {
				continue
			}

			currentNode, found := fetchNodeInHeap(nodesHeap, simulatedStates)
			// Check if node exists in open nodes (closed=false)
			if found && !currentNode.closed {
				if (parentNode.cost + action.cost) < currentNode.cost {
					currentNode.Action = action
					currentNode.world = simulatedStates
					currentNode.parentNode = parentNode
					currentNode.cost = parentNode.cost + action.cost
					currentNode.totalCost = parentNode.cost + action.cost + currentNode.heuristic
					currentNode.depth = parentNode.depth + 1

					// Fix heap position after cost update
					heap.Fix(&nodesHeap, currentNode.heapIndex)
				}
			} else if found && currentNode.closed {
				// Node was closed, reopen it with better cost
				if (parentNode.cost + action.cost) < currentNode.cost {
					currentNode.Action = action
					currentNode.world = simulatedStates
					currentNode.parentNode = parentNode
					currentNode.cost = parentNode.cost + action.cost
					currentNode.totalCost = parentNode.cost + action.cost + currentNode.heuristic
					currentNode.depth = parentNode.depth + 1
					currentNode.closed = false // Reopen

					// Fix heap position
					heap.Fix(&nodesHeap, currentNode.heapIndex)
				}
			} else {
				// New node
				heuristic := computeHeuristic(from, goal, simulatedStates)
				newNode := &node{
					Action:     action,
					world:      simulatedStates,
					parentNode: parentNode,
					cost:       parentNode.cost + action.cost,
					totalCost:  parentNode.cost + action.cost + heuristic,
					heuristic:  heuristic,
					depth:      parentNode.depth + 1,
					heapIndex:  -1,
					closed:     false,
				}
				heap.Push(&nodesHeap, newNode)
			}
		}
	}

	return Plan{}
}

// All the actions similar to initial world are useless:
// we consider they are not going towards the goal and are dead end
func getImpactingActions(from world, actions Actions) Actions {
	var availableActions Actions

	for _, action := range actions {
		if !action.effects.satisfyStates(from) {
			availableActions = append(availableActions, action)
		}
	}

	return availableActions
}

func fetchNodeInHeap(heap nodeHeap, w world) (*node, bool) {
	for _, n := range heap {
		if n.world.Check(w) {
			return n, true
		}
	}
	return nil, false
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

func simulateActionState(action *Action, w world) (world, bool) {
	/* If action effects implies no changes to current worldState,
	then avoid generating huge chunks of memory */
	if action.effects.satisfyStates(w) {
		return world{}, false
	}

	w.states = slices.Clone(w.states)
	err := action.effects.apply(&w)
	if err != nil {
		return world{}, false
	}

	return w, true
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

func countMissingGoal(goal goalInterface, w world) int {
	count := 0
	for _, condition := range goal.Conditions {
		if !condition.Check(w) {
			count++
		}
	}

	return count
}

/*
Improved heuristic using numeric distance calculation:
  - For each goal condition, calculate the numeric distance between current value and target
  - Sum all distances to get total heuristic
  - This provides much better guidance than simple binary satisfied/unsatisfied check
*/
func computeHeuristic(from world, goal goalInterface, w world) float32 {
	var totalDistance float32

	for _, condition := range goal.Conditions {
		key := condition.GetKey()
		stateIndex := w.states.GetIndex(key)

		if stateIndex >= 0 {
			// State exists, calculate actual distance
			state := w.states[stateIndex]
			distance := state.Distance(condition)
			totalDistance += distance
		} else {
			// State doesn't exist, use pessimistic estimate
			// If the condition is not satisfied and state doesn't exist, assume distance of 1
			if !condition.Check(w) {
				totalDistance += 1.0
			}
		}
	}

	return totalDistance
}
