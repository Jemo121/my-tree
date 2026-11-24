package domain

import (
	"math/rand/v2"
	"time"
)

// Season represents a calendar season.
type Season int

// Season values.
const (
	SeasonSpring Season = iota
	SeasonSummer
	SeasonAutumn
	SeasonWinter
)

// SeasonFromTime derives the season from a time.Time value.
func SeasonFromTime(t time.Time) Season {
	switch t.Month() {
	case 3, 4, 5:
		return SeasonSpring
	case 6, 7, 8:
		return SeasonSummer
	case 9, 10, 11:
		return SeasonAutumn
	default:
		return SeasonWinter
	}
}

// TimeOfDay represents a part of the day.
type TimeOfDay int

// TimeOfDay values.
const (
	TimeDawn TimeOfDay = iota
	TimeDay
	TimeDusk
	TimeNight
)

// TimeOfDayFromTime derives the time of day from a time.Time value.
func TimeOfDayFromTime(t time.Time) TimeOfDay {
	h := t.Hour()
	switch {
	case h >= 5 && h < 8:
		return TimeDawn
	case h >= 8 && h < 18:
		return TimeDay
	case h >= 18 && h < 21:
		return TimeDusk
	default:
		return TimeNight
	}
}

// WeatherType is a weather state.
type WeatherType int

// WeatherType values.
const (
	WeatherSunny WeatherType = iota
	WeatherCloudy
	WeatherRainy
	WeatherWindy
)

// WeatherIntensity is the strength of a weather state.
type WeatherIntensity int

// WeatherIntensity values.
const (
	IntensityLow WeatherIntensity = iota
	IntensityMedium
	IntensityHigh
)

// RainDrop is a single falling rain particle.
type RainDrop struct {
	Pos     Vec2
	Vel     Vec2
	Life    float64
	MaxLife float64
}

// CloudPuff is a visual cloud element.
type CloudPuff struct {
	Center Vec2
	Radius float64
	Alpha  float64
	Bounds Rect
	// clicks tracks how many times this cloud has been clicked.
	Clicks    int
	Vanishing bool
}

// WeatherState is the full WASM-internal weather state machine.
type WeatherState struct {
	Type      WeatherType
	Intensity WeatherIntensity
	Elapsed   float64 // seconds in current state
	Duration  float64 // total seconds for this state

	WindAngle    float64
	WindSpeed    float64
	WindStrength float64 // 0.0–1.0, updated each tick by TickWithCursor

	RainDrops  []RainDrop
	CloudPuffs []CloudPuff
	CloudCover float64

	ThunderTimer float64 // seconds until next thunder strike
	FlashElapsed float64 // seconds since last flash start (-1 = inactive)
}

type weatherTransition struct {
	Next   WeatherType
	Weight float64
}

var weatherTransitions = map[WeatherType][]weatherTransition{
	WeatherSunny:  {{WeatherCloudy, 0.6}, {WeatherWindy, 0.3}, {WeatherRainy, 0.1}},
	WeatherCloudy: {{WeatherRainy, 0.5}, {WeatherSunny, 0.3}, {WeatherWindy, 0.2}},
	WeatherRainy:  {{WeatherCloudy, 0.7}, {WeatherWindy, 0.2}, {WeatherSunny, 0.1}},
	WeatherWindy:  {{WeatherSunny, 0.4}, {WeatherCloudy, 0.4}, {WeatherRainy, 0.2}},
}

var weatherBaseDurations = map[WeatherType]float64{
	WeatherSunny:  600,
	WeatherCloudy: 450,
	WeatherRainy:  300,
	WeatherWindy:  400,
}

const (
	weatherMean   = 450.0
	weatherK      = 0.3
	weatherMinDur = 300.0
	weatherMaxDur = 1800.0
	thunderMinSec = 30.0
	thunderMaxSec = 120.0
)

// NewWeatherState creates a sunny low-intensity initial state.
func NewWeatherState() *WeatherState {
	w := &WeatherState{
		Type:         WeatherSunny,
		Intensity:    IntensityLow,
		Duration:     weatherBaseDurations[WeatherSunny],
		FlashElapsed: -1,
	}
	w.resetThunderTimer()
	return w
}

// Tick advances the weather state by deltaTime seconds.
func (w *WeatherState) Tick(dt float64) {
	w.Elapsed += dt

	if w.Type == WeatherRainy {
		w.tickRain(dt)
	}
	if w.Type == WeatherWindy {
		w.WindAngle += (rand.Float64()*0.04 - 0.02)
		w.WindSpeed = w.baseWindSpeed() * (0.8 + rand.Float64()*0.4)
	}
	if w.Type == WeatherRainy && w.Intensity == IntensityHigh {
		w.ThunderTimer -= dt
		if w.ThunderTimer <= 0 {
			w.FlashElapsed = 0
			w.resetThunderTimer()
		}
	}
	if w.FlashElapsed >= 0 {
		w.FlashElapsed += dt
		if w.FlashElapsed > 0.3 {
			w.FlashElapsed = -1
		}
	}

	if w.Elapsed >= w.Duration {
		w.transition()
	}
}

func (w *WeatherState) transition() {
	next := w.pickNextType()
	base := weatherBaseDurations[next]
	momentum := (w.Duration - weatherMean) * weatherK
	dur := base - momentum
	if dur < weatherMinDur {
		dur = weatherMinDur
	}
	if dur > weatherMaxDur {
		dur = weatherMaxDur
	}
	w.Type = next
	w.Intensity = w.pickNextIntensity()
	w.Duration = dur
	w.Elapsed = 0
	w.initStateParams()
}

func (w *WeatherState) pickNextType() WeatherType {
	opts := weatherTransitions[w.Type]
	total := 0.0
	for _, o := range opts {
		total += o.Weight
	}
	r := rand.Float64() * total
	acc := 0.0
	for _, o := range opts {
		acc += o.Weight
		if r <= acc {
			return o.Next
		}
	}
	return opts[0].Next
}

func (w *WeatherState) pickNextIntensity() WeatherIntensity {
	r := rand.Float64()
	switch {
	case r < 0.70:
		return w.Intensity
	case r < 0.85 && w.Intensity < IntensityHigh:
		return w.Intensity + 1
	case r >= 0.85 && w.Intensity > IntensityLow:
		return w.Intensity - 1
	default:
		return w.Intensity
	}
}

func (w *WeatherState) initStateParams() {
	switch w.Type {
	case WeatherWindy:
		w.WindAngle = rand.Float64() * 6.283
		w.WindSpeed = w.baseWindSpeed()
	case WeatherRainy:
		w.CloudCover = 0.8 + rand.Float64()*0.2
	case WeatherSunny:
		w.CloudCover = rand.Float64() * 0.1
		w.RainDrops = w.RainDrops[:0]
	case WeatherCloudy:
		w.CloudCover = 0.4 + float64(w.Intensity)*0.2
	}
	if w.Type == WeatherRainy && w.Intensity == IntensityHigh {
		w.resetThunderTimer()
	}
}

func (w *WeatherState) baseWindSpeed() float64 {
	return 80.0 + float64(w.Intensity)*60.0
}

func (w *WeatherState) resetThunderTimer() {
	w.ThunderTimer = thunderMinSec + rand.Float64()*(thunderMaxSec-thunderMinSec)
}

func (w *WeatherState) tickRain(dt float64) {
	spawnRate := 50.0 + float64(w.Intensity)*100.0
	spawnCount := int(spawnRate * dt)
	// rain drops are added by the render package which knows canvas dimensions
	_ = spawnCount
	for i := range w.RainDrops {
		d := &w.RainDrops[i]
		d.Pos = d.Pos.Add(d.Vel.Scale(dt))
		d.Life += dt
	}
	// remove dead drops
	live := w.RainDrops[:0]
	for _, d := range w.RainDrops {
		if d.Life < d.MaxLife {
			live = append(live, d)
		}
	}
	w.RainDrops = live
}

// IsFlashing reports whether a lightning flash is active.
func (w *WeatherState) IsFlashing() bool { return w.FlashElapsed >= 0 }

// FlashAlpha returns the current flash opacity (0–1).
func (w *WeatherState) FlashAlpha() float64 {
	if w.FlashElapsed < 0 {
		return 0
	}
	// double-flash pattern: two 80ms pulses 100ms apart
	t := w.FlashElapsed
	switch {
	case t < 0.08:
		return 0.4
	case t < 0.18:
		return 0
	case t < 0.26:
		return 0.3
	default:
		return 0
	}
}
