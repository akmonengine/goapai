package benchmark

import (
	"github.com/kelindar/goap"
	"strings"
	"testing"
)

// NewAction creates a new action from the given name, require and outcome.
func NewAction(name, require, outcome string) *action {
	return &action{
		name:    name,
		require: goap.StateOf(strings.Split(require, ",")...),
		outcome: goap.StateOf(strings.Split(outcome, ",")...),
	}
}

// action represents a single action that can be performed by the agent.
type action struct {
	name    string
	cost    int
	require *goap.State
	outcome *goap.State
}

// Simulate simulates the action and returns the required and outcome states.
func (a *action) Simulate(current *goap.State) (*goap.State, *goap.State) {
	return a.require, a.outcome
}

// Cost returns the cost of the action.
func (a *action) Cost() float32 {
	return 1
}

func BenchmarkKelindarGoap(b *testing.B) {
	init := goap.StateOf("attribute1=80", "!attribute2", "!attribute3")
	goal := goap.StateOf("attribute2>80")

	actions := []goap.Action{
		NewAction("action1", "attribute1>0", "attribute1-50,attribute2-5"),
		NewAction("action2", "attribute3<50", "attribute3+20,attribute2+10,attribute1+5"),
		NewAction("action3", "attribute3>30", "attribute3-30"),
	}

	for b.Loop() {
		_, err := goap.Plan(init, goal, actions)
		if err != nil {
			panic(err)
		}
	}

	b.ReportAllocs()
}
