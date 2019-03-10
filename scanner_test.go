package main

import (
	"testing"
	"time"
)

func TestMapStringToPoint(t *testing.T) {
	s := []string{"10", "37.969647", "23.722431", "1405588782"}
	outputPoint := Point{
		10,
		37.969647,
		23.722431,
		time.Unix(1405588782, 0)}

	if mapStringToPoint(s) != outputPoint {
		t.Error("Expected:", outputPoint)
	}
}
