package main

import "testing"

// data.csv has exactly 232050 points
const NUMBER_OF_POINTS = 232050

func TestFileSizeLoad(t *testing.T) {
	_, points := readCSV("data.csv")

	if len(points) != NUMBER_OF_POINTS {
		t.Error("File size is not 232050")
	}
}

func TestBSPCount(t *testing.T) {
	rect, points := readCSV("data.csv")

	// Add all points to the tree
	bsp := NewBSPTreeFromPoints(rect, &points)

	// Check if the tree has the correct number of nodes
	if bsp.size != NUMBER_OF_POINTS {
		t.Error("Tree does not have the correct number of nodes")
	}

	// Check if Querying the entire tree returns the correct number of points
	bspPoints := bsp.Query(rect)
	count := 0
	for i := 0; i < len(bspPoints); i++ {
		count += bspPoints[i].cnt
	}

	if count != NUMBER_OF_POINTS {
		t.Error("Query is not correct")
	}

	// Test bsp left + right size is correct
	s := bsp.left.size + bsp.right.size
	if s != bsp.size {
		t.Error("Size is not correct")
	}

	sl := bsp.left
	s = sl.left.size + sl.right.size
	if s != sl.size {
		t.Error("Size is not correct")
	}
}
