package main

import (
	"fmt"
	"sort"
	"time"
)

func main() {
	// All comments regarding time it takes to run the program are in debug mode.
	timeTemp := time.Now()

	rect, points := readCSV("data.csv") // Takes about 0.5s to read the file
	fmt.Println("Number of points:", len(points))
	// Starts a new binary space partition tree

	fmt.Println("Building BSP tree...")
	bsp := NewBSPTree(rect.x, rect.y, rect.w, rect.h) // 0.2s to build the tree
	for i := 0; i < len(points); i++ {
		bsp.Insert(&points[i])
	}

	// Print time
	fmt.Println("Time:", time.Since(timeTemp))

	// Query the tree
	fmt.Println("Querying tree...")
	timeTemp = time.Now()
	nPoints := 0
	for p := range bsp.Query(rect) {
		_ = p
		nPoints++
	}
	fmt.Println("Time:", time.Since(timeTemp))
	fmt.Println("Number of points in the tree:", nPoints)
	fmt.Println("Yep")
	clusters := dbscan(bsp, 0.0003, 5)

	pts := []Point{}
	for i := 0; i < len(clusters); i++ {
		pts = append(pts, clusters[i]...)
	}

	fmt.Println(len(clusters))
	fmt.Println("Original: ", len(points))
	fmt.Println("Returned: ", len(pts))

	// Save clusters to file

	// Sort points by size
	sort.Slice(clusters, func(i, j int) bool {
		return len(clusters[i]) > len(clusters[j])
	})

	// saveFolder := "clusters/"
	// for i := 0; i < len(clusters); i++ {
	// 	savePoints(saveFolder+fmt.Sprintf("cluster_%d.csv", i), clusters[i])
	// }
}
