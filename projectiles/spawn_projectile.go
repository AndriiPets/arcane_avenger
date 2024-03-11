package projectiles

import (
	"fmt"

	"github.com/solarlune/resolv"
)

type Spawner struct {
	ColorsAmmount map[string]int
	MaxAmmount    int
	Space         *resolv.Space
}

func NewSpawner(space *resolv.Space) *Spawner {
	c_map := map[string]int{
		"red":   0,
		"green": 0,
		"blue":  0,
	}
	s := &Spawner{
		ColorsAmmount: c_map,
		MaxAmmount:    10,
		Space:         space,
	}

	return s
}

func (s *Spawner) AddColor(color string) {

	s.ColorsAmmount[color] = s.MaxAmmount
}

func (s *Spawner) spend_color_point() {

	for key, point := range s.ColorsAmmount {
		if point > 0 {
			s.ColorsAmmount[key] = point - 1
		}
	}

}

func (s *Spawner) Spawn(position resolv.Vector, direction resolv.Vector, size float64) ProjectileInterface {

	var projectile ProjectileInterface

	switch r, g, b := s.ColorsAmmount["red"], s.ColorsAmmount["green"], s.ColorsAmmount["blue"]; r+g+b >= 0 {
	case b > 0 && r <= 0 && g <= 0:
		projectile = SpawnBlueProjectile(s.Space, position, direction, size)
		s.spend_color_point()
		fmt.Println(s.ColorsAmmount)
	default:
		fmt.Println(s.ColorsAmmount)
		projectile = SpawnNormalProjectile(s.Space, position, direction, size)
	}

	return projectile
}
