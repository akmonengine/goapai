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
	heapIndex  int // Index in the heap, needed for heap.Fix
}

// nodeHeap implements heap.Interface for a min-heap of nodes based on totalCost
type nodeHeap []*node

func (h nodeHeap) Len() int { return len(h) }

func (h nodeHeap) Less(i, j int) bool {
	return h[i].totalCost < h[j].totalCost
}

func (h nodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].heapIndex = i
	h[j].heapIndex = j
}

func (h *nodeHeap) Push(x interface{}) {
	n := x.(*node)
	n.heapIndex = len(*h)
	*h = append(*h, n)
}

func (h *nodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil      // avoid memory leak
	item.heapIndex = -1 // mark as removed
	*h = old[0 : n-1]
	return item
}

func astar(from world, goal goalInterface, actions Actions, maxDepth int) Plan {
	availableActions := getImpactingActions(from, actions)
	openNodes := make(map[uint64]*node)
	openNodesHeap := &nodeHeap{}
	heap.Init(openNodesHeap)
	closedNodes := make(map[uint64]*node)

	startNode := &node{
		Action: &Action{},
		world: world{
			Agent:  from.Agent,
			states: slices.Clone(from.states),
			hash:   from.hash,
		},
		parentNode: nil,
		cost:       0,
		totalCost:  0,
		heuristic:  0,
		depth:      0,
		heapIndex:  -1,
	}
	openNodes[startNode.world.hash] = startNode
	heap.Push(openNodesHeap, startNode)

	for openNodesHeap.Len() > 0 {
		parentNode := heap.Pop(openNodesHeap).(*node)
		delete(openNodes, parentNode.world.hash)

		if parentNode.depth > uint16(maxDepth) {
			closedNodes[parentNode.world.hash] = parentNode
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

			if currentNode, found := openNodes[simulatedStates.hash]; found {
				if (parentNode.cost + action.cost) < currentNode.cost {
					currentNode.Action = action
					currentNode.world = simulatedStates
					currentNode.parentNode = parentNode
					currentNode.cost = parentNode.cost + action.cost
					currentNode.totalCost = parentNode.cost + action.cost + currentNode.heuristic
					currentNode.depth = parentNode.depth + 1

					// Fix heap position after cost update
					heap.Fix(openNodesHeap, currentNode.heapIndex)
				}
			} else if currentNode, found := closedNodes[simulatedStates.hash]; found {
				if (parentNode.cost + action.cost) < currentNode.cost {
					currentNode.Action = action
					currentNode.world = simulatedStates
					currentNode.parentNode = parentNode
					currentNode.cost = parentNode.cost + action.cost
					currentNode.totalCost = parentNode.cost + action.cost + currentNode.heuristic
					currentNode.depth = parentNode.depth + 1

					openNodes[simulatedStates.hash] = currentNode
					delete(closedNodes, simulatedStates.hash)

					// Re-add to heap
					heap.Push(openNodesHeap, currentNode)
				}
			} else {
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
				}
				openNodes[simulatedStates.hash] = newNode
				heap.Push(openNodesHeap, newNode)
			}
		}

		closedNodes[parentNode.world.hash] = parentNode
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
A very simple (empiristic) model for h using:
  - how much required world are met

We try to be conservative and reduce the number of steps
*/
func computeHeuristic(from world, goal goalInterface, w world) float32 {
	missingGoalsCount := float32(countMissingGoal(goal, w))

	h := missingGoalsCount

	return h
}
