package projectiles

import (
	"image/color"

	"github.com/solarlune/resolv"
)

type ProjectileG struct {
	NormalProjectile
}

func SpawnProjectileG(space *resolv.Space, position resolv.Vector, direction resolv.Vector, color color.RGBA, size float64) *ProjectileG {

	p := &ProjectileG{
		NormalProjectile: NormalProjectile{
			Space:  space,
			Object: resolv.NewObject(position.X, position.Y, size, size),
			Size:   size,
			Speed:  direction,
			Alive:  true,
			Color:  color,
		},
	}

	p.Object.AddTags("projectile")
	space.Add(p.Object)

	return p
}