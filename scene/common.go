package scene

import (
	"ahasuerus/config"
	"ahasuerus/models"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type SceneProp string

const (
	CameraLevelMargin SceneProp = "cameraLevelMargin"
)

type SceneId int

const (
	Menu SceneId = iota
	Start
	Editor
)

var (
	WIDTH, HEIGHT = config.GetResolution()
	sceneMap      = make(map[SceneId]models.Scene, 0)
	lastScene     SceneId
)

func GetScene(id SceneId) models.Scene {
	scene, ok := sceneMap[id]
	if ok {
		return scene
	}

	sceneNames := map[SceneId]string{
		Start: "start",
	}

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawText("LOAD SCENE", int32(WIDTH)/3, int32(HEIGHT)/2, 90, rl.Gold)
	rl.EndDrawing()

	switch id {
	case Menu:
		scene = NewMenuScene()
	case Start:
		scene = NewGameScene(sceneNames[Start])
	case Editor:
		UnloadScene(Start)
		scene = NewEditScene(sceneNames[lastScene], lastScene)
	}

	if scene == nil {
		panic("scene not found")
	}

	if id != Editor {
		sceneMap[id] = scene
		lastScene = id
	}

	return scene
}

func UnloadScene(id SceneId) {
	scene, ok := sceneMap[id]
	if ok {
		scene.Unload()
		delete(sceneMap, id)
	} else {
		panic("scene to unload not found")
	}
}

func UnloadAllScenes() {
	for id, _ := range sceneMap {
		UnloadScene(id)
	}
}

func updateCameraSmooth(camera *rl.Camera2D, pos rl.Vector2, delta float32) {
	minSpeed := 60.0
	minEffectLength := 10
	fractionSpeed := 0.8

	diff := rl.Vector2Subtract(pos, camera.Target)
	length := rl.Vector2Length(diff)

	if length > float32(minEffectLength) {
		speed := float32(math.Max(fractionSpeed*float64(length), minSpeed))
		camera.Target = rl.Vector2Add(camera.Target, rl.Vector2Scale(diff, speed*delta/length))
	}
}

func updateCameraCenter(camera *rl.Camera2D, pos rl.Vector2) {
	camera.Target.X = pos.X
	camera.Target.Y = pos.Y
}
