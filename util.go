package main

import (
	urand "crypto/rand"
	"errors"
	"fmt"
	"math"
	"sort"
)

type SortableBoolSliceSlice struct {
	data       [][]bool
	serialised sort.StringSlice
}

func (s SortableBoolSliceSlice) Len() int           { return len(s.data) }
func (s SortableBoolSliceSlice) Less(i, j int) bool { return s.serialised.Less(i, j) }
func (s SortableBoolSliceSlice) Swap(i, j int) {
	s.serialised.Swap(i, j)
	s.data[i], s.data[j] = s.data[j], s.data[i]
}

func (s SortableBoolSliceSlice) initialiseStrings() {
	for i, _ := range s.data {
		s.serialised[i] = fmt.Sprint(s.data[i])
	}
}

func SortedBoolSlice(s [][]bool) [][]bool {
	sortable := SortableBoolSliceSlice{
		data:       s,
		serialised: sort.StringSlice(make([]string, len(s))),
	}
	sortable.initialiseStrings()
	sort.Sort(sortable)
	return sortable.data
}

func LogicalOr(first []bool, second []bool) (error, []bool) {
	if len(first) != len(second) {
		fmt.Println(first, second)
		return errors.New("Both bool slices must be the same length"), []bool{}
	}
	output := make([]bool, len(first))
	for i, _ := range first {
		if first[i] || second[i] {
			output[i] = true
		}
	}
	return nil, output
}

func FloorByte(in byte) int {
	return int(math.Floor(float64(in)))
}

func RoundInt(in float64) int {
	if in > 0 {
		return int(in + 0.5)
	} else {
		return int(in - 0.5)
	}
}

func BoolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

var randBuffer []byte = make([]byte, 100000)
var bufRead int = -1

func RandByte() byte {
	if bufRead == -1 || bufRead == 100000 {
		urand.Read(randBuffer)
		bufRead = 0
	}
	randomByte := randBuffer[bufRead]
	bufRead += 1
	return randomByte
}
