package main

import (
	"github.com/AndriiPets/ArcaneAvenger/enemies"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type WorldInterface interface {
	Init()
	Update()
	Draw(*ebiten.Image)
	GetPlayerPos(int, int) (float64, float64)
	GetColorInfo() map[string]int
	GetKillCount() int
	AddMonsterSpawner(sp *MonsterSpawner)
	GetSpace() *resolv.Space
	AddEnemy(enemies.EnemyInterface)
	Reset()
	IsPlayerDead() bool
	SetKillCount(count int)
}
