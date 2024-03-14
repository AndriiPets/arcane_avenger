package projectiles

import (
	"fmt"
	"image/color"
	"math"

	//"math"
	//"math/rand"

	"github.com/solarlune/resolv"
)

type Spawner struct {
	ColorsAmmount map[string]int
	MaxAmmount    int
	Space         *resolv.Space
	State         string
	Cooldowns     map[string]float64
	Cost          map[string]int
}

func NewSpawner(space *resolv.Space) *Spawner {
	c_map := map[string]int{
		"red":   0,
		"green": 0,
		"blue":  0,
	}

	cds_map := map[string]float64{
		"normal":     0.4,
		"red":        0.7,
		"green":      0.5,
		"blue":       0.3,
		"red_green":  0.6,
		"red_blue":   0.4,
		"green_blue": 0.4,
		"pre_rgb":    0.3,
		"rgb":        0.1,
	}

	cost_map := map[string]int{
		"normal":     0,
		"red":        2,
		"green":      1,
		"blue":       1,
		"red_green":  2,
		"red_blue":   3,
		"green_blue": 2,
		"pre_rgb":    4,
		"rgb":        4,
	}

	s := &Spawner{
		ColorsAmmount: c_map,
		MaxAmmount:    51,
		Space:         space,
		State:         "normal",
		Cooldowns:     cds_map,
		Cost:          cost_map,
	}

	return s
}

func (s *Spawner) AddColor(color string) {

	s.ColorsAmmount[color] += 17

	if s.ColorsAmmount[color] > s.MaxAmmount {
		s.ColorsAmmount[color] = s.MaxAmmount
	}
}

func (s *Spawner) Reset() {
	for key := range s.ColorsAmmount {
		s.ColorsAmmount[key] = 0
	}
}

func (s *Spawner) spend_color_point() {

	//random_c := []string{}

	for key, point := range s.ColorsAmmount {
		cost := s.Cost[s.State]
		if point >= cost {
			s.ColorsAmmount[key] -= cost
			//random_c = append(random_c, key)
		} else {
			s.ColorsAmmount[key] = 0
		}
	}

	//s.ColorsAmmount[random_c[rand.Intn(len(random_c))]] -= 1

}

func (s *Spawner) Spawn(position resolv.Vector, direction resolv.Vector, size float64) []ProjectileInterface {

	var projectile []ProjectileInterface

	factor := 255 / s.MaxAmmount

	rgb_limit := 34

	red := s.ColorsAmmount["red"] * factor
	green := s.ColorsAmmount["green"] * factor
	blue := s.ColorsAmmount["blue"] * factor

	color_p := color.RGBA{uint8(red), uint8(green), uint8(blue), 255}

	switch r, g, b := s.ColorsAmmount["red"], s.ColorsAmmount["green"], s.ColorsAmmount["blue"]; r+g+b >= 0 {
	case b > 0 && r <= 0 && g <= 0:
		//only blue
		s.State = "blue"
		projectile = append(projectile, SpawnProjectileB(s.Space, position, direction, color_p, size))
		s.spend_color_point()
		fmt.Println(s.ColorsAmmount)

	case r > 0 && b <= 0 && g <= 0:
		//only red
		s.State = "red"
		projectile = append(
			projectile,
			SpawnProjectileR(s.Space, position, direction, color_p, size),
		)
		s.spend_color_point()
		fmt.Println(s.ColorsAmmount)

	case g > 0 && b <= 0 && r <= 0:
		//only green
		s.State = "green"
		projectile = append(projectile, SpawnProjectileG(s.Space, position, direction, color_p, size))
		s.spend_color_point()
		fmt.Println(s.ColorsAmmount)

	case r > 0 && g > 0 && b <= 0:
		//red and green
		s.State = "red_green"
		projectile = append(projectile, SpawnProjectileRG(s.Space, position, direction, color_p, size))
		s.spend_color_point()
		fmt.Println(s.ColorsAmmount)

	case r > 0 && b > 0 && g <= 0:
		//red and blue
		s.State = "red_blue"
		projectile = append(projectile, SpawnProjectileRB(s.Space, position, direction, color_p, size))
		s.spend_color_point()
		fmt.Println(s.ColorsAmmount)

	case g > 0 && b > 0 && r <= 0:
		//green and blue
		s.State = "green_blue"
		projectile = append(projectile,
			SpawnProjectileGB(s.Space, position, direction, color.RGBA{uint8(red), uint8(blue), uint8(green), 255}, size),
			SpawnProjectileG(s.Space, position, direction, color_p, size),
		)

		s.spend_color_point()
		fmt.Println(s.ColorsAmmount)
	case g > 0 && b > 0 && r > 0:
		//not at rgb limit
		s.State = "pre_rgb"
		projectile = append(projectile,
			SpawnProjectileR(s.Space, position, direction, color_p, size),
			SpawnProjectileG(s.Space, position, direction.Rotate(math.Pi/3), color_p, size),
			SpawnProjectileB(s.Space, position, direction.Rotate(-math.Pi/3), color_p, size),
		)
		s.spend_color_point()
	case r > rgb_limit && g > rgb_limit && b > rgb_limit:
		//all colors
		s.State = "rgb"
		projectile = append(projectile,
			SpawnProjectileRGB(s.Space, position, direction, color_p, size+4),
			SpawnProjectileRGB(s.Space, position, direction.Rotate(math.Pi/2).Scale(3), color_p, size+4),
			SpawnProjectileRGB(s.Space, position, direction.Rotate(-math.Pi/2).Scale(3), color_p, size+4),
		)
		s.spend_color_point()
		fmt.Println(s.ColorsAmmount)

	default:
		//no color
		s.State = "normal"
		fmt.Println(s.ColorsAmmount)
		projectile = append(projectile, SpawnNormalProjectile(s.Space, position, direction, color_p, size))
	}

	return projectile
}
