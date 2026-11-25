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

func drawActiveMember(ctx *Ctx, m *domain.Member, cs config.ColorSet) {
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



// drawFruitGlyph renders the member's chosen fruit at radius r centred on the origin.
func drawFruitGlyph(ctx *Ctx, r float64, icon domain.MemberIcon) {
	switch icon {
	case domain.IconApple:
		fruitRound(ctx, r, "#c0392b", true)
	case domain.IconPear:
		fruitPear(ctx, r, "#a8c94b")
	case domain.IconQuince:
		fruitPear(ctx, r, "#e3b84a")
	case domain.IconCherry:
		fruitCherry(ctx, r, "#8b1a1a")
	case domain.IconPlum:
		fruitRound(ctx, r, "#6a0dad", false)
	case domain.IconPeach:
		fruitCreased(ctx, r, "#f4a460", "#e8724a")
	case domain.IconApricot:
		fruitCreased(ctx, r, "#f0a030", "#e07820")
	case domain.IconMango:
		fruitMango(ctx, r)
	case domain.IconAvocado:
		fruitAvocado(ctx, r)
	case domain.IconOrange:
		fruitRound(ctx, r, "#e67e22", true)
	case domain.IconLemon:
		fruitLemon(ctx, r, "#f9e547")
	case domain.IconLime:
		fruitLemon(ctx, r, "#7dc843")
	case domain.IconGrapefruit:
		fruitRound(ctx, r, "#f4a0a0", true)
	case domain.IconPomelo:
		fruitRound(ctx, r, "#d4e08a", true)
	case domain.IconCoconut:
		fruitCoconut(ctx, r)
	case domain.IconFig:
		fruitFig(ctx, r)
	case domain.IconDate:
		fruitDate(ctx, r)
	case domain.IconLychee:
		fruitLychee(ctx, r)
	case domain.IconPomegranate:
		fruitPomegranate(ctx, r)
	case domain.IconPersimmon:
		fruitPersimmon(ctx, r)
	default:
		fruitRound(ctx, r, "#c0392b", true) // fallback: apple
	}
}

// stem draws a short brown stem above the origin.
func stem(ctx *Ctx, r float64) {
	ctx.SetStroke("#5d4037")
	ctx.SetLineWidth(2)
	ctx.BeginPath()
	ctx.MoveTo(0, -r*0.85)
	ctx.LineTo(0, -r*0.5)
	ctx.Stroke()
}

// shine draws a small highlight circle at upper-left.
func shine(ctx *Ctx, r float64) {
	ctx.SetFill("rgba(255,255,255,0.28)")
	ctx.BeginPath()
	ctx.Arc(-r*0.28, -r*0.28, r*0.22, 0, math.Pi*2)
	ctx.Fill()
}

// fruitRound draws a simple round fruit. withDimple adds a top indentation.
func fruitRound(ctx *Ctx, r float64, color string, withDimple bool) {
	ctx.SetFill(color)
	ctx.BeginPath()
	ctx.Arc(0, r*0.08, r*0.82, 0, math.Pi*2)
	ctx.Fill()
	if withDimple {
		ctx.SetFill("rgba(0,0,0,0.12)")
		ctx.BeginPath()
		ctx.Arc(0, -r*0.6, r*0.18, 0, math.Pi*2)
		ctx.Fill()
	}
	stem(ctx, r)
	shine(ctx, r)
}

// fruitPear draws a teardrop pear/quince shape.
func fruitPear(ctx *Ctx, r float64, color string) {
	ctx.SetFill(color)
	ctx.BeginPath()
	ctx.MoveTo(0, -r*0.75)
	ctx.BezierCurveTo(r*0.6, -r*0.5, r*0.9, r*0.1, r*0.5, r*0.65)
	ctx.BezierCurveTo(r*0.2, r*0.9, -r*0.2, r*0.9, -r*0.5, r*0.65)
	ctx.BezierCurveTo(-r*0.9, r*0.1, -r*0.6, -r*0.5, 0, -r*0.75)
	ctx.ClosePath()
	ctx.Fill()
	stem(ctx, r)
	shine(ctx, r)
}

// fruitCherry draws a small cherry with a long stem.
func fruitCherry(ctx *Ctx, r float64, color string) {
	cr := r * 0.6 // cherry ball radius
	ctx.SetFill(color)
	ctx.BeginPath()
	ctx.Arc(0, r*0.3, cr, 0, math.Pi*2)
	ctx.Fill()
	// long stem curves up
	ctx.SetStroke("#5d4037")
	ctx.SetLineWidth(2)
	ctx.BeginPath()
	ctx.MoveTo(0, r*0.3-cr)
	ctx.BezierCurveTo(r*0.1, -r*0.2, r*0.3, -r*0.6, 0, -r*0.85)
	ctx.Stroke()
	shine(ctx, r*0.9)
}

// fruitCreased draws a peach/apricot — round with a vertical crease.
func fruitCreased(ctx *Ctx, r float64, body, crease string) {
	ctx.SetFill(body)
	ctx.BeginPath()
	ctx.Arc(0, r*0.08, r*0.82, 0, math.Pi*2)
	ctx.Fill()
	// crease line
	ctx.SetStroke(crease)
	ctx.SetLineWidth(1.5)
	ctx.BeginPath()
	ctx.MoveTo(0, -r*0.65)
	ctx.BezierCurveTo(r*0.08, 0, r*0.08, r*0.4, 0, r*0.72)
	ctx.Stroke()
	stem(ctx, r)
	shine(ctx, r)
}

// fruitMango draws a kidney-shaped mango.
func fruitMango(ctx *Ctx, r float64) {
	ctx.SetFill("#f9a825")
	ctx.BeginPath()
	ctx.MoveTo(0, -r*0.8)
	ctx.BezierCurveTo(r*0.9, -r*0.5, r*0.9, r*0.6, r*0.1, r*0.85)
	ctx.BezierCurveTo(-r*0.5, r*0.9, -r*0.9, r*0.4, -r*0.6, -r*0.3)
	ctx.BezierCurveTo(-r*0.4, -r*0.7, -r*0.1, -r*0.9, 0, -r*0.8)
	ctx.ClosePath()
	ctx.Fill()
	// blush
	ctx.SetFill("rgba(220,60,20,0.25)")
	ctx.BeginPath()
	ctx.Arc(r*0.3, 0, r*0.45, 0, math.Pi*2)
	ctx.Fill()
	stem(ctx, r)
	shine(ctx, r)
}

// fruitAvocado draws a dark pear-shaped avocado.
func fruitAvocado(ctx *Ctx, r float64) {
	ctx.SetFill("#2d5a1b")
	ctx.BeginPath()
	ctx.MoveTo(0, -r*0.85)
	ctx.BezierCurveTo(r*0.5, -r*0.5, r*0.85, r*0.1, r*0.5, r*0.7)
	ctx.BezierCurveTo(r*0.2, r*0.95, -r*0.2, r*0.95, -r*0.5, r*0.7)
	ctx.BezierCurveTo(-r*0.85, r*0.1, -r*0.5, -r*0.5, 0, -r*0.85)
	ctx.ClosePath()
	ctx.Fill()
	// seed
	ctx.SetFill("#8b4513")
	ctx.BeginPath()
	ctx.Arc(0, r*0.25, r*0.3, 0, math.Pi*2)
	ctx.Fill()
	stem(ctx, r)
}

// fruitLemon draws an oval with pointed ends.
func fruitLemon(ctx *Ctx, r float64, color string) {
	ctx.SetFill(color)
	ctx.BeginPath()
	ctx.MoveTo(0, -r*0.9)
	ctx.BezierCurveTo(r*0.7, -r*0.8, r*0.85, r*0.4, 0, r*0.9)
	ctx.BezierCurveTo(-r*0.85, r*0.4, -r*0.7, -r*0.8, 0, -r*0.9)
	ctx.ClosePath()
	ctx.Fill()
	stem(ctx, r)
	shine(ctx, r)
}

// fruitCoconut draws a brown oval with a fibrous ring.
func fruitCoconut(ctx *Ctx, r float64) {
	ctx.SetFill("#795548")
	ctx.BeginPath()
	ctx.Arc(0, 0, r*0.82, 0, math.Pi*2)
	ctx.Fill()
	ctx.SetStroke("#4e342e")
	ctx.SetLineWidth(1.5)
	ctx.BeginPath()
	ctx.Arc(0, 0, r*0.55, 0, math.Pi*2)
	ctx.Stroke()
}

// fruitFig draws a dark teardrop with a flat bottom.
func fruitFig(ctx *Ctx, r float64) {
	ctx.SetFill("#6a1b6e")
	ctx.BeginPath()
	ctx.MoveTo(0, -r*0.85)
	ctx.BezierCurveTo(r*0.6, -r*0.5, r*0.75, r*0.3, r*0.6, r*0.75)
	ctx.LineTo(-r*0.6, r*0.75)
	ctx.BezierCurveTo(-r*0.75, r*0.3, -r*0.6, -r*0.5, 0, -r*0.85)
	ctx.ClosePath()
	ctx.Fill()
	// flat bottom seam
	ctx.SetStroke("#4a0050")
	ctx.SetLineWidth(1)
	ctx.BeginPath()
	ctx.MoveTo(-r*0.4, r*0.75)
	ctx.LineTo(r*0.4, r*0.75)
	ctx.Stroke()
	stem(ctx, r)
}

// fruitDate draws an elongated amber oval.
func fruitDate(ctx *Ctx, r float64) {
	ctx.SetFill("#bf8040")
	ctx.BeginPath()
	ctx.MoveTo(0, -r*0.9)
	ctx.BezierCurveTo(r*0.5, -r*0.7, r*0.5, r*0.7, 0, r*0.9)
	ctx.BezierCurveTo(-r*0.5, r*0.7, -r*0.5, -r*0.7, 0, -r*0.9)
	ctx.ClosePath()
	ctx.Fill()
	stem(ctx, r)
	shine(ctx, r)
}

// fruitLychee draws a small red sphere with bumpy skin.
func fruitLychee(ctx *Ctx, r float64) {
	ctx.SetFill("#c0392b")
	ctx.BeginPath()
	ctx.Arc(0, 0, r*0.82, 0, math.Pi*2)
	ctx.Fill()
	// bumps — small darker circles on surface
	ctx.SetFill("rgba(0,0,0,0.15)")
	for i, angle := range []float64{0, 1, 2, 3, 4, 5} {
		a := angle * math.Pi / 3
		bx := math.Cos(a) * r * 0.52
		by := math.Sin(a) * r * 0.52
		ctx.BeginPath()
		ctx.Arc(bx, by, r*0.18, 0, math.Pi*2)
		ctx.Fill()
		_ = i
	}
	stem(ctx, r)
	shine(ctx, r)
}

// fruitPomegranate draws a round red fruit with a prominent crown.
func fruitPomegranate(ctx *Ctx, r float64) {
	ctx.SetFill("#8b1a1a")
	ctx.BeginPath()
	ctx.Arc(0, r*0.1, r*0.78, 0, math.Pi*2)
	ctx.Fill()
	// crown — 5 small spikes at top
	ctx.SetFill("#4a0a0a")
	for i := 0; i < 5; i++ {
		a := (float64(i)/5)*math.Pi - math.Pi*0.5
		px := math.Cos(a) * r * 0.35
		py := -r*0.55 + math.Sin(a)*r*0.1
		ctx.BeginPath()
		ctx.Arc(px, py, r*0.1, 0, math.Pi*2)
		ctx.Fill()
	}
	shine(ctx, r)
}

// fruitPersimmon draws an orange round fruit with a 4-petal calyx.
func fruitPersimmon(ctx *Ctx, r float64) {
	ctx.SetFill("#e8792a")
	ctx.BeginPath()
	ctx.Arc(0, r*0.08, r*0.82, 0, math.Pi*2)
	ctx.Fill()
	// 4-leaf calyx at top
	ctx.SetFill("#4caf50")
	for i := 0; i < 4; i++ {
		a := float64(i)*math.Pi/2 + math.Pi/4
		lx := math.Cos(a) * r * 0.38
		ly := -r*0.65 + math.Sin(a)*r*0.18
		ctx.BeginPath()
		ctx.Arc(lx, ly, r*0.2, 0, math.Pi*2)
		ctx.Fill()
	}
	shine(ctx, r)
}
