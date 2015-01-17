package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync"
)

var numRuns int = 0
var record *Record = NewRecord()

type Simulation struct {
	prey      []*Prey
	predators []*Predator
	dead      []*Prey
}

func (s *Simulation) RandomPopulation(numPred, numPrey int) {
	var wg sync.WaitGroup

	s.prey = make([]*Prey, numPrey)
	s.predators = make([]*Predator, numPred)

	for i := 0; i < numPrey; i++ {
		wg.Add(1)
		go func(id int) {
			genome := GenerateRandomGenome(InitialGenomeLength, ArtificialStartCodons)
			s.prey[id] = NewPrey(genome, false)
			wg.Done()
		}(i)
	}

	for j := 0; j < numPred; j++ {
		wg.Add(1)
		go func(id int) {
			genome := GenerateRandomGenome(InitialGenomeLength, ArtificialStartCodons)
			s.predators[id] = NewPredator(genome, false)
			wg.Done()
		}(j)
	}

	wg.Wait()
}

func (s *Simulation) InsertPredatorFromFile(filename string) {
	csvFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Failed to open genome file. ", err)
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1
	raw, err := reader.ReadAll()
	csvData := raw[0]
	if err != nil {
		fmt.Println("Failed to read CSV data. ", err)
	}
	genome := make([]byte, len(csvData))
	for i, d := range csvData {
		geneInt, err := strconv.Atoi(d)
		if err != nil {
			fmt.Println("Failed to parse gene to int", err)
		}
		genome[i] = byte(geneInt)
	}
	// fmt.Println(genome)

	s.predators = append(s.predators, NewPredator(genome, false))
}

func (s *Simulation) SavePredatorGenomes() {
	for prd, _ := range s.predators {
		genome := s.predators[prd].GetGenome()
		strGenome := make([]string, len(genome))
		for i, _ := range genome {
			strGenome[i] = strconv.Itoa(int(genome[i]))
		}

		csvFile, err := os.Create("genome/" + strconv.Itoa(prd) + ".csv")
		if err != nil {
			fmt.Println(err)
		}
		defer csvFile.Close()

		writer := csv.NewWriter(csvFile)
		err = writer.Write(strGenome)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		writer.Flush()
	}
}

// simulate the current simulation clonally for given steps
func (s *Simulation) SimulateHomogeneous(iterations int) {
	// TODO
}

// simulate the current simulation naturally for given steps
func (s *Simulation) SimulateHeterogeneous(iterations int) {
	var wg sync.WaitGroup
	var total int
	for iteration := 0; iteration < iterations; iteration++ {
		total = len(s.prey) + len(s.predators)
		wg.Add(total)
		for i, _ := range s.predators {
			go runAgentWG(s.predators[i], s.prey, s.predators, &wg)
		}
		for i, _ := range s.prey {
			go runAgentWG(s.prey[i], s.prey, s.predators, &wg)
		}
		wg.Wait()

		wg.Add(total)
		for i, _ := range s.predators {
			go stepAgentWG(s.predators[i], &wg)
		}
		for i, _ := range s.prey {
			go stepAgentWG(s.prey[i], &wg)
		}
		wg.Wait()
		s.processDeaths()
		s.RecordCurrentPositions()
		record.NewStep()
	}
	// record.WriteToFile(strconv.Itoa(numRuns))
	record.WriteToFile("1")
	// clear for next run
	record = NewRecord()
	numRuns += 1
}

func (s *Simulation) RecordCurrentPositions() {
	for i, _ := range s.predators {
		record.AddRecordToCurrentStep(s.predators[i])
	}
	for i, _ := range s.prey {
		record.AddRecordToCurrentStep(s.prey[i])
	}
}

func (s *Simulation) processDeaths() {
	deaths := make(map[int]bool)

	//TODO: implement quadtree? utilise view logic?
	for prd, _ := range s.predators {
		// handling time
		if EatCooldown && s.predators[prd].timeSinceKill < 10 {
			continue
		}
		// only check those nearby
		for pry, _ := range s.predators[prd].nearbyCache {
			preyLoc := s.prey[pry].GetLocation()
			distance := s.predators[prd].GetLocation().Subtract(preyLoc).Magnitude()
			// fmt.Printf("Distance from predator: %f \n", distance)
			if distance <= EatingDistance {
				if PredatorConfusion {
					denominator := math.Max(float64(len(s.predators[prd].viewCache)), 1)
					likelihood := 1 / denominator
					if rand.Float64() < likelihood {
						s.predators[prd].fitness += 1
						deaths[pry] = true
						s.predators[prd].timeSinceKill = 0
					}
				} else {
					s.predators[prd].fitness += 1
					s.predators[prd].timeSinceKill = 0
					deaths[pry] = true
				}
			}
		}
	}

	// get keys so we can sort
	keys := make([]int, len(deaths))
	i := 0
	for k, _ := range deaths {
		keys[i] = k
		i++
	}
	// descending order
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	for _, key := range keys {
		// add prey to dead list
		s.dead = append(s.dead, s.prey[key])
		// remove prey from slice
		plen := len(s.prey)
		s.prey[key] = s.prey[plen-1]
		s.prey = s.prey[:plen-1]
	}
}

func (s *Simulation) MoranSelectNextGeneration() {
	// add dead agents back to the population for eval
	s.prey = append(s.prey, s.dead...)
	// clear dead
	s.dead = nil
	s.prey = s.getMoranPreyGeneration()
	s.predators = s.getMoranPredatorGeneration()
}

func (s *Simulation) getMoranPreyGeneration() []*Prey {
	var totalPreyFitness float64 = math.SmallestNonzeroFloat64
	var highestFitness float64 = 0
	var fitness float64
	for i, _ := range s.prey {
		fitness = float64(s.prey[i].GetFitness())
		totalPreyFitness += fitness
		if fitness > highestFitness {
			highestFitness = fitness
		}
	}

	fmt.Printf("Total prey fitness: %f\n", totalPreyFitness)

	selectUniformly := false
	if int(totalPreyFitness) == 0 {
		fmt.Println("Selecting prey uniformly")
		selectUniformly = true
	}

	preyPopulationSize := len(s.prey)

	var genome []byte
	newPreyGeneration := make([]*Prey, preyPopulationSize)
	newPrey := 0
	for newPrey < preyPopulationSize {
		index := rand.Intn(preyPopulationSize)
		normFitness := float64(s.prey[index].GetFitness()) / highestFitness
		if rand.Float64() < normFitness || selectUniformly {
			genome = s.prey[index].GetGenome()
			newPreyGeneration[newPrey] = NewPrey(genome, true)
			newPrey += 1
		}
	}

	return newPreyGeneration
}

func (s *Simulation) getMoranPredatorGeneration() []*Predator {
	var totalPredatorFitness float64 = math.SmallestNonzeroFloat64
	var highestFitness float64 = 0
	var fitness float64
	for i, _ := range s.predators {
		// fmt.Printf("Predator %d fitness:\t%d\n", i, s.predators[i].GetFitness())
		fitness = float64(s.predators[i].GetFitness())
		totalPredatorFitness += fitness
		if fitness > highestFitness {
			highestFitness = fitness
		}
	}

	fmt.Printf("Total predator fitness: %f\n", totalPredatorFitness)

	selectUniformly := false
	if int(totalPredatorFitness) == 0 {
		fmt.Println("Selecting predators uniformly")
		selectUniformly = true
	}

	predPopulationSize := len(s.predators)

	var genome []byte
	newPredGeneration := make([]*Predator, predPopulationSize)
	newPredators := 0
	for newPredators < predPopulationSize {
		index := rand.Intn(predPopulationSize)
		normFitness := float64(s.predators[index].GetFitness()) / highestFitness
		if rand.Float64() < normFitness || selectUniformly {
			genome = s.predators[index].GetGenome()
			// fmt.Printf("Selecting agent with fitness %d\n", s.predators[index].GetFitness())
			newPredGeneration[newPredators] = NewPredator(genome, true)
			newPredators += 1
		}
	}

	// randomly introduce completely random genome
	// if rand.Float64() < 0.5 {
	// 	index := rand.Intn(predPopulationSize)
	// 	genome := GenerateRandomGenome(InitialGenomeLength, ArtificialStartCodons)
	// 	newPredGeneration[index] = NewPredator(genome, false)
	// }

	return newPredGeneration
}

func NewSimulation() *Simulation {
	return &Simulation{}
}
