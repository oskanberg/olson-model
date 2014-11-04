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

func (s *Prey) SetPosition(pos Position) {
	s.pos = &pos
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

func (s *Prey) canSee(target Agent) bool {
	differenceVector := s.GetLocation().Subtract(target.GetLocation())
	if differenceVector.Magnitude() > 100 {
		// too far away
		return false
	}
	dotProduct := s.GetDirection().Dot(differenceVector.Normalised())
	if dotProduct > AgentViewAngle {
		return true
	}
	return false
}

func (s *Prey) Run() {
	// read into sensors
	// run plgmn
	// update position
}

func (s *Prey) Step() {
	s.pos = s.posN
}

func NewPrey(genome []byte) *Prey {
	return &Prey{brain: DeserialiseGenome(genome)}
}
