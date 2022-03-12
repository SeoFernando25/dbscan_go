package main

import "math"

// Type point represents a point in 2D space
type Point struct {
	x float64
	y float64
}

func (p Point) Distance(q Point) float64 {
	dx := p.x - q.x
	dy := p.y - q.y
	return math.Sqrt(dx*dx + dy*dy)
}

// Calculates the pointAverage of a list of points
func pointAverage(points []Point) Point {
	var x, y float64
	for _, p := range points {
		x += p.x
		y += p.y
	}
	return Point{x / float64(len(points)), y / float64(len(points))}
}
