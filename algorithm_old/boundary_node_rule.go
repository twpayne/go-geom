package algorithm

type BoundaryNodeRule interface {
	// Tests whether a point that lies in boundaryCount
	// geometry component boundaries is considered to form
	// part of the boundary of the parent geometry.
	isInBoundary(boundaryCount int) bool
}

// A BoundaryNodeRule specifies that points are in the
// boundary of a lineal geometry iff
// the point lies on the boundary of an odd number
// of components.
//
// Under this rule LinearRings and closed
// LineStrings have an empty boundary.
//
// This is the rule specified by the <i>OGC SFS</i>,
// and is the default rule used in JTS.
type Mod2BoundaryNodeRule struct{}

var _ BoundaryNodeRule = Mod2BoundaryNodeRule{}

func (rule Mod2BoundaryNodeRule) isInBoundary(boundaryCount int) bool {
	return boundaryCount%2 == 1
}

//  A BoundaryNodeRule which specifies that any points which are endpoints
// of lineal components are in the boundary of the
// parent geometry.
// This corresponds to the "intuitive" topological definition
// of boundary.
// Under this rule LinearRings have a non-empty boundary
// (the common endpoint of the underlying LineString).
//
// This rule is useful when dealing with linear networks.
// For example, it can be used to check
// whether linear networks are correctly noded.
// The usual network topology constraint is that linear segments may touch only at endpoints.
// In the case of a segment touching a closed segment (ring) at one point,
// the Mod2 rule cannot distinguish between the permitted case of touching at the
// node point and the invalid case of touching at some other interior (non-node) point.
// The EndPoint rule does distinguish between these cases,
// so is more appropriate for use.
type EndPointBoundaryNodeRule struct{}

var _ BoundaryNodeRule = EndPointBoundaryNodeRule{}

func (rule EndPointBoundaryNodeRule) isInBoundary(boundaryCount int) bool {
	return boundaryCount > 0
}

// A BoundaryNodeRule which determines that only
// endpoints with valency greater than 1 are on the boundary.
// This corresponds to the boundary of a MultiLineString
// being all the "attached" endpoints, but not
// the "unattached" ones.
type MultiValentEndPointBoundaryNodeRule struct{}

var _ BoundaryNodeRule = MultiValentEndPointBoundaryNodeRule{}

func (rule MultiValentEndPointBoundaryNodeRule) isInBoundary(boundaryCount int) bool {
	return boundaryCount > 1
}

// A BoundaryNodeRule which determines that only
// endpoints with valency of exactly 1 are on the boundary.
// This corresponds to the boundary of a MultiLineString
// being all the "unattached" endpoints.
type MonoValentEndPointBoundaryNodeRule struct{}

var _ BoundaryNodeRule = MonoValentEndPointBoundaryNodeRule{}

func (rule MonoValentEndPointBoundaryNodeRule) isInBoundary(boundaryCount int) bool {
	return boundaryCount == 1
}
