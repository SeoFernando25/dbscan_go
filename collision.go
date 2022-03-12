package main

// Rect to rect collision detection
func rectIntersect(rect1 Rect, rect2 Rect) bool {
	left := rect1.x + rect1.w
	right := rect2.x + rect2.w
	top := rect1.y + rect1.h
	bottom := rect2.y + rect2.h
	return rect1.x <= right && left >= rect2.x && rect1.y <= bottom && top >= rect2.y
}

// Rect to Point collision detection
func rectPointIntersect(r Rect, p Point) bool {
	return r.x <= p.x && r.x+r.w >= p.x && r.y <= p.y && r.y+r.h >= p.y
}

// Point to Point collision detection
func pointIntersect(p1 Point, p2 Point) bool {
	return p1.x == p2.x && p1.y == p2.y
}
