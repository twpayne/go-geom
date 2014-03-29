package geom

import (
	"testing"

	. "launchpad.net/gocheck"
)

func Test(t *testing.T) { TestingT(t) }

type aliasesChecker struct {
	*CheckerInfo
}

var Aliases Checker = &aliasesChecker{
	&CheckerInfo{Name: "Aliases", Params: []string{"obtained", "expected"}},
}

func (checker *aliasesChecker) Check(params []interface{}, names []string) (result bool, error string) {
	x, ok := params[0].([]float64)
	if !ok {
		return false, "obtained value is not a []float64"
	}
	y, ok := params[1].([]float64)
	if !ok {
		return false, "expected value is not a []float64"
	}
	// http://golang.org/src/pkg/math/big/nat.go#L340
	return cap(x) > 0 && cap(y) > 0 && &x[0:cap(x)][cap(x)-1] == &y[0:cap(y)][cap(y)-1], ""
}
