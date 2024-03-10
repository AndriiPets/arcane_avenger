package projectiles

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type ProjectileInterface interface {
	Update()
	Draw(*ebiten.Image)
	IsAlive() bool
	GetObject() *resolv.Object
}
