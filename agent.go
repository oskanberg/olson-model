package main

import (
	urand "crypto/rand"
	"index/suffixarray"
	"math"
	"math/rand"
	"sync"
)

type Agent interface {
	GetLocation() *Vector2D
	GetDirection() *Vector2D
	SetRandomPosition(int, int)
	// simulate a single time step
	Run(*World)
	// update to the new position
	Step()
}

type Position struct {
	location  Vector2D
	direction Vector2D
}

type Prey struct {
	brain *PLGMN
	pos   Position
	posN  Position
}

func (s *Prey) GetLocation() *Vector2D {
	return &s.pos.location
}

func (s *Prey) GetDirection() *Vector2D {
	return &s.pos.direction
}

func (s *Prey) SetRandomPosition(maxWidth, maxHeight int) {
	newPos := &Position{
		location: Vector2D{
			x: float64(rand.Intn(maxWidth)),
			y: float64(rand.Intn(maxHeight)),
		},
		// TODO: random direction
		direction: Vector2D{
			x: 1,
			y: 0,
		},
	}
	s.pos = *newPos
	// by default, next timestep should be the same
	s.posN = *newPos
}

func (s *Prey) canSee(target Agent) (canSee bool, angle float64) {
	differenceVector := s.GetLocation()
	targetLoc := target.GetLocation()
	differenceVector = differenceVector.Subtract(targetLoc)
	if differenceVector.Magnitude() > 1000 {
		// too far away
		return false, 0
	}
	dotProduct := s.GetDirection().Dot(differenceVector.Normalised())
	if dotProduct > AgentViewAngle {
		return true, dotProduct
	}
	return false, 0
}

func (s *Prey) updatePosition(actuators []bool) {
	action := 0
	for i, v := range actuators {
		if v {
			action += int(math.Pow(2, float64(i)))
		}
	}
	switch action {
	case 0:
		// do nothing
		// fmt.Println("doing nothing")
	case 1:
		// turn right
		// fmt.Println("turning right")
		s.posN.direction = *s.GetDirection().Rotated(PreyTurnAmountRadians)
	case 2:
		// turn left
		// fmt.Println("turning left")
		s.posN.direction = *s.GetDirection().Rotated(-PreyTurnAmountRadians)
	case 3:
		// move straight ahead
		// TODO: implement variable speed
		// fmt.Println("moving straight")
		newLoc := s.GetLocation().Add(s.GetDirection())
		s.posN.location = *newLoc
	}
}

func (s *Prey) Run(w *World) {
	sensorValues := make([]bool, NumRetinaSlices*2)
	// read into first sensors (prey)
	for i, _ := range w.prey {
		// ignore itself
		if w.prey[i] != s {
			if b, a := s.canSee(w.prey[i]); b {
				// map to correct sensor
				// a can be negative
				sliceIndex := int(a/RetinaSliceWidthRadians) + (NumRetinaSlices / 2)
				sensorValues[sliceIndex] = true
			}
		}
	}
	// read into second set of sensors (predators)
	for i, _ := range w.predators {
		if b, a := s.canSee(w.predators[i]); b {
			sliceIndex := int(a/RetinaSliceWidthRadians) + NumRetinaSlices
			sensorValues[sliceIndex] = true
		}
	}
	actuators := s.brain.Run(sensorValues)
	// update positions
	s.updatePosition(actuators)
}

func (s *Prey) Step() {
	s.pos = s.posN
}

func NewPrey(genome []byte) *Prey {
	newPrey := &Prey{
		brain: DeserialiseGenome(genome),
		pos: Position{
			location:  Vector2D{0, 0},
			direction: Vector2D{1, 0},
		},
		posN: Position{
			location:  Vector2D{0, 0},
			direction: Vector2D{1, 0},
		},
	}
	newPrey.SetRandomPosition(500, 500)
	return newPrey
}

func DeserialiseGenome(genome []byte) *PLGMN {
	mn := NewPLGMN()

	index := suffixarray.New(genome)
	genomeStarts := index.Lookup([]byte{42, 213}, -1)

	// add each gate
	for _, start := range genomeStarts {
		mn.NewGateFromGenome(genome, start)
	}

	return mn
}

// Note: this is not guaranteed to create exact number of codons (due to overlap)
func GenerateRandomGenome(length int, artificialStartCodons int) []byte {
	genome := make([]byte, length)
	urand.Read(genome)
	for i := 0; i < artificialStartCodons; i++ {
		position := rand.Intn(length)
		genome[position] = 42
		genome[(position+1)%length] = 213
	}
	return genome
}

func runAgentWG(p Agent, world *World, wg *sync.WaitGroup) {
	p.Run(world)
	wg.Done()
}

func stepAgentWG(p Agent, wg *sync.WaitGroup) {
	p.Step()
	wg.Done()
}
