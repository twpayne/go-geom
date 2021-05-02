package lineintersector

import "testing"

func TestNonRobustLineIntersectionPointOnLine(t *testing.T) {
	executeLineIntersectionPointOnLineTest(t, NonRobustLineIntersector{})
}

func TestNonRobustLineIntersectionLines(t *testing.T) {
	executeLineIntersectionLinesTest(t, NonRobustLineIntersector{})
}
