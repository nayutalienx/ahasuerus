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

func (level *Level) Size() rl.Vector2 {
	var size rl.Vector2
	for _, point := range level.Images {
		for _, polygon := range point.BaseEditorItem.Polygons {
			for _, polygonPoint := range polygon.Points {
				if polygonPoint.X > size.X && polygonPoint.Y > size.Y {
					size = polygonPoint
				}
			}
		}
	}
	return size
}