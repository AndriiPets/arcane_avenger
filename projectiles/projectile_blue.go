package projectiles

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

type BlueProjectile struct {
	Space     *resolv.Space
	Type      string
	Object    *resolv.Object
	Size      float64
	Speed     resolv.Vector
	Alive     bool
	BounceNum int
}

func SpawnBlueProjectile(space *resolv.Space, position resolv.Vector, direction resolv.Vector, size float64) *BlueProjectile {

	p := &BlueProjectile{
		Space:     space,
		Object:    resolv.NewObject(position.X, position.Y, size, size),
		Size:      size,
		Speed:     direction,
		Alive:     true,
		BounceNum: 3,
	}

	p.Object.AddTags("projectile")
	space.Add(p.Object)

	return p
}

func (p *BlueProjectile) Update() {
	//p.Speed.Y += 0.1

	px := p.Speed.X * 4 //speed 4
	py := p.Speed.Y * 4

	if check := p.Object.Check(px, 0, "solid"); check != nil {

		px = check.ContactWithCell(check.Cells[0]).X

		if p.BounceNum > 0 {
			p.Speed.X *= -1
			p.BounceNum -= 1
		} else {
			p.Alive = false
		}
	}

	p.Object.Position.X += px

	if check := p.Object.Check(0, py, "solid"); check != nil {

		py = check.ContactWithCell(check.Cells[0]).Y

		if p.BounceNum > 0 {
			p.Speed.Y *= -1
			p.BounceNum -= 1
		} else {
			p.Alive = false
		}
	}

	p.Object.Position.Y += py

	p.Object.Update()
}

func (p *BlueProjectile) Draw(screen *ebiten.Image) {
	p_color := color.RGBA{30, 144, 255, 255} // blue
	vector.DrawFilledCircle(screen, float32(p.Object.Position.X), float32(p.Object.Position.Y), 4, p_color, false)
}

func (p *BlueProjectile) IsAlive() bool {
	return p.Alive
}

func (p *BlueProjectile) GetObject() *resolv.Object {
	return p.Object
}
