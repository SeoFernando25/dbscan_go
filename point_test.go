package main

import (
	"math"
	"math/rand"
	"testing"
)

func TestPointDistance(t *testing.T) {
	for i := 0; i < 100; i++ {
		rand1 := rand.Float64() * 100
		rand2 := rand.Float64() * 100
		slope := 1.0
		p1 := Point{rand1, rand1}
		p2 := Point{rand1 + slope*rand2, rand1 + slope*rand2}
		dx := p2.x - p1.x
		dy := p2.y - p1.y
		dist := math.Sqrt(dx*dx + dy*dy)
		if p1.Distance(p2) != dist {
			t.Error("Distance function is not correct")
		}
	}
}
