package powerups

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type PowerupInterface interface {
	Update()
	Draw(*ebiten.Image)
	IsPickedUp() bool
	GetObject() *resolv.Object
	GetColor() string
}
