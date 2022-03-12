package main

import (
	"testing"
)

func TestRectExpand(t *testing.T) {
	r := Rect{0, 0, 10, 10}
	r = r.Expand(5)
	if r.x != -5 || r.y != -5 || r.w != 20 || r.h != 20 {
		t.Errorf("Expand failed: %v", r)
	}

	r = Rect{0, 0, 10, 10}
	r = r.Expand(-1) // Make it smaller
	if r.x != 1 || r.y != 1 || r.w != 8 || r.h != 8 {
		t.Errorf("Expand failed: %v", r)
	}
}

func TestRectCentroid(t *testing.T) {
	r := Rect{0, 0, 10, 10}
	p := r.Centroid()
	if p.x != 5 || p.y != 5 {
		t.Errorf("Centroid failed: %v", p)
	}

	r = Rect{5, 5, 5, 5}
	p = r.Centroid()
	if p.x != 7.5 || p.y != 7.5 {
		t.Errorf("Centroid failed: %v", p)
	}
}

func TestMergeRect(t *testing.T) {
	r := Rect{0, 0, 10, 10}
	r2 := Rect{5, 5, 5, 5}
	// Merge should produce a rectangle that contains both (r1)
	r = r.Merge(r2)
	if r.x != 0 || r.y != 0 || r.w != 10 || r.h != 10 {
		t.Errorf("Merge failed: %v", r)
	}

	r = Rect{0, 0, 10, 10}
	r2 = Rect{10, 10, 5, 5}
	r = r.Merge(r2) // Expand r1 to contain r2
	if r.x != 0 || r.y != 0 || r.w != 15 || r.h != 15 {
		t.Errorf("Merge failed: %v", r)
	}
}
