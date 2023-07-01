package game

import (
	"ahasuerus/config"
	"ahasuerus/resources"
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
	rl.SetConfigFlags(rl.FlagVsyncHint)

	resources.LoadShaderCache(resources.SdfShader)
	literata := resources.LoadFont(resources.Literata)
	rl.SetTextureFilter(literata.Texture, rl.TextureFilterMode(rl.RL_TEXTURE_FILTER_BILINEAR))

	nextScene := scene.GetScene(scene.Menu)
	for nextScene != nil {
		nextScene = nextScene.Run()
	}

	resources.UnloadFont(resources.Literata)
	resources.UnloadShaderCache(resources.SdfShader)

	scene.UnloadAllScenes()

	rl.CloseAudioDevice()
}
