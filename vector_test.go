package main

import (
	"fmt"
	"math"
	"testing"
)

func TestRotate(t *testing.T) {
	vector := Vector2D{
		X: 1.0,
		Y: 0,
	}
	rotated := vector.Rotated(math.Pi)
	fmt.Println(vector, rotated)
	// some floating point errors
	if rotated.X != -1 || rotated.Y > 0.000000001 || rotated.Y < -0.000000000001 {
		t.Errorf("vector rotation failed. %f, %f", rotated.X, rotated.Y)
	}

	vector = Vector2D{
		X: -1.0,
		Y: 0,
	}
	rotated = vector.Rotated(math.Pi)
	fmt.Println(vector, rotated)
	// some floating point errors
	if rotated.X != 1 || rotated.Y > 0.000000001 || rotated.Y < -0.000000000001 {
		t.Errorf("vector rotation failed. %f, %f", rotated.X, rotated.Y)
	}

	vector = Vector2D{
		X: 1.0,
		Y: 1.0,
	}
	rotated = vector.Rotated(-math.Pi)
	fmt.Println(vector, rotated)
	// some floating point errors
	if rotated.X != -0.9999999999999999 || rotated.Y != -1.0000000000000002 {
		t.Errorf("vector rotation failed. %f, %f", rotated.X, rotated.Y)
	}
}

func TestDot(t *testing.T) {
	first := &Vector2D{
		X: 1.0,
		Y: 0,
	}
	second := &Vector2D{
		X: 1.0,
		Y: 1.0,
	}
	angle := first.Dot(second)
	fmt.Println(angle)

}
