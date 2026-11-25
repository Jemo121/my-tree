//go:build js && wasm

// Package render draws the canvas scene: tree, background, weather, animals, and UI.
package render

import (
	"math"

	"eguia.me/wasm/internal/domain"
)

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
	cr := r * 0.6
	ctx.SetFill(color)
	ctx.BeginPath()
	ctx.Arc(0, r*0.3, cr, 0, math.Pi*2)
	ctx.Fill()
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
	ctx.SetFill("rgba(0,0,0,0.15)")
	for i := 0; i < 6; i++ {
		a := float64(i) * math.Pi / 3
		ctx.BeginPath()
		ctx.Arc(math.Cos(a)*r*0.52, math.Sin(a)*r*0.52, r*0.18, 0, math.Pi*2)
		ctx.Fill()
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
	ctx.SetFill("#4a0a0a")
	for i := 0; i < 5; i++ {
		a := (float64(i)/5)*math.Pi - math.Pi*0.5
		ctx.BeginPath()
		ctx.Arc(math.Cos(a)*r*0.35, -r*0.55+math.Sin(a)*r*0.1, r*0.1, 0, math.Pi*2)
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
	ctx.SetFill("#4caf50")
	for i := 0; i < 4; i++ {
		a := float64(i)*math.Pi/2 + math.Pi/4
		ctx.BeginPath()
		ctx.Arc(math.Cos(a)*r*0.38, -r*0.65+math.Sin(a)*r*0.18, r*0.2, 0, math.Pi*2)
		ctx.Fill()
	}
	shine(ctx, r)
}
