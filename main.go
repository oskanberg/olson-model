package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func GenerateRandomPrey(number int) []*Prey {
	var newPrey *Prey
	var genome []byte
	var prey []*Prey
	for i := 0; i < number; i++ {
		genome = GenerateRandomGenome(100, 10)
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
		genome = GenerateRandomGenome(100, 10)
		newPredator = NewPredator(genome, false)
		predators = append(predators, newPredator)
	}
	return predators
}

func main() {
	// defer profile.Start(profile.CPUProfile).Stop()
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())

	simulation := NewSimulation()
	simulation.RandomPopulation(NumberOfPredators, NumberOfPrey)
	// for i := 0; i < NumberOfPredators; i++ {
	// 	simulation.InsertPredatorFromFile("genomes/startPredator.genome")
	// }

	for generation := 0; generation < TotalGenerations; generation++ {
		fmt.Println("Generation ", generation)
		// simulation.SimulateHomogeneous(TotalSimulationSteps)
		simulation.SimulateHeterogeneous(TotalSimulationSteps)
		simulation.MoranSelectNextGeneration()
	}

}
