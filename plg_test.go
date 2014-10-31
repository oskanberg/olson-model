package main

import (
	"testing"
)

func TestNormalise(t *testing.T) {
	plg := NewPLG()
	plg.AddInNode(&Node{false, false, 0})
	plg.AddOutNode(&Node{false, false, 1})
	plg.AddOutNode(&Node{false, false, 1})
	plg.transitionTable = []float64{123, 200, 91, 0, 1, 3, 4, 100}
	plg.NormaliseTransitionTable()
	sum := plg.transitionTable[0] + plg.transitionTable[1] + plg.transitionTable[2] + plg.transitionTable[3]
	if sum < 0.9999999999999999 || sum > 1.0 {
		t.Error("Table row sums to ", sum)
	}
	sum = plg.transitionTable[4] + plg.transitionTable[5] + plg.transitionTable[6] + plg.transitionTable[7]
	if sum < 0.9999999999999999 || sum > 1.0 {
		t.Error("Table row sums to ", sum)
	}
}
