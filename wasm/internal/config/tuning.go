package config

// Weather inverse-momentum tuning.
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

// Animal spawn weights per season index (spring=0, summer=1, autumn=2, winter=3).
var AnimalSpawnWeights = map[string][4]float64{
	"bird":      {0.8, 1.0, 0.6, 0.2},
	"squirrel":  {0.5, 0.7, 0.9, 0.3},
	"butterfly": {0.9, 0.8, 0.3, 0.0},
	"deer":      {0.3, 0.2, 0.5, 0.4},
}

// Shake decay — multiplied per frame tick.
const ShakeDecay = 0.85

// Rain drop terminal velocity (pixels/sec at 60fps equivalent).
const RainDropSpeed = 400.0

// Wind line drift speed (pixels/sec).
const WindLineSpeed = 200.0

// Birds ambient volume per weather intensity (index = intensity 0/1/2).
var BirdsAmbientVol = [3]float64{0.4, 0.2, 0.0}

// Breeze volume for cloudy weather (index = intensity 0/1/2).
var BreezeVol = [3]float64{0.0, 0.3, 0.6}
