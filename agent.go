package goapai

type Agent struct {
	actions Actions
	states  states
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

	states := states{
		Agent: &agent,
		data:  statesData{},
	}
	agent.states = states

	return agent
}

func SetState[T Numeric](agent *Agent, key StateKey, value T) {
	agent.states.data[key] = State[T]{
		Value: value,
	}
}

func SetStateBool(agent *Agent, key StateKey, value bool) {
	agent.states.data[key] = StateBool{
		Value: value,
	}
}

func SetStateString(agent *Agent, key StateKey, value string) {
	agent.states.data[key] = StateString{
		Value: value,
	}
}

func SetSensor[T Sensor](agent *Agent, name string, value T) {
	agent.sensors[name] = value
}
