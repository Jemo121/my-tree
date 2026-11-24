//go:build js && wasm

// Package render draws the canvas scene: tree, background, weather, animals, and UI.
package render

import (
	"math"

	"eguia.me/wasm/internal/domain"
)

// DrawAnimals renders all active animals in the scene.
// Called from the frame pipeline to maintain correct render layer order.
func DrawAnimals(ctx *Ctx, animals []*domain.Animal, tod domain.TimeOfDay) {
	for _, a := range animals {
		if a.Alpha <= 0 {
			continue
		}
		ctx.SetAlpha(a.Alpha)
		switch a.Type {
		case domain.AnimalBird:
			drawBird(ctx, a)
		case domain.AnimalSquirrel:
			drawSquirrel(ctx, a)
		case domain.AnimalButterfly:
			drawButterfly(ctx, a)
		case domain.AnimalFirefly:
			if tod == domain.TimeNight || tod == domain.TimeDusk {
				drawFirefly(ctx, a)
			}
		}
		ctx.SetAlpha(1)
	}
}

// DrawRareAnimals renders rare ground animals (Deer, Fox, SnowFox).
// Split from DrawAnimals to allow separate layer ordering.
func DrawRareAnimals(ctx *Ctx, animals []*domain.Animal) {
	for _, a := range animals {
		if a.Alpha <= 0 {
			continue
		}
		ctx.SetAlpha(a.Alpha)
		switch a.Type {
		case domain.AnimalDeer:
			drawDeer(ctx, a)
		case domain.AnimalFox:
			drawFox(ctx, a, false)
		case domain.AnimalSnowFox:
			drawFox(ctx, a, true)
		}
		ctx.SetAlpha(1)
	}
}

func drawBird(ctx *Ctx, a *domain.Animal) {
	x, y := a.Position.X, a.Position.Y
	// Simple two-arc wing shape
	ctx.SetStroke("#333333")
	ctx.SetLineWidth(1.5)
	ctx.BeginPath()
	ctx.MoveTo(x-6, y)
	ctx.QuadraticCurveTo(x-3, y-4, x, y)
	ctx.QuadraticCurveTo(x+3, y-4, x+6, y)
	ctx.Stroke()
}

func drawSquirrel(ctx *Ctx, a *domain.Animal) {
	x, y := a.Position.X, a.Position.Y
	// Body
	ctx.SetFill("#8B5E3C")
	ctx.BeginPath()
	ctx.Arc(x, y, 5, 0, math.Pi*2)
	ctx.Fill()
	// Tail arc
	ctx.SetStroke("#A0522D")
	ctx.SetLineWidth(3)
	ctx.BeginPath()
	ctx.MoveTo(x+3, y)
	ctx.QuadraticCurveTo(x+10, y-8, x+6, y-14)
	ctx.Stroke()
}

func drawButterfly(ctx *Ctx, a *domain.Animal) {
	x, y := a.Position.X, a.Position.Y
	// Two wing pairs using ellipse approximation with arcs
	// Upper wings
	ctx.SetFill("#E8A020")
	ctx.BeginPath()
	ctx.MoveTo(x, y)
	ctx.BezierCurveTo(x-10, y-12, x-14, y-4, x, y)
	ctx.Fill()
	ctx.BeginPath()
	ctx.MoveTo(x, y)
	ctx.BezierCurveTo(x+10, y-12, x+14, y-4, x, y)
	ctx.Fill()
	// Lower wings (smaller)
	ctx.SetFill("#C07010")
	ctx.BeginPath()
	ctx.MoveTo(x, y)
	ctx.BezierCurveTo(x-8, y+6, x-10, y+2, x, y)
	ctx.Fill()
	ctx.BeginPath()
	ctx.MoveTo(x, y)
	ctx.BezierCurveTo(x+8, y+6, x+10, y+2, x, y)
	ctx.Fill()
	// Body
	ctx.SetFill("#1A0A00")
	ctx.BeginPath()
	ctx.Arc(x, y, 2, 0, math.Pi*2)
	ctx.Fill()
}

func drawFirefly(ctx *Ctx, a *domain.Animal) {
	x, y := a.Position.X, a.Position.Y
	// Glow effect — only effective at night
	ctx.SetShadow(8, "#ADFF2F")
	ctx.SetFill("#ADFF2F")
	ctx.BeginPath()
	ctx.Arc(x, y, 2.5, 0, math.Pi*2)
	ctx.Fill()
	ctx.ClearShadow()
}
