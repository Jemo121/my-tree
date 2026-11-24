// Package config holds compile-time constants and tuning parameters for the renderer.
package config

import "math"

// Rendering and layout constants.
const (
	TargetFPS = 60
	FrameMS   = 1000.0 / TargetFPS
	TwoPi     = math.Pi * 2

	// Canvas layout proportions (fraction of canvas height)
	GroundFrac = 0.72 // Y position of ground line
	CrownFrac  = 0.20 // Crown top starts here
	RootFrac   = 0.88 // Root zone ends here

	// Tree trunk
	TrunkBaseWidth  = 0.04 // fraction of canvas width
	TrunkTopWidth   = 0.015
	TrunkHeightFrac = 0.52 // trunk occupies this fraction of canvas height

	// Member glyphs
	GlyphRadius       = 18.0 // canvas pixels
	GlyphRadiusActive = 20.0
	GlyphPadding      = 8.0

	// Sign fractions (used by layout algorithm)
	SignWidthFrac  = 0.12
	SignHeightFrac = 0.08

	// Sign absolute offsets for hit detection (pixels, relative to tree trunk centre)
	// Negative SignOffsetX places the sign to the left of the trunk.
	SignOffsetX = -220.0
	SignOffsetY = 40.0
	SignWidth   = 200.0
	SignHeight  = 70.0

	// Weather particles
	MaxRainDrops  = 300
	MaxWindLines  = 80
	MaxCloudPuffs = 6

	// Animals
	MaxBirds       = 5
	MaxSquirrels   = 2
	MaxButterflies = 4
	MaxDeer        = 1

	// Thunder timing
	ThunderMinSec = 30.0
	ThunderMaxSec = 120.0
	FlashDurMS    = 300.0
)
