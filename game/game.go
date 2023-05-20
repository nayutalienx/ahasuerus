package game

import (
	"ahasuerus/config"
	"ahasuerus/scene"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	APPLICATION = "ahasuerus"
)

var (
	WIDTH, HEIGHT = config.GetResolution()
	FPS           = config.GetFPS()
)

func Start() {
	rl.InitWindow(int32(WIDTH), int32(HEIGHT), APPLICATION)
	defer rl.CloseWindow()
	rl.SetTargetFPS(FPS)
	rl.InitAudioDevice()
	rl.SetConfigFlags(rl.FlagMsaa4xHint)

	nextScene := scene.NewMenuScene().Run()
	for nextScene != nil {
		nextScene = nextScene.Run()
	}

	rl.CloseAudioDevice()
}
