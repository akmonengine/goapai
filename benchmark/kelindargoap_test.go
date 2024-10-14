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
	b.StopTimer()

	init := goap.StateOf("hunger=80", "!food", "!tired")
	goal := goap.StateOf("food>80")

	actions := []goap.Action{
		NewAction("eat", "food>0", "hunger-50,food-5"),
		NewAction("forage", "tired<50", "tired+20,food+10,hunger+5"),
		NewAction("sleep", "tired>30", "tired-30"),
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := goap.Plan(init, goal, actions)
		if err != nil {
			panic(err)
		}
	}

	b.ReportAllocs()
}
