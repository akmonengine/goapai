package goapai

type Agent struct {
	actions Actions
	w       world
	sensors Sensors
	goals   Goals
}

type goalInterface struct {
	Conditions []ConditionInterface
	PriorityFn GoalPriorityFn
}
type GoalName string
type Goals map[GoalName]goalInterface

func CreateAgent(goals Goals, actions Actions) Agent {
	agent := Agent{
		actions: actions,
		goals:   goals,
		sensors: Sensors{},
	}

	states := world{
		Agent:  &agent,
		states: states{},
	}
	agent.w = states

	return agent
}

func SetState[T Numeric | bool | string](agent *Agent, key StateKey, value T) {
	agent.w.states = append(agent.w.states, State[T]{
		Key:   key,
		Value: value,
	})
}

func SetSensor[T Sensor](agent *Agent, name string, value T) {
	agent.sensors[name] = value
}
