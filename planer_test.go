package goapai

import "testing"

func TestPlan_GetTotalCost(t *testing.T) {
	tests := []struct {
		name string
		plan Plan
		want float64
	}{
		{"plan 1", Plan{}, 0.0},
		{"plan 2", Plan{{
			name: "action 1",
			cost: 1.0,
		}}, 1.0},
		{"plan 3", Plan{{
			name: "action 1",
			cost: 1.0,
		}, {
			name: "action 2",
			cost: 2.0,
		}}, 3.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.plan.GetTotalCost(); got != tt.want {
				t.Errorf("GetTotalCost() = %v, want %v", got, tt.want)
			}
		})
	}
}
