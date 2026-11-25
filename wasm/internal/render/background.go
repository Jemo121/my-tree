//go:build js && wasm

package render

import (
	"fmt"
	"math"

	"eguia.me/wasm/internal/config"
	"eguia.me/wasm/internal/domain"
)

// DrawBackground renders the sky gradient and ground.
func DrawBackground(ctx *Ctx, layout *domain.TreeLayout, cs config.ColorSet, weather *domain.WeatherState, tod domain.TimeOfDay) {
	cw, ch := ctx.W, ctx.H

	// Sky gradient
	sky := ctx.LinearGradient(0, 0, 0, layout.GroundY)
	AddColorStop(sky, 0, cs.SkyTop)
	AddColorStop(sky, 1, cs.SkyBottom)
	ctx.SetFillGradient(sky)
	ctx.FillRect(0, 0, cw, layout.GroundY)

	// Sun or moon
	drawCelestialBody(ctx, layout, cs, tod)

	// Ground
	ground := ctx.LinearGradient(0, layout.GroundY, 0, ch)
	AddColorStop(ground, 0, cs.Grass)
	AddColorStop(ground, 0.15, cs.Ground)
	AddColorStop(ground, 1, darken(cs.Ground, 0.6))
	ctx.SetFillGradient(ground)
	ctx.FillRect(0, layout.GroundY, cw, ch-layout.GroundY)

	// Fog overlay
	if cs.FogAlpha > 0 {
		fog := ctx.LinearGradient(0, layout.GroundY-ch*0.1, 0, layout.GroundY+ch*0.1)
		AddColorStop(fog, 0, "rgba(255,255,255,0)")
		AddColorStop(fog, 1, fmt.Sprintf("rgba(200,200,200,%.2f)", cs.FogAlpha))
		ctx.SetFillGradient(fog)
		ctx.FillRect(0, layout.GroundY-ch*0.1, cw, ch*0.2)
	}

	// Weather cloud puffs
	if weather.CloudCover > 0 {
		drawClouds(ctx, layout, weather, cs)
	}
}

func drawCelestialBody(ctx *Ctx, layout *domain.TreeLayout, cs config.ColorSet, tod domain.TimeOfDay) {
	cw, ch := ctx.W, ctx.H
	var x, y, r float64
	switch tod {
	case domain.TimeDawn:
		x, y, r = cw*0.15, layout.GroundY*0.85, ch*0.03
	case domain.TimeDay:
		x, y, r = cw*0.8, ch*0.12, ch*0.045
	case domain.TimeDusk:
		x, y, r = cw*0.85, layout.GroundY*0.9, ch*0.035
	case domain.TimeNight:
		x, y, r = cw*0.75, ch*0.1, ch*0.025
	}

	ctx.BeginPath()
	ctx.Arc(x, y, r, 0, math.Pi*2)
	ctx.SetFill(cs.MoonSun)
	ctx.SetShadow(r*2, cs.MoonSun)
	ctx.Fill()
	ctx.ClearShadow()

	// Stars at night
	if tod == domain.TimeNight {
		drawStars(ctx, cw, layout.GroundY)
	}
}

func drawStars(ctx *Ctx, cw, skyH float64) {
	ctx.SetFill("#ffffff")
	// deterministic star positions
	for i := 0; i < 60; i++ {
		x := pseudoRand(i*3) * cw
		y := pseudoRand(i*3+1) * skyH * 0.85
		r := pseudoRand(i*3+2)*1.5 + 0.5
		ctx.SetAlpha(0.4 + pseudoRand(i)*0.6)
		ctx.BeginPath()
		ctx.Arc(x, y, r, 0, math.Pi*2)
		ctx.Fill()
	}
	ctx.SetAlpha(1)
}

func drawClouds(ctx *Ctx, _ *domain.TreeLayout, weather *domain.WeatherState, _ config.ColorSet) {
	color := "#d8e4ee"
	if weather.Type == domain.WeatherRainy {
		color = "#6e7e8a"
	}
	for _, puff := range weather.CloudPuffs {
		if puff.Vanishing {
			continue
		}
		r := puff.Radius
		cx, cy := puff.Center.X, puff.Center.Y

		// Render each puff as 4 overlapping lobes — base + left + right + crown.
		// Together they form the billowing silhouette of a real cloud mass.
		ctx.SetFill(color)
		ctx.SetAlpha(puff.Alpha * 0.82)
		ctx.BeginPath()
		ctx.Arc(cx, cy, r, 0, math.Pi*2)
		ctx.Fill()

		ctx.SetAlpha(puff.Alpha * 0.75)
		ctx.BeginPath()
		ctx.Arc(cx-r*0.55, cy-r*0.2, r*0.72, 0, math.Pi*2)
		ctx.Fill()

		ctx.BeginPath()
		ctx.Arc(cx+r*0.55, cy-r*0.2, r*0.72, 0, math.Pi*2)
		ctx.Fill()

		ctx.SetAlpha(puff.Alpha * 0.70)
		ctx.BeginPath()
		ctx.Arc(cx, cy-r*0.45, r*0.60, 0, math.Pi*2)
		ctx.Fill()
	}
	ctx.SetAlpha(1)
}

// darken returns a hex color darkened by factor (0–1).
func darken(hex string, factor float64) string {
	if len(hex) == 7 && hex[0] == '#' {
		r := hexVal(hex[1:3])
		g := hexVal(hex[3:5])
		b := hexVal(hex[5:7])
		r2 := uint8(float64(r) * (1 - factor))
		g2 := uint8(float64(g) * (1 - factor))
		b2 := uint8(float64(b) * (1 - factor))
		return fmt.Sprintf("#%02x%02x%02x", r2, g2, b2)
	}
	return fmt.Sprintf("rgba(0,0,0,%.2f)", factor)
}

func hexVal(s string) uint8 {
	var v uint8
	for _, c := range s {
		v <<= 4
		switch {
		case c >= '0' && c <= '9':
			v |= uint8(c - '0')
		case c >= 'a' && c <= 'f':
			v |= uint8(c-'a') + 10
		case c >= 'A' && c <= 'F':
			v |= uint8(c-'A') + 10
		}
	}
	return v
}
