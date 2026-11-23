//go:build js && wasm

package render

import "eguia.me/wasm/internal/domain"

// ShakeKind is the type of shake to apply.
type ShakeKind int

const (
	ShakeTrunk ShakeKind = iota
	ShakeLocal
	ShakeEarth
)

// Shake is an in-progress shake effect at a point.
type Shake struct {
	Kind      ShakeKind
	Origin    domain.Vec2
	Magnitude float64 // current amplitude in canvas pixels
	OffsetX   float64
	OffsetY   float64
}

// ShakeSystem manages active shake effects.
type ShakeSystem struct {
	shakes []Shake
}

// Add starts a new shake at origin.
func (s *ShakeSystem) Add(kind ShakeKind, origin domain.Vec2, mag float64) {
	s.shakes = append(s.shakes, Shake{Kind: kind, Origin: origin, Magnitude: mag})
}

const shakeDecay = 0.82

// Tick advances all active shakes, decaying their magnitude.
func (s *ShakeSystem) Tick() {
	live := s.shakes[:0]
	for i := range s.shakes {
		sh := &s.shakes[i]
		sh.Magnitude *= shakeDecay
		if sh.Magnitude < 0.3 {
			continue
		}
		m := sh.Magnitude
		sh.OffsetX = (pseudoRand(i*2)*2 - 1) * m
		sh.OffsetY = (pseudoRand(i*2+1)*2 - 1) * m * 0.5
		live = append(live, *sh)
	}
	s.shakes = live
}

// Offset returns the total shake offset to apply to canvas translate.
func (s *ShakeSystem) Offset() (dx, dy float64) {
	for _, sh := range s.shakes {
		dx += sh.OffsetX
		dy += sh.OffsetY
	}
	return
}

// HasShake reports whether any shake is active.
func (s *ShakeSystem) HasShake() bool { return len(s.shakes) > 0 }

// pseudoRand returns a deterministic pseudo-random value in [0,1) for index i.
// Used to avoid importing math/rand in the render hot path.
func pseudoRand(i int) float64 {
	x := uint64(i)*6364136223846793005 + 1442695040888963407
	x ^= x >> 33
	x *= 0xff51afd7ed558ccd
	x ^= x >> 33
	return float64(x&0xfffff) / float64(0xfffff)
}
