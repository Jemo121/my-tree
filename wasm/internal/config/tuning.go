package config

// Tuning constants for weather inverse-momentum.
const (
	WeatherBaseSunny  = 600.0 // seconds
	WeatherBaseCloudy = 450.0
	WeatherBaseRainy  = 300.0
	WeatherBaseWindy  = 400.0
	WeatherMean       = 450.0
	WeatherMomentumK  = 0.3
	WeatherMinDur     = 300.0
	WeatherMaxDur     = 1800.0
)

// AnimalSpawnWeights maps season index (spring=0, summer=1, autumn=2, winter=3) to per-species spawn weights.
var AnimalSpawnWeights = map[string][4]float64{
	"bird":      {0.8, 1.0, 0.6, 0.2},
	"squirrel":  {0.5, 0.7, 0.9, 0.3},
	"butterfly": {0.9, 0.8, 0.3, 0.0},
	"deer":      {0.3, 0.2, 0.5, 0.4},
}

// ShakeDecay is the per-frame decay multiplier for shake effects.
const ShakeDecay = 0.85

// RainDropSpeed is the terminal velocity of rain drops in pixels per second.
const RainDropSpeed = 400.0

// WindLineSpeed is the drift speed of wind lines in pixels per second.
const WindLineSpeed = 200.0

// BirdsAmbientVol is the bird ambient volume per weather intensity (index = intensity 0/1/2).
var BirdsAmbientVol = [3]float64{0.4, 0.2, 0.0}

// BreezeVol is the breeze volume for cloudy weather (index = intensity 0/1/2).
var BreezeVol = [3]float64{0.0, 0.3, 0.6}
