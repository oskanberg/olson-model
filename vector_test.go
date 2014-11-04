package main

import (
	"math"
	"testing"
)

func TestRotate(t *testing.T) {
	vector := Vector2D{
		x: 1.0,
		y: 0,
	}
	vector.Rotate(math.Pi)
	// some floating point errors
	if vector.x != -1 || vector.y > 0.000000001 || vector.y < -0.000000000001 {
		t.Errorf("Vector rotation failed. %f, %f", vector.x, vector.y)
	}
}
