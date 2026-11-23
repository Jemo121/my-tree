package domain

// TreeNode is one node in the rendered tree structure.
type TreeNode struct {
	Member      *Member
	Children    []*TreeNode
	Parent      *TreeNode
	Position    Vec2
	BranchAngle float64
	BranchLen   float64
	Radius      float64 // trunk/branch thickness at this node
	Depth       int
}

// TreeLayout is the computed spatial representation of the tree.
type TreeLayout struct {
	Root          *TreeNode
	Nodes         []*TreeNode
	ActiveMembers []*Member
	InactMembers  []*Member

	TrunkBase   Vec2
	TrunkTop    Vec2
	TrunkBounds Rect
	CrownBounds Rect
	RootBounds  Rect
	SoilBounds  Rect
	SignBounds  Rect

	CanvasW float64
	CanvasH float64
	GroundY float64

	ActiveCount int
	InactCount  int
}

// MemberAt returns the member whose bounds contain p, or nil.
func (t *TreeLayout) MemberAt(p Vec2) *Member {
	for _, m := range t.ActiveMembers {
		if m.Bounds.Contains(p) {
			return m
		}
	}
	for _, m := range t.InactMembers {
		if m.Bounds.Contains(p) {
			return m
		}
	}
	return nil
}
