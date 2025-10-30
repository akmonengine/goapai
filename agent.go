// Package goapai implements a microlithic Goal-Oriented Action Planning (GOAP) system for AI agents.
//
// GOAP is a planning technique where agents dynamically generate action sequences to achieve
// goals based on the current world state. This implementation uses A* pathfinding to find
// optimal plans.
//
// # Key Concepts
//
// Agent: The central entity that maintains world state, goals, actions, and sensors.
//
// State: Key-value pairs representing the world state. States support numeric types,
// booleans, and strings, identified by compact StateKey (uint16) values.
//
// Action: Operations that modify the world state. Each action has preconditions (Conditions)
// and postconditions (Effects), plus a cost for pathfinding optimization.
//
// Goal: Desired world states with priority functions. The planner always works on the
// highest priority goal.
//
// Sensor: External data sources used in goal prioritization and procedural conditions
// without duplicating data during planning.
//
// # Example Usage
//
//	// Create an agent
//	agent := goapai.CreateAgent(
//	    goapai.Goals{
//	        "survive": {
//	            Conditions: goapai.Conditions{
//	                &goapai.Condition[int]{Key: 1, Value: 50, Operator: goapai.UPPER_OR_EQUAL},
//	            },
//	            PriorityFn: func(sensors goapai.Sensors) float32 {
//	                return 1.0
//	            },
//	        },
//	    },
//	    goapai.Actions{},
//	)
//
//	// Set initial state
//	goapai.SetState[int](&agent, 1, 20) // health = 20
//
//	// Add actions
//	agent.actions.AddAction("heal", 1.0, false,
//	    goapai.Conditions{},
//	    goapai.Effects{
//	        goapai.Effect[int]{Key: 1, Operator: goapai.ADD, Value: 50},
//	    },
//	)
//
//	// Generate plan
//	goalName, plan := goapai.GetPlan(agent, 10)
package goapai

// Agent represents an AI agent that uses GOAP (Goal-Oriented Action Planning) to make decisions.
//
// An agent maintains a world state, a set of goals with priorities, available actions,
// and sensors for external data. The agent uses A* pathfinding to generate optimal plans
// that achieve its highest priority goal.
type Agent struct {
	actions Actions
	w       world
	sensors Sensors
	goals   Goals
}

type goalInterface struct {
	Conditions []ConditionInterface
	PriorityFn GoalPriorityFn
}

// GoalName is a unique identifier for a goal.
type GoalName string

// Goals is a map of goal names to their definitions (conditions and priority function).
type Goals map[GoalName]goalInterface

// CreateAgent creates and initializes a new Agent with the given goals and actions.
//
// The agent is initialized with an empty world state and no sensors. Use SetState
// to initialize the world state and SetSensor to add sensor data.
func CreateAgent(goals Goals, actions Actions) Agent {
	agent := Agent{
		actions: actions,
		goals:   goals,
		sensors: Sensors{},
	}

	states := world{
		Agent:  &agent,
		states: states{},
	}
	agent.w = states

	return agent
}

// SetState adds or updates a state value in the agent's world state.
//
// State values can be numeric types (int, int8, uint8, uint64, float64), bool, or string.
// Each state is identified by a unique StateKey. Multiple calls with the same key will
// create duplicate states; this is generally not recommended.
//
// Example:
//
//	SetState[int](&agent, 1, 100)      // Set state key 1 to integer 100
//	SetState[bool](&agent, 2, true)    // Set state key 2 to boolean true
//	SetState[string](&agent, 3, "foo") // Set state key 3 to string "foo"
func SetState[T Numeric | bool | string](agent *Agent, key StateKey, value T) {
	agent.w.states = append(agent.w.states, State[T]{
		Key:   key,
		Value: value,
	})
}

// SetSensor adds or updates a sensor value for the agent.
//
// Sensors provide external data that can be used in goal priority functions and
// procedural conditions (ConditionFn) without duplicating data during planning.
// Unlike world state, sensors are not modified during plan simulation.
//
// Example:
//
//	SetSensor(&agent, "health", 100)
//	SetSensor(&agent, "enemy_visible", true)
func SetSensor[T Sensor](agent *Agent, name string, value T) {
	agent.sensors[name] = value
}
