package line_intersector_test

import (
	"github.com/twpayne/go-geom/algorithm/internal/line_intersector"
	"testing"
)

func TestNonRobustLineIntersectionPointOnLine(t *testing.T) {
	exectuteLineIntersectionPointOnLineTest(t, line_intersector.NonRobustLineIntersector{})
}

func TestNonRobustLineIntersectionLines(t *testing.T) {
	executeLineIntersectionLinesTest(t, line_intersector.NonRobustLineIntersector{})
}
