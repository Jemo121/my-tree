//go:build js && wasm

package render

import (
	"math"

	"eguia.me/wasm/internal/config"
	"eguia.me/wasm/internal/domain"
)

// DrawTree renders the trunk, crown mass, root zone, and rocks.
// shakeOff is applied only to the trunk and crown — the root zone is anchored to the ground.
func DrawTree(ctx *Ctx, layout *domain.TreeLayout, cs config.ColorSet, shakeOff domain.Vec2) {
	drawRootZone(ctx, layout, cs)

	ctx.Save()
	ctx.Translate(shakeOff.X, shakeOff.Y)
	drawTrunk(ctx, layout, cs)
	drawCrown(ctx, layout, cs)
	ctx.Restore()
}

// DrawMembers renders the member glyphs in crown (active) and roots (inactive).
func DrawMembers(ctx *Ctx, layout *domain.TreeLayout, cs config.ColorSet) {
	for _, m := range layout.ActiveMembers {
		drawActiveMember(ctx, m, cs)
	}
	for _, m := range layout.InactMembers {
		drawInactiveMember(ctx, m, cs)
	}
}

func drawTrunk(ctx *Ctx, layout *domain.TreeLayout, cs config.ColorSet) {
	bx, by := layout.TrunkBase.X, layout.TrunkBase.Y
	tx, ty := layout.TrunkTop.X, layout.TrunkTop.Y
	cw := layout.CanvasW
	bw := cw * config.TrunkBaseWidth
	tw := cw * config.TrunkTopWidth

	trunk := ctx.LinearGradient(bx-bw, 0, bx+bw, 0)
	AddColorStop(trunk, 0, cs.TrunkDark)
	AddColorStop(trunk, 0.35, cs.TrunkLight)
	AddColorStop(trunk, 0.65, cs.TrunkLight)
	AddColorStop(trunk, 1, cs.TrunkDark)
	ctx.SetFillGradient(trunk)

	// Natural trunk shape via bezier
	ctx.BeginPath()
	ctx.MoveTo(bx-bw/2, by)
	ctx.BezierCurveTo(bx-bw/2, by-30, tx-tw/2, ty+40, tx-tw/2, ty)
	ctx.LineTo(tx+tw/2, ty)
	ctx.BezierCurveTo(tx+tw/2, ty+40, bx+bw/2, by-30, bx+bw/2, by)
	ctx.ClosePath()
	ctx.Fill()

	// Bark texture lines
	ctx.SetStroke(cs.TrunkDark)
	ctx.SetAlpha(0.3)
	ctx.SetLineWidth(1)
	for i := 0; i < 6; i++ {
		t := float64(i+1) / 7.0
		lx := bx - bw/2*(1-t) + pseudoRand(i*5)*bw*0.3
		ctx.BeginPath()
		ctx.MoveTo(lx, by*(1-t)+ty*t)
		ctx.QuadraticCurveTo(lx+pseudoRand(i*5+1)*6-3, (by*(1-t)+ty*t)-20, lx+pseudoRand(i*5+2)*4-2, (by*(1-t)+ty*t)-35)
		ctx.Stroke()
	}
	ctx.SetAlpha(1)
}

func drawCrown(ctx *Ctx, layout *domain.TreeLayout, cs config.ColorSet) {
	cr := layout.CrownBounds
	cx := (cr.Min.X + cr.Max.X) / 2
	cy := (cr.Min.Y + cr.Max.Y) / 2
	r := (cr.Max.X - cr.Min.X) / 2

	// Number of puffs scales with active member count
	puffs := 7 + layout.ActiveCount*2
	if puffs > 24 {
		puffs = 24
	}

	for i := 0; i < puffs; i++ {
		angle := float64(i) * math.Pi * 2 / float64(puffs)
		spread := 0.6 + pseudoRand(i*4)*0.4
		pr := r * spread
		px := cx + pr*math.Cos(angle)
		py := cy + pr*math.Sin(angle)*0.75
		puffR := r*0.38 + pseudoRand(i*4+1)*r*0.18

		// alternate between primary and secondary leaf colours
		color := cs.LeafPri
		if i%3 == 1 {
			color = cs.LeafSec
		}

		ctx.SetFill(color)
		ctx.SetAlpha(0.85 + pseudoRand(i*4+2)*0.15)
		ctx.BeginPath()
		ctx.Arc(px, py, puffR, 0, math.Pi*2)
		ctx.Fill()
	}

	// Core fill to remove gaps
	ctx.SetFill(cs.LeafPri)
	ctx.SetAlpha(0.9)
	ctx.BeginPath()
	ctx.Arc(cx, cy, r*0.55, 0, math.Pi*2)
	ctx.Fill()
	ctx.SetAlpha(1)
}

func drawRootZone(ctx *Ctx, layout *domain.TreeLayout, cs config.ColorSet) {
	sb := layout.SoilBounds
	cx := (sb.Min.X + sb.Max.X) / 2
	w := sb.Max.X - sb.Min.X
	h := sb.Max.Y - sb.Min.Y

	// No filled rectangle — the background ground already paints this area.
	// A vignette gradient darkens the underground zone without hard edges.
	vignette := ctx.LinearGradient(0, sb.Min.Y, 0, sb.Max.Y)
	AddColorStop(vignette, 0, "rgba(0,0,0,0)")
	AddColorStop(vignette, 1, "rgba(0,0,0,0.35)")
	ctx.SetFillGradient(vignette)
	ctx.FillRect(0, sb.Min.Y, ctx.W, h)

	// Roots (sinuous lines radiating from trunk base)
	ctx.SetStroke(cs.RootCol)
	ctx.SetLineWidth(3)
	ctx.SetLineCap("round")
	for i := 0; i < 5+layout.InactCount; i++ {
		angle := math.Pi*0.15 + float64(i)*math.Pi*0.7/float64(4+layout.InactCount)
		rx := cx + math.Cos(math.Pi+angle)*w*0.4
		ry := sb.Min.Y + h*(0.4+pseudoRand(i*3)*0.4)
		ctx.BeginPath()
		ctx.MoveTo(cx, sb.Min.Y+5)
		ctx.BezierCurveTo(
			cx+pseudoRand(i*3+1)*20-10, sb.Min.Y+h*0.3,
			rx+pseudoRand(i*3+2)*15-7, ry-h*0.1,
			rx, ry,
		)
		ctx.SetAlpha(0.6 + pseudoRand(i*3)*0.3)
		ctx.Stroke()
	}
	ctx.SetAlpha(1)

	// Underground rocks for inactive members
	for i, m := range layout.InactMembers {
		drawRock(ctx, m.Position.X+pseudoRand(i*2)*10-5, m.Position.Y+15, cs)
	}
}

func drawRock(ctx *Ctx, x, y float64, cs config.ColorSet) {
	r := 8.0 + pseudoRand(int(x)*7)*4
	ctx.BeginPath()
	ctx.Arc(x, y, r, 0, math.Pi*2)
	ctx.SetFill(cs.RockCol)
	ctx.SetAlpha(0.7)
	ctx.Fill()
	ctx.SetAlpha(1)
}

func drawActiveMember(ctx *Ctx, m *domain.Member, _ config.ColorSet) {
	p := m.Position
	r := config.GlyphRadiusActive
	ctx.Save()
	ctx.Translate(p.X, p.Y)

	drawFruitGlyph(ctx, r, m.Icon)

	ctx.Restore()
}

func drawInactiveMember(ctx *Ctx, m *domain.Member, cs config.ColorSet) {
	p := m.Position
	r := config.GlyphRadius
	ctx.SetAlpha(0.45)
	ctx.BeginPath()
	ctx.Arc(p.X, p.Y, r, 0, math.Pi*2)
	ctx.SetFill(cs.RootCol)
	ctx.Fill()
	// small ring to distinguish
	ctx.SetStroke(cs.RockCol)
	ctx.SetLineWidth(1.5)
	ctx.Stroke()
	ctx.SetAlpha(1)
}
