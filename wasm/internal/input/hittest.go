package input

import (
	"eguia.me/wasm/internal/config"
	"eguia.me/wasm/internal/domain"
)

// HitSign returns true if (x, y) falls within the wooden sign bounds.
// signX, signY are the sign's top-left corner coordinates.
func HitSign(x, y, signX, signY float64) bool {
	return x >= signX && x <= signX+config.SignWidth &&
		y >= signY && y <= signY+config.SignHeight
}

const gridCell = 50.0

type gridKey struct{ X, Y int }

type hitEntry struct {
	target   domain.HitTarget
	priority int
}

// SpatialGrid is a bucketed spatial index for fast hit-testing.
type SpatialGrid struct {
	cells map[gridKey][]hitEntry
}

// NewSpatialGrid returns an empty grid.
func NewSpatialGrid() *SpatialGrid {
	return &SpatialGrid{cells: make(map[gridKey][]hitEntry)}
}

func (g *SpatialGrid) insert(e hitEntry) {
	r := e.target.Bounds
	x0, y0 := int(r.Min.X/gridCell), int(r.Min.Y/gridCell)
	x1, y1 := int(r.Max.X/gridCell), int(r.Max.Y/gridCell)
	for x := x0; x <= x1; x++ {
		for y := y0; y <= y1; y++ {
			k := gridKey{x, y}
			g.cells[k] = append(g.cells[k], e)
		}
	}
}

// Query returns the highest-priority target containing p, or a HitNone target.
func (g *SpatialGrid) Query(p domain.Vec2) domain.HitTarget {
	k := gridKey{int(p.X / gridCell), int(p.Y / gridCell)}
	entries := g.cells[k]
	best := domain.HitTarget{Type: domain.HitNone}
	bestPri := -1
	for _, e := range entries {
		if e.target.Bounds.Contains(p) && e.priority > bestPri {
			best = e.target
			bestPri = e.priority
		}
	}
	return best
}

// Build constructs a hit grid from the current layout and weather.
func Build(layout *domain.TreeLayout, weather *domain.WeatherState) *SpatialGrid {
	g := NewSpatialGrid()

	for _, m := range layout.ActiveMembers {
		g.insert(hitEntry{
			target:   domain.HitTarget{Type: domain.HitActiveMember, Member: m, Position: m.Position, Bounds: m.Bounds},
			priority: 100,
		})
	}
	for _, m := range layout.InactMembers {
		g.insert(hitEntry{
			target:   domain.HitTarget{Type: domain.HitInactiveMember, Member: m, Position: m.Position, Bounds: m.Bounds},
			priority: 90,
		})
	}
	g.insert(hitEntry{
		target:   domain.HitTarget{Type: domain.HitSign, Bounds: layout.SignBounds},
		priority: 80,
	})
	for i, puff := range weather.CloudPuffs {
		if puff.Vanishing {
			continue
		}
		g.insert(hitEntry{
			target:   domain.HitTarget{Type: domain.HitCloud, CloudIdx: i, Position: puff.Center, Bounds: puff.Bounds},
			priority: 70,
		})
	}
	g.insert(hitEntry{
		target:   domain.HitTarget{Type: domain.HitTrunk, Bounds: layout.TrunkBounds},
		priority: 60,
	})
	g.insert(hitEntry{
		target:   domain.HitTarget{Type: domain.HitRootSoil, Bounds: layout.SoilBounds},
		priority: 40,
	})

	return g
}
