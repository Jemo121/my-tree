package domain

// HitTargetType classifies what was hit by a click/tap.
type HitTargetType int

const (
	HitNone HitTargetType = iota
	HitActiveMember
	HitInactiveMember
	HitLeaf
	HitTrunk
	HitRootSoil
	HitCloud
	HitSign
)

// HitTarget is the result of a hit-test at a canvas point.
type HitTarget struct {
	Type     HitTargetType
	Member   *Member
	CloudIdx int
	Position Vec2
	Bounds   Rect
}

// ClickAction is what should happen when a hit target is activated.
type ClickAction int

const (
	ActionNone ClickAction = iota
	ActionNavigateMember
	ActionNavigateHidden
	ActionShakeLeaves
	ActionShakeTrunk
	ActionShakeRoots
	ActionSplashCloud
	ActionVanishCloud
	ActionShowSign
)

// Resolve maps a hit target to a click action and optional URL.
func (h *HitTarget) Resolve(weather *WeatherState) (ClickAction, string) {
	switch h.Type {
	case HitActiveMember:
		return ActionNavigateMember, "https://" + h.Member.Subdomain + ".eguia.me"
	case HitInactiveMember:
		return ActionNavigateHidden, "/roots/"
	case HitLeaf:
		return ActionShakeLeaves, ""
	case HitTrunk:
		return ActionShakeTrunk, ""
	case HitRootSoil:
		return ActionShakeRoots, ""
	case HitCloud:
		if weather.Type == WeatherRainy &&
			(weather.Intensity == IntensityMedium || weather.Intensity == IntensityHigh) {
			return ActionSplashCloud, ""
		}
		return ActionVanishCloud, ""
	case HitSign:
		return ActionShowSign, ""
	default:
		return ActionNone, ""
	}
}

// OverlaySpec describes a DOM <a> element to be positioned over the canvas.
type OverlaySpec struct {
	ID        string
	Href      string
	Rect      Rect
	Class     string
	AriaLabel string
}
