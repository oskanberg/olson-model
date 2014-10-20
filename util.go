package main

import "math"

func FloorByte(in byte) int {
	return int(math.Floor(float64(in)))
}

func RoundInt(in int) int {
	floatIn := float64(in)
	if in > 0 {
		return int(floatIn + 0.5)
	} else {
		return int(floatIn - 0.5)
	}
}
