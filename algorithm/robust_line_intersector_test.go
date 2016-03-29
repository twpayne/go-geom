package algorithm

import (
	"testing"
)

func TestRobustLineIntersectionPointOnLine(t *testing.T) {
	exectuteLineIntersectionPointOnLineTest(t, RobustLineIntersector{})
}

func TestRobustLineIntersectionLines(t *testing.T) {
	executeLineIntersectionLinesTest(t, RobustLineIntersector{})
}
