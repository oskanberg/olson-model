package main

import (
	urand "crypto/rand"
	"fmt"
	"index/suffixarray"
	"math/rand"
	"time"
)

func DeserialiseGenome(genome []byte) *PLGMN {
	mn := NewPLGMN()

	index := suffixarray.New(genome)
	genomeStarts := index.Lookup([]byte{42, 213}, -1)

	// add each gate
	for _, start := range genomeStarts {
		mn.NewGateFromGenome(genome, start)
	}

	return mn
}

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

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	genome := GenerateRandomGenome(32, 1)
	plgmn := DeserialiseGenome(genome)
	fmt.Print(plgmn.ToString())
	actuators := plgmn.Run([]bool{false, true, true, true, true, true, true, true, true, true, true, true, true})
	fmt.Println(actuators)
}
