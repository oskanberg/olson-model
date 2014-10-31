package main

import (
	"fmt"
	"math"
)

// PLGMN represents the whole MN of gates
type PLGMN struct {
	gates []*PLG
	nodes []Node
}

func (s *PLGMN) Run(sensorValues []bool) []bool {
	// just write to the first nodes
	for i, v := range sensorValues {
		s.nodes[i].SetState(v)
		// put new values into current timestep
		s.nodes[i].Step()
	}

	// activate each gate
	for i, _ := range s.gates {
		s.gates[i].Run()
	}

	// put new values into current timestep
	s.stepAllNodes()
	return s.getActuators()
}

func (s *PLGMN) stepAllNodes() {
	for i, _ := range s.nodes {
		s.nodes[i].Step()
	}
}

func (s *PLGMN) getActuators() []bool {
	actuatorVals := make([]bool, NumActuators)
	end := len(s.nodes)
	for i, v := range s.nodes[end-NumActuators : end] {
		actuatorVals[i] = v.GetState()
	}
	return actuatorVals
}

func (s *PLGMN) AddGate(newGate *PLG) {
	s.gates = append(s.gates, newGate)
}

func (s *PLGMN) ToString() string {
	str := "Begin PLGMN\n-------\n"
	str += fmt.Sprintf("%d gates\n", len(s.gates))
	str += fmt.Sprintf("%d nodes\n", len(s.nodes))
	str += "-------\n"
	for index, gate := range s.gates {
		str += fmt.Sprintf("-------\ngate %d\n-------\n", index)
		str += gate.ToString()
	}

	return str
}

func (s *PLGMN) NewGateFromGenome(genome []byte, startPosition int) {

	plg := NewPLG()

	genomeLen := len(genome)
	numNodes := len(s.nodes)

	// skip the start codon
	read := (startPosition + 2) % genomeLen

	// first byte is num of input nodes
	numberIn := FloorByte(genome[read] / (255 / MaximumInNodes))
	if numberIn < MinimumInNodes {
		numberIn = MinimumInNodes
	}
	read = (read + 1) % genomeLen

	// next byte is num of output nodes
	numberOut := FloorByte(genome[read] / (255 / MaximumOutNodes))
	if numberOut < MinimumOutNodes {
		numberOut = MinimumOutNodes
	}
	read = (read + 1) % genomeLen

	// next numberIn bytes are ids of input nodes
	for i := 0; i < numberIn; i++ {
		iterRead := (read + i) % genomeLen
		nodeIndex := RoundInt((int(genome[iterRead]) * numNodes) / 255)
		plg.AddInNode(&s.nodes[nodeIndex])
	}
	// there will always be MaximumInNodes spaces, if not vals
	read = (read + MaximumInNodes) % genomeLen

	// next numberOut bytes are ids of out nodes
	for i := 0; i < numberOut; i++ {
		iterRead := (read + i) % genomeLen
		nodeIndex := RoundInt((int(genome[iterRead]) * numNodes) / 255)
		plg.AddOutNode(&s.nodes[nodeIndex])
	}
	// there will always be MaximumOutNodes spaces, if not vals
	read = (read + MaximumOutNodes) % genomeLen

	//create float64 genome, for transition table
	genomeFloat64 := make([]float64, genomeLen)
	for i, v := range genome {
		genomeFloat64[i] = float64(v)
	}

	// read the probabilities for transition table
	// unfortunately can't do all at once since genome wraps
	leftToRead := int(math.Pow(2, float64(numberIn+numberOut)))
	// while the remainder wraps
	for read+leftToRead+1 > genomeLen {
		// read to end of genome
		plg.transitionTable = append(plg.transitionTable, genomeFloat64[read:]...)
		// subtract
		leftToRead -= genomeLen - read
		read = 0
	}

	plg.transitionTable = append(plg.transitionTable, genomeFloat64[read:read+leftToRead]...)
	plg.NormaliseTransitionTable()
	s.AddGate(plg)
}

func NewPLGMN() *PLGMN {
	nodes := make([]Node, NumTotalNodes)
	for i, _ := range nodes {
		nodes[i].SetId(i)
	}
	return &PLGMN{[]*PLG{}, nodes}
}
