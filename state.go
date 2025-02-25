package goapai

import (
	"fmt"
	"hash/fnv"
	"slices"
	"strconv"
)

type operator uint8

const (
	STATE_OPERATOR_EQUAL operator = iota
	STATE_OPERATOR_NOT_EQUAL
	STATE_OPERATOR_LOWER_OR_EQUAL
	STATE_OPERATOR_LOWER
	STATE_OPERATOR_UPPER_OR_EQUAL
	STATE_OPERATOR_UPPER
)

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
	~float32 | ~float64
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

	for _, data := range statesData {
		hash.Write([]byte(strconv.Itoa(int(data.GetKey()))))
		hash.Write([]byte(":"))
		hash.Write([]byte(fmt.Sprint(data.GetValue())))
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
		case STATE_OPERATOR_EQUAL:
			if state.Value == condition.Value {
				return true
			}
		case STATE_OPERATOR_NOT_EQUAL:
			if state.Value != condition.Value {
				return true
			}
		case STATE_OPERATOR_LOWER_OR_EQUAL:
			if state.Value <= condition.Value {
				return true
			}
		case STATE_OPERATOR_LOWER:
			if state.Value < condition.Value {
				return true
			}
		case STATE_OPERATOR_UPPER_OR_EQUAL:
			if state.Value >= condition.Value {
				return true
			}
		case STATE_OPERATOR_UPPER:
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
		case STATE_OPERATOR_EQUAL:
			if state.Value == conditionBool.Value {
				return true
			}
		case STATE_OPERATOR_NOT_EQUAL:
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
		case STATE_OPERATOR_EQUAL:
			if state.Value == conditionString.Value {
				return true
			}
		case STATE_OPERATOR_NOT_EQUAL:
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
