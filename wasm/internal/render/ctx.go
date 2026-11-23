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

func (c *Ctx) Save()    { c.v.Call("save") }
func (c *Ctx) Restore() { c.v.Call("restore") }

func (c *Ctx) ClearRect(x, y, w, h float64) { c.v.Call("clearRect", x, y, w, h) }
func (c *Ctx) FillRect(x, y, w, h float64)  { c.v.Call("fillRect", x, y, w, h) }
func (c *Ctx) StrokeRect(x, y, w, h float64) { c.v.Call("strokeRect", x, y, w, h) }

func (c *Ctx) BeginPath() { c.v.Call("beginPath") }
func (c *Ctx) ClosePath() { c.v.Call("closePath") }
func (c *Ctx) Fill()      { c.v.Call("fill") }
func (c *Ctx) Stroke()    { c.v.Call("stroke") }

func (c *Ctx) MoveTo(x, y float64)        { c.v.Call("moveTo", x, y) }
func (c *Ctx) LineTo(x, y float64)        { c.v.Call("lineTo", x, y) }
func (c *Ctx) Arc(x, y, r, sa, ea float64) { c.v.Call("arc", x, y, r, sa, ea) }
func (c *Ctx) ArcCCW(x, y, r, sa, ea float64) { c.v.Call("arc", x, y, r, sa, ea, true) }

func (c *Ctx) QuadraticCurveTo(cpx, cpy, x, y float64) {
	c.v.Call("quadraticCurveTo", cpx, cpy, x, y)
}

func (c *Ctx) BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	c.v.Call("bezierCurveTo", cp1x, cp1y, cp2x, cp2y, x, y)
}

func (c *Ctx) Translate(x, y float64) { c.v.Call("translate", x, y) }
func (c *Ctx) Rotate(angle float64)   { c.v.Call("rotate", angle) }
func (c *Ctx) Scale(x, y float64)     { c.v.Call("scale", x, y) }

func (c *Ctx) SetFill(color string)       { c.v.Set("fillStyle", color) }
func (c *Ctx) SetStroke(color string)     { c.v.Set("strokeStyle", color) }
func (c *Ctx) SetLineWidth(w float64)     { c.v.Set("lineWidth", w) }
func (c *Ctx) SetAlpha(a float64)         { c.v.Set("globalAlpha", a) }
func (c *Ctx) SetLineCap(cap string)      { c.v.Set("lineCap", cap) }
func (c *Ctx) SetFont(font string)        { c.v.Set("font", font) }
func (c *Ctx) SetTextAlign(align string)  { c.v.Set("textAlign", align) }
func (c *Ctx) SetTextBaseline(b string)   { c.v.Set("textBaseline", b) }
func (c *Ctx) SetComposite(op string)     { c.v.Set("globalCompositeOperation", op) }
func (c *Ctx) SetShadow(blur float64, color string) {
	c.v.Set("shadowBlur", blur)
	c.v.Set("shadowColor", color)
}
func (c *Ctx) ClearShadow() { c.v.Set("shadowBlur", 0); c.v.Set("shadowColor", "transparent") }

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
