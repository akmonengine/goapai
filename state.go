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

type Numeric interface {
	~int8 | ~int |
		~uint8 | ~uint64 |
		~float64
}

type StateInterface interface {
	Check(w world, key StateKey) bool
	GetKey() StateKey
	GetValue() any
	Store(w *world)
	GetHash() uint64
	Hash() uint64
}
type State[T Numeric | bool | string] struct {
	Key   StateKey
	Value T
	hash  uint64
}

type StateKey uint16

type states []StateInterface

type world struct {
	Agent  *Agent
	states states
	hash   uint64
}

func createState[T Numeric | bool | string](s State[T]) State[T] {
	s.hash = s.Hash()

	return s
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

// Check compares world and states2 by their hash.
func (world world) Check(world2 world) bool {
	return world.hash == world2.hash
}

func (statesData states) GetIndex(stateKey StateKey) int {
	for k, stateData := range statesData {
		if stateData.GetKey() == stateKey {
			return k
		}
	}

	return -1
}

// updateHashIncremental updates a hash by removing old state and adding new state
func updateHashIncremental(currentHash uint64, oldStateHash, newStateHash uint64) uint64 {
	currentHash ^= oldStateHash // Remove old
	currentHash ^= newStateHash // Add new

	return currentHash
}

type Sensor any
type Sensors map[string]Sensor

func (sensors Sensors) GetSensor(name string) Sensor {
	return sensors[name]
}

type ConditionInterface interface {
	GetKey() StateKey
	Check(w world) bool
}

type ConditionFn struct {
	Key      StateKey
	CheckFn  func(sensors Sensors) bool
	resolved bool
	valid    bool
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

type Condition[T Numeric] struct {
	Key      StateKey
	Value    T
	Operator operator
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

type ConditionBool struct {
	Key      StateKey
	Value    bool
	Operator operator
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

type ConditionString struct {
	Key      StateKey
	Value    string
	Operator operator
}

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

type Conditions []ConditionInterface

func (conditions Conditions) Check(w world) bool {
	for _, condition := range conditions {
		if !condition.Check(w) {
			return false
		}
	}

	return true
}
