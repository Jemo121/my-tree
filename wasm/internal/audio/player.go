//go:build js && wasm

package audio

// Player is a stub for Phase 1. The Go Beep WAV generation and Web Audio
// playback are implemented in Phase 2 (PR-5). All calls are no-ops until then.
type Player struct{}

// NewPlayer creates a no-op audio player.
// Replace with the full implementation in PR-5 (audio generation + embedding).
func NewPlayer() *Player { return &Player{} }

// Play is a no-op stub.
func (p *Player) Play(_ SoundID, _, _ float64) {}

// Resume must be called on the first user interaction (browser autoplay policy).
func (p *Player) Resume() {}
