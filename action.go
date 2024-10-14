package goapai

import (
	"fmt"
	"maps"
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
	effectFn   EffectFn
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

	if _, ok := states.data[effect.Key]; !ok {
		return false
	}
	if _, ok := states.data[effect.Key].(State[T]); !ok {
		return false
	}

	s := states.data[effect.Key].(State[T])

	return s.Value == effect.Value
}

func (effect Effect[T]) apply(data statesData) error {
	if _, ok := data[effect.Key]; !ok {
		if slices.Contains([]arithmetic{EFFECT_ARITHMETIC_SET, EFFECT_ARITHMETIC_ADD}, effect.Operator) {
			data[effect.Key] = State[T]{Value: effect.Value}
			return nil
		} else if slices.Contains([]arithmetic{EFFECT_ARITHMETIC_SUBSTRACT}, effect.Operator) {
			data[effect.Key] = State[T]{Value: -effect.Value}
			return nil
		}
		return fmt.Errorf("data does not exist")
	}
	if _, ok := data[effect.Key].(State[T]); !ok {
		return fmt.Errorf("type does not match")
	}

	state := data[effect.Key].(State[T])
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

	data[effect.Key] = state

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

	if _, ok := states.data[effectBool.Key]; !ok {
		return false
	}
	if _, ok := states.data[effectBool.Key].(StateBool); !ok {
		return false
	}

	s := states.data[effectBool.Key].(StateBool)

	return s.Value == effectBool.Value
}

func (effectBool EffectBool) apply(data statesData) error {
	if effectBool.Operator != EFFECT_ARITHMETIC_SET {
		return fmt.Errorf("operation %v not allowed on bool type", effectBool.Operator)
	}

	if _, ok := data[effectBool.Key]; !ok {
		data[effectBool.Key] = StateBool{Value: effectBool.Value}
		return nil
	}
	if _, ok := data[effectBool.Key].(StateBool); !ok {
		return fmt.Errorf("type does not match")
	}

	state := data[effectBool.Key].(StateBool)
	state.Value = effectBool.Value
	data[effectBool.Key] = state

	return nil
}

type EffectString struct {
	Key      StateKey
	Value    string
	Operator arithmetic
}

func (effectString EffectString) check(states states) bool {
	if _, ok := states.data[effectString.Key]; !ok {
		return false
	}
	if _, ok := states.data[effectString.Key].(StateString); !ok {
		return false
	}

	s := states.data[effectString.Key].(StateString)

	return s.Value == effectString.Value
}

func (effectString EffectString) apply(data statesData) error {
	if !slices.Contains([]arithmetic{EFFECT_ARITHMETIC_SET, EFFECT_ARITHMETIC_ADD}, effectString.Operator) {
		return fmt.Errorf("arithmetic operation %v not allowed on string type", effectString.Operator)
	}

	if _, ok := data[effectString.Key]; !ok {
		data[effectString.Key] = StateString{Value: effectString.Value}
		return nil
	}
	if _, ok := data[effectString.Key].(StateString); !ok {
		return fmt.Errorf("type does not match")
	}

	state := data[effectString.Key].(StateString)
	switch effectString.Operator {
	case EFFECT_ARITHMETIC_SET:
		state.Value = effectString.Value
	case EFFECT_ARITHMETIC_ADD:
		state.Value = fmt.Sprint(state.Value, effectString.Value)
	}
	data[effectString.Key] = state

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

func (effects Effects) apply(states states) statesData {
	data := maps.Clone(states.data)

	for _, effect := range effects {
		effect.apply(data)
	}

	return data
}
