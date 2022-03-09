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
	defer file.Close()
	if err != nil {
		fmt.Println("Error:", err)
		return Rect{}, nil
	}

	// Create a new list
	list := make([]Point, 0)
	minx, miny, maxx, maxy := math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64

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

		// Update the bounding box
		if p.x < minx {
			minx = p.x
		}
		if p.y < miny {
			miny = p.y
		}
		if p.x > maxx {
			maxx = x
		}
		if p.y > maxy {
			maxy = p.y
		}

		list = append(list, p)
	}

	return Rect{minx - 1, miny - 1, maxx - minx + 2, maxy - miny + 2}, list
}

// Writes the list of "list of points" to a CSV file, eg:
// ClusterId,Latitude,Longitude,Size
// 1,37.78,122.45,1
func writeCSV(filename string, clusters [][]Point) {
	// Sort the clusters by length of each item
	sort.Slice(clusters, func(i, j int) bool {
		return len(clusters[i]) > len(clusters[j])
	})

	// Open the file
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Write the header
	file.WriteString("ClusterId,Latitude,Longitude,Size\n")

	// Write the clusters
	clusterId := 0
	for _, cluster := range clusters {
		// Get the cluster
		// Get the average point
		p := pointAverage(cluster)
		// Write the cluster to the file
		file.WriteString(fmt.Sprintf("%d,%f,%f,%d\n", clusterId, p.y, p.x, len(cluster)))
		clusterId++
	}
}

func writeClusterPoints(filename string, clusters [][]Point) {
	// Sort the clusters by length of each item
	sort.Slice(clusters, func(i, j int) bool {
		return len(clusters[i]) > len(clusters[j])
	})

	// Open the file
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Write the header
	file.WriteString("ClusterId,Latitude,Longitude\n")
	clusterId := 1
	for _, cluster := range clusters {
		for _, p := range cluster {
			// Write the cluster to the file
			file.WriteString(fmt.Sprintf("%d,%f,%f\n", clusterId, p.y, p.x))
		}
		clusterId++
	}
}

func savePoints(filename string, points []Point) {

	// Open the file
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Write the header
	file.WriteString("Id,Latitude,Longitude\n")

	// Write the clusters
	clusterId := 0
	for _, p := range points {
		// Get the cluster
		// Get the average point
		// Write the cluster to the file
		file.WriteString(fmt.Sprintf("%d,%f,%f\n", clusterId, p.y, p.x))
		clusterId++
	}
}
