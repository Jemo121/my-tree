package config

import "eguia.me/wasm/internal/domain"

// ColorSet holds all colors needed for one (season, time, weather) combination.
type ColorSet struct {
	SkyTop, SkyBottom string
	Ground, Grass     string
	TrunkDark, TrunkLight string
	LeafPri, LeafSec  string
	FruitCol, FlowerCol, TwigCol string
	RootCol, RockCol  string
	Accent, MoonSun   string
	FogAlpha          float64
}

// shared dusk / night palettes reused across seasons
var (
	sharedDawnPalette = ColorSet{
		SkyTop: "#ff7043", SkyBottom: "#ffcc80", Ground: "#2d1b0e",
		Grass: "#4a9e2f", TrunkDark: "#2e1a04", TrunkLight: "#5a3010",
		LeafPri: "#388e3c", LeafSec: "#66bb6a", FruitCol: "#e65100",
		FlowerCol: "#e91e63", TwigCol: "#795548", RootCol: "#5d4037",
		RockCol: "#78909c", Accent: "#ffcc80", MoonSun: "#ff7043", FogAlpha: 0.1,
	}
	sharedDuskPalette = ColorSet{
		SkyTop: "#7b1fa2", SkyBottom: "#ef6c00", Ground: "#1a0e05",
		Grass: "#2e7d32", TrunkDark: "#1a0a00", TrunkLight: "#4a2800",
		LeafPri: "#1b5e20", LeafSec: "#43a047", FruitCol: "#bf360c",
		FlowerCol: "#880e4f", TwigCol: "#4e342e", RootCol: "#3e2723",
		RockCol: "#546e7a", Accent: "#ef6c00", MoonSun: "#ef6c00", FogAlpha: 0.15,
	}
	sharedNightPalette = ColorSet{
		SkyTop: "#0a0a2e", SkyBottom: "#1a1a4e", Ground: "#0d0805",
		Grass: "#1a3d12", TrunkDark: "#0d0500", TrunkLight: "#2a1200",
		LeafPri: "#0d3d10", LeafSec: "#1b5e20", FruitCol: "#7f3100",
		FlowerCol: "#4a0033", TwigCol: "#2e1f1a", RootCol: "#1a0d08",
		RockCol: "#37474f", Accent: "#7fff7f", MoonSun: "#fffde7", FogAlpha: 0.05,
	}
)

// Palette maps (Season, TimeOfDay) → ColorSet.
// Weather modifiers are applied on top at render time.
var Palette = map[domain.Season]map[domain.TimeOfDay]ColorSet{
	domain.SeasonSpring: {
		domain.TimeDay: {
			SkyTop: "#87ceeb", SkyBottom: "#c5e8f0", Ground: "#3d2b1a",
			Grass: "#5aad3f", TrunkDark: "#3b2008", TrunkLight: "#6b3d1a",
			LeafPri: "#4caf50", LeafSec: "#81c784", FruitCol: "#ff8f00",
			FlowerCol: "#f48fb1", TwigCol: "#8d6e63", RootCol: "#6d4c41",
			RockCol: "#90a4ae", Accent: "#7fff7f", MoonSun: "#fff9c4", FogAlpha: 0,
		},
		domain.TimeDawn:  sharedDawnPalette,
		domain.TimeDusk:  sharedDuskPalette,
		domain.TimeNight: sharedNightPalette,
	},
	domain.SeasonSummer: {
		domain.TimeDay: {
			SkyTop: "#1565c0", SkyBottom: "#90caf9", Ground: "#3d2b1a",
			Grass: "#2e7d32", TrunkDark: "#3b2008", TrunkLight: "#6b3d1a",
			LeafPri: "#1b5e20", LeafSec: "#388e3c", FruitCol: "#e53935",
			FlowerCol: "#ffd600", TwigCol: "#8d6e63", RootCol: "#6d4c41",
			RockCol: "#90a4ae", Accent: "#7fff7f", MoonSun: "#fff59d", FogAlpha: 0,
		},
		domain.TimeDawn:  sharedDawnPalette,
		domain.TimeDusk:  sharedDuskPalette,
		domain.TimeNight: sharedNightPalette,
	},
	domain.SeasonAutumn: {
		domain.TimeDay: {
			SkyTop: "#78909c", SkyBottom: "#b0bec5", Ground: "#2d1a09",
			Grass: "#795548", TrunkDark: "#3b2008", TrunkLight: "#6b3d1a",
			LeafPri: "#e65100", LeafSec: "#ff8f00", FruitCol: "#bf360c",
			FlowerCol: "#ffa000", TwigCol: "#5d4037", RootCol: "#4e342e",
			RockCol: "#78909c", Accent: "#ffcc80", MoonSun: "#ffecb3", FogAlpha: 0.08,
		},
		domain.TimeDawn: {
			SkyTop: "#bf360c", SkyBottom: "#ffcc02", Ground: "#1a0e05",
			Grass: "#4e342e", TrunkDark: "#1a0a00", TrunkLight: "#3e1a00",
			LeafPri: "#bf360c", LeafSec: "#e64a19", FruitCol: "#7f3100",
			FlowerCol: "#e65100", TwigCol: "#3e2723", RootCol: "#2e1b0e",
			RockCol: "#607d8b", Accent: "#ffcc80", MoonSun: "#bf360c", FogAlpha: 0.12,
		},
		domain.TimeDusk:  sharedDuskPalette,
		domain.TimeNight: sharedNightPalette,
	},
	domain.SeasonWinter: {
		domain.TimeDay: {
			SkyTop: "#b0bec5", SkyBottom: "#eceff1", Ground: "#1a1a1a",
			Grass: "#37474f", TrunkDark: "#212121", TrunkLight: "#424242",
			LeafPri: "#37474f", LeafSec: "#546e7a", FruitCol: "#c62828",
			FlowerCol: "#e0e0e0", TwigCol: "#424242", RootCol: "#212121",
			RockCol: "#90a4ae", Accent: "#e0e0e0", MoonSun: "#fffde7", FogAlpha: 0.2,
		},
		domain.TimeDawn:  sharedDawnPalette,
		domain.TimeDusk:  sharedDuskPalette,
		domain.TimeNight: sharedNightPalette,
	},
}

// Get returns the ColorSet for the given season and time of day.
// Falls back to spring/day if the combination is missing.
func Get(s domain.Season, t domain.TimeOfDay) ColorSet {
	if times, ok := Palette[s]; ok {
		if cs, ok := times[t]; ok {
			return cs
		}
	}
	return Palette[domain.SeasonSpring][domain.TimeDay]
}
