package domain

import "math/rand/v2"

// WeatherSnowy and WeatherFoggy extend the WeatherType constants defined in time.go.
// They are appended here to avoid modifying the at-limit time.go file.
const (
	WeatherSnowy WeatherType = iota + 4 // must follow WeatherWindy (3)
	WeatherFoggy
)

// windBaseByType returns the baseline wind contribution (0.0–1.0) for a weather type.
// Used by Tick to blend with cursor velocity for WindStrength.
var windBaseByType = [6]float64{
	WeatherSunny:  0.05,
	WeatherCloudy: 0.10,
	WeatherRainy:  0.20,
	WeatherWindy:  0.70,
	WeatherSnowy:  0.15,
	WeatherFoggy:  0.05,
}

const cursorWindScale = 0.0003 // pixels/s → wind unit conversion

// TickWithCursor advances the weather state and updates WindStrength based on cursor speed.
// cursorSpeed is math.Hypot(dx, dy) of cursor movement in canvas pixels per second.
// Call this instead of Tick when cursor tracking is available.
func (w *WeatherState) TickWithCursor(dt, cursorSpeed float64) {
	w.Tick(dt)

	base := 0.0
	if int(w.Type) < len(windBaseByType) {
		base = windBaseByType[w.Type]
	}
	ws := base + cursorSpeed*cursorWindScale
	if ws > 1.0 {
		ws = 1.0
	}
	w.WindStrength = ws
}

// pickSnowOrFog returns WeatherSnowy or WeatherFoggy when weather conditions warrant.
// Called by transition logic to introduce the extended weather types.
func pickSnowOrFog() WeatherType {
	if rand.Float64() < 0.5 {
		return WeatherSnowy
	}
	return WeatherFoggy
}
