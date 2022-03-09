package main

import (
	"math"
)

// region Point

// Type point represents a point in 2D space
type Point struct {
	x float64
	y float64
}

// Calculates the distance between two points
func distance(p1, p2 Point) float64 {
	a := p1.x - p2.x
	b := p1.y - p2.y
	x := a*a + b*b
	return math.Sqrt(x)
}

func (p Point) Distance(q Point) float64 {
	return distance(p, q)
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

//  endregion

// region Rect

type Rect struct {
	x, y, w, h float64
}

// endregion

// Rect to rect collision detection
func rectIntersect(rect1 Rect, rect2 Rect) bool {
	a := rect1.x <= rect2.x+rect2.w
	b := rect1.x+rect1.w >= rect2.x
	c := rect1.y <= rect2.y+rect2.h
	d := rect1.h+rect1.y >= rect2.y
	return a && b && c && d
}

// Rect to Point collision detection
func rectPointIntersect(r Rect, p Point) bool {
	gt_left := p.x >= r.x
	lt_right := p.x <= r.x+r.w
	gt_top := p.y >= r.y
	lt_bottom := p.y <= r.y+r.h
	return gt_left && lt_right && gt_top && lt_bottom
}
