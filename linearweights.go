package main

import (
	// "fmt"
	"math"
	"math/rand"
)

// PLGMN represents the whole MN of gates
type LinearWeights struct {
	noopWeights    []byte
	leftWeights    []byte
	rightWeights   []byte
	forwardWeights []byte
}

func (s *LinearWeights) Reset() {
	// nothing to reset
}

func (s *LinearWeights) calculateWeights(weights []byte, sensors []bool) int {
	numSensors := len(sensors)
	// allow overflow
	var total byte = 0
	for i := 0; i < numSensors; i++ {
		if sensors[i] {
			total += weights[i]
		}
	}
	return int(total)
}

func (s *LinearWeights) Run(sensorValues []bool) []bool {
	n := s.calculateWeights(s.noopWeights, sensorValues)
	l := s.calculateWeights(s.leftWeights, sensorValues)
	r := s.calculateWeights(s.rightWeights, sensorValues)
	fw := s.calculateWeights(s.forwardWeights, sensorValues)

	total := n + l + r + fw

	if total == 0 {
		// compare final weights of each
		n = int(s.noopWeights[len(s.noopWeights)-1])
		l = int(s.leftWeights[len(s.leftWeights)-1])
		r = int(s.rightWeights[len(s.rightWeights)-1])
		fw = int(s.forwardWeights[len(s.forwardWeights)-1])
		total = int(math.Max(float64(n+l+r+fw), 1))
	}

	selection := rand.Intn(total)

	if selection < r {
		// right
		return []bool{true, false}
	} else if selection < r+l {
		// left
		return []bool{false, true}
	} else if selection < r+l+fw {
		// forward
		return []bool{true, true}
	} else {
		// noop
		return []bool{false, false}
	}
}

//TODO
func (s *LinearWeights) ToString() string {
	return ""
}

func (s *LinearWeights) AddWeightsFromGenome(genome []byte, start int) {
	genomeLen := len(genome)
	numWeights := len(s.leftWeights)
	// skip the start codon
	read := (start + 2) % genomeLen
	for i := 0; i < numWeights; i++ {
		s.leftWeights[i] += genome[(read+i)%genomeLen]
		s.rightWeights[i] += genome[(read+numWeights+i)%genomeLen]
		s.forwardWeights[i] += genome[(read+(numWeights*2)+i)%genomeLen]
	}
}

func NewLinearWeights() *LinearWeights {
	retinaSlices := NumRetinaSlices * 2
	lw := &LinearWeights{
		noopWeights:    make([]byte, retinaSlices+1),
		leftWeights:    make([]byte, retinaSlices+1),
		rightWeights:   make([]byte, retinaSlices+1),
		forwardWeights: make([]byte, retinaSlices+1),
	}
	// initialise to zero so we can just add weights
	for i := 0; i < retinaSlices+1; i++ {
		lw.noopWeights[i] = 0
		lw.leftWeights[i] = 0
		lw.rightWeights[i] = 0
		lw.forwardWeights[i] = 0
	}
	return lw
}
