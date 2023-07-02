package scene

import (
	"ahasuerus/config"
	"ahasuerus/models"
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type SceneId string

const (
	Undefined SceneId = ""
	Menu SceneId = "menu"
	Start SceneId = "start"
	Editor SceneId = "editor"
	Close SceneId = "close"
	Level1 SceneId = "level1"
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

	drawLoadScene()

	switch id {
	case Menu:
		scene = NewMenuScene()
	case Editor:
		UnloadScene(lastScene)
		scene = NewEditScene(string(lastScene), lastScene)
	default:
		UnloadScene(lastScene)
		scene = NewGameScene(string(id))
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

func drawLoadScene() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.EndDrawing()
}

func UnloadScene(id SceneId) {
	scene, ok := sceneMap[id]
	if ok {
		scene.Unload()
		delete(sceneMap, id)
	} else {
		for i := 0; i < 10; i++ {
			fmt.Println("scene to unload not found")
		}
	}
}

func UnloadAllScenes() {
	for id, _ := range sceneMap {
		UnloadScene(id)
	}
}

func updateCameraWithMode(
	camera *rl.Camera2D,
	pos rl.Vector2,
	delta float32) {

	minSpeed := 0.0
	minEffectLength := 10
	fractionSpeed := 0.0

	minSpeed = 30.0
	fractionSpeed = 5.0

	diff := rl.Vector2Subtract(pos, camera.Target)
	length := rl.Vector2Length(diff)

	if length > float32(minEffectLength) {
		speed := float32(math.Max(fractionSpeed*float64(length), minSpeed))
		camera.Target = rl.Vector2Add(camera.Target, rl.Vector2Scale(diff, speed*delta/length))
	}
}

func updateCameraCenter(camera *rl.Camera2D, pos rl.Vector2, delta float32) {
	camera.Target.X = pos.X
	camera.Target.Y = pos.Y
}
