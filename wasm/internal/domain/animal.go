package domain

// AnimalType identifies the kind of animal.
type AnimalType int

const (
	AnimalBird AnimalType = iota
	AnimalSquirrel
	AnimalDeer
	AnimalButterfly
)

// AnimalState is the current behaviour mode of an animal.
type AnimalState int

const (
	AnimalIdle AnimalState = iota
	AnimalMoving
	AnimalFleeing
	AnimalPerching
)

// Animal is a single living creature on the tree scene.
type Animal struct {
	ID          int
	Type        AnimalType
	State       AnimalState
	Position    Vec2
	Velocity    Vec2
	Target      Vec2
	SpriteFrame int
	FrameTimer  float64
	PerchNode   *TreeNode
	GroundY     float64
	SeasonMask  [4]bool // indexed by Season constant
	Alpha       float64 // 0–1 for fade in/out
}

// SeasonFor reports whether this animal should appear in the given season.
func (a *Animal) SeasonFor(s Season) bool {
	return a.SeasonMask[int(s)]
}
