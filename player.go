package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

type Player struct {
	Object          *resolv.Object
	Speed           resolv.Vector
	FacingRight     bool
	Direction       resolv.Vector
	WeaponPositionX float32
	WeaponPositionY float32
	Health          int
	Time            float64

	//damage time
	Vulnerable    bool
	HitTime       float64
	InvicibleTime float64

	//spawn projectile
	SpawnProjectile    bool
	ProjectileCooldown float64
	SpawnTime          float64
	OnCooldown         bool

	Color color.RGBA
}

func NewPlayer(space *resolv.Space) *Player {

	p := &Player{
		Object:             resolv.NewObject(32, 128, 16, 16),
		FacingRight:        true,
		Direction:          resolv.NewVector(0, 0),
		Health:             3,
		Color:              color.RGBA{0, 225, 0, 225},
		InvicibleTime:      1.0,
		ProjectileCooldown: 0.4,
		Vulnerable:         true,
	}

	space.Add(p.Object)
	p.Object.AddTags("player", "entity")

	return p
}

func (p *Player) Update(cursorX, cursorY float64) {
	dx, dy := 0.0, 0.0
	speed := 2.0

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		dy = -speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		dy += speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		dx = -speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		dx += speed
	}

	if col := p.Object.Check(dx, 0); col != nil {

		if col.HasTags("solid", "entity") {
			dx = col.ContactWithCell(col.Cells[0]).X
		}

		if col.HasTags("projectile", "enemy") {
			p.take_damage()
		}
	}

	p.Object.Position.X += dx

	if col := p.Object.Check(0, dy); col != nil {

		if col.HasTags("solid", "entity") {
			dy = col.ContactWithCell(col.Cells[0]).Y
		}

		if col.HasTags("projectile", "enemy") {
			p.take_damage()
		}

	}

	//update direction unit vector
	p.Object.Position.Y += dy

	p.Object.Update()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		p.spawn_projectile()
	}

	p.Time += 1.0 / 60.0

	p.update_cooldowns()

	//calculate weapon direction vector

	mouse_vec := resolv.NewVector(cursorX-12, cursorY-12)
	player_vec := resolv.NewVector(p.Object.Position.X, p.Object.Position.Y)

	projectile_unit := mouse_vec.Sub(player_vec)
	projectile_unit = projectile_unit.Unit()
	p.Direction = projectile_unit

}

func (p *Player) take_damage() {

	if p.Vulnerable {

		p.Health -= 1
		p.HitTime = p.Time

		p.Vulnerable = false
		fmt.Println(p.Health, p.Vulnerable)

	}

}

func (p *Player) spawn_projectile() {

	if !p.OnCooldown {
		p.SpawnProjectile = true

		p.SpawnTime = p.Time
		p.OnCooldown = true
	}
}

func (p *Player) update_cooldowns() {

	//health cooldown
	if !p.Vulnerable {

		p.Color = color.RGBA{225, 0, 0, 255}

		if p.Time-p.HitTime >= p.InvicibleTime {

			fmt.Println(p.Time - p.HitTime)
			p.Vulnerable = true
			p.Color = color.RGBA{0, 225, 0, 225}
		}
	}

	//weapon cooldown
	if p.OnCooldown {

		if p.Time-p.SpawnTime >= p.ProjectileCooldown {

			fmt.Println("cooldown end", p.Time)
			p.OnCooldown = false

		}
	}
}

func (p *Player) Draw(screen *ebiten.Image) {

	posX, posY := float32(p.Object.Position.X), float32(p.Object.Position.Y)
	sizeX, sizeY := float32(p.Object.Size.X), float32(p.Object.Size.Y)
	var weapon_vector resolv.Vector

	if p.OnCooldown {
		weapon_vector = p.Direction.Scale(10)
	} else {
		weapon_vector = p.Direction.Scale(20)
	}

	centerX, CenterY := posX+(sizeX/2), posY+(sizeY/2)
	weapon_position_x, weapon_position_y := centerX+float32(weapon_vector.X), CenterY+float32(weapon_vector.Y)

	p.WeaponPositionX = weapon_position_x
	p.WeaponPositionY = weapon_position_y

	weapon_color := color.RGBA{255, 255, 255, 255}

	vector.DrawFilledRect(
		screen,
		posX,
		posY,
		sizeX,
		sizeY,
		p.Color, false)

	vector.DrawFilledCircle(screen, p.WeaponPositionX, p.WeaponPositionY, 7, weapon_color, false)
}
