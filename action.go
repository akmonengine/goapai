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

// Action represents a single action that an agent can perform to modify the world state.
//
// An action has preconditions (conditions) that must be met before it can be executed,
// and postconditions (effects) that describe how it modifies the world state.
// Actions have a cost that is used by the A* algorithm to find the optimal plan.
type Action struct {
	name       string
	cost       float32
	repeatable bool
	conditions Conditions
	effects    Effects
}

// Actions is a collection of Action pointers.
type Actions []*Action

// AddAction creates a new Action and appends it to the Actions collection.
//
// Parameters:
//   - name: unique identifier for the action
//   - cost: numeric cost used by pathfinding (lower costs are preferred)
//   - repeatable: if false, the action can only be used once per plan
//   - conditions: preconditions that must be satisfied before the action can be executed
//   - effects: postconditions that describe how the action modifies the world state
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

// GetName returns the action's name identifier.
func (action *Action) GetName() string {
	return action.name
}

// GetEffects returns the action's effects (postconditions).
func (action *Action) GetEffects() Effects {
	return action.effects
}

// EffectInterface defines the interface that all effect types must implement.
// Effects describe how an action modifies the world state.
type EffectInterface interface {
	GetKey() StateKey
	check(w world) bool
	apply(w *world) error
}

// Effect represents a numeric state modification for types constrained by Numeric.
//
// It supports arithmetic operators (SET, ADD, SUBTRACT, MULTIPLY, DIVIDE) to modify
// numeric state values. The effect is applied when an action is executed during planning.
type Effect[T Numeric] struct {
	Key      StateKey   // State key to modify
	Operator arithmetic // Arithmetic operation to perform
	Value    T          // Value to use in the operation
}

// GetKey returns the state key that this effect modifies.
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
			w.states = append(w.states, State[T]{Key: effect.Key, Value: effect.Value})
			return nil
		} else if slices.Contains([]arithmetic{SUBSTRACT}, effect.Operator) {
			w.states = append(w.states, State[T]{Key: effect.Key, Value: -effect.Value})
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

// EffectBool represents a boolean state modification.
//
// Only the SET operator is allowed for boolean effects. Attempting to use other
// operators (ADD, SUBTRACT, etc.) will result in an error when the effect is applied.
type EffectBool struct {
	Key      StateKey   // State key to modify
	Value    bool       // Boolean value to set
	Operator arithmetic // Must be SET
}

// GetKey returns the state key that this effect modifies.
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
		w.states = append(w.states, State[bool]{Key: effectBool.Key, Value: effectBool.Value})
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

// EffectString represents a string state modification.
//
// Supports SET (replace string) and ADD (concatenate) operators. Other operators
// (SUBTRACT, MULTIPLY, DIVIDE) are not allowed and will result in an error.
type EffectString struct {
	Key      StateKey   // State key to modify
	Value    string     // String value to use
	Operator arithmetic // Allowed: SET, ADD
}

// GetKey returns the state key that this effect modifies.
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
		w.states = append(w.states, State[string]{Key: effectString.Key, Value: effectString.Value})
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

// EffectFn is a function type for custom procedural effects that directly modify the agent.
// This allows for effects that cannot be expressed through simple state modifications.
type EffectFn func(agent *Agent)

// Effects is a collection of EffectInterface implementations that describe how
// an action modifies the world state.
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
