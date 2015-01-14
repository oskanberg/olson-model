package main

import (
	"math"
	"math/rand"
)

type Prey struct {
	fitness int
	genome  []byte
	brain   *PLGMN
	pos     Position
	posN    Position
}

func (s *Prey) GetGenome() []byte {
	return s.genome
}

func (s *Prey) GetFitness() int {
	return s.fitness
}

func (s *Prey) GetLocation() *Vector2D {
	return &s.pos.Location
}

func (s *Prey) GetDirection() *Vector2D {
	return &s.pos.Direction
}

func (s *Prey) SetRandomPosition(maxWidth, maxHeight int) {
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

func (s *Prey) canSee(target Agent) (canSee bool, angle float64) {
	differenceVector := s.GetLocation().Subtract(target.GetLocation())
	if differenceVector.Magnitude() > PreyViewDistance {
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
		s.posN.Direction = *s.GetDirection().Rotated(PreyTurnAmountRadians)
	case 2:
		// turn left
		// fmt.Println("turning left")
		s.posN.Direction = *s.GetDirection().Rotated(-PreyTurnAmountRadians)
	case 3:
		// move straight ahead
		// fmt.Println("moving straight")
		newLoc := s.GetLocation().Add(s.GetDirection()).Wrap(SimulationSpaceSize, SimulationSpaceSize)
		s.posN.Location = *newLoc
	}
}

func (s *Prey) Run(prey []*Prey, predators []*Predator) {
	// fmt.Println(s)
	// update fitness
	s.fitness += 1

	sensorValues := make([]bool, NumRetinaSlices*2)
	// read into first sensors (prey)
	for i, _ := range prey {
		// ignore itself
		if prey[i] != s {
			if b, a := s.canSee(prey[i]); b {
				// map to correct sensor
				// a can be negative

				sliceIndex := int(a/RetinaSliceWidthRadians) + (NumRetinaSlices / 2)
				sensorValues[sliceIndex] = true
			}
		}
	}
	// read into second set of sensors (predators)
	for i, _ := range predators {
		if b, a := s.canSee(predators[i]); b {
			sliceIndex := int(a/RetinaSliceWidthRadians) + (NumRetinaSlices / 2) + NumRetinaSlices
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
