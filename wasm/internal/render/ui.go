//go:build js && wasm

package render

import (
	"syscall/js"

	"eguia.me/wasm/assets"
	"eguia.me/wasm/internal/domain"
)

// SignSwinger tracks a simple rotation oscillation on sign click (~600 ms decay).
type SignSwinger struct {
	Angle float64 // current rotation in radians
	vel   float64 // angular velocity
}

const signSwingDecay = 0.88

// Trigger starts a swing with an initial kick.
func (s *SignSwinger) Trigger() { s.vel = 0.18 }

// Tick advances the swing physics (call each frame).
func (s *SignSwinger) Tick() {
	s.vel *= signSwingDecay
	s.Angle += s.vel
	s.Angle *= signSwingDecay
}

// DrawSign renders the wooden sign element on the canvas.
// swing is the current rotation angle in radians (from SignSwinger); 0 means upright.
func DrawSign(ctx *Ctx, layout *domain.TreeLayout, swing float64) {
	sb := layout.SignBounds
	w := sb.Max.X - sb.Min.X
	h := sb.Max.Y - sb.Min.Y
	cx := (sb.Min.X + sb.Max.X) / 2
	postH := layout.GroundY - sb.Max.Y
	postW := 7.0 // px per post

	// Two posts at ~25% inset from each side of the board.
	leftPost := sb.Min.X + w*0.22
	rightPost := sb.Max.X - w*0.22
	ctx.SetFill("#5c3a1e")
	ctx.FillRect(leftPost, sb.Max.Y, postW, postH)
	ctx.FillRect(rightPost, sb.Max.Y, postW, postH)

	// Board — rotates around the post tops on swing.
	ctx.Save()
	ctx.Translate(cx, sb.Max.Y)
	ctx.Rotate(swing)
	ctx.Translate(-cx, -sb.Max.Y)

	ctx.BeginPath()
	ctx.MoveTo(sb.Min.X+4, sb.Min.Y)
	ctx.LineTo(sb.Max.X-4, sb.Min.Y)
	ctx.LineTo(sb.Max.X, sb.Min.Y+4)
	ctx.LineTo(sb.Max.X, sb.Max.Y-4)
	ctx.LineTo(sb.Max.X-4, sb.Max.Y)
	ctx.LineTo(sb.Min.X+4, sb.Max.Y)
	ctx.LineTo(sb.Min.X, sb.Max.Y-4)
	ctx.LineTo(sb.Min.X, sb.Min.Y+4)
	ctx.ClosePath()
	ctx.SetFill("#7a4f2a")
	ctx.Fill()
	ctx.SetStroke("#3e2208")
	ctx.SetLineWidth(1.5)
	ctx.Stroke()

	// Text — small "?" to invite clicking.
	ctx.SetFill("#f5e6c8")
	ctx.SetFont("bold " + fontSize(h*0.5) + "px Georgia, serif")
	ctx.SetTextAlign("center")
	ctx.SetTextBaseline("middle")
	ctx.FillText("?", cx, (sb.Min.Y+sb.Max.Y)/2)
	ctx.Restore()
	_ = w
}

// UpdateSignPanel pushes the current member counts and i18n labels into the DOM panel.
// cfg holds localised strings sourced from Hugo data attributes.
func UpdateSignPanel(layout *domain.TreeLayout, cfg assets.SignConfig) {
	doc := js.Global().Get("document")
	el := doc.Call("getElementById", "sign-active")
	if !el.IsNull() {
		el.Set("textContent", js.ValueOf(layout.ActiveCount))
	}
	el = doc.Call("getElementById", "sign-roots")
	if !el.IsNull() {
		el.Set("textContent", js.ValueOf(layout.InactCount))
	}
	// Push localised labels into the panel so the DOM template can reference them.
	el = doc.Call("getElementById", "sign-label-active")
	if !el.IsNull() {
		el.Set("textContent", cfg.LabelActive)
	}
	el = doc.Call("getElementById", "sign-label-roots")
	if !el.IsNull() {
		el.Set("textContent", cfg.LabelRoots)
	}
}

// ShowSignPanel removes the hidden attribute from the sign panel element.
func ShowSignPanel() {
	doc := js.Global().Get("document")
	panel := doc.Call("getElementById", "sign-panel")
	if !panel.IsNull() {
		panel.Call("removeAttribute", "hidden")
	}
}

// HideSignPanel adds the hidden attribute to the sign panel element.
func HideSignPanel() {
	doc := js.Global().Get("document")
	panel := doc.Call("getElementById", "sign-panel")
	if !panel.IsNull() {
		panel.Call("setAttribute", "hidden", "")
	}
}

// SetCSSVar sets a CSS custom property on :root.
func SetCSSVar(name, value string) {
	js.Global().Get("document").Get("documentElement").
		Get("style").Call("setProperty", name, value)
}

// RemoveLoadingOverlay removes the #tree-loading DOM element.
func RemoveLoadingOverlay() {
	doc := js.Global().Get("document")
	el := doc.Call("getElementById", "tree-loading")
	if !el.IsNull() {
		el.Call("remove")
	}
}

func fontSize(h float64) string {
	if h < 10 {
		return "10"
	}
	if h > 32 {
		h = 32
	}
	return js.ValueOf(int(h)).String()
}
