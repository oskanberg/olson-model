package main

import "math"

type Vector2D struct {
	x, y float64
}

func (s *Vector2D) Normalised() *Vector2D {
	mag := s.Magnitude()
	return &Vector2D{
		x: s.x / mag,
		y: s.y / mag,
	}
}

func (s *Vector2D) Magnitude() float64 {
	return math.Sqrt((s.x * s.x) + (s.y * s.y))
}

func (s *Vector2D) Subtract(v *Vector2D) *Vector2D {
	return &Vector2D{
		x: s.x - v.x,
		y: s.y - v.y,
	}
}

func (s *Vector2D) Add(v *Vector2D) *Vector2D {
	return &Vector2D{
		x: s.x + v.x,
		y: s.y + v.y,
	}
}

func (s *Vector2D) Dot(v *Vector2D) float64 {
	return (s.x * v.x) + (s.y * v.y)
}

func (s *Vector2D) Rotated(radians float64) *Vector2D {
	cr := math.Cos(radians)
	sr := math.Sin(radians)
	return &Vector2D{
		x: (s.x * cr) - (s.y * sr),
		y: (s.x * sr) + (s.y * cr),
	}
}

func (s *Vector2D) Multiplied(by float64) *Vector2D {
	return &Vector2D{
		x: s.x * by,
		y: s.y * by,
	}
}
