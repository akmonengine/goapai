# GoapAI - Goal Oriented Action Planning
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/akmonengine/goapai)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Reference](https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://pkg.go.dev/github.com/akmonengine/goapai)
[![Go Report Card](https://goreportcard.com/badge/github.com/akmonengine/goapai)](https://goreportcard.com/report/github.com/akmonengine/goapai)
![Tests](https://img.shields.io/github/actions/workflow/status/akmonengine/goapai/code_coverage.yml?label=tests)
![Codecov](https://img.shields.io/codecov/c/github/akmonengine/goapai)
![GitHub Issues or Pull Requests](https://img.shields.io/github/issues/akmonengine/goapai)
![GitHub Issues or Pull Requests](https://img.shields.io/github/issues-pr/akmonengine/goapai)

GoapAI is a Go implementation of microlithic GOAP for AI agents in game development.

GOAP is a common tool to write autonomous AI agents, like NPCs. Those agents have goals and actions available. They can automatically
decide which is their current most important goal, and the best (ordered) actions to achieve it.

Similarly to a FSM (Finite State Machine), we end up with a graph of all possible States. But contrary to a FSM,
we do not have to manually define all the relational possibilities between States:
the A* pathfinding algorithm generates this graph on the fly.

GOAP was initially developed by Jeff Orkin for the game F.E.A.R., and has been proven viable for decision making AI.
Other algorithms that you might find useful are Behavioral Trees (BT) and Hierarchical Task Network (HTN).

## Monolithic vs Microlithic
Monolithic implementations of GOAP map one Action to one given Goal.

Microlithic are the most common version: one Action has a list of Effects, and no directed Goal.
It is the pathfinding algorithm that finds itself the chain of Actions to reach a given goal, depending of each Actions' effect.

One of the advantage of Microlithic GOAP is that it is considered to generate new and unpredicted behaviors,
where Monolithic uses predefined behaviors.
One of its disadvantage is the computation requirements, and exponential explosion if the Actions, their Conditions and Effects
are not well scoped.

## Features
- Multi-types States (numerics, bool, string), Conditions and Effects
- Relational operators (==, !=, <=, <, >=, >) for Conditions
- Algorithm operators (=, +, -, *, /) for Effects
- Usage of uint16 typed names for States (type StateKey), instead of the more common strings, to reduce the memory footprint
- Possibility to integrate custom States, Conditions & Effects through interface, for a better representation of your world
- Procedural preconditions through ConditionFn. These have access to your entity through Sensors, and are resolved once per planning request.
Because it is not registered in the worldState, it is a good tool to rely on to reduce the memory usage related to the standard GOAP algorithm. 
It does not duplicate a huge temporary worldState for each Effect.
Using this, you can most of the time use the worldState only as Actions' effects, and not initialize it with hundred of data that would not be
used anyway for a specific goal.
- Repeatable Actions. Non repeated Actions (default configuration) can hugely improve the performances of the algorithm.
But repeatable Actions can be a requirement for your goal (e.g. the AI needs 10 apples, the action "pick apple" gives one,
then this action should be repeated 10 times).
- A* forward implementation
- Floating Cost property on Actions: this allows a simple heuristic calculation in the A* path traveling,
for a better representation of your world in your Actions.
- Configurable Depth Limit to avoid generating plans of a hundred Actions

## Basic Usage
- First we need an Agent, to apply the AI on:
```go
type Entity struct {
    agent      goapai.Agent
    isHungry bool
    hasWood bool
}

entity := Entity{
    attributes: Attributes{
        isHungry: true,
        hasWood: false,
    },
}

entity.agent = goapai.CreateAgent(goals, actions)
```

- When you have a complex model and lot of data that does not requires to be mapped to the worldState (data not evolving during the Plan), you can use Sensors.
Sensors are not required, but they can have a huge performance impact so use them whenever you can:
```go
goapai.SetSensor(&entity.agent, "entity", &entity)
```

- Create a list of actions available for your agent. Each Action can be configured:
  - Its repeatability, can slower drastically the algorithm if true.
  - Its cost when executed, so that the pathfinding algorithm will select the solution with the smallest cost.
  - Its conditions, can be applied using the worldState attributes or a ConditionFn to check the entity properties directly.
  - Its effects, applied to the worldState.
```go
actions := goapai.Actions{}

actions.AddAction("fetch apple", 2, true, goapai.Conditions{}, goapai.Effects{
    goapai.Effect[bool]{
        Key:      ATTRIBUTE_HAS_APPLE,
        Value:    true,
    },
})

actions.AddAction("eat", 1, true, goapai.Conditions{
        goapai.ConditionBool{Key: ATTRIBUTE_HAS_APPLE, Value: true},
        goapai.ConditionBool{Key: ATTRIBUTE_HUNGRY, Value: true},
}, goapai.Effects{
    goapai.Effect[bool]{
        Key:      ATTRIBUTE_HUNGRY,
        Value:    false,
    },
    goapai.Effect[bool]{
        Key:      ATTRIBUTE_HAS_APPLE,
        Value:    false,
    },
})
```

- Create all the available goals for your agent. Each goal can be configured:
  - Its conditions to be met, so that the goal is considered achieved.
  - Its priority function calculation, so that the Planner can choose the most important goal to work on.
  You can use Sensors for a better definition of the priority depending on your data.
```go
goals := goapai.Goals{
    "eat": {
        Conditions: []goapai.ConditionInterface{
            goapai.ConditionBool{Key: ATTRIBUTE_HUNGRY, Value: false},
        },
        PriorityFn: func(sensors goapai.Sensors) float64 {
            if sensors.GetSensor("entity").(*Entity).attributes.isHungry {
                return 1.0
            }

            return 0.0
        },
    },
    "get wood": {
        Conditions: []goapai.ConditionInterface{
            goapai.Condition[int]{Key: ATTRIBUTE_HAS_WOOD, Value: true},
        },
        PriorityFn: func(sensors goapai.Sensors) float64 {
            if !sensors.GetSensor("entity").(*Entity).attributes.hasWood {
                return 0.5
            }

            return 0.0
        },
    },
}
entity.agent = goapai.CreateAgent(goals, actions)
```

- Set your initial WorldState for your agent:
```go
goapai.SetStateBool(&entity.agent, ATTRIBUTE_HUNGRY, entity.isHungry)
goapai.SetStateBool(&entity.agent, ATTRIBUTE_HAS_WOOD, entity.hasWood)
```

- Search the best Goal and the best Plan for it.
The maxDepth argument defines the maximum number of steps acceptable to achieve the Goal:
```go
goalName, plan := goapai.GetPlan(entity.agent, 10)
```
It returns the GoalName and the structure Plan being a slice of all the ordered Actions required for the Goal.

Multiple types are available for your conditions, states and effects:
```go
goapai.State[T Numeric]
goapai.StateBool
goapai.StateString

goapai.Condition[T Numeric]
goapai.ConditionBool
goapai.ConditionString
goapai.ConditionFn

goapai.Effect[T Numeric]
goapai.EffectBool
goapai.EffectString
```

Depending on your requirements, the number of Agents and the number of Actions,
you can either call goapai.GetPlan() every game loop or once per N frame, or only once an Action is resolved.
GOAP needs to be benchmarked and monitored regularly because of exponential risks with the WorldState.
Though if well scoped you can manage hundred of Actions for 200µs per Agent.

## What's next?
- The simulateActionState functions, to check the effect of an action on a node, takes up to 40% of CPU and 40% of memory.
We need to refactorize this part, or find another logical path.
- Heuristic calculation in A* is done poorly, we need a better algorithm to improve performances.
- Benchmark a backward implementation like D*, to improve performances.

## Sources
- https://web.archive.org/web/20230912145018/http://alumni.media.mit.edu/~jorkin/goap.html
- https://www.gamedevs.org/uploads/three-states-plan-ai-of-fear.pdf
- https://www.youtube.com/watch?v=gm7K68663rA
- https://excaliburjs.com/blog/goal-oriented-action-planning

## Contributing Guidelines

See [how to contribute](CONTRIBUTING.md).

## Licence
This project is distributed under the [Apache 2.0 licence](LICENCE.md).
