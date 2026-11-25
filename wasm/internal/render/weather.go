//go:build js && wasm

package render

import (
	"math"
	"math/rand/v2"

	"eguia.me/wasm/internal/config"
	"eguia.me/wasm/internal/domain"
)

// SpawnRainDrops adds new rain particles relative to canvas dimensions.
// Called each frame when weather is rainy.
func SpawnRainDrops(w *domain.WeatherState, cw, ch float64) {
	spawnRate := 50 + int(w.Intensity)*100
	for i := 0; i < spawnRate/60; i++ {
		if len(w.RainDrops) >= config.MaxRainDrops {
			break
		}
		x := rand.Float64() * cw
		angle := math.Pi*0.1 + w.WindAngle*0.3
		speed := config.RainDropSpeed * (0.8 + rand.Float64()*0.4)
		groundY := ch * config.GroundFrac
		velY := math.Cos(angle) * speed
		w.RainDrops = append(w.RainDrops, domain.RainDrop{
			Pos:     domain.Vec2{X: x, Y: -10},
			Vel:     domain.Vec2{X: math.Sin(angle) * speed * 0.3, Y: velY},
			MaxLife: (groundY + 10) / velY,
		})
	}
}

// DrawWeather renders active weather effects over the scene.
func DrawWeather(ctx *Ctx, weather *domain.WeatherState, cs config.ColorSet) {
	switch weather.Type {
	case domain.WeatherRainy:
		drawRain(ctx, weather, cs)
	case domain.WeatherWindy:
		drawWind(ctx, weather)
	}
}

func drawRain(ctx *Ctx, w *domain.WeatherState, _ config.ColorSet) {
	intensity := 0.4 + float64(w.Intensity)*0.2
	ctx.SetAlpha(intensity)
	ctx.SetStroke("rgba(174, 214, 241, 0.75)")
	dropLen := 8.0 + float64(w.Intensity)*6
	ctx.SetLineWidth(1.0 + float64(w.Intensity)*0.5)

	for _, d := range w.RainDrops {
		ctx.BeginPath()
		ctx.MoveTo(d.Pos.X, d.Pos.Y)
		ctx.LineTo(d.Pos.X+d.Vel.X*0.05, d.Pos.Y-dropLen)
		ctx.Stroke()
	}
	ctx.SetAlpha(1)
}

func drawWind(ctx *Ctx, w *domain.WeatherState) {
	ctx.SetAlpha(0.12 + float64(w.Intensity)*0.08)
	ctx.SetStroke("#ffffff")
	ctx.SetLineWidth(1)

	for i := 0; i < config.MaxWindLines; i++ {
		x := pseudoRand(i*3) * ctx.W
		y := pseudoRand(i*3+1) * ctx.H * 0.9
		l := 20 + pseudoRand(i*3+2)*60
		ctx.BeginPath()
		ctx.MoveTo(x, y)
		ctx.LineTo(x+math.Cos(w.WindAngle)*l, y+math.Sin(w.WindAngle)*l*0.3)
		ctx.Stroke()
	}
	ctx.SetAlpha(1)
}

// UpdateCloudPuffs initialises or repopulates cloud puffs on the weather state.
// Count and opacity scale with weather intensity.
func UpdateCloudPuffs(w *domain.WeatherState, cw, ch float64) {
	live := w.CloudPuffs[:0]
	for _, p := range w.CloudPuffs {
		if !p.Vanishing {
			live = append(live, p)
		}
	}
	w.CloudPuffs = live
	if len(w.CloudPuffs) > 0 {
		return
	}

	// Scale puff count with intensity: Low=2, Medium=4, Heavy=6.
	count := 2 + int(w.Intensity)*2
	if w.Type == domain.WeatherRainy {
		count += 2 // rainy always adds more cloud cover
	}
	baseAlpha := 0.55 + float64(w.Intensity)*0.15 // Low=0.55, Med=0.70, High=0.85

	// Drift speed: 8–22 px/s, direction from wind angle.
	driftSpeed := 8.0 + float64(w.Intensity)*7.0
	for i := 0; i < count; i++ {
		x := cw * (0.08 + rand.Float64()*0.84)
		y := ch * (0.04 + rand.Float64()*0.18)
		r := cw*0.045 + rand.Float64()*cw*0.055
		speed := driftSpeed * (0.7 + rand.Float64()*0.6)
		puff := domain.CloudPuff{
			Center: domain.Vec2{X: x, Y: y},
			Vel:    domain.Vec2{X: math.Cos(w.WindAngle) * speed, Y: math.Sin(w.WindAngle) * speed * 0.1},
			Radius: r,
			Alpha:  baseAlpha + rand.Float64()*0.1,
		}
		puff.Bounds = domain.RectFromCenter(puff.Center, r*1.5, r*1.2)
		w.CloudPuffs = append(w.CloudPuffs, puff)
	}
}

// VanishCloud marks a cloud for removal or splashes it.
func VanishCloud(w *domain.WeatherState, idx int) {
	if idx < 0 || idx >= len(w.CloudPuffs) {
		return
	}
	w.CloudPuffs[idx].Clicks++
	if w.CloudPuffs[idx].Clicks >= 3 {
		w.CloudPuffs[idx].Vanishing = true
	}
}
