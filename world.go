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
	"github.com/AndriiPets/ArcaneAvenger/particles"
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
	MosterSpawners    []*MonsterSpawner
	BloodSpawner      *particles.BloodParticleSpawner
	ProjectileSpawner *projectiles.Spawner
	KillCount         int
}

func NewWorld(g *Game) *World {

	p := make(map[uuid.UUID]projectiles.ProjectileInterface)
	pow := make(map[uuid.UUID]powerups.PowerupInterface)
	en := make(map[uuid.UUID]enemies.EnemyInterface)

	w := &World{Game: g, Projectiles: p, Powerups: pow, Enemies: en, KillCount: 0}
	w.Init()
	return w
}

func (w *World) Init() {

	gw := float64(w.Game.Width) * 2
	gh := float64(w.Game.Height) * 2

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
	w.BloodSpawner = particles.NewBloodParticleSpawner()
}

func (w *World) Update() {

	w.Player.Update(w.Game.Camera.cursorX, w.Game.Camera.cursorY)

	w.Player.SetPojectileCooldown(w.ProjectileSpawner.Cooldowns[w.ProjectileSpawner.State])

	//spawn projectile
	if w.Player.SpawnProjectile {
		fmt.Println("Hit")

		spread := 0.2 - rand.Float64()*0.4
		//spread *= 0.1

		proj := w.ProjectileSpawner.Spawn(
			resolv.NewVector(float64(w.Player.WeaponPositionX),
				float64(w.Player.WeaponPositionY)),
			w.Player.Direction.Rotate(spread),
			4)

		for _, p := range proj {
			w.Projectiles[uuid.New()] = p
		}

		//fmt.Println(len(w.Projectiles))
		w.Player.SpawnProjectile = false
	}

	//if inpututil.IsKeyJustPressed(ebiten.KeyE) {
	//	p := powerups.SpawnColorPowerup(
	//		w.Space,
	//		"blue",
	//		resolv.NewVector(float64(w.Player.WeaponPositionX+20), float64(w.Player.WeaponPositionY+20)),
	//		4,
	//	)
	//	w.Powerups[uuid.New()] = p
	//}

	//if inpututil.IsKeyJustPressed(ebiten.KeyF) {
	//	e := enemies.NewEnemy(w.Space, resolv.NewVector(float64(w.Player.WeaponPositionX+20), float64(w.Player.WeaponPositionY+20)))
	//	w.Enemies[uuid.New()] = e
	//}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		w.Game.Camera.Rotation += 1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		w.Game.Camera.Reset()
	}

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

		if e.HitPlayer() {

			w.Player.take_damage()
			e.HitPlayerComplete()

		}

		if !e.IsVunerable() {
			chance := rand.Intn(100)

			if chance <= 5 {
				w.BloodSpawner.Spawn(e.GetObject().Position, e.GetColorRGBA())
			}

		}

		if !e.IsAlive() {

			w.KillCount += 1

			if e.DeathDrop() {

				p := powerups.SpawnColorPowerup(
					w.Space,
					e.GetColor(),
					resolv.NewVector(float64(e.GetObject().Right()), float64(e.GetObject().Bottom())),
					4,
				)
				w.Powerups[uuid.New()] = p

			}

			delete(w.Enemies, key)
			w.Space.Remove(e.GetObject())
		}

		e.Update(w.Player.Object)
	}

	//update blood
	w.BloodSpawner.Update()

	//update monster spawner
	for _, ms := range w.MosterSpawners {
		ms.Update()
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

func (w *World) Reset() {

	for key, p := range w.Powerups {
		//player picked up the powerup
		delete(w.Powerups, key)
		w.Space.Remove(p.GetObject())
	}

	for key, e := range w.Enemies {
		delete(w.Enemies, key)
		w.Space.Remove(e.GetObject())
	}

	for key, p := range w.Projectiles {
		delete(w.Projectiles, key)
		w.Space.Remove(p.GetObject())
	}

	w.ProjectileSpawner.Reset()

	w.Player.Reset()
}

func (w *World) AddMonsterSpawner(sp *MonsterSpawner) {
	w.MosterSpawners = append(w.MosterSpawners, sp)
}

func (w *World) IsPlayerDead() bool {
	if w.Player.Dead {
		w.Player.Dead = false
		return true
	} else {
		return false
	}
}

func (w *World) GetColorInfo() map[string]int {
	return w.ProjectileSpawner.ColorsAmmount
}

func (w *World) GetKillCount() int {
	return w.KillCount
}

func (w *World) SetKillCount(count int) {
	w.KillCount = count
}

func (w *World) GetSpace() *resolv.Space {
	return w.Space
}

func (w *World) AddEnemy(e enemies.EnemyInterface) {
	w.Enemies[uuid.New()] = e
}

func (w *World) GetPlayerPos(screenX, screenY int) (float64, float64) {
	cx := w.Player.Object.Position.X - (float64(screenX) / 2)
	cy := w.Player.Object.Position.Y - (float64(screenY) / 2)
	//fmt.Println(w.Player.Object.Position.X, w.Player.Object.Position.Y)
	return cx, cy
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

	//blood draw
	w.BloodSpawner.Draw(screen)

	//draw player
	w.Player.Draw(screen)

	//draw powerups
	for _, p := range w.Powerups {
		p.Draw(screen)
	}

	//draw enemies
	for _, e := range w.Enemies {
		e.Draw(screen)
	}

	//draw projectiles
	for _, p := range w.Projectiles {
		p.Draw(screen)
	}

	if w.Game.Debug {
		w.Game.DebugDraw(screen, w.Space)
	}

}
