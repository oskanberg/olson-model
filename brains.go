package main

import (
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
		mn.NewPLGFromGenome(genome, start)
	}

	return mn
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	genome := []byte{250, 100, 4, 163, 42, 213, 7, 4, 2, 42, 213, 163, 95}
	plgmn := DeserialiseGenome(genome)
	fmt.Print(plgmn.ToString())
	plgmn.Run([]bool{false, true, true, true, true, true, true, true, true, true, true, true, true})
}
