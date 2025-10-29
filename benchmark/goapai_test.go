package benchmark

import (
	"fmt"
	"goapai"
	"os"
	"runtime/pprof"
	"runtime/trace"
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
		&goapai.Condition[int]{Key: ATTRIBUTE_1, Operator: goapai.UPPER, Value: 0},
	}, goapai.Effects{
		goapai.Effect[int]{Key: ATTRIBUTE_1, Operator: goapai.SUBSTRACT, Value: 50},
		goapai.Effect[int]{Key: ATTRIBUTE_2, Operator: goapai.SUBSTRACT, Value: 5},
	})
	actions.AddAction("action2", 1, true, goapai.Conditions{
		&goapai.Condition[int]{Key: ATTRIBUTE_3, Operator: goapai.LOWER, Value: 50},
	}, goapai.Effects{
		goapai.Effect[int]{Key: ATTRIBUTE_3, Operator: goapai.ADD, Value: 20},
		goapai.Effect[int]{Key: ATTRIBUTE_2, Operator: goapai.ADD, Value: 10},
		goapai.Effect[int]{Key: ATTRIBUTE_1, Operator: goapai.ADD, Value: 5},
	})
	actions.AddAction("action3", 1, true, goapai.Conditions{
		&goapai.Condition[int]{Key: ATTRIBUTE_3, Operator: goapai.UPPER, Value: 30},
	}, goapai.Effects{
		goapai.Effect[int]{Key: ATTRIBUTE_3, Operator: goapai.SUBSTRACT, Value: 30},
	})

	entity := Entity{attributes: Attributes{}}

	goals := goapai.Goals{
		"goal1": {
			Conditions: goapai.Conditions{
				&goapai.Condition[int]{Key: ATTRIBUTE_2, Value: 80, Operator: goapai.UPPER},
			},
			PriorityFn: func(sensors goapai.Sensors) float32 {
				return 1.0
			},
		},
		"goal2": {
			Conditions: goapai.Conditions{
				&goapai.Condition[int]{Key: ATTRIBUTE_2, Value: 100, Operator: goapai.EQUAL},
			},
			PriorityFn: func(sensors goapai.Sensors) float32 {
				return 0.0
			},
		},
	}

	entity.agent = goapai.CreateAgent(goals, actions)
	goapai.SetSensor(&entity.agent, "entity", &entity)
	goapai.SetState[int](&entity.agent, ATTRIBUTE_2, 0)
	goapai.SetState[int](&entity.agent, ATTRIBUTE_1, 80)
	goapai.SetState[int](&entity.agent, ATTRIBUTE_3, 0)

	//Write to the trace file.
	f, _ := os.Create("trace.out")
	fcpu, _ := os.Create(`cpu.prof`)
	fheap, _ := os.Create(`heap.prof`)

	pprof.StartCPUProfile(fcpu)
	pprof.WriteHeapProfile(fheap)
	trace.Start(f)

	var lastPlan goapai.Plan
	for b.Loop() {
		//goapai.GetPlan(entity.agent, 15)
		_, lastPlan = goapai.GetPlan(entity.agent, 15)
	}
	fmt.Println("Actions for plan", len(lastPlan))
	for _, j := range lastPlan {
		fmt.Printf("		- %v\n", j.GetName())
	}

	defer f.Close()
	defer fcpu.Close()
	defer fheap.Close()

	trace.Stop()
	pprof.StopCPUProfile()

	b.ReportAllocs()
}
