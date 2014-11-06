package main

import "math"

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

func Booltobyte(b bool) byte {
	if b {
		return 1
	}
	return 0
}
