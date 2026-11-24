//go:build js && wasm

package input

import (
	"fmt"
	"syscall/js"

	"eguia.me/wasm/internal/domain"
)

// OverlayManager synchronises DOM <a> elements over the canvas for member links.
type OverlayManager struct {
	doc       js.Value
	container js.Value
	canvas    js.Value
	overlays  map[string]js.Value
	// canvas CSS rect (updated each frame)
	cxOff, cyOff, scale float64
}

// NewOverlayManager creates the overlay container div and injects it after the canvas.
func NewOverlayManager(canvasID string) *OverlayManager {
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", canvasID)

	container := doc.Call("createElement", "div")
	container.Set("id", "member-overlays")
	container.Get("style").Set("position", "absolute")
	container.Get("style").Set("inset", "0")
	container.Get("style").Set("pointerEvents", "none")
	container.Get("style").Set("zIndex", "10")

	// Replace existing #member-overlays if present
	existing := doc.Call("getElementById", "member-overlays")
	if !existing.IsNull() {
		existing.Get("parentNode").Call("replaceChild", container, existing)
	} else {
		canvas.Get("parentNode").Call("insertBefore", container, canvas.Get("nextSibling"))
	}

	return &OverlayManager{
		doc:       doc,
		container: container,
		canvas:    canvas,
		overlays:  make(map[string]js.Value),
	}
}

// Sync reconciles DOM <a> tags with the given overlay specs.
func (m *OverlayManager) Sync(specs []domain.OverlaySpec) {
	rect := m.canvas.Call("getBoundingClientRect")
	cssW := rect.Get("width").Float()
	cw := m.canvas.Get("width").Float()
	if cssW == 0 || cw == 0 {
		return
	}
	m.scale = cw / cssW
	m.cxOff = rect.Get("left").Float()
	m.cyOff = rect.Get("top").Float()

	seen := make(map[string]bool, len(specs))
	for _, spec := range specs {
		seen[spec.ID] = true
		el, exists := m.overlays[spec.ID]
		if !exists {
			el = m.doc.Call("createElement", "a")
			el.Set("id", spec.ID)
			el.Set("className", spec.Class)
			el.Set("href", spec.Href)
			el.Set("aria-label", spec.AriaLabel)
			st := el.Get("style")
			st.Set("position", "absolute")
			st.Set("pointerEvents", "auto")
			st.Set("display", "block")
			st.Set("borderRadius", "50%")
			m.container.Call("appendChild", el)
			m.overlays[spec.ID] = el
		} else {
			el.Set("href", spec.Href)
		}
		m.position(el, spec.Rect)
	}

	for id, el := range m.overlays {
		if !seen[id] {
			m.container.Call("removeChild", el)
			delete(m.overlays, id)
		}
	}
}

func (m *OverlayManager) position(el js.Value, r domain.Rect) {
	left := r.Min.X/m.scale + m.cxOff
	top := r.Min.Y/m.scale + m.cyOff
	w := (r.Max.X - r.Min.X) / m.scale
	h := (r.Max.Y - r.Min.Y) / m.scale
	st := el.Get("style")
	st.Set("left", fmt.Sprintf("%.1fpx", left))
	st.Set("top", fmt.Sprintf("%.1fpx", top))
	st.Set("width", fmt.Sprintf("%.1fpx", w))
	st.Set("height", fmt.Sprintf("%.1fpx", h))
}

// BuildSpecs converts the tree layout into overlay specs.
func BuildSpecs(layout *domain.TreeLayout) []domain.OverlaySpec {
	specs := make([]domain.OverlaySpec, 0, layout.ActiveCount+layout.InactCount+1)
	for _, m := range layout.ActiveMembers {
		specs = append(specs, domain.OverlaySpec{
			ID:        "member-" + m.Subdomain,
			Href:      "https://" + m.Subdomain + ".eguia.me",
			Rect:      m.Bounds,
			Class:     "member-overlay member-active",
			AriaLabel: "Visit " + m.Name + "'s page",
		})
	}
	for _, m := range layout.InactMembers {
		specs = append(specs, domain.OverlaySpec{
			ID:        "member-hidden-" + m.Subdomain,
			Href:      "/roots/",
			Rect:      m.Bounds,
			Class:     "member-overlay member-inactive",
			AriaLabel: m.Name + " — private profile",
		})
	}
	specs = append(specs, domain.OverlaySpec{
		ID:        "sign-trigger",
		Href:      "#sign",
		Rect:      layout.SignBounds,
		Class:     "member-overlay sign-overlay",
		AriaLabel: "Family information",
	})
	return specs
}
