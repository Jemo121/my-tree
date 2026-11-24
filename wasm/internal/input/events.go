//go:build js && wasm

package input

import (
	"syscall/js"

	"eguia.me/wasm/internal/domain"
)

// EventState accumulates input events between frames.
type EventState struct {
	Clicks    []domain.Vec2
	PanDelta  domain.Vec2
	ZoomDelta float64
	Resized   bool

	// pan/zoom internals
	lastTouch  domain.Vec2
	pinchDist0 float64
	dragging   bool
}

// ViewTransform tracks pan and zoom applied to the canvas draw context.
type ViewTransform struct {
	PanX, PanY float64
	Zoom       float64
}

// NewViewTransform returns the identity transform.
func NewViewTransform() *ViewTransform { return &ViewTransform{Zoom: 1} }

// Apply pan/zoom delta from accumulated events.
func (vt *ViewTransform) Apply(es *EventState, cw, ch float64) {
	vt.PanX += es.PanDelta.X
	vt.PanY += es.PanDelta.Y
	if es.ZoomDelta != 0 {
		factor := 1 + es.ZoomDelta*0.001
		vt.Zoom *= factor
		if vt.Zoom < 0.4 {
			vt.Zoom = 0.4
		}
		if vt.Zoom > 4 {
			vt.Zoom = 4
		}
	}
	// clamp pan so tree stays on screen
	maxPanX := cw * (vt.Zoom - 1) / 2
	if vt.PanX < -maxPanX {
		vt.PanX = -maxPanX
	}
	if vt.PanX > maxPanX {
		vt.PanX = maxPanX
	}
	maxPanY := ch * (vt.Zoom - 1) / 2
	if vt.PanY < -maxPanY {
		vt.PanY = -maxPanY
	}
	if vt.PanY > maxPanY {
		vt.PanY = maxPanY
	}
}

// Register attaches all canvas event listeners and returns the EventState.
func Register(canvasID string) *EventState {
	es := &EventState{}
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", canvasID)
	win := js.Global()

	// Click
	canvas.Call("addEventListener", "click", js.FuncOf(func(_ js.Value, args []js.Value) any {
		e := args[0]
		rect := canvas.Call("getBoundingClientRect")
		dpr := win.Get("devicePixelRatio").Float()
		if dpr == 0 {
			dpr = 1
		}
		x := (e.Get("clientX").Float() - rect.Get("left").Float()) * dpr
		y := (e.Get("clientY").Float() - rect.Get("top").Float()) * dpr
		es.Clicks = append(es.Clicks, domain.Vec2{X: x, Y: y})
		return nil
	}))

	// Wheel zoom
	canvas.Call("addEventListener", "wheel", js.FuncOf(func(_ js.Value, args []js.Value) any {
		e := args[0]
		e.Call("preventDefault")
		es.ZoomDelta -= e.Get("deltaY").Float()
		return nil
	}), map[string]any{"passive": false})

	// Mouse drag pan
	canvas.Call("addEventListener", "mousedown", js.FuncOf(func(_ js.Value, args []js.Value) any {
		e := args[0]
		es.lastTouch = domain.Vec2{X: e.Get("clientX").Float(), Y: e.Get("clientY").Float()}
		es.dragging = true
		return nil
	}))
	doc.Call("addEventListener", "mousemove", js.FuncOf(func(_ js.Value, args []js.Value) any {
		if !es.dragging {
			return nil
		}
		e := args[0]
		cx := e.Get("clientX").Float()
		cy := e.Get("clientY").Float()
		es.PanDelta.X += cx - es.lastTouch.X
		es.PanDelta.Y += cy - es.lastTouch.Y
		es.lastTouch = domain.Vec2{X: cx, Y: cy}
		return nil
	}))
	doc.Call("addEventListener", "mouseup", js.FuncOf(func(_ js.Value, _ []js.Value) any {
		es.dragging = false
		return nil
	}))

	// Window resize
	win.Call("addEventListener", "resize", js.FuncOf(func(_ js.Value, _ []js.Value) any {
		es.Resized = true
		return nil
	}))

	return es
}

// Drain resets per-frame consumed state.
func (es *EventState) Drain() {
	es.Clicks = es.Clicks[:0]
	es.PanDelta = domain.Vec2{}
	es.ZoomDelta = 0
	es.Resized = false
}
