package enemies

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type EnemyInterface interface {
	Update()
	Draw(*ebiten.Image)
	IsAlive() bool
	GetObject() *resolv.Object
}
