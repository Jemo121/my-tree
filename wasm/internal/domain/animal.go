// Package domain defines the core types for the family tree scene: members, animals, weather, and time.
package domain

// AnimalType identifies the kind of animal.
type AnimalType int

// AnimalType values.
const (
	AnimalBird      AnimalType = iota
	AnimalSquirrel             // tree-bound, year-round
	AnimalDeer                 // rare ground appearance, slow walk across
	AnimalButterfly            // passive drift across canvas, any time of day
	AnimalFirefly              // night-only, glow/pulse effect, slow drift
	AnimalFox                  // rare ground appearance, non-winter
	AnimalSnowFox              // rare ground appearance, winter only (white coat)
)

// SpawnProb constants (per tick).
const (
	SpawnProbButterfly = 0.002  // per tick
	SpawnProbFirefly   = 0.003  // per tick, night only
	SpawnProbDeer      = 0.0001 // rare; at most one active at a time
	SpawnProbFox       = 0.0001 // rare; at most one Fox/SnowFox active at a time
)

// CanSpawnFirefly returns true only during nighttime.
func CanSpawnFirefly(isNight bool) bool { return isNight }

// CanSpawnFox returns true outside of winter season.
func CanSpawnFox(season Season) bool { return season != SeasonWinter }

// CanSpawnSnowFox returns true only during winter.
func CanSpawnSnowFox(season Season) bool { return season == SeasonWinter }

// AnimalState is the current behaviour mode of an animal.
type AnimalState int

// AnimalState values.
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
