package projectiles

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"

	"github.com/AndriiPets/ArcaneAvenger/particles"
)

type ProjectileB struct {
	Space     *resolv.Space
	Type      string
	Object    *resolv.Object
	Size      float64
	Speed     resolv.Vector
	Alive     bool
	Color     color.RGBA
	BounceNum int
	Particle  *particles.ParticleSpawner
}

func SpawnProjectileB(space *resolv.Space, position resolv.Vector, direction resolv.Vector, color color.RGBA, size float64) *ProjectileB {

	p := &ProjectileB{
		Space:     space,
		Object:    resolv.NewObject(position.X, position.Y, size, size),
		Size:      size,
		Speed:     direction,
		Alive:     true,
		BounceNum: 3,
		Color:     color,
	}

	p.Particle = particles.NewParticleSpawner(p.Object, p.Color, "circle", true)

	p.Object.AddTags("projectile")
	space.Add(p.Object)

	return p
}

func (p *ProjectileB) Update() {
	//p.Speed.Y += 0.1

	px := p.Speed.X * 6 //speed 4
	py := p.Speed.Y * 6

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
}

func (p *ProjectileB) check_bounce(vert bool) {

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

func (p *ProjectileB) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, float32(p.Object.Position.X), float32(p.Object.Position.Y), 4, p.Color, false)
	p.Particle.Draw(screen)
}

func (p *ProjectileB) IsAlive() bool {
	return p.Alive
}

func (p *ProjectileB) GetObject() *resolv.Object {
	return p.Object
}
