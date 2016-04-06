package algorithm_test

import (
	"github.com/twpayne/go-geom/algorithm"
	"testing"
)

func TestNonRobustLineIntersectionPointOnLine(t *testing.T) {
	exectuteLineIntersectionPointOnLineTest(t, algorithm.NonRobustLineIntersector{})
}

func TestNonRobustLineIntersectionLines(t *testing.T) {
	executeLineIntersectionLinesTest(t, algorithm.NonRobustLineIntersector{})
}
