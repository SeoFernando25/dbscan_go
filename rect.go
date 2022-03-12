package main

import "math"

type Rect struct {
	x, y, w, h float64
}

// Returns the center point of the rect
func (r Rect) Centroid() Point {
	return Point{r.x + r.w/2, r.y + r.h/2}
}

// Merges two rectangles
func (r Rect) Merge(other Rect) Rect {
	return Rect{
		x: math.Min(r.x, other.x),
		y: math.Min(r.y, other.y),
		w: math.Max(r.x+r.w, other.x+other.w) - math.Min(r.x, other.x),
		h: math.Max(r.y+r.h, other.y+other.h) - math.Min(r.y, other.y),
	}
}

// Add a padding to the rect
func (r Rect) Expand(amount float64) Rect {
	return Rect{
		x: r.x - amount,
		y: r.y - amount,
		w: r.w + amount*2,
		h: r.h + amount*2,
	}
}
