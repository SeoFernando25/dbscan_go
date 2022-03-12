package main

import (
	"math/rand"
	"testing"
)

func TestPointCollision(t *testing.T) {
	for i := 0; i < 100; i++ {
		rand1 := rand.Float64() * 100
		p1 := Point{rand1, rand1}
		p2 := Point{rand1, rand1}
		if !pointIntersect(p1, p2) {
			t.Error("Point to Point collision is not correct")
		}
	}

	for i := 0; i < 100; i++ {
		rand1 := rand.Float64() * 100
		rand2 := 200 + rand.Float64()*100
		p1 := Point{rand1, rand1}
		p2 := Point{rand1, rand1 + rand2}
		if pointIntersect(p1, p2) {
			t.Error("Point to Point collision is not correct")
		}
	}
}

func TestRectPointCollision(t *testing.T) {
	// Should collide
	r := Rect{0, 0, 10, 10}
	p := Point{0, 0}
	if !rectPointIntersect(r, p) { // Points intersect on edge
		t.Error("Rect to Point collision is not correct")
	}

	p = Point{5, 5}
	if !rectPointIntersect(r, p) { // Point inside
		t.Error("Rect to Point collision is not correct")
	}

	p = Point{10, 10}
	if !rectPointIntersect(r, p) { // Point on edge
		t.Error("Rect to Point collision is not correct")
	}

	// Should not collide
	p = Point{-1, -1}
	if rectPointIntersect(r, p) { // Point outside
		t.Error("Rect to Point collision is not correct")
	}

	p = Point{-1, 5}
	if rectPointIntersect(r, p) { // Point outside x axis
		t.Error("Rect to Point collision is not correct")
	}

	p = Point{5, -1}
	if rectPointIntersect(r, p) { // Point outside y axis
		t.Error("Rect to Point collision is not correct")
	}
}

func TestRectCollision(t *testing.T) {
	// Should collide
	r1 := Rect{0, 0, 10, 10}
	r2 := Rect{0, 0, 10, 10}

	if !rectIntersect(r1, r2) { // Same rect
		t.Error("Rect to Rect collision is not correct")
	}

	r2 = Rect{0, 10, 10, 10}
	if !rectIntersect(r1, r2) { // Touching on edge
		t.Error("Rect to Rect collision is not correct")
	}

	r2 = Rect{-99999, -99999, 99999 * 2, 99999 * 2}
	if !rectIntersect(r1, r2) { // Way bigger rect
		t.Error("Rect to Rect collision is not correct")
	}

	r2 = Rect{5, 5, 10, 10}
	if !rectIntersect(r1, r2) { // Partially inside
		t.Error("Rect to Rect collision is not correct")
	}

	// Should not collide
	r2 = Rect{-9999, -9999, 10, 10}
	if rectIntersect(r1, r2) { // Way outside
		t.Error("Rect to Rect collision is not correct")
	}

	r2 = Rect{10.0001, 10.0001, 0.001, 0.001}
	if rectIntersect(r1, r2) { // Almost touching
		t.Error("Rect to Rect collision is not correct")
	}

}
