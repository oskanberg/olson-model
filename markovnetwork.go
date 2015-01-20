package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

// PLGMN represents the whole MN of gates
type MarkovNetwork struct {
	gates      []*ProbabilisticGate
	nodes      []bool
	statistics map[string][]bool
}

func (s *MarkovNetwork) Reset() {
	for i, _ := range s.nodes {
		s.nodes[i] = false
	}
}

func (s *MarkovNetwork) AddGateFromGenome(genome []byte, startPosition int) {
	gate := NewProbabilisticGate()
	gate.LoadFromGenome(genome, startPosition)
	s.gates = append(s.gates, gate)
}

func (s *MarkovNetwork) Run(sensorValues []bool) []bool {
	padded := make([]bool, len(s.nodes))
	copy(padded, sensorValues)
	err, lOr := LogicalOr(s.nodes, padded)
	if err != nil {
		panic(err)
	}
	s.nodes = lOr
	for i, _ := range s.gates {
		err, lOr = LogicalOr(s.nodes, s.gates[i].Run(padded))
		if err != nil {
			panic(err)
		}
		s.nodes = lOr
	}

	s.statistics[fmt.Sprint(sensorValues)] = s.nodes[NumTotalNodes-NumActuators : NumTotalNodes]
	// actuators are just the last nodes
	return s.nodes[NumTotalNodes-NumActuators : NumTotalNodes]
}

//TODO
func (s *MarkovNetwork) ToString() string {
	var description string = "Contains " + fmt.Sprint(len(s.gates)) + " gates.\n"

	readFrom := make(map[byte]bool)
	for i, _ := range s.gates {
		for _, v := range s.gates[i].inputIndexes {
			readFrom[v] = true
		}
	}

	description += "Read from nodes: "
	for k, _ := range readFrom {
		description += fmt.Sprint(k) + ", "
	}

	description += "\n"
	for i, _ := range s.gates {
		description += "Gate " + fmt.Sprint(i) + " reads from "
		for _, v := range s.gates[i].inputIndexes {
			description += fmt.Sprint(v) + ", "
		}
		description += "\n"
	}

	writeTo := make(map[byte]bool)
	for i, _ := range s.gates {
		for _, v := range s.gates[i].outputIndexes {
			writeTo[v] = true
		}
	}

	description += "\nWrite to nodes: "
	for k, _ := range writeTo {
		description += fmt.Sprint(k) + ", "
	}

	description += "\n"
	for i, _ := range s.gates {
		description += "Gate " + fmt.Sprint(i) + " writes to "
		for _, v := range s.gates[i].outputIndexes {
			if v >= (NumTotalNodes - NumActuators) {
				description += "**************"
			}
			description += fmt.Sprint(v) + ", "
		}
		description += "\n"
	}

	return description
}

// func (s *MarkovNetwork) PrintStatistics() {
// 	keys := make([]string, len(s.statistics))
// 	i := 0
// 	for k, _ := range s.statistics {
// 		keys[i] = k
// 		i++
// 	}
// 	sort.Strings(keys)
// 	for _, k := range keys {
// 		if s.statistics[k][0] && !s.statistics[k][1] {
// 			fmt.Println("\n\n\n", k, "\n", s.statistics[k])
// 		}
// 	}
// }

func NewMarkovNetwork() *MarkovNetwork {
	mn := &MarkovNetwork{
		nodes: make([]bool, NumTotalNodes),
		// statistics: make(map[string][]bool),
	}
	if RigMarkovNetwork {
		// artificially ensure that some gates output to both actuators
		genome := GenerateRandomGenome(InitialGenomeLength, 0)

		// 4 + Max
		genome[4+MaximumInNodes] = NumTotalNodes - 1
		genome[5+MaximumInNodes] = NumTotalNodes - 1
		genome[6+MaximumInNodes] = NumTotalNodes - 2
		genome[4+MaximumInNodes+MaximumOutNodes] = NumTotalNodes - 1
		genome[5+MaximumInNodes+MaximumOutNodes] = NumTotalNodes - 2
		genome[6+MaximumInNodes+MaximumOutNodes] = NumTotalNodes - 2
		mn.AddGateFromGenome(genome, 0)
		mn.AddGateFromGenome(genome, 1)
	}
	return mn
}

func NewProbabilisticGate() *ProbabilisticGate {
	return &ProbabilisticGate{
		transitionTable: make(map[string][]byte),
	}
}

type ProbabilisticGate struct {
	inputIndexes  []byte
	outputIndexes []byte

	// bit of a hack - convert input to string
	transitionTable map[string][]byte

	outputCombinationCache [][]bool
}

func (s *ProbabilisticGate) LoadFromGenome(genome []byte, startPosition int) {
	genomeLen := len(genome)

	// skip the start codon
	read := (startPosition + 2) % genomeLen

	// number of inputs to gate
	numIn := genome[read] % (MaximumInNodes + 1)
	numIn = byte(math.Max(float64(numIn), MinimumInNodes))
	s.inputIndexes = make([]byte, numIn)

	read = (read + 1) % genomeLen

	// number of outputs from gate
	numOut := genome[read] % (MaximumOutNodes + 1)
	numOut = byte(math.Max(float64(numOut), MinimumOutNodes))
	s.outputIndexes = make([]byte, numOut)

	for i := byte(0); i < numIn; i++ {
		read = (read + 1) % genomeLen
		s.inputIndexes[i] = genome[read] % (NumTotalNodes - NumActuators)
	}

	// skip over rest of (poss redundant) indexes
	read = (read + (MaximumInNodes - int(numIn)) + 1) % genomeLen

	for i := byte(0); i < numOut; i++ {
		read = (read + 1) % genomeLen
		s.outputIndexes[i] = genome[read] % NumTotalNodes
	}

	// skip over rest of (poss redundant) indexes
	read = (read + (MaximumOutNodes - int(numOut)) + 1) % genomeLen

	// generate list of string representations of inputs, alphabetic
	formattedCombinations := make([]string, int(math.Pow(2, float64(numIn))))
	boolAlphabet := InterfaceSlice([]bool{true, false})

	var formatted int = 0
	for combination := range GenerateCombinations(boolAlphabet, int(numIn)) {
		formattedCombinations[formatted] = fmt.Sprint(combination)
		formatted += 1
	}
	sort.Strings(formattedCombinations)

	// read into the transition table
	var outputLikelihoods []byte
	var numOutputCombinations int = int(math.Pow(2, float64(numOut)))
	// for each possible input combination
	for _, sc := range formattedCombinations {
		// for each possibe output combination
		outputLikelihoods = make([]byte, numOutputCombinations)
		for i := 0; i < numOutputCombinations; i++ {
			outputLikelihoods[i] = genome[read]
			read = (read + 1) % genomeLen
		}
		s.transitionTable[sc] = outputLikelihoods
	}

	// cache output possibilities for later
	s.outputCombinationCache = make([][]bool, numOutputCombinations)
	var poss int = 0
	for combination := range GenerateCombinations(boolAlphabet, int(numOut)) {
		// convert back from interface{}
		var c []bool = make([]bool, len(combination))
		for i, _ := range combination {
			c[i] = combination[i].(bool)
		}
		s.outputCombinationCache[poss] = c
		poss += 1
	}
	s.outputCombinationCache = SortedBoolSlice(s.outputCombinationCache)

}

func (s *ProbabilisticGate) Run(inputNodes []bool) []bool {
	input := make([]bool, len(s.inputIndexes))
	for i, _ := range input {
		input[i] = inputNodes[s.inputIndexes[i]]
	}

	likelihoods := s.transitionTable[fmt.Sprint(input)]
	var total int = 0
	for _, v := range likelihoods {
		total += int(v)
	}

	if total <= 0 {
		return make([]bool, NumTotalNodes)
	}

	r := rand.Intn(total)
	cumulative := int(likelihoods[0])
	var choice int = 0
	for r > int(cumulative) {
		choice += 1
		cumulative += int(likelihoods[choice])
	}

	output := make([]bool, NumTotalNodes)
	for i, v := range s.outputIndexes {
		value := s.outputCombinationCache[choice][i]
		output[v] = value
	}

	return output
}
