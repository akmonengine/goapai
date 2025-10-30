package goapai

import (
	"math"
)

type operator uint8

const (
	EQUAL operator = iota
	NOT_EQUAL
	LOWER_OR_EQUAL
	LOWER
	UPPER_OR_EQUAL
	UPPER
)

// Numeric is a constraint that defines the numeric types supported by generic State and Condition.
// Supported types are: int8, int, uint8, uint64, and float64.
type Numeric interface {
	~int8 | ~int |
		~uint8 | ~uint64 |
		~float64
}

// StateInterface defines the interface that all state types must implement.
// States represent key-value pairs in the world state, with support for hashing and distance calculation.
type StateInterface interface {
	Check(w world, key StateKey) bool
	GetKey() StateKey
	GetValue() any
	Store(w *world)
	GetHash() uint64
	Hash() uint64
	Distance(condition ConditionInterface) float32
}

// State represents a single key-value pair in the world state.
//
// States can hold numeric types (constrained by Numeric), bool, or string values.
// Each state is identified by a unique StateKey and includes a cached hash for performance.
type State[T Numeric | bool | string] struct {
	Key   StateKey // Unique identifier for this state
	Value T        // The state's value
	hash  uint64   // Cached hash value for fast comparison
}

// StateKey is a compact 16-bit unsigned integer used to identify states.
// Using uint16 instead of strings reduces memory usage and improves performance.
type StateKey uint16

type states []StateInterface

type world struct {
	Agent  *Agent
	states states
	hash   uint64
}

// Check compares world and states2 by their hash.
func (world world) Check(world2 world) bool {
	return world.hash == world2.hash
}

func (state State[T]) GetKey() StateKey {
	return state.Key
}

func (state State[T]) Check(w world, key StateKey) bool {
	k := w.states.GetIndex(key)
	if k < 0 {
		return false
	}
	s := w.states[k]
	if agentState, ok := s.(State[T]); ok {
		if agentState.Value == state.Value {
			return true
		}
	}

	return false
}

func (state State[T]) GetValue() any {
	return state.Value
}

func (state State[T]) Store(w *world) {
	oldHash := state.hash
	state.hash = state.Hash()
	w.hash = updateHashIncremental(w.hash, oldHash, state.hash)
	k := w.states.GetIndex(state.Key)
	if k < 0 {
		w.states = append(w.states, state)
	} else {
		w.states[k] = state
	}
}

func (state State[T]) GetHash() uint64 {
	return state.hash
}

// Hash returns a unique hash for this state using a fast multiplicative hash
// It implements a fast inline multiplicative hash
// Uses prime multipliers for good distribution without allocations
func (state State[T]) Hash() uint64 {
	const (
		prime1 uint64 = 11400714819323198485 // Large prime for key
		prime2 uint64 = 14029467366897019727 // Second prime for value
	)

	// Start with key
	hash := uint64(state.Key) * prime1

	// Mix in value based on type
	switch v := any(state.Value).(type) {
	case int8:
		hash ^= uint64(v) * prime2
	case int:
		hash ^= uint64(v) * prime2
	case uint8:
		hash ^= uint64(v) * prime2
	case uint64:
		hash ^= v * prime2
	case float64:
		hash ^= math.Float64bits(v) * prime2
	case bool:
		if v {
			hash ^= prime2
		}
	case string:
		// For strings, hash each byte
		for i := 0; i < len(v); i++ {
			hash = hash*prime2 ^ uint64(v[i])
		}
	}

	return hash
}

// updateHashIncremental updates a hash by removing old state and adding new state
func updateHashIncremental(currentHash uint64, oldStateHash, newStateHash uint64) uint64 {
	currentHash ^= oldStateHash // Remove old
	currentHash ^= newStateHash // Add new

	return currentHash
}

func (statesData states) GetIndex(stateKey StateKey) int {
	for k, stateData := range statesData {
		if stateData.GetKey() == stateKey {
			return k
		}
	}

	return -1
}

// Sensor is an alias for any type, used to store external data accessed by goal priority
// functions and procedural conditions.
type Sensor any

// Sensors is a map of sensor names to their values, providing external data to the agent
// without duplicating it in the world state during planning.
type Sensors map[string]Sensor

// GetSensor retrieves a sensor value by name.
// Returns nil if the sensor doesn't exist.
func (sensors Sensors) GetSensor(name string) Sensor {
	return sensors[name]
}

// ConditionInterface defines the interface that all condition types must implement.
// Conditions are preconditions that must be satisfied for actions or goals.
type ConditionInterface interface {
	GetKey() StateKey
	Check(w world) bool
}

// ConditionFn represents a procedural condition that evaluates against sensor data.
//
// Unlike state-based conditions, ConditionFn uses a custom function to check sensors.
// The result is cached after the first evaluation to avoid redundant computation during planning.
//
// Example:
//
//	condition := &ConditionFn{
//	    Key: 100,
//	    CheckFn: func(sensors Sensors) bool {
//	        health := sensors["health"].(int)
//	        return health < 50
//	    },
//	}
type ConditionFn struct {
	Key      StateKey                 // Unique identifier for this condition
	CheckFn  func(sensors Sensors) bool // Function that evaluates the condition
	resolved bool                       // Whether the condition has been evaluated
	valid    bool                       // Cached result of the evaluation
}

func (conditionFn *ConditionFn) GetKey() StateKey {
	return conditionFn.Key
}

func (conditionFn *ConditionFn) Check(w world) bool {
	if !conditionFn.resolved {
		conditionFn.valid = conditionFn.CheckFn(w.Agent.sensors)
		conditionFn.resolved = true
	}

	return conditionFn.valid
}

// Condition represents a numeric state-based condition with comparison operators.
//
// Conditions check if a state value satisfies a comparison (EQUAL, UPPER, LOWER, etc.)
// against a target value. Supported types are constrained by the Numeric interface.
//
// Example:
//
//	// Check if state key 1 is greater than or equal to 100
//	condition := &Condition[int]{
//	    Key:      1,
//	    Value:    100,
//	    Operator: UPPER_OR_EQUAL,
//	}
type Condition[T Numeric] struct {
	Key      StateKey // State key to check
	Value    T        // Target value to compare against
	Operator operator // Comparison operator (EQUAL, UPPER, LOWER, etc.)
}

func (condition *Condition[T]) GetKey() StateKey {
	return condition.Key
}

func (condition *Condition[T]) Check(w world) bool {
	k := w.states.GetIndex(condition.Key)
	if k < 0 {
		return false
	}
	s := w.states[k]
	if state, ok := s.(State[T]); ok {
		switch condition.Operator {
		case EQUAL:
			if state.Value == condition.Value {
				return true
			}
		case NOT_EQUAL:
			if state.Value != condition.Value {
				return true
			}
		case LOWER_OR_EQUAL:
			if state.Value <= condition.Value {
				return true
			}
		case LOWER:
			if state.Value < condition.Value {
				return true
			}
		case UPPER_OR_EQUAL:
			if state.Value >= condition.Value {
				return true
			}
		case UPPER:
			if state.Value > condition.Value {
				return true
			}
		}
	}

	return false
}

// ConditionBool represents a boolean state-based condition.
//
// Only EQUAL and NOT_EQUAL operators are supported for boolean conditions.
// Other operators will cause Check to return false.
//
// Example:
//
//	// Check if state key 2 is true
//	condition := &ConditionBool{
//	    Key:      2,
//	    Value:    true,
//	    Operator: EQUAL,
//	}
type ConditionBool struct {
	Key      StateKey // State key to check
	Value    bool     // Target boolean value
	Operator operator // Allowed: EQUAL, NOT_EQUAL
}

func (conditionBool *ConditionBool) GetKey() StateKey {
	return conditionBool.Key
}

func (conditionBool *ConditionBool) Check(w world) bool {
	k := w.states.GetIndex(conditionBool.Key)
	if k < 0 {
		return false
	}
	s := w.states[k]
	if state, ok := s.(State[bool]); ok {
		switch conditionBool.Operator {
		case EQUAL:
			if state.Value == conditionBool.Value {
				return true
			}
		case NOT_EQUAL:
			if state.Value != conditionBool.Value {
				return true
			}
		default:
			return false
		}
	}

	return false
}

// ConditionString represents a string state-based condition.
//
// Only EQUAL and NOT_EQUAL operators are supported for string conditions.
// Other operators will cause Check to return false.
//
// Example:
//
//	// Check if state key 3 equals "ready"
//	condition := &ConditionString{
//	    Key:      3,
//	    Value:    "ready",
//	    Operator: EQUAL,
//	}
type ConditionString struct {
	Key      StateKey // State key to check
	Value    string   // Target string value
	Operator operator // Allowed: EQUAL, NOT_EQUAL
}

// GetKey returns the state key that this condition checks.
func (conditionString *ConditionString) GetKey() StateKey {
	return conditionString.Key
}

func (conditionString *ConditionString) Check(w world) bool {
	k := w.states.GetIndex(conditionString.Key)
	if k < 0 {
		return false
	}
	s := w.states[k]
	if state, ok := s.(State[string]); ok {
		switch conditionString.Operator {
		case EQUAL:
			if state.Value == conditionString.Value {
				return true
			}
		case NOT_EQUAL:
			if state.Value != conditionString.Value {
				return true
			}
		default:
			return false
		}
	}

	return false
}

// Conditions is a collection of ConditionInterface implementations that must all be satisfied.
type Conditions []ConditionInterface

// Check returns true if all conditions in the collection are satisfied in the given world state.
// Returns true for an empty condition list (vacuous truth).
func (conditions Conditions) Check(w world) bool {
	for _, condition := range conditions {
		if !condition.Check(w) {
			return false
		}
	}

	return true
}
