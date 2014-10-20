package main

import (
	"fmt"
	"index/suffixarray"
	"math"
)

// Node represents a node in the PLGMN
type Node struct {
	t, tNext bool
	id       int
}

func (s *Node) GetState() bool {
	return s.t
}

func (s *Node) SetState(state bool) {
	s.tNext = state
}

func (s *Node) Step() {
	s.t = s.tNext
}

// PLG represents a single gate in the MN
type PLG struct {
	in              []*Node
	out             []*Node
	transitionTable []byte
}

func (s *PLG) AddInNode(node *Node) {
	s.in = append(s.in, node)
}

func (s *PLG) AddOutNode(node *Node) {
	s.out = append(s.out, node)
}

func (s *PLG) ToString() string {
	str := fmt.Sprintf("%d inputs\n", len(s.in))
	for _, node := range s.in {
		str += fmt.Sprintf("in node %d\n", node.id)
	}
	str += fmt.Sprintf("%d outputs\n", len(s.out))
	for _, node := range s.out {
		str += fmt.Sprintf("out node %d\n", node.id)
	}
	return str
}

// PLGMN represents the whole MN of gates
type PLGMN struct {
	gates []*PLG
	nodes []*Node
}

func (s *PLGMN) AddGate(newGate *PLG) {
	s.gates = append(s.gates, newGate)
}

func (s *PLGMN) GenerateNodes(numNodes int) {
	// TODO: this feels very inefficient
	for i := 0; i < numNodes; i++ {
		s.nodes = append(s.nodes, NewNode(i))
	}
}

func (s *PLGMN) ToString() string {
	str := "Begin PLGMN\n-------\n"
	str += fmt.Sprintf("%d gates\n", len(s.gates))
	str += fmt.Sprintf("%d nodes\n", len(s.nodes))
	str += "-------\n"
	for index, gate := range s.gates {
		str += fmt.Sprintf("gate %d\n-------\n", index)
		str += gate.ToString()
	}

	return str
}

func (s *PLGMN) NewPLGFromGenome(genome []byte, startPosition int) {
	plg := NewPLG()

	genomeLen := len(genome)
	numNodes := len(s.nodes)

	// skip the start codon
	read := (startPosition + 2) % genomeLen

	// first byte is num of input nodes
	numberIn := FloorByte(genome[read] / (255 / MaximumInNodes))
	read = (read + 1) % genomeLen

	// next byte is num of output nodes
	numberOut := FloorByte(genome[read] / (255 / MaximumOutNodes))
	read = (read + 1) % genomeLen

	// next numberIn bytes are ids of input nodes
	for i := 0; i < numberIn; i++ {
		nodeIndex := RoundInt((int(genome[read+i]) * numNodes) / 255)
		plg.AddInNode(s.nodes[nodeIndex])
	}
	// there will always be MaximumInNodes spaces, if not vals
	read = (read + MaximumInNodes) % genomeLen

	// next numberOut bytes are ids of out nodes
	for i := 0; i < numberOut; i++ {
		nodeIndex := RoundInt((int(genome[read+i]) * numNodes) / 255)
		plg.AddOutNode(s.nodes[nodeIndex])
	}
	// there will always be MaximumOutNodes spaces, if not vals
	read = (read + MaximumOutNodes) % genomeLen

	// read the probabilities for transition table
	// unfortunately can't do all at once since genome wraps
	numProbabilities := int(math.Pow(2, float64(numberIn+numberOut)))
	end := read + numProbabilities
	if end > genomeLen {
		plg.transitionTable = append(plg.transitionTable, genome[read:genomeLen]...)
		plg.transitionTable = append(plg.transitionTable, genome[0:end-genomeLen]...)
	} else {
		plg.transitionTable = append(plg.transitionTable, genome[read:end]...)
	}
	s.AddGate(plg)
}

func NewNode(index int) *Node {
	return &Node{id: index}
}

func NewPLG() *PLG {
	return &PLG{[]*Node{}, []*Node{}, []byte{}}
}

func NewPLGMN() *PLGMN {
	return &PLGMN{[]*PLG{}, []*Node{}}
}

func DeserialiseGenome(genome []byte) *PLGMN {
	mn := NewPLGMN()

	index := suffixarray.New(genome)
	genomeStarts := index.Lookup([]byte{42, 213}, -1)

	// create all the Nodes
	mn.GenerateNodes(len(genomeStarts))

	// add each gate
	for start := range genomeStarts {
		mn.NewPLGFromGenome(genome, start)
	}

	return mn
}

func main() {
	genome := []byte{0, 3, 4, 163, 243, 53, 7, 4, 2, 42, 213, 163, 95}
	plgmn := DeserialiseGenome(genome)
	fmt.Print(plgmn.ToString())
}
