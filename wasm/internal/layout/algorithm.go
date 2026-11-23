package layout

import (
	"math"

	"eguia.me/wasm/internal/config"
	"eguia.me/wasm/internal/domain"
)

// Build computes the full TreeLayout from a list of members and canvas dimensions.
func Build(members []*domain.Member, cw, ch float64) *domain.TreeLayout {
	layout := &domain.TreeLayout{
		CanvasW: cw,
		CanvasH: ch,
		GroundY: ch * config.GroundFrac,
	}

	// Classify members
	for _, m := range members {
		if m.IsActive() {
			layout.ActiveMembers = append(layout.ActiveMembers, m)
		} else {
			layout.InactMembers = append(layout.InactMembers, m)
		}
	}
	layout.ActiveCount = len(layout.ActiveMembers)
	layout.InactCount = len(layout.InactMembers)

	// Trunk geometry
	trunkBaseX := cw / 2
	trunkBaseY := layout.GroundY + ch*0.02
	trunkTopY := ch * config.CrownFrac
	layout.TrunkBase = domain.Vec2{X: trunkBaseX, Y: trunkBaseY}
	layout.TrunkTop = domain.Vec2{X: trunkBaseX, Y: trunkTopY}
	hw := cw * config.TrunkBaseWidth / 2
	layout.TrunkBounds = domain.Rect{
		Min: domain.Vec2{X: trunkBaseX - hw, Y: trunkTopY},
		Max: domain.Vec2{X: trunkBaseX + hw, Y: trunkBaseY},
	}

	// Crown bounds
	crownRadius := crownSize(layout.ActiveCount, cw, ch)
	cx, cy := trunkBaseX, trunkTopY+crownRadius*0.6
	layout.CrownBounds = domain.Rect{
		Min: domain.Vec2{X: cx - crownRadius, Y: cy - crownRadius},
		Max: domain.Vec2{X: cx + crownRadius, Y: cy + crownRadius},
	}

	// Root bounds (below ground)
	rootRadius := rootSize(layout.InactCount, cw, ch)
	ry := trunkBaseY + ch*0.04
	layout.RootBounds = domain.Rect{
		Min: domain.Vec2{X: cx - rootRadius, Y: ry},
		Max: domain.Vec2{X: cx + rootRadius, Y: ry + rootRadius*1.5},
	}
	layout.SoilBounds = domain.Rect{
		Min: domain.Vec2{X: cx - cw*0.35, Y: trunkBaseY},
		Max: domain.Vec2{X: cx + cw*0.35, Y: ch * config.RootFrac},
	}

	// Wooden sign (right of trunk by default)
	signW, signH := cw*config.SignWidthFrac, ch*config.SignHeightFrac
	signX := trunkBaseX + cw*0.2
	signY := trunkBaseY - signH - ch*0.05
	layout.SignBounds = domain.Rect{
		Min: domain.Vec2{X: signX, Y: signY},
		Max: domain.Vec2{X: signX + signW, Y: signY + signH},
	}

	// Place active members in crown
	placeActiveMembers(layout, cx, cy, crownRadius)

	// Place inactive members in roots
	placeInactiveMembers(layout, cx, ry, rootRadius)

	return layout
}

// crownSize returns the crown radius based on member count and canvas dimensions.
func crownSize(n int, cw, ch float64) float64 {
	base := ch * 0.22
	growth := math.Log(float64(n+2)) * ch * 0.06
	max := cw * 0.42
	r := base + growth
	if r > max {
		r = max
	}
	return r
}

// rootSize returns the root spread radius.
func rootSize(n int, cw, ch float64) float64 {
	base := ch * 0.10
	growth := math.Log(float64(n+2)) * ch * 0.03
	max := cw * 0.28
	r := base + growth
	if r > max {
		r = max
	}
	return r
}

// placeActiveMembers distributes active members concentrically within the crown.
func placeActiveMembers(l *domain.TreeLayout, cx, cy, r float64) {
	members := l.ActiveMembers
	if len(members) == 0 {
		return
	}

	// Group by generation
	byGen := groupByGeneration(members)
	maxGen := 0
	for g := range byGen {
		if g > maxGen {
			maxGen = g
		}
	}

	for gen, group := range byGen {
		if len(group) == 0 {
			continue
		}
		// Innermost generation closest to trunk; outermost near edge
		frac := 0.3
		if maxGen > 1 {
			frac = 0.3 + 0.65*float64(gen-1)/float64(maxGen-1)
		}
		ringR := r * frac
		spread(group, cx, cy, ringR, config.GlyphRadius)
	}
}

// placeInactiveMembers distributes inactive members in the root zone.
func placeInactiveMembers(l *domain.TreeLayout, cx, startY, r float64) {
	members := l.InactMembers
	if len(members) == 0 {
		return
	}
	spread(members, cx, startY+r*0.5, r*0.7, config.GlyphRadius)
}

// spread arranges members evenly around a circle of radius r centred at (cx, cy).
func spread(members []*domain.Member, cx, cy, r, glyphR float64) {
	n := len(members)
	for i, m := range members {
		angle := math.Pi*1.5 + float64(i)*math.Pi*2/float64(n)
		if n == 1 {
			angle = math.Pi * 1.5
		}
		px := cx + r*math.Cos(angle)
		py := cy + r*math.Sin(angle)
		m.Position = domain.Vec2{X: px, Y: py}
		m.Bounds = domain.RectFromCenter(m.Position, glyphR+config.GlyphPadding, glyphR+config.GlyphPadding)
	}
}

// groupByGeneration groups members by their Generation field.
func groupByGeneration(members []*domain.Member) map[int][]*domain.Member {
	out := make(map[int][]*domain.Member)
	for _, m := range members {
		out[m.Generation] = append(out[m.Generation], m)
	}
	return out
}
