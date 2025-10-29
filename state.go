package goapai

import (
	"encoding/binary"
	"hash/fnv"
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
	Hash() uint64
}
type State[T Numeric | bool | string] struct {
	Key   StateKey
	Value T
}

type StateKey uint16

type statesData []StateInterface

type world struct {
	Agent *Agent
	data  statesData
	hash  uint64
}

func (state State[T]) GetKey() StateKey {
	return state.Key
}

func (state State[T]) Check(w world, key StateKey) bool {
	k := w.data.GetIndex(key)
	if k < 0 {
		return false
	}
	s := w.data[k]
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

// Hash returns a unique hash for this state using FNV-64a
func (state State[T]) Hash() uint64 {
	h := fnv.New64a()

	// Write the key
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(state.Key))
	h.Write(buf)

	// Write the value based on type
	buf = make([]byte, 8)
	switch v := any(state.Value).(type) {
	case int8:
		binary.LittleEndian.PutUint64(buf, uint64(v))
	case int:
		binary.LittleEndian.PutUint64(buf, uint64(v))
	case uint8:
		binary.LittleEndian.PutUint64(buf, uint64(v))
	case uint64:
		binary.LittleEndian.PutUint64(buf, v)
	case float64:
		binary.LittleEndian.PutUint64(buf, math.Float64bits(v))
	case bool:
		if v {
			binary.LittleEndian.PutUint64(buf, 1)
		} else {
			binary.LittleEndian.PutUint64(buf, 0)
		}
	case string:
		h.Write([]byte(v))
		return h.Sum64()
	}
	h.Write(buf)

	return h.Sum64()
}

// Check compares world and states2 by their hash.
func (world world) Check(world2 world) bool {
	return world.hash == world2.hash
}

func (statesData statesData) GetIndex(stateKey StateKey) int {
	for k, stateData := range statesData {
		if stateData.GetKey() == stateKey {
			return k
		}
	}

	return -1
}

// hashStates computes the initial hash using XOR of individual state hashes
// This is O(n) but only called once when creating initial state
func (statesData statesData) hashStates() uint64 {
	var hash uint64 = 0
	for _, state := range statesData {
		hash ^= state.Hash() // XOR for incremental updates
	}
	return hash
}

// updateHashIncremental updates a hash by removing old state and adding new state
// This is O(1) - the key optimization
func updateHashIncremental(currentHash uint64, oldState, newState StateInterface) uint64 {
	if oldState != nil {
		currentHash ^= oldState.Hash() // Remove old
	}
	if newState != nil {
		currentHash ^= newState.Hash() // Add new
	}
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
	k := w.data.GetIndex(condition.Key)
	if k < 0 {
		return false
	}
	s := w.data[k]
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
	k := w.data.GetIndex(conditionBool.Key)
	if k < 0 {
		return false
	}
	s := w.data[k]
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
	k := w.data.GetIndex(conditionString.Key)
	if k < 0 {
		return false
	}
	s := w.data[k]
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
