package main

import (
	//"fmt"
	"fmt"
	"image/color"
	"math/rand"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"

	"github.com/AndriiPets/ArcaneAvenger/enemies"
	"github.com/AndriiPets/ArcaneAvenger/powerups"
	"github.com/AndriiPets/ArcaneAvenger/projectiles"
)

type World struct {
	Game              *Game
	Space             *resolv.Space
	Player            *Player
	Projectiles       map[uuid.UUID]projectiles.ProjectileInterface
	Powerups          map[uuid.UUID]powerups.PowerupInterface
	Enemies           map[uuid.UUID]enemies.EnemyInterface
	ProjectileSpawner *projectiles.Spawner
}

func NewWorld(g *Game) *World {

	p := make(map[uuid.UUID]projectiles.ProjectileInterface)
	pow := make(map[uuid.UUID]powerups.PowerupInterface)
	en := make(map[uuid.UUID]enemies.EnemyInterface)

	w := &World{Game: g, Projectiles: p, Powerups: pow, Enemies: en}
	w.Init()
	return w
}

func (w *World) Init() {

	gw := float64(w.Game.Width)
	gh := float64(w.Game.Height)

	cell_size := 8

	w.Space = resolv.NewSpace(int(gw), int(gh), cell_size, cell_size)

	// Construct geometry
	geometry := []*resolv.Object{

		resolv.NewObject(0, 0, 16, gh),
		resolv.NewObject(gw-16, 0, 16, gh),
		resolv.NewObject(0, 0, gw, 16),
		resolv.NewObject(0, gh-24, gw, 32),
		resolv.NewObject(0, gh-24, gw, 32),

		//resolv.NewObject(200, -160, 16, gh),
	}

	w.Space.Add(geometry...)

	for _, o := range w.Space.Objects() {
		o.AddTags("solid")
	}

	w.Player = NewPlayer(w.Space)
	w.ProjectileSpawner = projectiles.NewSpawner(w.Space)
}

func (w *World) Update() {

	w.Player.Update()

	//spawn projectile
	if w.Player.SpawnProjectile {
		fmt.Println("Hit")

		spread := 0.2 - rand.Float64()*0.4
		//spread *= 0.1

		p := w.ProjectileSpawner.Spawn(
			resolv.NewVector(float64(w.Player.WeaponPositionX),
				float64(w.Player.WeaponPositionY)),
			w.Player.Direction.Rotate(spread),
			4)

		w.Projectiles[uuid.New()] = p
		//fmt.Println(len(w.Projectiles))
		w.Player.SpawnProjectile = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		p := powerups.SpawnColorPowerup(
			w.Space,
			"blue",
			resolv.NewVector(float64(w.Player.WeaponPositionX+20), float64(w.Player.WeaponPositionY+20)),
			4,
		)
		w.Powerups[uuid.New()] = p
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		e := enemies.NewEnemy(w.Space, resolv.NewVector(float64(w.Player.WeaponPositionX+20), float64(w.Player.WeaponPositionY+20)))
		w.Enemies[uuid.New()] = e
	}
	//player picked up the powerup

	//update projectiles
	for key, p := range w.Projectiles {
		if !p.IsAlive() {
			delete(w.Projectiles, key)
			w.Space.Remove(p.GetObject())
		}

		p.Update()
	}

	//update enemies
	for key, e := range w.Enemies {
		if !e.IsAlive() {
			delete(w.Enemies, key)
			w.Space.Remove(e.GetObject())
		}

		e.Update()
	}

	//update powerups
	for key, p := range w.Powerups {
		//player picked up the powerup
		if p.IsPickedUp() {

			w.fill_color(p.GetColor())

			delete(w.Powerups, key)
			w.Space.Remove(p.GetObject())
		}

		p.Update()

	}
}

func (w *World) fill_color(color string) {
	w.ProjectileSpawner.AddColor(color)
}

func (w *World) Draw(screen *ebiten.Image) {

	//draw map
	for _, o := range w.Space.Objects() {
		drawColor := color.RGBA{60, 60, 60, 255}
		if o.HasTags("solid") {
			vector.DrawFilledRect(
				screen,
				float32(o.Position.X),
				float32(o.Position.Y),
				float32(o.Size.X),
				float32(o.Size.Y),
				drawColor, false)
		}

	}

	//draw player
	w.Player.Draw(screen)

	//draw projectiles
	for _, p := range w.Projectiles {
		p.Draw(screen)
	}

	//draw powerups
	for _, p := range w.Powerups {
		p.Draw(screen)
	}

	//draw enemies
	for _, e := range w.Enemies {
		e.Draw(screen)
	}

	//draw aim circle
	mouseX, mouseY := ebiten.CursorPosition()
	mx, my := float32(mouseX), float32(mouseY)

	aimColor := color.RGBA{0, 225, 0, 225}

	vector.StrokeCircle(screen, mx, my, 12, 2, aimColor, false)
	vector.DrawFilledCircle(screen, mx, my, 2, aimColor, false)

	if w.Game.Debug {
		w.Game.DebugDraw(screen, w.Space)
	}

}
