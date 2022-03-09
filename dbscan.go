package main

import (
	"sync"
)

// Function dbscan takes list of point and return list of list of points
func dbscan(bsp *BSPTree, epsilon float64, minPts int) [][]Point {
	clusters := sync.Map{} // Maps Point to Int
	merges := sync.Map{}   // Maps Int to map of int
	var wg sync.WaitGroup

	clusterId := 0
	// For each point in the tree
	for p := range bsp.Iterate() {
		clusterId++
		// If point was already assigned to a cluster, skip it
		wg.Add(1)
		go func(p *Point, clusterId int) {
			toMerge := make(map[int]bool)

			// // Query neighbors (including the point itself)
			r := Rect{p.x - epsilon, p.y - epsilon, epsilon * 2, epsilon * 2}
			toVisit := []Point{}
			for n := range bsp.Query(r) {
				toVisit = append(toVisit, *n)
			}

			// Recursively visit neighbors
			for len(toVisit) > 0 {
				current := toVisit[0]
				toVisit = toVisit[1:]

				if val, ok := clusters.Load(current); ok {
					toMerge[val.(int)] = true
					continue
				}

				// Query neighbors (including the point itself)
				r := Rect{current.x - epsilon, current.y - epsilon, epsilon * 2, epsilon * 2}
				for n := range bsp.Query(r) {
					if n.Distance(current) <= epsilon && n.x != p.x && n.y != p.y {
						toVisit = append(toVisit, *n)
					}
				}
				clusters.Store(current, clusterId)
			}

			merges.Store(clusterId, toMerge)
			wg.Done()
		}(p, clusterId)
	}
	wg.Wait()
	// Printlen of cluster[0]
	// fmt.Println(len(clusters))

	// Collect the values of clusterMap into a list
	ret := make([][]Point, 0)

	clusters.Range(func(key, value interface{}) bool {
		ret = append(ret, []Point{key.(Point)})
		return true

	})

	// for _, v := range clusters {
	// 	if len(v) >= minPts {
	// 		ret = append(ret, v)
	// 	}
	// }

	// Print size of clusters
	return ret
}
