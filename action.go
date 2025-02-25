package goapai

import (
	"fmt"
	"slices"
)

type arithmetic uint8

const (
	EFFECT_ARITHMETIC_SET arithmetic = iota
	EFFECT_ARITHMETIC_ADD
	EFFECT_ARITHMETIC_SUBSTRACT
	EFFECT_ARITHMETIC_MULTIPLY
	EFFECT_ARITHMETIC_DIVIDE
)

type Action struct {
	name       string
	cost       float64
	repeatable bool
	conditions Conditions
	effects    Effects
}
type Actions []*Action

func (actions *Actions) AddAction(name string, cost float64, repeatable bool, conditions Conditions, effects Effects) {
	action := Action{
		name:       name,
		cost:       cost,
		repeatable: repeatable,
		conditions: conditions,
		effects:    effects,
	}

	*actions = append(*actions, &action)
}

func (action *Action) GetName() string {
	return action.name
}

func (action *Action) GetEffects() Effects {
	return action.effects
}

type EffectInterface interface {
	check(states states) bool
	apply(data statesData) error
}

type Effect[T Numeric] struct {
	Key      StateKey
	Value    T
	Operator arithmetic
}

func (effect Effect[T]) check(states states) bool {
	// Other operators than '=' mean the effect will have an impact of the states
	if effect.Operator != EFFECT_ARITHMETIC_SET {
		return false
	}

	k := states.data.GetIndex(effect.Key)
	if k < 0 {
		return false
	}
	s := states.data[k]

	if _, ok := s.(State[T]); !ok {
		return false
	}

	return s.(State[T]).Value == effect.Value
}

func (effect Effect[T]) apply(data statesData) error {
	k := data.GetIndex(effect.Key)
	if k < 0 {
		if slices.Contains([]arithmetic{EFFECT_ARITHMETIC_SET, EFFECT_ARITHMETIC_ADD}, effect.Operator) {
			data = append(data, State[T]{Value: effect.Value})
			return nil
		} else if slices.Contains([]arithmetic{EFFECT_ARITHMETIC_SUBSTRACT}, effect.Operator) {
			data = append(data, State[T]{Value: -effect.Value})
			return nil
		}
		return fmt.Errorf("data does not exist")
	}
	s := data[k]

	if _, ok := s.(State[T]); !ok {
		return fmt.Errorf("type does not match")
	}

	state := s.(State[T])
	switch effect.Operator {
	case EFFECT_ARITHMETIC_SET:
		state.Value = effect.Value
	case EFFECT_ARITHMETIC_ADD:
		state.Value += effect.Value
	case EFFECT_ARITHMETIC_SUBSTRACT:
		state.Value -= effect.Value
	case EFFECT_ARITHMETIC_MULTIPLY:
		state.Value *= effect.Value
	case EFFECT_ARITHMETIC_DIVIDE:
		state.Value /= effect.Value
	}

	data[k] = state

	return nil
}

type EffectBool struct {
	Key      StateKey
	Value    bool
	Operator arithmetic
}

func (effectBool EffectBool) check(states states) bool {
	// Other operators than '=' is not allowed
	if effectBool.Operator != EFFECT_ARITHMETIC_SET {
		return false
	}

	k := states.data.GetIndex(effectBool.Key)
	if k < 0 {
		return false
	}
	if _, ok := states.data[k].(State[bool]); !ok {
		return false
	}

	s := states.data[k].(State[bool])

	return s.Value == effectBool.Value
}

func (effectBool EffectBool) apply(data statesData) error {
	if effectBool.Operator != EFFECT_ARITHMETIC_SET {
		return fmt.Errorf("operation %v not allowed on bool type", effectBool.Operator)
	}

	k := data.GetIndex(effectBool.Key)
	if k < 0 {
		data = append(data, State[bool]{Value: effectBool.Value})
		return nil
	}
	if _, ok := data[k].(State[bool]); !ok {
		return fmt.Errorf("type does not match")
	}

	state := data[k].(State[bool])
	state.Value = effectBool.Value
	data[k] = state

	return nil
}

type EffectString struct {
	Key      StateKey
	Value    string
	Operator arithmetic
}

func (effectString EffectString) check(states states) bool {
	k := states.data.GetIndex(effectString.Key)
	if k < 0 {
		return false
	}
	if _, ok := states.data[k].(State[string]); !ok {
		return false
	}

	s := states.data[k].(State[string])

	return s.Value == effectString.Value
}

func (effectString EffectString) apply(data statesData) error {
	if !slices.Contains([]arithmetic{EFFECT_ARITHMETIC_SET, EFFECT_ARITHMETIC_ADD}, effectString.Operator) {
		return fmt.Errorf("arithmetic operation %v not allowed on string type", effectString.Operator)
	}

	k := data.GetIndex(effectString.Key)
	if k < 0 {
		data = append(data, State[string]{Value: effectString.Value})
		return nil
	}
	if _, ok := data[k].(State[string]); !ok {
		return fmt.Errorf("type does not match")
	}

	state := data[k].(State[string])
	switch effectString.Operator {
	case EFFECT_ARITHMETIC_SET:
		state.Value = effectString.Value
	case EFFECT_ARITHMETIC_ADD:
		state.Value = fmt.Sprint(state.Value, effectString.Value)
	}
	data[k] = state

	return nil
}

type EffectFn func(agent *Agent)

type Effects []EffectInterface

// If all the effects already exist in states,
// it is probably not a good path
func (effects Effects) satisfyStates(states states) bool {
	for _, effect := range effects {
		if !effect.check(states) {
			return false
		}
	}

	return true
}

func (effects Effects) apply(states states) (statesData, error) {
	data := slices.Clone(states.data)

	for _, effect := range effects {
		err := effect.apply(data)

		if err != nil {
			return nil, err
		}
	}

	return data, nil
}
