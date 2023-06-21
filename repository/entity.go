package repository

import (
	"ahasuerus/models"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const dataId = "data"

type Level struct {
	Name               string
	Characters         []models.Npc
	Lights             []models.Light
	CollissionHitboxes []models.CollisionHitbox
	Images             []models.Image
	ParticleSources    []models.ParticleSource

	CameraPos          rl.Vector2
	CameraStartEndMove rl.Vector2
	PlayerPos          rl.Vector2

	MusicTheme        string
	MusicThemeReverse        string
}

func GetLevel(levelName string) Level {
	var level Level
	err := db.Read(levelName, dataId, &level)
	if err != nil {
		panic(err)
	}
	level.Name = levelName
	return level
}

func (level *Level) SaveLevel() {
	err := db.Write(level.Name, dataId, level)
	if err != nil {
		panic(err)
	}
}
