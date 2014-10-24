package main

import (
	"testing"
)

func TestNormalise(t *testing.T) {
	plg := NewPLG()
	plg.AddInNode(&Node{false, false, 0})
	plg.AddOutNode(&Node{false, false, 1})
	plg.transitionTable = []float64{123, 200, 91, 0}
	plg.NormaliseTransitionTable()
	if plg.transitionTable[0]+plg.transitionTable[1] != 1.0 {
		t.Error("Table row sums to ", plg.transitionTable[0]+plg.transitionTable[1])
	}
	if plg.transitionTable[2]+plg.transitionTable[3] != 1.0 {
		t.Error("Table row sums to ", plg.transitionTable[2]+plg.transitionTable[3])
	}
}
