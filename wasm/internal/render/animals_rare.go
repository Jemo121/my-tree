//go:build js && wasm

package render

import (
	"math"

	"eguia.me/wasm/internal/domain"
)

// drawDeer renders a stylised deer: tan body, head, and antler lines.
func drawDeer(ctx *Ctx, a *domain.Animal) {
	x, y := a.Position.X, a.Position.Y

	// Body — elongated rectangle approximation
	ctx.SetFill("#C8A060")
	ctx.BeginPath()
	ctx.MoveTo(x-14, y-6)
	ctx.LineTo(x+14, y-6)
	ctx.LineTo(x+14, y+6)
	ctx.LineTo(x-14, y+6)
	ctx.ClosePath()
	ctx.Fill()

	// Head
	ctx.BeginPath()
	ctx.Arc(x+16, y-8, 5, 0, math.Pi*2)
	ctx.Fill()

	// Legs (four lines below body)
	ctx.SetStroke("#A07840")
	ctx.SetLineWidth(2)
	for _, lx := range [4]float64{x - 10, x - 4, x + 4, x + 10} {
		ctx.BeginPath()
		ctx.MoveTo(lx, y+6)
		ctx.LineTo(lx, y+16)
		ctx.Stroke()
	}

	// Antlers
	ctx.SetStroke("#7A5020")
	ctx.SetLineWidth(1.5)
	// Left antler
	ctx.BeginPath()
	ctx.MoveTo(x+13, y-13)
	ctx.LineTo(x+10, y-22)
	ctx.MoveTo(x+10, y-18)
	ctx.LineTo(x+6, y-22)
	ctx.Stroke()
	// Right antler
	ctx.BeginPath()
	ctx.MoveTo(x+19, y-13)
	ctx.LineTo(x+22, y-22)
	ctx.MoveTo(x+22, y-18)
	ctx.LineTo(x+26, y-22)
	ctx.Stroke()
}

// drawFox renders a fox. snow=true uses white/grey winter coat (SnowFox).
func drawFox(ctx *Ctx, a *domain.Animal, snow bool) {
	x, y := a.Position.X, a.Position.Y

	bodyColor := "#D2531A"
	accentColor := "#FFFFFF"
	if snow {
		bodyColor = "#F0F0F0"
		accentColor = "#AAAAAA"
	}

	// Body
	ctx.SetFill(bodyColor)
	ctx.BeginPath()
	ctx.Arc(x, y, 7, 0, math.Pi*2)
	ctx.Fill()

	// Head
	ctx.BeginPath()
	ctx.Arc(x+9, y-4, 5, 0, math.Pi*2)
	ctx.Fill()

	// Ears (two small triangles)
	ctx.BeginPath()
	ctx.MoveTo(x+6, y-9)
	ctx.LineTo(x+8, y-15)
	ctx.LineTo(x+10, y-9)
	ctx.ClosePath()
	ctx.Fill()
	ctx.BeginPath()
	ctx.MoveTo(x+10, y-9)
	ctx.LineTo(x+12, y-15)
	ctx.LineTo(x+14, y-9)
	ctx.ClosePath()
	ctx.Fill()

	// White chest patch
	ctx.SetFill(accentColor)
	ctx.BeginPath()
	ctx.Arc(x+9, y-2, 2.5, 0, math.Pi*2)
	ctx.Fill()

	// Bushy tail
	ctx.SetStroke(bodyColor)
	ctx.SetLineWidth(4)
	ctx.BeginPath()
	ctx.MoveTo(x-6, y)
	ctx.QuadraticCurveTo(x-16, y-6, x-12, y-14)
	ctx.Stroke()
	// Tail tip
	ctx.SetFill(accentColor)
	ctx.BeginPath()
	ctx.Arc(x-12, y-14, 3, 0, math.Pi*2)
	ctx.Fill()
}
