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
	String() string
}
type State[T Numeric] struct {
	Value T
}

type StateBool struct {
	Value bool
}

type StateString struct {
	Value string
}

type StateKey uint16

func (stateKey StateKey) String() string {
	return strconv.Itoa(int(stateKey))
}

type statesData map[StateKey]StateInterface

type states struct {
	Agent *Agent
	data  statesData
	hash  uint64
}

func (state State[T]) Check(states states, key StateKey) bool {
	if s, ok := states.data[key]; ok {
		if agentState, ok := s.(State[T]); ok {
			if agentState.Value == state.Value {
				return true
			}
		}
	}

	return false
}

func (state State[T]) String() string {
	return fmt.Sprint(state.Value)
}

func (stateBool StateBool) Check(states states, key StateKey) bool {
	return stateBool == states.data[key]
}

// String transforms the boolean value to string.
func (stateBool StateBool) String() string {
	return strconv.FormatBool(stateBool.Value)
}

func (stateString StateString) Check(states states, key StateKey) bool {
	if s, ok := states.data[key]; ok {
		if agentState, ok := s.(StateString); ok {
			if agentState.Value == stateString.Value {
				return true
			}
		}
	}

	return false
}

// String returns the string value.
func (stateString StateString) String() string {
	return stateString.Value
}

// Check compares states and states2 by their hash.
func (states states) Check(states2 states) bool {
	return states.hash == states2.hash
}

func (statesData statesData) hashStates() uint64 {
	hash := fnv.New64()
	keys := make([]StateKey, 0, len(statesData))

	for k := range statesData {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for _, k := range keys {
		hash.Write([]byte(k.String()))
		hash.Write([]byte(":"))
		hash.Write([]byte(statesData[k].String()))
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
	if s, ok := states.data[condition.Key]; ok {
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
	if s, ok := states.data[conditionBool.Key]; ok {
		if state, ok := s.(StateBool); ok {
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
	if s, ok := states.data[conditionString.Key]; ok {
		if state, ok := s.(StateString); ok {
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
