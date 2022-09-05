package lineintersector

import (
	"reflect"
	"runtime/debug"
	"testing"

	"github.com/twpayne/go-geom/geomtest"
	"github.com/twpayne/go-geom/xy/lineintersection"
)

func TestRobustLineIntersectionPointOnLine(t *testing.T) {
	executeLineIntersectionPointOnLineTest(t, RobustLineIntersector{})
}

func TestRobustLineIntersectionLines(t *testing.T) {
	executeLineIntersectionLinesTest(t, RobustLineIntersector{})

	// extra tests that have different values for robust intersector than non-robust

	for i, tc := range robustLineOnLineIntersectionData {
		doLineIntersectsLineTest(t, RobustLineIntersector{}, i, tc)
	}
}

func executeLineIntersectionPointOnLineTest(t *testing.T, intersectionStrategy Strategy) {
	t.Helper()
	for i, tc := range pointOnLineIntersectionTestData {
		calculatedResult := PointIntersectsLine(intersectionStrategy, tc.P, tc.LineEnd1, tc.LineEnd2)
		if !reflect.DeepEqual(tc.Result, calculatedResult) {
			t.Errorf("Test '%v' failed: expected \n%v was \n%v", i+1, tc.Result, calculatedResult)
		}
	}
}

func executeLineIntersectionLinesTest(t *testing.T, intersectionStrategy Strategy) {
	t.Helper()
	for i, tc := range lineOnLineIntersectionTestData {
		doLineIntersectsLineTest(t, intersectionStrategy, i, tc)
	}
}

func doLineIntersectsLineTest(t *testing.T, intersectionStrategy Strategy, i int, tc lineIntersectsLinesTestData) {
	t.Helper()
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("%T - An error occurred during Test '%v' (%v): %v\n%s", intersectionStrategy, i+1, tc.Desc, err, debug.Stack())
		}
	}()

	calculatedResult := LineIntersectsLine(intersectionStrategy, tc.P1, tc.P2, tc.P3, tc.P4)

	if !lineInteresectionResultEqualsRel(calculatedResult, tc.Result, 1e-3) {
		t.Errorf("%T - Test '%v' (%v) failed: expected \n%v but was \n%v", intersectionStrategy, i+1, tc.Desc, tc.Result, calculatedResult)
	}
}

func lineInteresectionResultEqualsRel(r1, r2 lineintersection.Result, epsilon float64) bool {
	if r1.Type() != r2.Type() {
		return false
	}
	i1 := r1.Intersection()
	i2 := r2.Intersection()
	if len(i1) != len(i2) {
		return false
	}
	for i := range i1 {
		if !geomtest.CoordsEqualRel(i1[i], i2[i], epsilon) {
			return false
		}
	}
	return true
}
