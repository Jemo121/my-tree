//go:build js && wasm

package render

import (
	"syscall/js"

	"eguia.me/wasm/internal/domain"
)

// DrawSign renders the wooden sign element on the canvas.
func DrawSign(ctx *Ctx, layout *domain.TreeLayout) {
	sb := layout.SignBounds
	w := sb.Max.X - sb.Min.X
	h := sb.Max.Y - sb.Min.Y
	cx := (sb.Min.X + sb.Max.X) / 2

	// Post
	postX := cx - 2
	ctx.SetFill("#5c3a1e")
	ctx.FillRect(postX, sb.Max.Y, 4, layout.GroundY-sb.Max.Y)

	// Board
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

	// Text — small "?" to invite clicking
	ctx.SetFill("#f5e6c8")
	ctx.SetFont("bold " + fontSize(h*0.5) + "px Georgia, serif")
	ctx.SetTextAlign("center")
	ctx.SetTextBaseline("middle")
	ctx.FillText("?", cx, (sb.Min.Y+sb.Max.Y)/2)
	_ = w
}

// UpdateSignPanel pushes the current member counts into the DOM panel.
func UpdateSignPanel(layout *domain.TreeLayout) {
	doc := js.Global().Get("document")
	el := doc.Call("getElementById", "sign-active")
	if !el.IsNull() {
		el.Set("textContent", js.ValueOf(layout.ActiveCount))
	}
	el = doc.Call("getElementById", "sign-roots")
	if !el.IsNull() {
		el.Set("textContent", js.ValueOf(layout.InactCount))
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
