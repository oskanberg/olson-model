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
		inputState = inputState | BoolToByte(s.in[i].GetState())
	}

	// fmt.Println("input state")
	// fmt.Println(inputState)

	rowLength := int(math.Pow(2, float64(len(s.out))))
	start := int(inputState) * rowLength
	row := s.transitionTable[start : start+rowLength]

	point := rand.Float64()
	total := float64(0)
	index := -1

	// keep iterating until total is >= rand number
	for total < point {
		index += 1
		total += row[index]
	}

	// fmt.Println("out state")
	// fmt.Println(index)

	// index now represents outnodes in binary
	for i, _ := range s.out {
		// if lsb is 1, set state true
		if index&1 == 1 {
			// fmt.Println("writing true to ", s.out[i].GetId())
			s.out[i].SetState(true)
		} else {
			s.out[i].SetState(false)
		}
		index = index >> 1
	}

}

func (s *PLG) NormaliseTransitionTable() {
	// row length = num output conditions = 2 ^ num outnodes
	rowLength := int(math.Pow(2, float64(len(s.out))))
	numRows := len(s.transitionTable) / rowLength

	// for each row
	for i := 0; i < numRows; i++ {
		index := i * rowLength
		row := s.transitionTable[index : index+rowLength]
		var rowTotal float64 = 0
		for _, val := range row {
			rowTotal += val
		}
		for i, val := range row {
			row[i] = float64((val) / (rowTotal))
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

	rowLen := int(math.Pow(2, float64(len(s.out))))
	numRows := int(math.Pow(2, float64(len(s.in))))

	str += "table:\n"
	// for each row
	for i := 0; i < numRows; i++ {
		row := s.transitionTable[i*rowLen : (i+1)*rowLen]
		str += fmt.Sprintln(row)
	}

	str += "\n"
	str += fmt.Sprintln(s.transitionTable)
	str += "\n"
	return str
}

func NewPLG() *PLG {
	return &PLG{[]*Node{}, []*Node{}, []float64{}}
}
