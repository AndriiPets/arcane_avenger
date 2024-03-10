package projectiles

import "github.com/solarlune/resolv"

func Spawner(space *resolv.Space, p_type string, position resolv.Vector, direction resolv.Vector, size float64) ProjectileInterface {

	var projectile ProjectileInterface

	switch t := p_type; t {
	case "normal":
		projectile = SpawnNormalProjectile(space, position, direction, size)
	case "blue":
		projectile = SpawnBlueProjectile(space, position, direction, size)
	default:
		projectile = SpawnNormalProjectile(space, position, direction, size)
	}

	return projectile
}
