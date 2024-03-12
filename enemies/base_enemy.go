package enemies

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

type BaseEnemy struct {
	Space     *resolv.Space
	Object    *resolv.Object
	Health    int
	Color     color.RGBA
	Direction resolv.Vector
	Alive     bool
	Time      float64
	ColorName string

	//damage time
	Vulnerable    bool
	HitTime       float64
	InvicibleTime float64

	//player damage
	PlayerDamaged bool
}

func NewEnemy(space *resolv.Space, position resolv.Vector) EnemyInterface {

	colors := [3]color.RGBA{{30, 144, 255, 255}, {225, 0, 0, 255}, {11, 156, 49, 255}}
	names := [3]string{"blue", "red", "green"}
	pick := rand.Intn(len(colors))

	e := &BaseEnemy{
		Space:         space,
		Object:        resolv.NewObject(position.X, position.Y, 16, 16),
		Health:        3,
		Color:         colors[pick],
		InvicibleTime: 1.0,
		Vulnerable:    true,
		Alive:         true,
		ColorName:     names[pick],
	}

	e.Direction = e.GetMovementVector()

	e.Object.AddTags("enemy", "entity")
	e.Space.Add(e.Object)

	fmt.Println("Enemy Created!")

	return e
}

func (e *BaseEnemy) GetMovementVector() resolv.Vector {
	v := resolv.NewVector(e.Object.Position.X, e.Object.Position.Y).Unit()
	v.Rotate(rand.Float64())

	return v
}

func (e *BaseEnemy) take_damage() {

	if e.Vulnerable {

		e.Health -= 1
		e.HitTime = e.Time

		e.Vulnerable = false
		fmt.Println(e.Health, e.Vulnerable)

	}

}

func (e *BaseEnemy) update_cooldowns() {

	//health cooldown
	if !e.Vulnerable {

		e.Color = color.RGBA{225, 0, 0, 255}

		if e.Time-e.HitTime >= e.InvicibleTime {

			fmt.Println(e.Time - e.HitTime)
			e.Vulnerable = true
			e.Color = color.RGBA{0, 225, 0, 225}
		}
	}
}

func (e *BaseEnemy) update_direction(p *resolv.Object) {

	player_vec := resolv.NewVector(p.Position.X, p.Position.Y)
	enemy_vec := resolv.NewVector(e.Object.Position.X, e.Object.Position.Y)

	unit := player_vec.Sub(enemy_vec).Unit()

	e.Direction = unit
}

func (e *BaseEnemy) Update(p *resolv.Object) {

	px, py := e.Direction.X, e.Direction.Y

	if col := e.Object.Check(px, 0); col != nil {

		if col.HasTags("solid", "entity") {
			px = col.ContactWithCell(col.Cells[0]).X
		}

		if col.HasTags("player") {
			e.PlayerDamaged = true
		}

		if col.HasTags("projectile") {
			e.take_damage()
		}
	}

	e.Object.Position.X += px

	if col := e.Object.Check(0, py); col != nil {

		if col.HasTags("solid", "entity") {
			py = col.ContactWithCell(col.Cells[0]).Y
		}

		if col.HasTags("player") {
			e.PlayerDamaged = true
		}

		if col.HasTags("projectile") {
			e.take_damage()
		}
	}

	e.Object.Position.Y += py

	e.Object.Update()
	e.update_direction(p)

	e.Time += 1.0 / 60.0

	e.update_cooldowns()
	e.check_death()

}

func (e *BaseEnemy) IsAlive() bool {
	return e.Alive
}

func (e *BaseEnemy) HitPlayer() bool {
	return e.PlayerDamaged
}

func (e *BaseEnemy) HitPlayerComplete() {
	e.PlayerDamaged = false
}

func (e *BaseEnemy) DeathDrop() bool {

	chance := rand.Intn(10)
	if chance >= 5 {
		return true
	}
	return false
}

func (e *BaseEnemy) check_death() {
	if e.Health <= 0 {
		e.Alive = false
	}
}

func (e *BaseEnemy) GetObject() *resolv.Object {
	return e.Object
}

func (e *BaseEnemy) GetColor() string {
	return e.ColorName
}

func (e *BaseEnemy) Draw(screen *ebiten.Image) {

	posX, posY := float32(e.Object.Position.X), float32(e.Object.Position.Y)
	sizeX, sizeY := float32(e.Object.Size.X), float32(e.Object.Size.Y)

	vector.DrawFilledRect(
		screen,
		posX,
		posY,
		sizeX,
		sizeY,
		e.Color, false)

}
