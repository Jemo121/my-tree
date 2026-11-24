//go:build js && wasm

package render

import (
	"fmt"
	"syscall/js"
)

// Ctx is a thin wrapper around a CanvasRenderingContext2D JS object.
type Ctx struct {
	v  js.Value
	el js.Value
	W  float64
	H  float64
}

// NewCtx retrieves the canvas element by ID and returns its 2D context.
func NewCtx(canvasID string) *Ctx {
	doc := js.Global().Get("document")
	el := doc.Call("getElementById", canvasID)
	ctx := el.Call("getContext", "2d")
	return &Ctx{v: ctx, el: el}
}

// Resize sets the canvas backing-store resolution to match the CSS display size.
func (c *Ctx) Resize() {
	win := js.Global()
	w := win.Get("innerWidth").Float()
	h := win.Get("innerHeight").Float()
	dpr := win.Get("devicePixelRatio").Float()
	if dpr == 0 {
		dpr = 1
	}
	c.W = w * dpr
	c.H = h * dpr
	c.el.Set("width", c.W)
	c.el.Set("height", c.H)
	c.el.Get("style").Set("width", fmt.Sprintf("%gpx", w))
	c.el.Get("style").Set("height", fmt.Sprintf("%gpx", h))
}

// Save pushes the current canvas state onto the state stack.
func (c *Ctx) Save() { c.v.Call("save") }

// Restore pops the most recently saved canvas state.
func (c *Ctx) Restore() { c.v.Call("restore") }

// ClearRect erases the given rectangle to transparent black.
func (c *Ctx) ClearRect(x, y, w, h float64) { c.v.Call("clearRect", x, y, w, h) }

// FillRect fills the given rectangle with the current fill style.
func (c *Ctx) FillRect(x, y, w, h float64) { c.v.Call("fillRect", x, y, w, h) }

// StrokeRect strokes the outline of the given rectangle.
func (c *Ctx) StrokeRect(x, y, w, h float64) { c.v.Call("strokeRect", x, y, w, h) }

// BeginPath starts a new path, discarding any previous sub-paths.
func (c *Ctx) BeginPath() { c.v.Call("beginPath") }

// ClosePath adds a straight line from the current point to the path start.
func (c *Ctx) ClosePath() { c.v.Call("closePath") }

// Fill fills the current path with the current fill style.
func (c *Ctx) Fill() { c.v.Call("fill") }

// Stroke strokes the current path with the current stroke style.
func (c *Ctx) Stroke() { c.v.Call("stroke") }

// MoveTo moves the current point to (x, y) without drawing.
func (c *Ctx) MoveTo(x, y float64) { c.v.Call("moveTo", x, y) }

// LineTo draws a straight line from the current point to (x, y).
func (c *Ctx) LineTo(x, y float64) { c.v.Call("lineTo", x, y) }

// Arc draws a circular arc centred at (x, y) from startAngle to endAngle.
func (c *Ctx) Arc(x, y, r, sa, ea float64) { c.v.Call("arc", x, y, r, sa, ea) }

// ArcCCW draws a counter-clockwise arc centred at (x, y).
func (c *Ctx) ArcCCW(x, y, r, sa, ea float64) { c.v.Call("arc", x, y, r, sa, ea, true) }

// QuadraticCurveTo draws a quadratic Bézier curve to (x, y) using control point (cpx, cpy).
func (c *Ctx) QuadraticCurveTo(cpx, cpy, x, y float64) {
	c.v.Call("quadraticCurveTo", cpx, cpy, x, y)
}

// BezierCurveTo draws a cubic Bézier curve to (x, y) using two control points.
func (c *Ctx) BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	c.v.Call("bezierCurveTo", cp1x, cp1y, cp2x, cp2y, x, y)
}

// Translate moves the canvas origin by (x, y).
func (c *Ctx) Translate(x, y float64) { c.v.Call("translate", x, y) }

// Rotate rotates the canvas by angle radians around the origin.
func (c *Ctx) Rotate(angle float64) { c.v.Call("rotate", angle) }

// Scale applies a scaling transform of (x, y) to the canvas.
func (c *Ctx) Scale(x, y float64) { c.v.Call("scale", x, y) }

// SetFill sets the fill style to the given CSS color string.
func (c *Ctx) SetFill(color string) { c.v.Set("fillStyle", color) }

// SetStroke sets the stroke style to the given CSS color string.
func (c *Ctx) SetStroke(color string) { c.v.Set("strokeStyle", color) }

// SetLineWidth sets the stroke line width in canvas pixels.
func (c *Ctx) SetLineWidth(w float64) { c.v.Set("lineWidth", w) }

// SetAlpha sets the global opacity (0.0 transparent – 1.0 opaque).
func (c *Ctx) SetAlpha(a float64) { c.v.Set("globalAlpha", a) }

// SetLineCap sets the shape of line endpoints ("butt", "round", or "square").
func (c *Ctx) SetLineCap(style string) { c.v.Set("lineCap", style) }

// SetFont sets the font for text operations (CSS font string).
func (c *Ctx) SetFont(font string) { c.v.Set("font", font) }

// SetTextAlign sets horizontal text alignment ("left", "center", or "right").
func (c *Ctx) SetTextAlign(align string) { c.v.Set("textAlign", align) }

// SetTextBaseline sets the text baseline ("top", "middle", or "bottom").
func (c *Ctx) SetTextBaseline(b string) { c.v.Set("textBaseline", b) }

// SetComposite sets the global composite operation (e.g. "source-over").
func (c *Ctx) SetComposite(op string) { c.v.Set("globalCompositeOperation", op) }

// SetShadow sets the shadow blur radius and colour.
func (c *Ctx) SetShadow(blur float64, color string) {
	c.v.Set("shadowBlur", blur)
	c.v.Set("shadowColor", color)
}

// ClearShadow removes the drop shadow.
func (c *Ctx) ClearShadow() { c.v.Set("shadowBlur", 0); c.v.Set("shadowColor", "transparent") }

// FillText draws text at (x, y) using the current fill style and font.
func (c *Ctx) FillText(text string, x, y float64) { c.v.Call("fillText", text, x, y) }

// LinearGradient creates a linear gradient and returns a js.Value to pass to SetFill.
func (c *Ctx) LinearGradient(x0, y0, x1, y1 float64) js.Value {
	return c.v.Call("createLinearGradient", x0, y0, x1, y1)
}

// AddColorStop adds a stop to a gradient value.
func AddColorStop(grad js.Value, offset float64, color string) {
	grad.Call("addColorStop", offset, color)
}

// SetFillGradient sets a gradient as the fill style.
func (c *Ctx) SetFillGradient(grad js.Value) { c.v.Set("fillStyle", grad) }
