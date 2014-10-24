package main

import (
	"fmt"
	"testing"
)

func TestNewGateFromGenome(t *testing.T) {
	//                       in out in in in inoutoutoutout p1 p2 p3    p4 p5 p6 p7 p8
	genome := []byte{42, 213, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 100, 200, 1, 1, 1, 1}
	plgmn := NewPLGMN()
	plgmn.NewGateFromGenome(genome, 0)
	fmt.Println(plgmn.gates[0].ToString())
}
