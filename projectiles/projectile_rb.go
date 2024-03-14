package projectiles

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"

	"github.com/AndriiPets/ArcaneAvenger/particles"
)

type ProjectileRB struct {
	Space     *resolv.Space
	Type      string
	Object    *resolv.Object
	Size      float64
	Speed     resolv.Vector
	Alive     bool
	Color     color.RGBA
	Particle  *particles.ParticleSpawner
	BounceNum int
	//Boom
	Time            float64
	BoomTime        float64
	BoomDuration    float64
	BoomObj         *resolv.Object
	BoomActiv       bool
	ExplosionRadius float64
}

func SpawnProjectileRB(space *resolv.Space, position resolv.Vector, direction resolv.Vector, color color.RGBA, size float64) *ProjectileRB {

	p := &ProjectileRB{
		Space:           space,
		Object:          resolv.NewObject(position.X, position.Y, size, size),
		Size:            size,
		Speed:           direction,
		Alive:           true,
		Color:           color,
		BoomDuration:    0.2,
		BoomActiv:       false,
		ExplosionRadius: 50,
		BounceNum:       4,
	}

	p.Particle = particles.NewParticleSpawner(p.Object, p.Color, "square", true)

	p.Object.AddTags("projectile")
	space.Add(p.Object)

	return p
}

func (p *ProjectileRB) Update() {
	//p.Speed.Y += 0.1

	px := p.Speed.X * 6
	py := p.Speed.Y * 6

	if p.BoomActiv {
		px = 0
		py = 0
	}
	p.Particle.Update(resolv.NewVector(px, py))

	if check := p.Object.Check(px, 0, "solid"); check != nil {

		px = check.ContactWithCell(check.Cells[0]).X
		p.check_bounce(false)
	}

	p.Object.Position.X += px

	if check := p.Object.Check(0, py, "solid"); check != nil {

		py = check.ContactWithCell(check.Cells[0]).Y
		p.check_bounce(true)
	}

	p.Object.Position.Y += py

	p.Object.Update()

	if p.BoomActiv {

		if p.Time-p.BoomTime >= p.BoomDuration {
			p.Space.Remove(p.BoomObj)
			p.BoomActiv = false
			p.Alive = false
		}
	}

	p.Time += 1.0 / 60.0

}

func (p *ProjectileRB) check_bounce(vert bool) {

	if p.BounceNum > 0 {

		if vert {
			p.Speed.Y *= -1
		} else {
			p.Speed.X *= -1
		}

		p.BounceNum -= 1
	} else {
		p.Alive = false
	}

}

func (p *ProjectileRB) Boom() {

	if !p.BoomActiv {

		p.BoomTime = p.Time

		center := p.Object.Center()

		boom_obj := resolv.NewObject(center.X-p.ExplosionRadius, center.Y-p.ExplosionRadius, p.ExplosionRadius*2, p.ExplosionRadius*2)
		boom_obj.AddTags("explosion")
		p.BoomObj = boom_obj
		p.Space.Add(boom_obj)

		p.BoomActiv = true
	}
}

func (p *ProjectileRB) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, float32(p.Object.Position.X), float32(p.Object.Position.Y), 4, p.Color, false)
	p.Particle.Draw(screen)

	if p.BoomActiv {
		vector.DrawFilledCircle(screen, float32(p.Object.Position.X), float32(p.Object.Position.Y), float32(p.ExplosionRadius), p.Color, false)
	}
}

func (p *ProjectileRB) IsAlive() bool {
	return p.Alive
}

func (p *ProjectileRB) GetObject() *resolv.Object {
	return p.Object
}
