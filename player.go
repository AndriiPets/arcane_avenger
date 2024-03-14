package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
	"github.com/yohamta/ganim8/v2"

	"github.com/AndriiPets/ArcaneAvenger/resources/images"
)

var (
	sprite_image *ebiten.Image
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

	IsMoving bool

	Color color.RGBA

	Animation   *ganim8.Animation
	SpriteImage *ebiten.Image

	Dead bool
}

func NewPlayer(space *resolv.Space) *Player {

	p := &Player{
		Object:             resolv.NewObject(550, 128, 16, 16),
		FacingRight:        true,
		Direction:          resolv.NewVector(0, 0),
		Health:             4,
		Color:              color.RGBA{0, 225, 0, 225},
		InvicibleTime:      1.0,
		ProjectileCooldown: 0.3,
		Vulnerable:         true,
		IsMoving:           false,
		Dead:               false,
	}

	space.Add(p.Object)
	p.Object.AddTags("player", "entity")
	p.setup_animations()

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

		if col.HasTags("projectile", "enemy", "explosion") {
			p.take_damage()
		}
	}

	p.Object.Position.X += dx

	if col := p.Object.Check(0, dy); col != nil {

		if col.HasTags("solid", "entity") {
			dy = col.ContactWithCell(col.Cells[0]).Y
		}

		if col.HasTags("projectile", "enemy", "explosion") {
			p.take_damage()
		}

	}

	//update direction unit vector
	p.Object.Position.Y += dy

	p.Object.Update()

	p.check_movement(dx, dy)
	p.Animation.Update()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		p.spawn_projectile()
	}

	p.Time += 1.0 / 60.0

	p.update_cooldowns()

	if p.Health <= 0 {
		p.Dead = true
	}

	//calculate weapon direction vector

	mouse_vec := resolv.NewVector(cursorX-12, cursorY-12)
	player_vec := resolv.NewVector(p.Object.Position.X, p.Object.Position.Y)

	projectile_unit := mouse_vec.Sub(player_vec)
	projectile_unit = projectile_unit.Unit()
	p.Direction = projectile_unit

}

func (p *Player) check_movement(x, y float64) {

	if x == 0.0 && y == 0.0 {
		p.IsMoving = false
	} else {
		p.IsMoving = true
	}
}

func (p *Player) Reset() {
	p.Health = 4
	p.Object.Position.X = 550
	p.Object.Position.Y = 128
}

func (p *Player) setup_animations() {

	img, _, err := image.Decode(bytes.NewReader(images.MainSprite))
	if err != nil {
		log.Fatal(err)
	}

	sprite_image := ebiten.NewImageFromImage(img)
	p.SpriteImage = sprite_image

	g16 := ganim8.NewGrid(16, 16, 48, 80)

	p.Animation = ganim8.New(sprite_image, g16.Frames(1, 1, 2, 1, 3, 1), time.Millisecond*300)
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

func (p *Player) SetPojectileCooldown(cooldown float64) {
	p.ProjectileCooldown = cooldown
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

	//weapon_color := color.RGBA{255, 255, 255, 255}

	//vector.DrawFilledRect(
	//	screen,
	//	posX,
	//	posY,
	//	sizeX,
	//	sizeY,
	//	p.Color, false)

	if p.IsMoving {
		p.Animation.Draw(screen, ganim8.DrawOpts(float64(posX), float64(posY)))
	} else {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(posX), float64(posY))

		screen.DrawImage(p.SpriteImage.SubImage(image.Rect(16, 0, 32, 16)).(*ebiten.Image), op)
	}

	//vector.DrawFilledCircle(screen, p.WeaponPositionX, p.WeaponPositionY, 7, weapon_color, false)

	op := &ebiten.DrawImageOptions{}

	angle := math.Atan2(float64(p.WeaponPositionX), float64(p.WeaponPositionY))
	degrees := 180 * angle / math.Pi

	op.GeoM.Translate(float64(-16/2), float64(-16/2))
	op.GeoM.Rotate(math.Round(degrees))
	op.GeoM.Translate(float64(p.WeaponPositionX), float64(p.WeaponPositionY))

	screen.DrawImage(p.SpriteImage.SubImage(image.Rect(16, 16, 32, 32)).(*ebiten.Image), op)

	//healthbar frame
	vector.DrawFilledRect(
		screen,
		posX+1,
		posY-11,
		float32(p.Health)*4+2,
		6,
		color.RGBA{255, 255, 255, 255}, false)

	//healthbar
	vector.DrawFilledRect(
		screen,
		posX+2,
		posY-10,
		float32(p.Health)*4,
		4,
		color.RGBA{255, 0, 0, 255}, false)
}
