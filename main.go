package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

func main() {
	// Fancy progress bar :)
	// (or deadlock sanity check)
	progressI := 0
	progressBar := []string{
		"░█████████",
		"█░████████",
		"██░███████",
		"███░██████",
		"████░█████",
		"█████░████",
		"██████░███",
		"███████░██",
		"████████░█",
		"█████████░"}

	// Defaults
	inputFile := "./data.csv"
	epsilon := 0.0003
	minPts := 5
	threadN := runtime.NumCPU() // Number of CPU cores
	maxJobSize := 1_000

	// If no arguments are given, print help
	if len(os.Args) == 1 {
		fmt.Println("Usage:   ./dbscan <input_file> <epsilon> <minPts> <maxJobSize> <threadN>")
		fmt.Println("Example: ./dbscan ./data.csv   0.0003    5        1000         12")
		fmt.Println("Note:    maxJobSize is the maximum number of points that can be processed by a single job in the thread pool")
		fmt.Println("         threadN defaults to the number of logical cores on the machine (so you probably can leave it empty)")
		fmt.Println("         Other than that, the values are defaulted to the example above")
		fmt.Println("         If you're not feeling like going for a coffee break, you can try using a smaller epsilon or maxJobSize")
	}

	// Try get input file
	if len(os.Args) > 1 {
		inputFile = os.Args[1]
	}
	// Try get epsilon
	if len(os.Args) > 2 {
		epsilon, _ = strconv.ParseFloat(os.Args[2], 64)
	}
	// Try get minPts
	if len(os.Args) > 3 {
		minPts, _ = strconv.Atoi(os.Args[3])
	}
	// Try get maxJobSize
	if len(os.Args) > 4 {
		maxJobSize, _ = strconv.Atoi(os.Args[4])
	}
	// Try get threadN
	if len(os.Args) > 5 {
		threadN, _ = strconv.Atoi(os.Args[5])
	}

	// Print settings
	fmt.Println()
	fmt.Println("Current settings:")
	fmt.Println("Input file:", inputFile)
	fmt.Println("Epsilon:", epsilon)
	fmt.Println("MinPts:", minPts)
	fmt.Println("MaxJobSize:", maxJobSize)
	fmt.Println("ThreadN:", threadN)
	fmt.Println()

	startT := time.Now() // For benchmark only
	checkPointT := startT

	// Read the CSV file and return a list of points and a bounding box
	fmt.Println("Reading file...")
	rect, points := readCSV(inputFile)
	// Starts a new binary space partition for speed-up querying
	fmt.Println("Building BSP tree...")
	bsp := NewBSPTreeFromPoints(rect, &points)
	fmt.Println("Starting DBSCAN...")

	clustersResult := []Cluster{}
	for cluster := range dbscanParallel(bsp, epsilon, maxJobSize, threadN) {
		// Calculate cluster size
		clusterSize := 0
		for i := 0; i < len(cluster.points); i++ {
			clusterSize += cluster.points[i].cnt // cnt is the number of points in the coordinate
		}

		// If cluster size greater than or equal to minPts, add to result
		if clusterSize >= minPts {
			clustersResult = append(clustersResult, cluster)
		}

		// Print progress
		fmt.Print(progressBar[progressI], " Clusters: ", len(clustersResult), "\033[G") // Print progress bar and move cursor to beginning of line
		progressI++
		if progressI >= len(progressBar) {
			progressI = 0
		}
	}
	fmt.Println("▓▓▓▓▓▓▓▓▓▓ Clusters found:", len(clustersResult), "| ΔT:", time.Since(startT), " + 0s |")
	checkPointT = time.Now()

	// Merge clusters
	fmt.Println("Merging clusters...")

	// Fake progress bar
	progressI = 0
	done := false
	go func() {
		for {
			fmt.Print(progressBar[progressI], " Merging...\033[G") // Print progress bar and move cursor to beginning of line
			progressI++
			time.Sleep(time.Millisecond * 100)
			// Wrap progress bar
			if progressI >= len(progressBar) {
				progressI = 0
			}
			if done {
				break
			}
		}
	}()

	clustersResult = mergeClusters(clustersResult, epsilon)
	done = true // Stop fake progress bar
	time.Sleep(time.Millisecond * 250)
	// Print len of merged clusters
	fmt.Println("▓▓▓▓▓▓▓▓▓▓ Merged clusters:", len(clustersResult), "| ΔT:", time.Since(startT), " + ", time.Since(checkPointT), "|")
	checkPointT = time.Now()

	fmt.Println("Saving results...")
	// Write clusters to file
	writeCSV("./clusters.csv", clustersResult)
	// Write points to file
	writeClusterPoints("./points.csv", clustersResult)
	fmt.Println("Total elapsed time:", time.Since(startT))
}
