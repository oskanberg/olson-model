package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Prey struct {
	fitness int
	genome  []byte
	brain   Brain
	pos     Position
	posN    Position
	sensors string
}

func (s *Prey) GetSensors() string {
	return s.sensors
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

func (s *Prey) CanSee(target Agent) (canSee bool, angle float64) {
	differenceVector := target.GetLocation().Subtract(s.GetLocation())
	if differenceVector.Magnitude() > PreyViewDistance {
		// too far away
		return false, 0
	}
	dotProduct := s.GetDirection().Dot(differenceVector.Normalised())
	if dotProduct > AgentViewAngle {
		pdp := s.GetDirection().X*differenceVector.Y - s.GetDirection().Y*differenceVector.X
		if pdp > 0 {
			dotProduct = 1 + (1 - dotProduct)
		}
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
			if b, a := s.CanSee(prey[i]); b {
				// map to correct sensor
				// a is a number from AgentViewAngle to 1 + (1 - AgentViewAngle)
				sliceIndex := int((a - AgentViewAngle) / RetinaSliceWidth)
				sensorValues[sliceIndex] = true
			}
		}
	}
	// read into second set of sensors (predators)
	for i, _ := range predators {
		if b, a := s.CanSee(predators[i]); b {
			// map to correct sensor
			// a is a number from AgentViewAngle to 1 + (1 - AgentViewAngle)
			sliceIndex := int((a-AgentViewAngle)/RetinaSliceWidth) + NumRetinaSlices
			sensorValues[sliceIndex] = true
		}
	}
	s.sensors = fmt.Sprint(sensorValues)
	actuators := s.brain.Run(sensorValues)
	// update positions
	s.updatePosition(actuators)
}

func (s *Prey) Step() {
	s.pos = s.posN
}
