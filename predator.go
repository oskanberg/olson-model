package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Predator struct {
	fitness       int
	genome        []byte
	brain         Brain
	pos           Position
	posN          Position
	sensors       string
	nearbyCache   []Agent
	viewCache     []Agent
	timeSinceKill int
}

func (s *Predator) Reset() {
	s.SetRandomPosition(SimulationSpaceSize, SimulationSpaceSize)
	s.brain.Reset()
	s.timeSinceKill = PreyHeadStart
}

func (s *Predator) GetSensors() string {
	return s.sensors
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

func (s *Predator) CanSee(target Agent) (canSee bool, angle float64) {
	differenceVector := target.GetLocation().Subtract(s.GetLocation())
	if differenceVector.Magnitude() > PredatorViewDistance {
		// too far away
		return false, 0
	}
	s.nearbyCache = append(s.nearbyCache, target)
	dotProduct := s.GetDirection().Dot(differenceVector.Normalised())
	if dotProduct > CosHalfAgentViewAngle {
		// sometimes float precision error causes Acos to panic
		if dotProduct > 1 {
			dotProduct = 1
		}
		angle := math.Acos(dotProduct)
		pdp := s.GetDirection().X*differenceVector.Y - s.GetDirection().Y*differenceVector.X
		if pdp > 0 {
			angle = -angle
		}
		return true, angle
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
		newLoc := s.GetLocation().Add(s.GetDirection().Multiplied(PredatorSpeedMultiplier)).Wrap(SimulationSpaceSize, SimulationSpaceSize)
		s.posN.Location = *newLoc
	}
}

func (s *Predator) Run(prey []*Prey, predators []*Predator) {
	s.viewCache = nil
	s.nearbyCache = nil
	s.timeSinceKill += 1
	sensorValues := make([]bool, NumRetinaSlices*2)
	// read into first sensors (prey)
	for i, _ := range prey {
		if b, a := s.CanSee(prey[i]); b {
			s.viewCache = append(s.viewCache, prey[i])
			// map to correct sensor
			// a is a number from -(AgentViewAngleRadians/2) to AgentViewAngleRadians/2
			sliceIndex := int((a + HalfAgentViewAngleRadians) / RetinaSliceWidth)
			sensorValues[sliceIndex] = true
		}
	}
	// read into second set of sensors (predators)
	for i, _ := range predators {
		if predators[i] != s {
			if b, a := s.CanSee(predators[i]); b {
				// map to correct sensor
				//a is a number from -(AgentViewAngleRadians/2) to AgentViewAngleRadians/2
				sliceIndex := int((a+HalfAgentViewAngleRadians)/RetinaSliceWidth) + NumRetinaSlices
				sensorValues[sliceIndex] = true
			}
		}
	}
	s.sensors = fmt.Sprint(sensorValues)
	actuators := s.brain.Run(sensorValues)
	// update positions
	s.updatePosition(actuators)
}

func (s *Predator) Step() {
	s.pos = s.posN
}

// func (s *Predator) PrintStatistics() {
// 	s.brain.PrintStatistics()
// }

func (s *Predator) ToString() string {
	return s.brain.ToString()
}
