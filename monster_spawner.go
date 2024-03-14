package main

import (
	"fmt"

	"github.com/AndriiPets/ArcaneAvenger/enemies"
	"github.com/solarlune/resolv"
)

type MonsterSpawner struct {
	Position     resolv.Vector
	World        WorldInterface
	Time         float64
	WaveCooldown float64
	SpawnTime    float64
	WaveReady    bool
	WaveNum      int
	LevelUp      float64
}

func NewMonsterSpawner(wld WorldInterface, pos resolv.Vector) *MonsterSpawner {

	ms := &MonsterSpawner{
		Position:     pos,
		World:        wld,
		Time:         0.0,
		WaveCooldown: 10.0,
		WaveReady:    true,
		WaveNum:      4,
		LevelUp:      20.0,
	}

	return ms
}

func (ms *MonsterSpawner) Update() {
	ms.spawn_wave()
	ms.Time += 1.0 / 60.0

	ms.update_cooldown()

	if ms.Time > ms.LevelUp {
		ms.WaveCooldown -= 1
		ms.LevelUp += ms.Time
	}
}

func (ms *MonsterSpawner) spawn_wave() {

	if ms.WaveReady {
		fmt.Println(ms.WaveReady)

		positions := map[int]resolv.Vector{
			1: resolv.NewVector(ms.Position.X+5, ms.Position.Y),
			2: resolv.NewVector(ms.Position.X-5, ms.Position.Y),
			3: resolv.NewVector(ms.Position.X, ms.Position.Y-5),
			4: resolv.NewVector(ms.Position.X, ms.Position.Y+5),
		}

		for i := 0; i <= ms.WaveNum; i++ {
			e := enemies.NewEnemy(ms.World.GetSpace(), positions[i])

			ms.World.AddEnemy(e)
		}

		ms.SpawnTime = ms.Time
		ms.WaveReady = false
	}
}

func (ms *MonsterSpawner) update_cooldown() {

	if !ms.WaveReady {

		if ms.Time-ms.SpawnTime >= ms.WaveCooldown {

			ms.WaveReady = true
		}
	}
}
