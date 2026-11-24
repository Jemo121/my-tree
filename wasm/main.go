//go:build js && wasm

// Package main is the entry point for the eguia.me family tree Go WASM module.
package main

import (
	"syscall/js"
	"time"

	"eguia.me/wasm/assets"
	"eguia.me/wasm/internal/config"
	"eguia.me/wasm/internal/domain"
	"eguia.me/wasm/internal/input"
	"eguia.me/wasm/internal/layout"
	"eguia.me/wasm/internal/render"
)

func main() {
	members, err := loadMembers()
	if err != nil {
		js.Global().Get("console").Call("error", "Failed to load members: "+err.Error())
		return
	}

	ctx := render.NewCtx("tree-canvas")
	ctx.Resize()

	treeLayout := layout.Build(members, ctx.W, ctx.H)
	weather := domain.NewWeatherState()
	shakes := &render.ShakeSystem{}
	signSwing := &render.SignSwinger{}
	particles := &render.ParticleSystem{}
	events := input.Register("tree-canvas")
	overlays := input.NewOverlayManager("tree-canvas")

	signCfg := assets.LoadSignConfig()

	render.UpdateCloudPuffs(weather, ctx.W, ctx.H)
	render.UpdateSignPanel(treeLayout, signCfg)
	bindCloseButton()
	overlays.Sync(input.BuildSpecs(treeLayout))
	render.RemoveLoadingOverlay()

	var last float64
	var frameCount int

	var tick js.Func
	tick = js.FuncOf(func(_ js.Value, args []js.Value) any {
		now := args[0].Float()
		dt := 0.016
		if last > 0 {
			dt = (now - last) / 1000.0
			if dt > 0.1 {
				dt = 0.1
			}
		}
		last = now
		frameCount++

		if events.Resized {
			ctx.Resize()
			treeLayout = layout.Build(members, ctx.W, ctx.H)
			render.UpdateCloudPuffs(weather, ctx.W, ctx.H)
		}

		weather.Tick(dt)
		shakes.Tick()
		signSwing.Tick()
		particles.Tick(dt)
		render.SpawnRainDrops(weather, ctx.W, ctx.H)
		if len(weather.CloudPuffs) == 0 {
			render.UpdateCloudPuffs(weather, ctx.W, ctx.H)
		}

		now2 := time.Now()
		season := domain.SeasonFromTime(now2)
		tod := domain.TimeOfDayFromTime(now2)
		cs := config.Get(season, tod)

		render.SetCSSVar("--sky-top", cs.SkyTop)
		render.SetCSSVar("--sky-bottom", cs.SkyBottom)

		// Seasonal particle emission
		crown := treeLayout.CrownBounds
		crownCX := (crown.Min.X + crown.Max.X) / 2
		crownCY := (crown.Min.Y + crown.Max.Y) / 2
		switch season {
		case domain.SeasonAutumn:
			particles.EmitLeaves(crownCX, crownCY, weather.WindStrength)
		case domain.SeasonWinter:
			doubleSnow := weather.Type == domain.WeatherSnowy
			particles.EmitSnow(ctx.W, weather.WindStrength, doubleSnow)
		case domain.SeasonSpring:
			particles.TriggerBudPop(season, nil)
		default:
			particles.TriggerBudPop(season, nil) // arms re-trigger for next spring
		}

		grid := input.Build(treeLayout, weather)
		for _, click := range events.Clicks {
			handleClick(click, grid, treeLayout, weather, shakes, signSwing, signCfg)
		}

		dx, dy := shakes.Offset()
		shakeOff := domain.Vec2{X: dx, Y: dy}

		ctx.ClearRect(0, 0, ctx.W, ctx.H)
		render.DrawBackground(ctx, treeLayout, cs, weather, tod)
		render.DrawTree(ctx, treeLayout, cs, shakeOff)
		particles.Draw(ctx)
		render.DrawMembers(ctx, treeLayout, cs)
		render.DrawWeather(ctx, weather, cs)
		render.DrawSign(ctx, treeLayout)
		render.DrawFlash(ctx, weather)

		if frameCount%3 == 0 {
			overlays.Sync(input.BuildSpecs(treeLayout))
		}

		events.Drain()
		js.Global().Call("requestAnimationFrame", tick)
		return nil
	})

	js.Global().Call("requestAnimationFrame", tick)
	select {}
}

func loadMembers() ([]*domain.Member, error) {
	raw, err := assets.LoadMembers()
	if err != nil {
		return nil, err
	}
	out := make([]*domain.Member, 0, len(raw))
	for i, r := range raw {
		m := &domain.Member{
			Name:       r.Name,
			Subdomain:  r.Subdomain,
			Template:   r.Template,
			Generation: r.Generation,
			Parent:     r.Parent,
			Status:     domain.MemberStatus(r.Status),
			Icon:       domain.MemberIcon(r.Icon),
			NodeID:     i,
		}
		out = append(out, m)
	}
	return out, nil
}

func handleClick(
	click domain.Vec2,
	grid *input.SpatialGrid,
	treeLayout *domain.TreeLayout,
	weather *domain.WeatherState,
	shakes *render.ShakeSystem,
	signSwing *render.SignSwinger,
	signCfg assets.SignConfig,
) {
	target := grid.Query(click)
	action, url := target.Resolve(weather)

	switch action {
	case domain.ActionNavigateMember, domain.ActionNavigateHidden:
		js.Global().Get("location").Set("href", url)
	case domain.ActionShakeTrunk:
		shakes.Add(render.ShakeTrunk, treeLayout.TrunkBase, 8)
	case domain.ActionShakeLeaves:
		shakes.Add(render.ShakeLocal, click, 4)
	case domain.ActionShakeRoots:
		shakes.Add(render.ShakeEarth, click, 5)
	case domain.ActionSplashCloud, domain.ActionVanishCloud:
		render.VanishCloud(weather, target.CloudIdx)
	case domain.ActionShowSign:
		signSwing.Trigger()
		render.UpdateSignPanel(treeLayout, signCfg)
		render.ShowSignPanel()
	}
}

func bindCloseButton() {
	btn := js.Global().Get("document").Call("getElementById", "sign-panel-close")
	if btn.IsNull() {
		return
	}
	btn.Call("addEventListener", "click", js.FuncOf(func(_ js.Value, _ []js.Value) any {
		render.HideSignPanel()
		return nil
	}))
}
