package main

type World struct {
	prey      []Agent
	predators []Agent
}

func (s *World) AddPrey(p *Prey) {
	s.prey = append(s.prey, p)
}

func (s *World) GenerateRandomPrey(number int) {
	var newPrey *Prey
	var genome []byte
	for i := 0; i < number; i++ {
		genome = GenerateRandomGenome(100, 10)
		newPrey = NewPrey(genome)
		s.AddPrey(newPrey)
	}
}

func NewWorld() *World {
	return &World{}
}
