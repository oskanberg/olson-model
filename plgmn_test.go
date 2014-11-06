package main

import "testing"

func TestNewGateFromGenome(t *testing.T) {
	//                       in out in in in inoutoutoutout p1 p2 p3    p4 p5 p6 p7 p8
	genome := []byte{42, 213, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1}
	plgmn := NewPLGMN()
	plgmn.NewGateFromGenome(genome, 0)
	newGate := plgmn.gates[0]
	if len(newGate.in) != MinimumInNodes {
		t.Error("Generated gate does not have correct number of input nodes")
	}
	if len(newGate.out) != MinimumOutNodes {
		t.Error("Generated gate does not have correct number of output nodes")
	}
	// fmt.Println(newGate.ToString())
}
