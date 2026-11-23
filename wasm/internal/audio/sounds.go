package audio

// SoundID identifies a pre-generated WAV sound.
type SoundID int

const (
	SoundKnock SoundID = iota // trunk click — knocking wood
	SoundRainLow              // rain ambient, low intensity (looping)
	SoundRainMed              // rain ambient, medium intensity (looping)
	SoundRainHigh             // rain ambient, high intensity (looping)
	SoundThunder              // distant thunder crack + rumble
	SoundWindLow              // wind ambient, low (looping)
	SoundWindMed              // wind ambient, medium (looping)
	SoundWindHigh             // wind ambient, high (looping)
	SoundChirp                // bird chirp
	SoundRustle               // leaf / flower click
	SoundThud                 // soil / root click
	SoundDrop                 // cloud click (single water drop)
	SoundBirdsAmbient         // birdsong loop for sunny low/med
	SoundBreeze               // soft breeze loop for cloudy med/high
	soundCount
)

// String returns the filename (without extension) for each SoundID.
func (s SoundID) String() string {
	return [soundCount]string{
		"knock", "rain_low", "rain_med", "rain_high", "thunder",
		"wind_low", "wind_med", "wind_high", "chirp", "rustle",
		"thud", "drop", "birds_ambient", "breeze",
	}[s]
}
