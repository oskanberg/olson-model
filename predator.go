package main

import (
	"math"
	"math/rand"
)

type Predator struct {
	fitness int
	genome  []byte
	brain   *PLGMN
	pos     Position
	posN    Position
}

func (s *Predator) GetGenome() []byte {
	return s.genome
}

func (s *Predator) GetFitness() int {
	return s.fitness
}

func (s *Predator) GetLocation() *Vector2D {
	return &s.pos.Location
}

func (s *Predator) GetDirection() *Vector2D {
	return &s.pos.Direction
}

func (s *Predator) SetRandomPosition(maxWidth, maxHeight int) {
	newPos := &Position{
		Location: Vector2D{
			X: float64(rand.Intn(maxWidth)),
			Y: float64(rand.Intn(maxHeight)),
		},
		Direction: *NewRandomUnitVector(),
	}
	s.pos = *newPos
	// by default, next timestep should be the same
	s.posN = *newPos
}

func (s *Predator) canSee(target Agent) (canSee bool, angle float64) {
	differenceVector := s.GetLocation().Subtract(target.GetLocation())
	if differenceVector.Magnitude() > PredatorViewDistance {
		// too far away
		return false, 0
	}
	dotProduct := s.GetDirection().Dot(differenceVector.Normalised())
	if dotProduct > AgentViewAngle {
		return true, dotProduct
	}
	return false, 0
}

func (s *Predator) updatePosition(actuators []bool) {
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
		s.posN.Direction = *s.GetDirection().Rotated(PredatorTurnAmountRadians)
	case 2:
		// turn left
		// fmt.Println("turning left")
		s.posN.Direction = *s.GetDirection().Rotated(-PredatorTurnAmountRadians)
	case 3:
		// move straight ahead
		// fmt.Println("moving straight")
		newLoc := s.GetLocation().Add(s.GetDirection()).Wrap(SimulationSpaceSize, SimulationSpaceSize)
		s.posN.Location = *newLoc
	}
}

func (s *Predator) Run(prey []*Prey, predators []*Predator) {

	sensorValues := make([]bool, NumRetinaSlices*2)
	// read into first sensors (prey)
	for i, _ := range prey {
		if b, a := s.canSee(prey[i]); b {
			// map to correct sensor
			// a can be negative

			sliceIndex := int(a/RetinaSliceWidthRadians) + (NumRetinaSlices / 2)
			sensorValues[sliceIndex] = true
		}
	}
	actuators := s.brain.Run(sensorValues)
	// fmt.Println(actuators)
	// update positions
	s.updatePosition(actuators)
}

func (s *Predator) Step() {
	s.pos = s.posN
}
