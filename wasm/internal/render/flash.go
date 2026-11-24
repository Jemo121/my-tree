//go:build js && wasm

package render

import "eguia.me/wasm/internal/domain"

// DrawFlash renders the lightning flash overlay if active.
func DrawFlash(ctx *Ctx, weather *domain.WeatherState) {
	a := weather.FlashAlpha()
	if a <= 0 {
		return
	}
	ctx.SetAlpha(a)
	ctx.SetFill("rgba(200,220,255,1)")
	ctx.FillRect(0, 0, ctx.W, ctx.H)
	ctx.SetAlpha(1)
}
