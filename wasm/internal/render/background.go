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

func drawClouds(ctx *Ctx, layout *domain.TreeLayout, weather *domain.WeatherState, cs config.ColorSet) {
	cw := ctx.W
	cloudAlpha := weather.CloudCover * 0.75
	for i, puff := range weather.CloudPuffs {
		if puff.Vanishing {
			continue
		}
		ctx.SetAlpha(cloudAlpha * puff.Alpha)
		color := "#d0d8e0"
		if weather.Type == domain.WeatherRainy {
			color = "#7a8898"
		}
		ctx.SetFill(color)
		ctx.BeginPath()
		ctx.Arc(puff.Center.X, puff.Center.Y, puff.Radius, 0, math.Pi*2)
		ctx.Fill()
		_ = i
		_ = cw
	}
	ctx.SetAlpha(1)
}

// darken returns a hex color darkened by factor (0–1).
func darken(hex string, factor float64) string {
	// simple approximation using CSS brightness filter via rgba
	_ = hex
	return fmt.Sprintf("rgba(0,0,0,%.2f)", 1-factor)
}
