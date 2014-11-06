package main

import (
	"fmt"
	"math"
	"testing"
)

func TestRotate(t *testing.T) {
	vector := Vector2D{
		x: 1.0,
		y: 0,
	}
	rotated := vector.Rotated(math.Pi)
	fmt.Println(vector, rotated)
	// some floating point errors
	if rotated.x != -1 || rotated.y > 0.000000001 || rotated.y < -0.000000000001 {
		t.Errorf("vector rotation failed. %f, %f", rotated.x, rotated.y)
	}

	vector = Vector2D{
		x: -1.0,
		y: 0,
	}
	rotated = vector.Rotated(math.Pi)
	fmt.Println(vector, rotated)
	// some floating point errors
	if rotated.x != 1 || rotated.y > 0.000000001 || rotated.y < -0.000000000001 {
		t.Errorf("vector rotation failed. %f, %f", rotated.x, rotated.y)
	}

	vector = Vector2D{
		x: 1.0,
		y: 1.0,
	}
	rotated = vector.Rotated(-math.Pi)
	fmt.Println(vector, rotated)
	// some floating point errors
	if rotated.x != -0.9999999999999999 || rotated.y != -1.0000000000000002 {
		t.Errorf("vector rotation failed. %f, %f", rotated.x, rotated.y)
	}
}

func TestDot(t *testing.T) {
	first := &Vector2D{
		x: 1.0,
		y: 0,
	}
	second := &Vector2D{
		x: 1.0,
		y: 1.0,
	}
	angle := first.Dot(second)
	fmt.Println(angle)

}
