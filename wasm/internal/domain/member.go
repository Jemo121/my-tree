package domain

// MemberStatus represents a family member's presence on the tree.
type MemberStatus string

// MemberStatus values.
const (
	StatusActive   MemberStatus = "active"
	StatusInactive MemberStatus = "inactive"
)

// MemberIcon identifies which fruit glyph represents the member in the crown.
// Each member chooses their own fruit. Only tree fruits are valid.
type MemberIcon string

// MemberIcon values — pome, stone, citrus, tropical.
const (
	// Pome
	IconApple  MemberIcon = "apple"
	IconPear   MemberIcon = "pear"
	IconQuince MemberIcon = "quince"

	// Stone fruits
	IconCherry  MemberIcon = "cherry"
	IconPlum    MemberIcon = "plum"
	IconPeach   MemberIcon = "peach"
	IconApricot MemberIcon = "apricot"
	IconMango   MemberIcon = "mango"
	IconAvocado MemberIcon = "avocado"

	// Citrus
	IconOrange     MemberIcon = "orange"
	IconLemon      MemberIcon = "lemon"
	IconLime       MemberIcon = "lime"
	IconGrapefruit MemberIcon = "grapefruit"
	IconPomelo     MemberIcon = "pomelo"

	// Tropical / Mediterranean
	IconCoconut     MemberIcon = "coconut"
	IconFig         MemberIcon = "fig"
	IconDate        MemberIcon = "date"
	IconLychee      MemberIcon = "lychee"
	IconPomegranate MemberIcon = "pomegranate"
	IconPersimmon   MemberIcon = "persimmon"
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
