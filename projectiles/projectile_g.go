package projectiles

import (
	"fmt"
	"image/color"
	"math"

	"github.com/AndriiPets/ArcaneAvenger/particles"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

type ProjectileG struct {
	Space    *resolv.Space
	Type     string
	Object   *resolv.Object
	Size     float64
	Speed    resolv.Vector
	Alive    bool
	Color    color.RGBA
	Time     float64
	Particle *particles.ParticleSpawner
	//
	Amplitude float64 // Controls the height of the sine wave
	Frequency float64 // Controls the number of cycles per unit time
	Phase     float64
}

func SpawnProjectileG(space *resolv.Space, position resolv.Vector, direction resolv.Vector, color color.RGBA, size float64) *ProjectileG {

	p := &ProjectileG{
		Space:  space,
		Object: resolv.NewObject(position.X, position.Y, size, size),
		Size:   size,
		Speed:  direction,
		Alive:  true,
		Color:  color,
		//
		Amplitude: 5,
		Frequency: math.Pi,
		Phase:     0,
	}

	p.Particle = particles.NewParticleSpawner(p.Object, p.Color, "stroke_circle", true)

	p.Object.AddTags("projectile")
	space.Add(p.Object)
	fmt.Println("New projectile Green")

	return p
}

func (p *ProjectileG) Update() {
	//p.Speed.Y += 0.1

	px := p.Speed.X * 6
	py := p.Speed.Y * 6
	p.Particle.Update(resolv.NewVector(px, py))

	if check := p.Object.Check(px, 0, "solid", "entity"); check != nil {

		px = check.ContactWithCell(check.Cells[0]).X
		p.Alive = false
	}

	p.Object.Position.X += px

	if check := p.Object.Check(0, py, "solid", "entity"); check != nil {

		py = check.ContactWithCell(check.Cells[0]).Y
		p.Alive = false
	}
	//sine pattern
	yOffset := p.Amplitude * math.Sin(2*math.Pi*p.Frequency*p.Time+p.Phase)
	py += yOffset

	p.Object.Position.Y += py

	p.Object.Update()

	p.Time += 1.0 / 60.0

}

func (p *ProjectileG) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, float32(p.Object.Position.X), float32(p.Object.Position.Y), 4, p.Color, false)
	p.Particle.Draw(screen)
}

func (p *ProjectileG) spawn_add_projectiles() {
	angle1 := 90 * 180 / math.Pi
	angle2 := math.Pi

	SpawnProjectileG(p.Space, p.Object.Position.Scale(8), p.Speed.Rotate(angle1), p.Color, 4)
	SpawnProjectileG(p.Space, p.Object.Position.Scale(8), p.Speed.Rotate(angle2), p.Color, 4)
}

func (p *ProjectileG) IsAlive() bool {
	return p.Alive
}

func (p *ProjectileG) GetObject() *resolv.Object {
	return p.Object
}
