//go:build js && wasm

package render

import (
	"math"
	"math/rand/v2"

	"eguia.me/wasm/internal/domain"
)

// Particle is a single ephemeral visual element (leaf, snowflake, bud).
type Particle struct {
	Pos   domain.Vec2
	Vel   domain.Vec2
	Life  float64 // seconds alive
	Max   float64 // lifetime in seconds
	Size  float64
	Color string
	// budScale is used by BudPop particles: scale from 0 → 1.2 → 1.0
	budScale float64
	isBud    bool
}

// alive reports whether this particle should continue to be updated.
func (p *Particle) alive() bool { return p.Life < p.Max }

// alpha returns the current opacity (fades out over last 20% of life).
func (p *Particle) alpha() float64 {
	t := p.Life / p.Max
	if t < 0.8 {
		return 1.0
	}
	return 1.0 - (t-0.8)/0.2
}

// ParticleSystem manages leaf, snow, and bud-pop particles.
type ParticleSystem struct {
	leaves     []Particle
	snowflakes []Particle
	buds       []Particle
	// budsFired tracks the last season a bud-pop was fired (avoid re-firing same spring).
	budsFiredSeason domain.Season
	budsArmed       bool // set to true when transitioning away from Spring; rearms trigger
}

const (
	maxLeaves     = 40
	maxSnowflakes = 60
	maxBuds       = 24
)

// Tick advances all active particles.
func (ps *ParticleSystem) Tick(dt float64) {
	tickParticles(ps.leaves, dt)
	tickParticles(ps.snowflakes, dt)
	ps.leaves = cullDead(ps.leaves)
	ps.snowflakes = cullDead(ps.snowflakes)
	tickBuds(ps.buds, dt)
	ps.buds = cullDead(ps.buds)
}

func tickParticles(pp []Particle, dt float64) {
	for i := range pp {
		p := &pp[i]
		p.Pos = p.Pos.Add(p.Vel.Scale(dt))
		p.Life += dt
	}
}

func tickBuds(pp []Particle, dt float64) {
	for i := range pp {
		p := &pp[i]
		p.Life += dt
		t := p.Life / p.Max
		switch {
		case t < 0.4:
			p.budScale = t / 0.4 * 1.2
		case t < 0.6:
			p.budScale = 1.2 - (t-0.4)/0.2*0.2
		default:
			p.budScale = 1.0
		}
	}
}

func cullDead(pp []Particle) []Particle {
	live := pp[:0]
	for _, p := range pp {
		if p.alive() {
			live = append(live, p)
		}
	}
	return live
}

// EmitLeaves spawns 0-1 leaf particles per call during Autumn.
// crownRadius controls the emission spread across the crown area.
func (ps *ParticleSystem) EmitLeaves(crownX, crownY, crownRadius, windStrength float64) {
	if len(ps.leaves) >= maxLeaves {
		return
	}
	leafColors := [4]string{"#D2691E", "#FF8C00", "#DAA520", "#B22222"}
	if rand.IntN(2) == 0 {
		return // 50% chance to skip — keeps density natural at reduced call rate
	}
	// Distribute across the crown silhouette, not just the centre.
	angle := rand.Float64() * math.Pi * 2
	r := rand.Float64() * crownRadius * 0.85
	ps.leaves = append(ps.leaves, Particle{
		Pos:   domain.Vec2{X: crownX + r*math.Cos(angle), Y: crownY + r*math.Sin(angle)*0.6},
		Vel:   domain.Vec2{X: windStrength*20 + (rand.Float64()-0.5)*18, Y: 8 + rand.Float64()*10},
		Max:   4.0 + rand.Float64()*3.0,
		Size:  3 + rand.Float64()*3,
		Color: leafColors[rand.IntN(4)],
	})
}

// EmitSnow spawns snowflake particles. doubleRate doubles emission for WeatherSnowy.
func (ps *ParticleSystem) EmitSnow(cw, windStrength float64, doubleRate bool) {
	limit := maxSnowflakes
	rate := 1
	if doubleRate {
		rate = 2
	}
	for range rate {
		if len(ps.snowflakes) >= limit {
			break
		}
		ps.snowflakes = append(ps.snowflakes, Particle{
			Pos:   domain.Vec2{X: rand.Float64() * cw, Y: -4},
			Vel:   domain.Vec2{X: windStrength*15 + (rand.Float64()-0.5)*8, Y: 18 + rand.Float64()*12},
			Max:   3.0 + rand.Float64()*2.0,
			Size:  2 + rand.Float64()*2,
			Color: "#E8F4FC",
		})
	}
}

// TriggerBudPop fires the one-time spring bud-pop if it has not fired this spring.
// branchTips is a slice of canvas coordinates at branch endpoints.
func (ps *ParticleSystem) TriggerBudPop(season domain.Season, branchTips []domain.Vec2) {
	if season != domain.SeasonSpring {
		ps.budsArmed = true // leaving Spring re-arms for next spring
		return
	}
	if !ps.budsArmed {
		return
	}
	ps.budsArmed = false
	ps.buds = ps.buds[:0]
	for _, tip := range branchTips {
		ps.buds = append(ps.buds, Particle{
			Pos:      tip,
			Max:      0.4,
			Size:     5,
			Color:    "#6DBF67",
			isBud:    true,
			budScale: 0,
		})
		if len(ps.buds) >= maxBuds {
			break
		}
	}
}

// Draw renders all particles onto the canvas.
func (ps *ParticleSystem) Draw(ctx *Ctx) {
	drawLeaves(ctx, ps.leaves)
	drawSnowflakes(ctx, ps.snowflakes)
	drawBuds(ctx, ps.buds)
}

func drawLeaves(ctx *Ctx, pp []Particle) {
	for i := range pp {
		p := &pp[i]
		if !p.alive() {
			continue
		}
		ctx.SetAlpha(p.alpha())
		ctx.SetFill(p.Color)
		ctx.Save()
		ctx.Translate(p.Pos.X, p.Pos.Y)
		ctx.Rotate(p.Life * 2.5)
		ctx.BeginPath()
		ctx.MoveTo(0, -p.Size)
		ctx.BezierCurveTo(p.Size, -p.Size*0.5, p.Size, p.Size*0.5, 0, p.Size)
		ctx.BezierCurveTo(-p.Size, p.Size*0.5, -p.Size, -p.Size*0.5, 0, -p.Size)
		ctx.Fill()
		ctx.Restore()
	}
	ctx.SetAlpha(1)
}

func drawSnowflakes(ctx *Ctx, pp []Particle) {
	ctx.SetFill("#E8F4FC")
	for i := range pp {
		p := &pp[i]
		if !p.alive() {
			continue
		}
		ctx.SetAlpha(p.alpha() * 0.85)
		ctx.BeginPath()
		ctx.Arc(p.Pos.X, p.Pos.Y, p.Size, 0, math.Pi*2)
		ctx.Fill()
	}
	ctx.SetAlpha(1)
}

func drawBuds(ctx *Ctx, pp []Particle) {
	ctx.SetFill("#6DBF67")
	for i := range pp {
		p := &pp[i]
		if p.budScale <= 0 {
			continue
		}
		ctx.SetAlpha(0.9)
		ctx.Save()
		ctx.Translate(p.Pos.X, p.Pos.Y)
		ctx.Scale(p.budScale, p.budScale)
		ctx.BeginPath()
		ctx.Arc(0, 0, p.Size, 0, math.Pi*2)
		ctx.Fill()
		ctx.Restore()
	}
	ctx.SetAlpha(1)
}
