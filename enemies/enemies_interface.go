package enemies

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type EnemyInterface interface {
	Update(*resolv.Object)
	Draw(*ebiten.Image)
	IsAlive() bool
	GetObject() *resolv.Object
	DeathDrop() bool
	GetColor() string
	HitPlayer() bool
	HitPlayerComplete()
}
