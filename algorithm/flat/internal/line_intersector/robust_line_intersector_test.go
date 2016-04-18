package line_intersector_test

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/algorithm/flat/internal/line_intersector"
	"reflect"
	"runtime/debug"
	"testing"
)

func TestRobustLineIntersectionPointOnLine(t *testing.T) {
	exectuteLineIntersectionPointOnLineTest(t, line_intersector.RobustLineIntersector{})
}

func TestRobustLineIntersectionLines(t *testing.T) {
	executeLineIntersectionLinesTest(t, line_intersector.RobustLineIntersector{})

	// extra tests that have different values for robust intersector than non-robust

	for i, tc := range line_intersector.ROBUST_LINE_ON_LINE_INTERSECTION_DATA {
		doLineIntersectsLineTest(t, line_intersector.RobustLineIntersector{}, i, tc)
	}
}

func exectuteLineIntersectionPointOnLineTest(t *testing.T, intersectionStrategy line_intersector.LineIntersectionStrategy) {
	for i, tc := range line_intersector.POINT_ON_LINE_INTERSECTION_TEST_DATA {
		intersector := line_intersector.LineIntersector{Layout: geom.XY, Strategy: intersectionStrategy}
		calculatedResult := intersector.PointIntersectsLine(tc.P, tc.LineEnd1, tc.LineEnd2)
		if !reflect.DeepEqual(tc.Result, calculatedResult) {
			t.Errorf("Test '%v' failed: expected \n%v was \n%v", i+1, tc.Result, calculatedResult)
		}
	}
}

func executeLineIntersectionLinesTest(t *testing.T, intersectionStrategy line_intersector.LineIntersectionStrategy) {
	for i, tc := range line_intersector.LINE_ON_LINE_INTERSECTION_TEST_DATA {
		doLineIntersectsLineTest(t, intersectionStrategy, i, tc)
	}
}

func doLineIntersectsLineTest(t *testing.T, intersectionStrategy line_intersector.LineIntersectionStrategy, i int, tc line_intersector.LineIntersectsLinesTestData) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("%T - An error occurred during Test '%v' (%v): %v\n%s", intersectionStrategy, i+1, tc.Desc, err, debug.Stack())
		}
	}()

	intersector := line_intersector.LineIntersector{Layout: geom.XY, Strategy: intersectionStrategy}
	calculatedResult := intersector.LineIntersectsLine(tc.P1, tc.P2, tc.P3, tc.P4)

	if !reflect.DeepEqual(calculatedResult, tc.Result) {
		t.Errorf("%T - Test '%v' (%v) failed: expected \n%v but was \n%v", intersectionStrategy, i+1, tc.Desc, tc.Result, calculatedResult)
	}
}
