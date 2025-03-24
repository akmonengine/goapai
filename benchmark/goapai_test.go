package benchmark

import (
	"fmt"
	"goapai"
	"testing"
)

const (
	ATTRIBUTE_1 = iota
	ATTRIBUTE_2
	ATTRIBUTE_3
)

type Attributes struct{}
type Entity struct {
	agent      goapai.Agent
	attributes Attributes
}

func BenchmarkGoapAI(b *testing.B) {
	actions := goapai.Actions{}

	actions.AddAction("action1", 1, true, goapai.Conditions{
		&goapai.Condition[int]{Key: ATTRIBUTE_1, Value: 0, Operator: goapai.STATE_OPERATOR_UPPER},
	}, goapai.Effects{
		goapai.Effect[int]{
			Key:      ATTRIBUTE_1,
			Value:    50,
			Operator: goapai.EFFECT_ARITHMETIC_SUBSTRACT,
		},
		goapai.Effect[int]{
			Key:      ATTRIBUTE_2,
			Value:    5,
			Operator: goapai.EFFECT_ARITHMETIC_SUBSTRACT,
		},
	})
	actions.AddAction("action2", 1, true, goapai.Conditions{
		&goapai.Condition[int]{Key: ATTRIBUTE_3, Value: 50, Operator: goapai.STATE_OPERATOR_LOWER},
	}, goapai.Effects{
		goapai.Effect[int]{
			Key:      ATTRIBUTE_3,
			Value:    20,
			Operator: goapai.EFFECT_ARITHMETIC_ADD,
		},
		goapai.Effect[int]{
			Key:      ATTRIBUTE_2,
			Value:    10,
			Operator: goapai.EFFECT_ARITHMETIC_ADD,
		},
		goapai.Effect[int]{
			Key:      ATTRIBUTE_1,
			Value:    5,
			Operator: goapai.EFFECT_ARITHMETIC_ADD,
		},
	})
	actions.AddAction("action3", 1, true, goapai.Conditions{
		&goapai.Condition[int]{Key: ATTRIBUTE_3, Value: 30, Operator: goapai.STATE_OPERATOR_UPPER},
	}, goapai.Effects{
		goapai.Effect[int]{
			Key:      ATTRIBUTE_3,
			Value:    30,
			Operator: goapai.EFFECT_ARITHMETIC_SUBSTRACT,
		},
	})

	entity := Entity{attributes: Attributes{}}

	goals := goapai.Goals{
		"goal1": {
			Conditions: goapai.Conditions{
				&goapai.Condition[int]{Key: ATTRIBUTE_2, Value: 80, Operator: goapai.STATE_OPERATOR_UPPER},
			},
			PriorityFn: func(sensors goapai.Sensors) float64 {
				return 1.0
			},
		},
		"goal2": {
			Conditions: goapai.Conditions{
				&goapai.Condition[int]{Key: ATTRIBUTE_2, Value: 100, Operator: goapai.STATE_OPERATOR_EQUAL},
			},
			PriorityFn: func(sensors goapai.Sensors) float64 {
				return 0.0
			},
		},
	}

	entity.agent = goapai.CreateAgent(goals, actions)
	goapai.SetSensor(&entity.agent, "entity", &entity)
	goapai.SetState[int](&entity.agent, ATTRIBUTE_2, 0)
	goapai.SetState[int](&entity.agent, ATTRIBUTE_1, 80)
	goapai.SetState[int](&entity.agent, ATTRIBUTE_3, 0)

	// Write to the trace file.
	//f, _ := os.Create("trace.out")
	//fcpu, _ := os.Create(`cpu.prof`)
	//fheap, _ := os.Create(`heap.prof`)
	//
	//pprof.StartCPUProfile(fcpu)
	//pprof.WriteHeapProfile(fheap)
	//trace.Start(f)

	var lastPlan goapai.Plan
	for b.Loop() {
		//goapai.GetPlan(entity.agent, 15)
		_, lastPlan = goapai.GetPlan(entity.agent, 15)
	}
	fmt.Println(len(lastPlan))

	//defer f.Close()
	//defer fcpu.Close()
	//defer fheap.Close()
	//
	//trace.Stop()
	//pprof.StopCPUProfile()

	b.ReportAllocs()
}
