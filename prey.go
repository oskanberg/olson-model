package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Prey struct {
	fitness    int
	genome     []byte
	brain      Brain
	pos        Position
	posN       Position
	sensors    string
	nearbyPrey int
}

func (s *Prey) Reset() {
	s.SetRandomPosition(SimulationSpaceSize, SimulationSpaceSize)
	s.brain.Reset()
	s.nearbyPrey = 0
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

func (s *Prey) CanSee(target Agent) (canSee bool, angle float64, distance float64) {
	differenceVector := target.GetLocation().Subtract(s.GetLocation())
	if differenceVector.Magnitude() > PreyViewDistance {
		// too far away
		return false, 0, 0
	}
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
		return true, angle, differenceVector.Magnitude()
	}
	return false, 0, 0
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
			if b, a, d := s.CanSee(prey[i]); b {
				if d <= 30 {
					s.nearbyPrey++
				}
				// map to correct sensor
				sliceIndex := int((a + HalfAgentViewAngleRadians) / RetinaSliceWidth)
				sensorValues[sliceIndex] = true
			}
		}
	}
	// read into second set of sensors (predators)
	for i, _ := range predators {
		if b, a, _ := s.CanSee(predators[i]); b {
			// map to correct sensor
			sliceIndex := int((a+HalfAgentViewAngleRadians)/RetinaSliceWidth) + NumRetinaSlices
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
