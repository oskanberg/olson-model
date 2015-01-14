package main

import (
	urand "crypto/rand"
	"math"
)

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
