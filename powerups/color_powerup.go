package powerups

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

type ColorPowerup struct {
	Color    string
	Space    *resolv.Space
	Object   *resolv.Object
	ColorMap map[string]color.RGBA
	PickedUp bool
}

func SpawnColorPowerup(space *resolv.Space, powerup_color string, position resolv.Vector, size float64) PowerupInterface {

	c_map := map[string]color.RGBA{
		"blue":  {30, 144, 255, 255},
		"red":   {225, 0, 0, 255},
		"green": {11, 156, 49, 255},
	}

	p := &ColorPowerup{
		Color:    powerup_color,
		Space:    space,
		Object:   resolv.NewObject(position.X, position.Y, size, size),
		PickedUp: false,
		ColorMap: c_map,
	}

	p.Object.AddTags("powerup")
	p.Space.Add(p.Object)

	return p
}

func (p *ColorPowerup) Update() {

	if p.Object.SharesCellsTags("player") {
		fmt.Println("Powerup Hit!")
		p.PickedUp = true
	}

}

func (p *ColorPowerup) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, float32(p.Object.Position.X), float32(p.Object.Position.Y), 5, color.RGBA{255, 255, 255, 225}, false)
	vector.DrawFilledCircle(screen, float32(p.Object.Position.X), float32(p.Object.Position.Y), 4, p.ColorMap[p.Color], false)
}

func (p *ColorPowerup) IsPickedUp() bool {
	return p.PickedUp
}

func (p *ColorPowerup) GetObject() *resolv.Object {
	return p.Object
}

func (p *ColorPowerup) GetColor() string {
	return p.Color
}
