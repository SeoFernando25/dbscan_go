package main

import (
	"fmt"
	"math"
)

// BSPTree is a 2d spatial index of Point objects.

type BSPTreePoint struct {
	*Point
	cnt int
}

type BSPTree struct {
	cnt   int // How many points are there in this node?
	size  int // How many points are in the tree in total?
	rect  Rect
	point *Point
	left  *BSPTree
	right *BSPTree
}

// Create a new BSPTree
func NewBSPTree(x, y, w, h float64) *BSPTree {
	return &BSPTree{
		rect: Rect{x, y, w, h},
	}
}

// Create a new BSPTree
func NewBSPTreeFromPoints(r Rect, points *[]Point) *BSPTree {
	tree := NewBSPTree(r.x, r.y, r.w, r.h)
	for i := 0; i < len(*points); i++ {
		p := &(*points)[i]
		tree.Insert(p)
	}
	return tree
}

// Tree insert
func (q *BSPTree) Insert(p *Point) {
	q.size++
	if q.point == nil && q.left == nil && q.right == nil { // Try normal insert
		q.point = p
		q.cnt = 1
	} else if q.left != nil && q.right != nil { // Find closes quadrant
		distLeft := p.Distance(q.left.rect.Centroid())
		distRight := p.Distance(q.right.rect.Centroid())
		if distLeft < distRight {
			q.left.Insert(p)
		} else {
			q.right.Insert(p)
		}
	} else if q.point != nil && q.point.Distance(*p) == 0 { // If point is in the exact same place, add it to the tree
		q.cnt++
	} else { // Subdivide
		q.Subdivide(p)
	}
}

// Subdivide tree while adding point
func (q *BSPTree) Subdivide(p *Point) {
	// Initialize the quadrants
	// If rect is vertical rectangle split vertically, else split horizontally
	ratio := q.rect.w / q.rect.h
	if ratio >= 1 { // Split vertically
		w := q.rect.w / 2
		q.left = NewBSPTree(q.rect.x, q.rect.y, w, q.rect.h)
		q.right = NewBSPTree(q.rect.x+w, q.rect.y, w, q.rect.h)
	} else { // Split horizontally
		h := q.rect.h / 2
		q.left = NewBSPTree(q.rect.x, q.rect.y, q.rect.w, h)
		q.right = NewBSPTree(q.rect.x, q.rect.y+h, q.rect.w, h)
	}

	// Add points to their respective quadrants
	distLeft := q.point.Distance(q.left.rect.Centroid())
	distRight := q.point.Distance(q.right.rect.Centroid())

	var toInsert *BSPTree
	if distLeft < distRight {
		toInsert = q.left
	} else {
		toInsert = q.right
	}

	toInsert.point = q.point
	toInsert.cnt = q.cnt
	toInsert.size = q.cnt // Subtract the point we just added (we are inserting a new point)

	distLeft = p.Distance(q.left.rect.Centroid())
	distRight = p.Distance(q.right.rect.Centroid())

	if distLeft < distRight {
		q.left.Insert(p)
	} else {
		q.right.Insert(p)
	}

	q.point = nil // Clear the point (it's been inserted into the children)
	q.cnt = 0
}

func (q *BSPTree) Rebuild(p *Point) {
	fmt.Println("Warning: rebuilding tree!!!")
	// Get all points in the tree
	points := q.Query(q.rect)

	// Get bounding box of all points including the new point
	var minX, minY, maxX, maxY float64 = math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64
	allPoints := make([]*Point, len(points))
	allPoints = append(allPoints, p)
	for _, bspPoint := range points {
		for i := 0; i < bspPoint.cnt; i++ {
			allPoints = append(allPoints, bspPoint.Point)
		}
	}
	for _, p := range allPoints {
		if p.x < minX {
			minX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
		if p.x > maxX {
			maxX = p.x
		}
		if p.y > maxY {
			maxY = p.y
		}
	}

	// Create new bounding box with padding
	q.rect = Rect{minX - 1, minY - 1, maxX - minX + 2, maxY - minY + 2}
	q.point = nil
	q.left = nil
	q.right = nil

	// Re-add all points to the tree
	for _, pts := range allPoints {
		q.Insert(pts)
	}
}

func (q *BSPTree) Query(r Rect) []BSPTreePoint {
	c := q.QueryAsync(r)
	var points []BSPTreePoint
	for p := range c {
		points = append(points, p)
	}
	return points
}

func (q *BSPTree) QueryChan(r Rect, c chan BSPTreePoint) {
	if q == nil {
		return
	}

	q.left.QueryChan(r, c)
	q.right.QueryChan(r, c)

	if q.point != nil && rectPointIntersect(r, *q.point) {
		c <- BSPTreePoint{q.point, q.cnt}
	}
}

// Returns an iterator over a query
func (q *BSPTree) QueryAsync(r Rect) <-chan BSPTreePoint {
	c := make(chan BSPTreePoint, q.size)

	go func() {
		q.QueryChan(r, c)
		close(c)
	}()

	return c
}

// Iterate over all points
func (q *BSPTree) Iterate() <-chan BSPTreePoint {
	return q.QueryAsync(q.rect)
}
