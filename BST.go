package main

import (
	"fmt"
	"math"
	"sync"
)

// BSPTree is a 2d spatial index of Point objects.
type BSPTree struct {
	size       int
	subdivided bool
	rect       Rect
	points     []*Point
	left       *BSPTree
	right      *BSPTree
}

// Create a new BSPTree
func NewBSPTree(x, y, w, h float64) *BSPTree {
	return &BSPTree{
		rect:       Rect{x, y, w, h},
		subdivided: false,
		points:     []*Point{},
	}
}

// Tree insert
func (q *BSPTree) Insert(p *Point) {
	q.size++
	// If already subdivided, insert into appropriate branch
	if q.subdivided {
		if q.left.Contains(*p) {
			q.left.Insert(p)
		} else if q.right.Contains(*p) {
			q.right.Insert(p)
		} else {
			q.Rebuild(p)
		}
		return
	}

	// Otherwise, if point is too close to existing point, add to list
	// Otherwise, subdivide
	if len(q.points) == 0 {
		q.points = append(q.points, p)
	} else {
		delta := q.points[0].Distance(*p)
		const epsilon = 0.001
		q.points = append(q.points, p)
		if delta > epsilon {
			q.Subdivide(p)
		}
	}
}

// Subdivide tree while adding point
func (q *BSPTree) Subdivide(p *Point) {
	// Initialize the quadrants
	// If rect is vertical rectangle split vertically, else split horizontally
	ratio := q.rect.w / q.rect.h
	if ratio > 1 {
		q.left = NewBSPTree(q.rect.x, q.rect.y, q.rect.w/2, q.rect.h)
		q.right = NewBSPTree(q.rect.x+q.rect.w/2, q.rect.y, q.rect.w/2, q.rect.h)
	} else {
		q.left = NewBSPTree(q.rect.x, q.rect.y, q.rect.w, q.rect.h/2)
		q.right = NewBSPTree(q.rect.x, q.rect.y+q.rect.h/2, q.rect.w, q.rect.h/2)
	}

	// Add points to their respective quadrants
	for _, pts := range q.points {
		if q.left.Contains(*pts) {
			q.left.Insert(pts)
		} else {
			q.right.Insert(pts)
		}
	}

	q.subdivided = true
	q.points = []*Point{}
}

// Check if point is inside tree rect
func (q *BSPTree) Contains(p Point) bool {
	gt_left := p.x >= q.rect.x
	lt_right := p.x < q.rect.x+q.rect.w
	gt_top := p.y >= q.rect.y
	lt_bottom := p.y < q.rect.y+q.rect.h
	return gt_left && lt_right && gt_top && lt_bottom
}

// Get all the points in the tree
func (q *BSPTree) GetPoints() []*Point {
	if q.subdivided {
		return append(q.left.GetPoints(), q.right.GetPoints()...)
	}
	return q.points
}

func (q *BSPTree) Rebuild(p *Point) {
	fmt.Println("Warning: rebuilding tree!!!")
	// Get all points in the tree
	points := q.GetPoints()

	// Get bounding box of all points including the new point
	var minx, miny, maxx, maxy float64 = math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64
	allPoints := []*Point{p}
	allPoints = append(allPoints, points...)
	for _, p := range allPoints {
		if p.x < minx {
			minx = p.x
		}
		if p.y < miny {
			miny = p.y
		}
		if p.x > maxx {
			maxx = p.x
		}
		if p.y > maxy {
			maxy = p.y
		}
	}

	// Create new bounding box with padding
	q.rect = Rect{minx - 1, miny - 1, maxx - minx + 2, maxy - miny + 2}
	q.points = []*Point{}
	q.subdivided = false
	q.left = nil
	q.right = nil

	// Re-add all points to the tree
	for _, pts := range points {
		q.Insert(pts)
	}
	q.Insert(p)
}

func (q *BSPTree) QueryImpl(r Rect, c chan *Point, g *sync.WaitGroup) {
	if q.subdivided {
		g.Add(2)
		q.left.QueryImpl(r, c, g)
		q.right.QueryImpl(r, c, g)
	} else {
		for _, pts := range q.points {
			if rectPointIntersect(r, *pts) {
				c <- pts
			}
		}
	}
}

// Query the tree for points within a given rect (unstable)
func (q *BSPTree) Query(r Rect) <-chan *Point {
	c := make(chan *Point, q.size)

	go func() {
		var wg sync.WaitGroup
		q.QueryImpl(r, c, &wg)
		close(c)
	}()

	return c
}

// Iterate over all points
func (q *BSPTree) Iterate() <-chan *Point {
	return q.Query(q.rect)
}
