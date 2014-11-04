package main

import "math/rand"

type Agent interface {
	GetLocation() *Vector2D
	GetDirection() *Vector2D
	SetPosition()
	SetRandomPosition()
	// simulate a single time step
	Run()
	// update to the new position
	Step()
}

type Position struct {
	location  Vector2D
	direction Vector2D
}

type Prey struct {
	brain *PLGMN
	pos   *Position
	posN  *Position
}

func (s *Prey) SetLocation(loc *Vector2D) {
	s.pos.location = *loc
}

func (s *Prey) SetDirection(direction *Vector2D) {
	s.pos.location = *direction
}

func (s *Prey) GetLocation() *Vector2D {
	return &s.pos.location
}

func (s *Prey) GetDirection() *Vector2D {
	return &s.pos.direction
}

func (s *Prey) SetRandomPosition(maxWidth, maxHeight int) {
	newPos := Position{
		location: Vector2D{
			x: float64(rand.Intn(maxWidth)),
			y: float64(rand.Intn(maxHeight)),
		},
		// TODO: random direction
		direction: Vector2D{
			x: 0,
			y: 0,
		},
	}
	s.SetPosition(newPos)
}

func (s *Prey) canSee(target Agent) (canSee bool, angle float64) {
	differenceVector := s.GetLocation().Subtract(target.GetLocation())
	if differenceVector.Magnitude() > 100 {
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
	action := byte(0)
	for _, v := range actuators {
		action = action << 1
		if v {
			action = action & 1
		}
	}
	switch action {
	case 0:
		// do nothing
	case 1:
		// turn right
		s.GetDirection().Rotate(PreyTurnAmountRadians)
	case 2:
		// turn left
		s.GetDirection().Rotate(-PreyTurnAmountRadians)
	case 3:
		// move straight ahead
		// TODO: implement variable speed
		newPos := s.GetLocation().Add(s.GetDirection())
		s.SetPosition(newPos)
	}
}

func (s *Prey) Run(w *World) {
	sensorValues := make([]bool, NumRetinaSlices*2)
	// read into first sensors (prey)
	for _, agent := range w.prey {
		if b, a := s.canSee(agent); b {
			sliceIndex := int(a / RetinaSliceWidthRadians)
			sensorValues[sliceIndex] = true
		}
	}
	// read into second set of sensors (predators)
	for _, agent := range w.predators {
		if b, a := s.canSee(agent); b {
			sliceIndex := int(a/RetinaSliceWidthRadians) + NumRetinaSlices
			sensorValues[sliceIndex] = true
		}
	}
	actuators := s.brain.Run(sensorValues)
	// update position
	s.updatePosition(actuators)
}

func (s *Prey) Step() {
	s.pos = s.posN
}

func NewPrey(genome []byte) *Prey {
	return &Prey{brain: DeserialiseGenome(genome)}
}
