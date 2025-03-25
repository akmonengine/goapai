package goapai

import (
	"encoding/binary"
	"hash/fnv"
	"slices"
	"strconv"
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
	Check(states states, key StateKey) bool
	GetKey() StateKey
	GetValue() any
}
type State[T Numeric | bool | string] struct {
	Key   StateKey
	Value T
}

type StateKey uint16

type statesData []StateInterface

type states struct {
	Agent *Agent
	data  statesData
	hash  uint64
}

func (state State[T]) GetKey() StateKey {
	return state.Key
}

func (state State[T]) Check(states states, key StateKey) bool {
	k := states.data.GetIndex(key)
	if k < 0 {
		return false
	}
	s := states.data[k]
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

// Check compares states and states2 by their hash.
func (states states) Check(states2 states) bool {
	return states.hash == states2.hash
}

func (statesData statesData) GetIndex(stateKey StateKey) int {
	for k, stateData := range statesData {
		if stateData.GetKey() == stateKey {
			return k
		}
	}

	return -1
}

func (statesData statesData) sort() {
	slices.SortFunc(statesData, func(a, b StateInterface) int {
		if a.GetKey() > b.GetKey() {
			return 1
		} else if a.GetKey() < b.GetKey() {
			return -1
		}

		return 0
	})
}

func (statesData statesData) hashStates() uint64 {
	hash := fnv.New64()

	buf := make([]byte, binary.MaxVarintLen64)
	for _, data := range statesData {
		n := binary.PutVarint(buf, int64(data.GetKey()))
		hash.Write(buf[:n])
		hash.Write([]byte(":"))

		switch v := data.GetValue().(type) {
		case int8:
			n = binary.PutVarint(buf, int64(v))
			hash.Write(buf[:n])
		case int:
			n = binary.PutVarint(buf, int64(v))
			hash.Write(buf[:n])
		case uint8:
			n = binary.PutUvarint(buf, uint64(v))
			hash.Write(buf[:n])
		case uint64:
			n = binary.PutUvarint(buf, v)
			hash.Write(buf[:n])
		case float64:
			hash.Write([]byte(strconv.FormatFloat(v, 'f', -1, 64)))
		case string:
			hash.Write([]byte(v))
		case []byte:
			hash.Write(v)
		default:
			binary.Write(hash, binary.LittleEndian, data.GetValue())
		}
		hash.Write([]byte(";"))
	}

	return hash.Sum64()
}

type Sensor any
type Sensors map[string]Sensor

func (sensors Sensors) GetSensor(name string) Sensor {
	return sensors[name]
}

type ConditionInterface interface {
	GetKey() StateKey
	Check(states states) bool
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

func (conditionFn *ConditionFn) Check(states states) bool {
	if !conditionFn.resolved {
		conditionFn.valid = conditionFn.CheckFn(states.Agent.sensors)
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

func (condition *Condition[T]) Check(states states) bool {
	k := states.data.GetIndex(condition.Key)
	if k < 0 {
		return false
	}
	s := states.data[k]
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

func (conditionBool *ConditionBool) Check(states states) bool {
	k := states.data.GetIndex(conditionBool.Key)
	if k < 0 {
		return false
	}
	s := states.data[k]
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

func (conditionString *ConditionString) Check(states states) bool {
	k := states.data.GetIndex(conditionString.Key)
	if k < 0 {
		return false
	}
	s := states.data[k]
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

func (conditions Conditions) Check(states states) bool {
	for _, condition := range conditions {
		if !condition.Check(states) {
			return false
		}
	}

	return true
}
