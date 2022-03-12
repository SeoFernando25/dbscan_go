package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Reads the CSV file and returns a list of points and a bounding box
func readCSV(filename string) (Rect, []Point) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return Rect{}, nil
	}
	defer file.Close()

	// Create a new list
	list := make([]Point, 0)

	readFirstLine := false
	minX, minY, maxX, maxY := math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	// Skip the first line
	scanner.Scan()
	for scanner.Scan() {
		// Save the 8th and 9th fields (x-y coordinates) as a point
		fields := strings.Split(scanner.Text(), ",")
		x, _ := strconv.ParseFloat(fields[8], 64)
		y, _ := strconv.ParseFloat(fields[9], 64)
		scale := 1.0
		p := Point{x * scale, y * scale}

		if !readFirstLine {
			readFirstLine = true
			minX = p.x
			minY = p.y
			maxX = p.x
			maxY = p.y
		} else {
			// Update the bounding box
			if p.x <= minX {
				minX = p.x
			}
			if p.y <= minY {
				minY = p.y
			}
			if p.x >= maxX {
				maxX = p.x
			}
			if p.y >= maxY {
				maxY = p.y
			}
		}

		list = append(list, p)
	}

	return Rect{minX, minY, maxX - minX, maxY - minY}, list
}

// Writes the list of clusters to a CSV file
func writeCSV(filename string, clusters []Cluster) {
	// Sort the clusters by length of each item
	sort.Slice(clusters, func(i, j int) bool {
		return len(clusters[i].points) > len(clusters[j].points)
	})

	// Open the file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Write the header
	file.WriteString("ClusterId,Latitude,Longitude,Size\n")

	// Write the clusters
	clusterId := 0
	for _, cluster := range clusters {
		// Get the cluster
		// Get the average point
		points := []Point{}
		clusterSize := 0
		for _, p := range cluster.points {
			points = append(points, *p.Point)
			clusterSize += p.cnt
		}
		p := pointAverage(points)
		// Write the cluster to the file

		file.WriteString(fmt.Sprintf("%d,%f,%f,%d\n", clusterId, p.y, p.x, clusterSize))
		clusterId++
	}
}

// Saves the cluster points from a list of clusters to a CSV file
func writeClusterPoints(filename string, clusters []Cluster) {
	// Sort the clusters by length of each item
	sort.Slice(clusters, func(i, j int) bool {
		return len(clusters[i].points) > len(clusters[j].points)
	})

	// Open the file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Write the header
	file.WriteString("ClusterId,Latitude,Longitude\n")
	clusterId := 1
	for _, cluster := range clusters {
		for _, p := range cluster.points {
			// Write the cluster to the file
			file.WriteString(fmt.Sprintf("%d,%f,%f\n", clusterId, p.y, p.x))
		}
		clusterId++
	}
}
