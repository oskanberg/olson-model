package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
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
		if Method == "Homogenous" {
			SimulateHomogeneous(simulation)
		} else if Method == "Hetrogenous" {
			SimulateHetrogenous(simulation)
		}
	}
	if SavePredators {
		simulation.SavePredatorGenomes()
	}
	if SavePrey {
		simulation.SavePreyGenomes()
	}
}

func SimulateHetrogenous(s *Simulation) {
	for i := 0; i < RoundsPerGeneration; i++ {
		s.Simulate(TotalSimulationSteps)
		s.ResetPopulation()
		AppendRecordFloat([]float64{s.meanNearbyPrey}, "meannearby.csv")
		fmt.Println("Mean nearby prey:", s.meanNearbyPrey)
		s.meanNearbyPrey = 0
	}
	s.MoranSelectNextGeneration()
}

func SimulateHomogeneous(s *Simulation) {
	var wg sync.WaitGroup
	var predatorGenome []byte
	var preyGenome []byte

	var meanNearbyPrey float64

	for _, prd := range s.predators {
		predatorGenome = prd.GetGenome()
		wg.Add(1)
		go func() {
			simulation := NewSimulation()
			for i := 0; i < NumberOfPredatorClones; i++ {
				simulation.AddPredatorFromGenome(predatorGenome)
			}

			for _, pry := range s.prey {
				preyGenome = pry.GetGenome()
				for i := 0; i < NumberOfPreyClones; i++ {
					simulation.AddPreyFromGenome(preyGenome)
				}

				for i := 0; i < RoundsPerGeneration; i++ {
					simulation.Simulate(TotalSimulationSteps)
					simulation.ResetPopulation()
				}

				// copy prey fitness back
				var newFitness int
				for _, sPry := range simulation.GetPrey() {
					newFitness += sPry.GetFitness()
				}
				newFitness = newFitness / NumberOfPreyClones
				// fmt.Println("prey updating with fitness", newFitness)
				pry.fitness = newFitness

				newFitness = 0
				for _, sPrd := range simulation.GetPredators() {
					newFitness += sPrd.GetFitness()
				}
				newFitness = int(math.Ceil(float64(newFitness) / float64(NumberOfPredatorClones)))
				// fmt.Println("predator updating with fitness", newFitness)
				// += to other prey adversaries
				prd.fitness += newFitness
				meanNearbyPrey += simulation.meanNearbyPrey
			}
			// prd.fitness /= len(s.prey)
			wg.Done()
		}()
	}

	wg.Wait()

	meanNearbyPrey /= float64(len(s.prey) * len(s.predators))
	fmt.Println("Mean nearby prey", meanNearbyPrey)

	s.MoranSelectNextGeneration()
}
