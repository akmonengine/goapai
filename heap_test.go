package goapai

import "testing"

func TestNodeHeap_Less(t *testing.T) {
	tests := []struct {
		name     string
		nodes    nodeHeap
		i, j     int
		wantLess bool
	}{
		{
			name: "first node has lower cost",
			nodes: nodeHeap{
				{totalCost: 5.0},
				{totalCost: 10.0},
			},
			i:        0,
			j:        1,
			wantLess: true,
		},
		{
			name: "first node has higher cost",
			nodes: nodeHeap{
				{totalCost: 15.0},
				{totalCost: 10.0},
			},
			i:        0,
			j:        1,
			wantLess: false,
		},
		{
			name: "equal costs",
			nodes: nodeHeap{
				{totalCost: 10.0},
				{totalCost: 10.0},
			},
			i:        0,
			j:        1,
			wantLess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.nodes.Less(tt.i, tt.j)
			if got != tt.wantLess {
				t.Errorf("Less(%d, %d) = %v, want %v", tt.i, tt.j, got, tt.wantLess)
			}
		})
	}
}

func TestNodeHeap_Swap(t *testing.T) {
	tests := []struct {
		name         string
		initialNodes []*node
		i, j         int
		wantI, wantJ int
		wantIIndex   int
		wantJIndex   int
	}{
		{
			name: "swap first and second",
			initialNodes: []*node{
				{totalCost: 5.0, heapIndex: 0},
				{totalCost: 10.0, heapIndex: 1},
			},
			i:          0,
			j:          1,
			wantI:      1,
			wantJ:      0,
			wantIIndex: 0,
			wantJIndex: 1,
		},
		{
			name: "swap first and third",
			initialNodes: []*node{
				{totalCost: 5.0, heapIndex: 0},
				{totalCost: 10.0, heapIndex: 1},
				{totalCost: 15.0, heapIndex: 2},
			},
			i:          0,
			j:          2,
			wantI:      2,
			wantJ:      0,
			wantIIndex: 0,
			wantJIndex: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := nodeHeap(tt.initialNodes)
			originalI := h[tt.i]
			originalJ := h[tt.j]

			h.Swap(tt.i, tt.j)

			// Verify nodes were swapped
			if h[tt.i] != originalJ {
				t.Errorf("After swap, h[%d] should be originalJ", tt.i)
			}
			if h[tt.j] != originalI {
				t.Errorf("After swap, h[%d] should be originalI", tt.j)
			}

			// Verify heapIndex was updated
			if h[tt.i].heapIndex != tt.wantIIndex {
				t.Errorf("h[%d].heapIndex = %d, want %d", tt.i, h[tt.i].heapIndex, tt.wantIIndex)
			}
			if h[tt.j].heapIndex != tt.wantJIndex {
				t.Errorf("h[%d].heapIndex = %d, want %d", tt.j, h[tt.j].heapIndex, tt.wantJIndex)
			}
		})
	}
}

func TestNodeHeap_Push(t *testing.T) {
	tests := []struct {
		name         string
		initialNodes []*node
		pushNode     *node
		wantLen      int
		wantIndex    int
	}{
		{
			name:         "push to empty heap",
			initialNodes: []*node{},
			pushNode:     &node{totalCost: 5.0},
			wantLen:      1,
			wantIndex:    0,
		},
		{
			name: "push to non-empty heap",
			initialNodes: []*node{
				{totalCost: 5.0, heapIndex: 0},
				{totalCost: 10.0, heapIndex: 1},
			},
			pushNode:  &node{totalCost: 15.0},
			wantLen:   3,
			wantIndex: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := nodeHeap(tt.initialNodes)
			h.Push(tt.pushNode)

			if len(h) != tt.wantLen {
				t.Errorf("Len() = %d, want %d", len(h), tt.wantLen)
			}

			if tt.pushNode.heapIndex != tt.wantIndex {
				t.Errorf("heapIndex = %d, want %d", tt.pushNode.heapIndex, tt.wantIndex)
			}

			if h[len(h)-1] != tt.pushNode {
				t.Error("Pushed node should be at the end of heap")
			}
		})
	}
}

func TestNodeHeap_Pop(t *testing.T) {
	tests := []struct {
		name         string
		initialNodes []*node
		wantNode     *node
		wantLen      int
	}{
		{
			name: "pop from heap with one element",
			initialNodes: []*node{
				{totalCost: 5.0, heapIndex: 0},
			},
			wantNode: &node{totalCost: 5.0, heapIndex: -1},
			wantLen:  0,
		},
		{
			name: "pop from heap with multiple elements",
			initialNodes: []*node{
				{totalCost: 5.0, heapIndex: 0},
				{totalCost: 10.0, heapIndex: 1},
				{totalCost: 15.0, heapIndex: 2},
			},
			wantNode: &node{totalCost: 15.0, heapIndex: -1},
			wantLen:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := nodeHeap(tt.initialNodes)
			originalLast := h[len(h)-1]

			popped := h.Pop().(*node)

			if len(h) != tt.wantLen {
				t.Errorf("Len() = %d, want %d", len(h), tt.wantLen)
			}

			if popped.heapIndex != -1 {
				t.Errorf("Popped node heapIndex = %d, want -1", popped.heapIndex)
			}

			if popped != originalLast {
				t.Error("Pop should return the last element")
			}
		})
	}
}
