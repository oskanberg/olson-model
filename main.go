package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"time"
)

func GenerateRandomPrey(number int) []*Prey {
	var newPrey *Prey
	var genome []byte
	var prey []*Prey
	for i := 0; i < number; i++ {
		genome = GenerateRandomGenome(InitialGenomeLength, ArtificialStartCodons)
		newPrey = NewPrey(genome, false)
		prey = append(prey, newPrey)
	}
	return prey
}

func GenerateRandomPredators(number int) []*Predator {
	var newPredator *Predator
	var genome []byte
	var predators []*Predator
	for i := 0; i < number; i++ {
		genome = GenerateRandomGenome(InitialGenomeLength, ArtificialStartCodons)
		newPredator = NewPredator(genome, false)
		predators = append(predators, newPredator)
	}
	return predators
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// defer profile.Start(profile.CPUProfile).Stop()
	runtime.GOMAXPROCS(runtime.NumCPU())

	simulation := NewSimulation()
	if SeedPredators {
		numPredators := int(math.Max(float64(NumberOfPredators-1), 0))
		simulation.RandomPopulation(numPredators, NumberOfPrey)
		// for i := 0; i < NumberOfPredators; i++ {
		simulation.InsertPredatorFromFile("genome/predator/" + Model + "/0.genome")
		// }
	} else {
		simulation.RandomPopulation(NumberOfPredators, NumberOfPrey)
	}

	for generation := 0; generation < TotalGenerations; generation++ {
		fmt.Println("Generation ", generation)
		SimulateHetrogenous(simulation)
		simulation.MoranSelectNextGeneration()
	}
	if SavePredators {
		simulation.SavePredatorGenomes()
	}
	if SavePrey {
		simulation.SavePreyGenomes()
	}
}

func SimulateHetrogenous(s *Simulation) {
	s.Simulate(TotalSimulationSteps, RoundsPerGeneration)
}

func SimulateHomogeneous() {

}
