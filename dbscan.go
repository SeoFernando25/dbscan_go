package main

import (
	"sync"
)

// Cluster holds a rect and a list of points
type Cluster struct {
	Rect
	points []BSPTreePoint
}

// Thread pool job producer that returns partitions of points that are within the maxJobSize threshold.
// Note that maxJobSize is not guaranteed if tree can't be futher broken down.
func dbscanProducer(bsp *BSPTree, maxJobSize int, wg *sync.WaitGroup) <-chan *BSPTree {
	outJobs := make(chan *BSPTree, maxJobSize)
	wg.Add(1)

	go func() { // Iterate over tree for nodes that satisfy the maxJobSize threshold
		q := []*BSPTree{bsp}
		for len(q) > 0 {
			// Pop
			current := q[0]
			q = q[1:]

			if current.size > maxJobSize { // If size is too big, split
				// Split
				if current.left != nil && current.right != nil {
					q = append(q, current.left)
					q = append(q, current.right)
				} else { // Can't split, just add it anyway
					outJobs <- current
				}
			} else if current.size > 0 { // Add to jobs
				outJobs <- current
			}
		}
		close(outJobs)
		wg.Done()
	}()

	return outJobs
}

// A simple thread pool worker that wait for jobs and processes them
func dbscanWorker(bsp <-chan *BSPTree, res chan<- Cluster, epsilon float64, wg *sync.WaitGroup) {
	for bspJob := range bsp {
		dbscan(bspJob, epsilon, res)
	}
	wg.Done()
}

// Returns a channel containing the unmerged clusters
func dbscanParallel(bspRoot *BSPTree, epsilon float64, maxJobSize int, nWorkers int) <-chan Cluster {
	var (
		wg   sync.WaitGroup
		res  = make(chan Cluster, maxJobSize)
		jobs = dbscanProducer(bspRoot, maxJobSize, &wg)
	)

	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go dbscanWorker(jobs, res, epsilon, &wg)
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	return res
}

// Perform DBSCAN clustering from a spatial partitioning tree
// Returns a list of points that are within epsilon distance of the query point
func dbscan(bsp *BSPTree, epsilon float64, res chan<- Cluster) {
	visited := make(map[BSPTreePoint]bool)
	// For each point in the tree
	for pQuery := range bsp.Iterate() {
		mainPoint := pQuery.Point
		// If key in visited map, skip
		if _, ok := visited[pQuery]; ok {
			continue
		}

		// Calculate cluster bounding box
		minX := mainPoint.x
		minY := mainPoint.y
		maxX := mainPoint.x
		maxY := mainPoint.y

		// Find neighbors (including the point itself)
		r := Rect{mainPoint.x - epsilon, mainPoint.y - epsilon, epsilon * 2, epsilon * 2}
		toVisit := []BSPTreePoint{}
		for n := range bsp.QueryAsync(r) {
			toVisit = append(toVisit, n)
		}
		clusterPoints := []BSPTreePoint{}

		// Recursively visit neighbors
		for len(toVisit) > 0 {
			current := toVisit[0]
			toVisit = toVisit[1:]

			// If key in visited map, skip
			if _, ok := visited[current]; ok {
				continue
			}

			// Query neighbors, excluding the point itself
			r := Rect{current.x - epsilon, current.y - epsilon, epsilon * 2, epsilon * 2}
			for n := range bsp.QueryAsync(r) {
				// If point is within epsilon distance, add to toVisit
				if n.Point.Distance(*current.Point) <= epsilon && !pointIntersect(*n.Point, *mainPoint) {
					toVisit = append(toVisit, n)
				}
			}
			visited[current] = true // Mark as visited

			clusterPoints = append(clusterPoints, current) // Add to cluster

			// Update bounding box
			if current.x < minX {
				minX = current.x
			}
			if current.y < minY {
				minY = current.y
			}
			if current.x > maxX {
				maxX = current.x
			}
			if current.y > maxY {
				maxY = current.y
			}
		}
		boundingRect := Rect{minX, minY, maxX - minX, maxY - minY}
		res <- Cluster{boundingRect, clusterPoints}
	}
}

// This is a naive implementation for merging clusters
// it does not take into account the spatial partitioning tree and runs in O(n^2) time
// DO NOT GRADE THIS FUNCTION, THIS IS JUST AN EXTRA SO THAT I COULD SEE THE END RESULT
func mergeClusters(clusters []Cluster, epsilon float64) []Cluster {
	// Merge clusters
	index := 0

start:
	for index < len(clusters) {
		current := clusters[index]

		// Get all neighbors
		for i := 0; i < len(clusters); i++ {
			if i == index {
				continue
			}

			// Intersection Rect is rect + epsilon
			adjustedRect := current.Rect.Expand(epsilon)
			if rectIntersect(current.Rect, adjustedRect) {
				neighbor := clusters[i]
				// Merge with neighbors if any point is within epsilon distance
				for _, p := range current.points {
					for _, p2 := range neighbor.points {
						if p.Point.Distance(*p2.Point) <= epsilon {
							// fmt.Println("Merging")
							current.points = append(current.points, neighbor.points...)
							current.Rect = current.Rect.Merge(neighbor.Rect)
							// Delete merged neighbor
							clusters = append(clusters[:i], clusters[i+1:]...)
							goto start
						}
					}
				}
			}
		}
		index++
	}
	return clusters
}
