package main

import (
	urand "crypto/rand"
	"index/suffixarray"
	"math/rand"
	"sync"
)

type Brain interface {
	Run([]bool) []bool
	ToString() string
	Reset()
	// PrintStatistics
}

type Agent interface {
	GetSensors() string
	GetFitness() int
	GetGenome() []byte
	GetLocation() *Vector2D
	GetDirection() *Vector2D
	SetRandomPosition(int, int)
	CanSee(Agent) (bool, float64, float64)
	// simulate a single time step
	Run([]*Prey, []*Predator)
	// update to the new position
	Step()
}

type Position struct {
	Location  Vector2D
	Direction Vector2D
}

func Mutate(genome []byte) []byte {
	for i, _ := range genome {
		if rand.Float64() < MutationRate {
			genome[i] = RandByte()
		}
	}
	if len(genome) < 20000 && rand.Float64() < DuplicationLikelihood {
		genomeLen := len(genome)
		start := rand.Intn(genomeLen)
		end := rand.Intn(genomeLen-start) + start
		insert := rand.Intn(genomeLen)
		// insert
		// fmt.Println("Copying section")
		suffix := append(genome[start:end], genome[insert:]...)
		genome = append(genome[:insert], suffix...)
	}
	if len(genome) > 1000 && rand.Float64() < DeletionLikelihood {
		genomeLen := len(genome)
		start := rand.Intn(genomeLen)
		end := rand.Intn(genomeLen-start) + start
		// delete
		// fmt.Println("Deleting section")
		genome = append(genome[:start], genome[end:]...)
	}
	return genome
}

func NewPrey(genomeO []byte, mutate bool) *Prey {
	genome := make([]byte, len(genomeO))
	copy(genome, genomeO)
	if mutate {
		genome = Mutate(genome)
	}
	newPrey := &Prey{
		fitness: 0,
		genome:  genome,
		pos: Position{
			Location:  Vector2D{0, 0},
			Direction: Vector2D{1, 0},
		},
		posN: Position{
			Location:  Vector2D{0, 0},
			Direction: Vector2D{1, 0},
		},
		nearbyPrey: 0,
	}
	if Model == "MarkovNetwork" {
		newPrey.brain = DeserialiseGenomeMarkovNetwork(genome)
	} else if Model == "Olson" {
		newPrey.brain = DeserialiseGenomePLGMN(genome)
	} else if Model == "LinearWeights" {
		newPrey.brain = DeserialiseGenomeLinearWeights(genome)
	}
	newPrey.SetRandomPosition(SimulationSpaceSize, SimulationSpaceSize)
	return newPrey
}

func NewPredator(genomeO []byte, mutate bool) *Predator {
	genome := make([]byte, len(genomeO))
	copy(genome, genomeO)
	if mutate {
		genome = Mutate(genome)
	}
	newPredator := &Predator{
		fitness: 0,
		genome:  genome,
		pos: Position{
			Location:  Vector2D{0, 0},
			Direction: Vector2D{1, 0},
		},
		posN: Position{
			Location:  Vector2D{0, 0},
			Direction: Vector2D{1, 0},
		},
		nearbyCache:   make([]Agent, 0),
		timeSinceKill: -PreyHeadStart,
	}
	if Model == "MarkovNetwork" {
		newPredator.brain = DeserialiseGenomeMarkovNetwork(genome)
	} else if Model == "Olson" {
		newPredator.brain = DeserialiseGenomePLGMN(genome)
	} else if Model == "LinearWeights" {
		newPredator.brain = DeserialiseGenomeLinearWeights(genome)
	}
	newPredator.SetRandomPosition(SimulationSpaceSize, SimulationSpaceSize)
	return newPredator
}

func DeserialiseGenomePLGMN(genome []byte) *PLGMN {
	mn := NewPLGMN()

	index := suffixarray.New(genome)
	genomeStarts := index.Lookup([]byte{42, 213}, -1)

	// add each gate
	for _, start := range genomeStarts {
		mn.NewGateFromGenome(genome, start)
	}

	return mn
}

func DeserialiseGenomeLinearWeights(genome []byte) *LinearWeights {
	lw := NewLinearWeights()

	index := suffixarray.New(genome)
	genomeStarts := index.Lookup([]byte{42, 213}, -1)

	// add each gate
	for _, start := range genomeStarts {
		lw.AddWeightsFromGenome(genome, start)
	}

	return lw
}

func DeserialiseGenomeMarkovNetwork(genome []byte) *MarkovNetwork {
	mn := NewMarkovNetwork()

	index := suffixarray.New(genome)
	genomeStarts := index.Lookup([]byte{42, 213}, -1)

	// add each gate
	for _, start := range genomeStarts {
		mn.AddGateFromGenome(genome, start)
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

func runAgentWG(p Agent, prey []*Prey, predators []*Predator, wg *sync.WaitGroup) {
	p.Run(prey, predators)
	wg.Done()
}

func stepAgentWG(p Agent, wg *sync.WaitGroup) {
	p.Step()
	wg.Done()
}
