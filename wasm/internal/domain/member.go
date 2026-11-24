package domain

// MemberStatus represents a family member's presence on the tree.
type MemberStatus string

// MemberStatus values.
const (
	StatusActive   MemberStatus = "active"
	StatusInactive MemberStatus = "inactive"
)

// MemberIcon determines which glyph represents the member in the crown.
type MemberIcon string

// MemberIcon values.
const (
	IconLeaf   MemberIcon = "leaf"
	IconFruit  MemberIcon = "fruit"
	IconFlower MemberIcon = "flower"
	IconTwig   MemberIcon = "twig"
)

// Member holds one family member's data plus runtime position state.
type Member struct {
	Name       string
	Subdomain  string
	Template   string
	Generation int
	Parent     *string
	Status     MemberStatus
	Icon       MemberIcon

	// Runtime-only: populated by the layout algorithm.
	NodeID   int
	Position Vec2
	Bounds   Rect
}

// IsActive reports whether this member is in the crown (active).
func (m *Member) IsActive() bool { return m.Status == StatusActive }
