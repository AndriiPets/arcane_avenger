package projectiles

import (
	"fmt"
	"image/color"
	"math/rand"

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
		MaxAmmount:    51,
		Space:         space,
	}

	return s
}

func (s *Spawner) AddColor(color string) {

	s.ColorsAmmount[color] += 17

	if s.ColorsAmmount[color] > s.MaxAmmount {
		s.ColorsAmmount[color] = s.MaxAmmount
	}
}

func (s *Spawner) spend_color_point() {

	random_c := []string{}

	for key, point := range s.ColorsAmmount {
		if point > 0 {
			random_c = append(random_c, key)
		}
	}

	s.ColorsAmmount[random_c[rand.Intn(len(random_c))]] -= 1

}

func (s *Spawner) Spawn(position resolv.Vector, direction resolv.Vector, size float64) ProjectileInterface {

	var projectile ProjectileInterface

	factor := 255 / s.MaxAmmount

	red := s.ColorsAmmount["red"] * factor
	green := s.ColorsAmmount["green"] * factor
	blue := s.ColorsAmmount["blue"] * factor

	color := color.RGBA{uint8(red), uint8(green), uint8(blue), 255}

	switch r, g, b := s.ColorsAmmount["red"], s.ColorsAmmount["green"], s.ColorsAmmount["blue"]; r+g+b >= 0 {
	case b > 0 && r <= 0 && g <= 0:
		//only blue
		projectile = SpawnProjectileB(s.Space, position, direction, color, size)
		s.spend_color_point()
		fmt.Println(s.ColorsAmmount)

	case r > 0 && b <= 0 && g <= 0:
		//only red
		projectile = SpawnProjectileR(s.Space, position, direction, color, size)
		s.spend_color_point()
		fmt.Println(s.ColorsAmmount)

	case g > 0 && b <= 0 && r <= 0:
		//only green
		projectile = SpawnProjectileG(s.Space, position, direction, color, size)
		s.spend_color_point()
		fmt.Println(s.ColorsAmmount)

	case r > 0 && g > 0 && b <= 0:
		//red and green
		projectile = SpawnProjectileRG(s.Space, position, direction, color, size)
		s.spend_color_point()
		fmt.Println(s.ColorsAmmount)

	case r > 0 && b > 0 && g <= 0:
		//red and blue
		projectile = SpawnProjectileRB(s.Space, position, direction, color, size)
		s.spend_color_point()
		fmt.Println(s.ColorsAmmount)

	case g > 0 && b > 0 && r <= 0:
		projectile = SpawnProjectileGB(s.Space, position, direction, color, size)
		s.spend_color_point()
		fmt.Println(s.ColorsAmmount)

	case r > 0 && g > 0 && b > 0:
		//all colors
		projectile = SpawnProjectileRGB(s.Space, position, direction, color, size)
		s.spend_color_point()
		fmt.Println(s.ColorsAmmount)

	default:
		//no color
		fmt.Println(s.ColorsAmmount)
		projectile = SpawnNormalProjectile(s.Space, position, direction, color, size)
	}

	return projectile
}
