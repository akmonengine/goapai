package goapai

import (
	"fmt"
	"slices"
)

type arithmetic uint8

const (
	SET arithmetic = iota
	ADD
	SUBSTRACT
	MULTIPLY
	DIVIDE
)

type Action struct {
	name       string
	cost       float32
	repeatable bool
	conditions Conditions
	effects    Effects
}
type Actions []*Action

func (actions *Actions) AddAction(name string, cost float32, repeatable bool, conditions Conditions, effects Effects) {
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
	GetKey() StateKey
	check(w world) bool
	apply(w *world) error
}

type Effect[T Numeric] struct {
	Key      StateKey
	Operator arithmetic
	Value    T
}

func (effect Effect[T]) GetKey() StateKey {
	return effect.Key
}

func (effect Effect[T]) check(w world) bool {
	// Other operators than '=' mean the effect will have an impact of the world
	if effect.Operator != SET {
		return false
	}

	k := w.states.GetIndex(effect.Key)
	if k < 0 {
		return false
	}
	s := w.states[k]

	if _, ok := s.(State[T]); !ok {
		return false
	}

	return s.(State[T]).Value == effect.Value
}

func (effect Effect[T]) apply(w *world) error {
	k := w.states.GetIndex(effect.Key)
	if k < 0 {
		if slices.Contains([]arithmetic{SET, ADD}, effect.Operator) {
			w.states = append(w.states, State[T]{Value: effect.Value})
			return nil
		} else if slices.Contains([]arithmetic{SUBSTRACT}, effect.Operator) {
			w.states = append(w.states, State[T]{Value: -effect.Value})
			return nil
		}
		return fmt.Errorf("w does not exist")
	}
	if _, ok := w.states[k].(State[T]); !ok {
		return fmt.Errorf("type does not match")
	}

	state := w.states[k].(State[T])
	switch effect.Operator {
	case SET:
		state.Value = effect.Value
	case ADD:
		state.Value += effect.Value
	case SUBSTRACT:
		state.Value -= effect.Value
	case MULTIPLY:
		state.Value *= effect.Value
	case DIVIDE:
		state.Value /= effect.Value
	}

	state.Store(w)

	return nil
}

type EffectBool struct {
	Key      StateKey
	Value    bool
	Operator arithmetic
}

func (effectBool EffectBool) GetKey() StateKey {
	return effectBool.Key
}

func (effectBool EffectBool) check(w world) bool {
	// Other operators than '=' is not allowed
	if effectBool.Operator != SET {
		return false
	}

	k := w.states.GetIndex(effectBool.Key)
	if k < 0 {
		return false
	}
	if _, ok := w.states[k].(State[bool]); !ok {
		return false
	}

	s := w.states[k].(State[bool])

	return s.Value == effectBool.Value
}

func (effectBool EffectBool) apply(w *world) error {
	if effectBool.Operator != SET {
		return fmt.Errorf("operation %v not allowed on bool type", effectBool.Operator)
	}

	k := w.states.GetIndex(effectBool.Key)
	if k < 0 {
		w.states = append(w.states, State[bool]{Value: effectBool.Value})
		return nil
	}
	if _, ok := w.states[k].(State[bool]); !ok {
		return fmt.Errorf("type does not match")
	}

	state := w.states[k].(State[bool])
	state.Value = effectBool.Value

	state.Store(w)

	return nil
}

type EffectString struct {
	Key      StateKey
	Value    string
	Operator arithmetic
}

func (effectString EffectString) GetKey() StateKey {
	return effectString.Key
}

func (effectString EffectString) check(w world) bool {
	k := w.states.GetIndex(effectString.Key)
	if k < 0 {
		return false
	}
	if _, ok := w.states[k].(State[string]); !ok {
		return false
	}

	s := w.states[k].(State[string])

	return s.Value == effectString.Value
}

func (effectString EffectString) apply(w *world) error {
	if !slices.Contains([]arithmetic{SET, ADD}, effectString.Operator) {
		return fmt.Errorf("arithmetic operation %v not allowed on string type", effectString.Operator)
	}

	k := w.states.GetIndex(effectString.Key)
	if k < 0 {
		w.states = append(w.states, State[string]{Value: effectString.Value})
		return nil
	}
	if _, ok := w.states[k].(State[string]); !ok {
		return fmt.Errorf("type does not match")
	}

	state := w.states[k].(State[string])
	switch effectString.Operator {
	case SET:
		state.Value = effectString.Value
	case ADD:
		state.Value = fmt.Sprint(state.Value, effectString.Value)
	}

	state.Store(w)

	return nil
}

type EffectFn func(agent *Agent)

type Effects []EffectInterface

// If all the effects already exist in world,
// it is probably not a good path
func (effects Effects) satisfyStates(w world) bool {
	for _, effect := range effects {
		if !effect.check(w) {
			return false
		}
	}

	return true
}

func (effects Effects) apply(w *world) error {
	for _, effect := range effects {
		err := effect.apply(w)

		if err != nil {
			return err
		}
	}

	return nil
}
