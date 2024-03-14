package particles

import (
	"image/color"
	"math/rand"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

type Particle struct {
	Scale     float32
	Position  resolv.Vector
	SpawnTime float64
	Type      string //circle square stroke_circle stroke_square
}

type ParticleSpawner struct {
	Object    *resolv.Object
	Direction resolv.Vector
	Color     color.RGBA
	Time      float64
	Particles map[uuid.UUID]*Particle
	Cooldown  float64
	SpawnTime float64
	Tick      bool
	Type      string
	Wobble    bool
}

func NewParticleSpawner(obj *resolv.Object, color color.RGBA, p_type string, wobble bool) *ParticleSpawner {

	p_map := make(map[uuid.UUID]*Particle)

	ps := &ParticleSpawner{
		Object:    obj,
		Color:     color,
		Cooldown:  0.02,
		Particles: p_map,
		Tick:      true,
		Type:      p_type,
		Wobble:    wobble,
	}

	ps.Direction = resolv.NewVector(ps.Object.Position.X, ps.Object.Position.Y)

	return ps
}

func (ps *ParticleSpawner) Update(v resolv.Vector) {

	//for _, p := range ps.Particles {

	//	p.Position = p.Position.Add(v)
	//}

	if ps.Tick {

		parent_vec := resolv.NewVector(ps.Object.Position.X, ps.Object.Position.Y)

		if ps.Wobble {

			spread := 0.02 - rand.Float64()*0.04
			parent_vec = resolv.NewVector(ps.Object.Position.X, ps.Object.Position.Y).Rotate(spread)

		}

		pos := parent_vec

		for key, p := range ps.Particles {

			if p.Scale < 1 {
				delete(ps.Particles, key)
			}

			p.Scale -= 0.5
		}

		p := &Particle{
			Scale:     6,
			Position:  pos,
			SpawnTime: ps.Time,
			Type:      ps.Type,
		}

		ps.SpawnTime = ps.Time

		ps.Particles[uuid.New()] = p

	}

	ps.Time += 1.0 / 60.0
	ps.ParticleTick()

}

func (ps *ParticleSpawner) ParticleTick() {

	if ps.Time-ps.SpawnTime >= ps.Cooldown {
		ps.Tick = true
	} else {
		ps.Tick = false
	}

}

func (ps *ParticleSpawner) Draw(screen *ebiten.Image) {

	for _, p := range ps.Particles {

		switch t := p.Type; t {

		case "circle":
			vector.DrawFilledCircle(screen, float32(p.Position.X), float32(p.Position.Y), p.Scale, ps.Color, false)

		case "square":
			vector.DrawFilledRect(screen, float32(p.Position.X), float32(p.Position.Y), p.Scale, p.Scale, ps.Color, false)

		case "stroke_circle":
			vector.StrokeCircle(screen, float32(p.Position.X), float32(p.Position.Y), p.Scale, 2, ps.Color, false)

		case "stroke_square":
			vector.StrokeRect(screen, float32(p.Position.X), float32(p.Position.Y), p.Scale, p.Scale, 2, ps.Color, false)

		default:
			vector.DrawFilledCircle(screen, float32(p.Position.X), float32(p.Position.Y), p.Scale, ps.Color, false)
		}

	}
}
