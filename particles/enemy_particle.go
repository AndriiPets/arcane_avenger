package particles

import (
	"image/color"
	"math/rand"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

type BloodParticle struct {
	Scale     float32
	Position  resolv.Vector
	SpawnTime float64
	Color     color.RGBA
}

type BloodParticleSpawner struct {
	Time      float64
	Particles map[uuid.UUID]*BloodParticle
	SpawnTime float64
	BloodTime float64
}

func NewBloodParticleSpawner() *BloodParticleSpawner {

	p_map := make(map[uuid.UUID]*BloodParticle)

	ps := &BloodParticleSpawner{
		Particles: p_map,
		BloodTime: 15.0,
	}

	return ps
}

func (bp *BloodParticleSpawner) Update() {

	for key, b := range bp.Particles {
		if bp.Time-b.SpawnTime >= bp.BloodTime {
			delete(bp.Particles, key)
		}
	}

	bp.Time += 1.0 / 60.0
}

func (bp *BloodParticleSpawner) Draw(screen *ebiten.Image) {

	for _, b := range bp.Particles {
		vector.DrawFilledCircle(screen, float32(b.Position.X), float32(b.Position.Y), b.Scale, b.Color, false)
	}
}

func (bp *BloodParticleSpawner) Spawn(position resolv.Vector, b_color color.RGBA) {

	//spread := 0.2 - rand.Float64()*0.4
	size := rand.Intn(8)
	distance := rand.Float64()

	pos := position

	blood := &BloodParticle{
		Scale:     float32(size),
		Position:  pos.Scale(distance),
		Color:     b_color,
		SpawnTime: bp.Time,
	}

	bp.Particles[uuid.New()] = blood

}
