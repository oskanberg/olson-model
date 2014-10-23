package main

import (
	"fmt"
	"math"
	"math/rand"
)

// PLG represents a single gate in the MN
type PLG struct {
	in              []*Node
	out             []*Node
	transitionTable []float64
}

func (s *PLG) AddInNode(node *Node) {
	s.in = append(s.in, node)
}

func (s *PLG) AddOutNode(node *Node) {
	s.out = append(s.out, node)
}

func (s *PLG) Run() {
	var inputState byte = 0
	for i, _ := range s.in {
		inputState = inputState << 1
		inputState = inputState | Booltobyte(s.in[i].GetState())
	}

	rowLength := byte(math.Pow(2, float64(len(s.out))))
	start := inputState * rowLength
	row := s.transitionTable[start : start+rowLength]

	point := rand.Float64()
	total := float64(0)
	index := -1

	// keep iterating until total is >= rand number
	for total < point {
		index += 1
		total += row[index]
	}

	// index now represents outnodes in binary
	for i, _ := range s.out {
		// if lsb is 1, set state true
		s.out[i].SetState(index&1 == 1)
		index = index >> 1
	}

}

func (s *PLG) NormaliseTransitionTable(numberIn int, numberOut int) {
	rowLength := int(math.Pow(2, float64(len(s.out))))

	// for each row
	for i, rows := float64(0), math.Pow(2, float64(numberIn)); i < rows; i++ {
		index := int(i) * rowLength
		row := s.transitionTable[index : index+int(math.Pow(2, float64(numberOut)))]
		var rowTotal float64 = 0
		for _, val := range row {
			rowTotal += (val + 1)
		}
		for i, val := range row {
			row[i] = float64((val + 1) / (rowTotal))
		}
	}
}

func (s *PLG) ToString() string {
	str := fmt.Sprintf("%d inputs\n", len(s.in))
	for _, node := range s.in {
		str += fmt.Sprintf("in node %d\n", node.GetId())
	}
	str += fmt.Sprintf("%d outputs\n", len(s.out))
	for _, node := range s.out {
		str += fmt.Sprintf("out node %d\n", node.GetId())
	}

	numberIn := len(s.in)
	numberOut := len(s.out)
	str += "table:\n"
	// for each row
	for i, rows := float64(0), math.Pow(2, float64(numberIn)); i < rows; i++ {
		index := int(i) * numberIn
		row := s.transitionTable[index : index+int(math.Pow(2, float64(numberOut)))]
		str += fmt.Sprintln(row)
	}
	return str
}

func NewPLG() *PLG {
	return &PLG{[]*Node{}, []*Node{}, []float64{}}
}
