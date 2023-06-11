package scene

import (
	"ahasuerus/config"
	"ahasuerus/models"
	"math"
	"math/rand"
	"strings"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type SceneProp string

const (
	StartCameraFollowPos SceneProp = "startCameraFollowPos"
	EndCameraFollowPos   SceneProp = "endCameraFollowPos"
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

	for i := 0; i < rand.Intn(5); i++ {
		drawLoadScene(randomFunLoadMessage(), 20, time.Second/5)
	}

	switch id {
	case Menu:
		scene = NewMenuScene()
	case Start:
		scene = NewGameScene(sceneNames[Start])
	case Editor:
		UnloadScene(Start)
		scene = NewEditScene(sceneNames[lastScene], lastScene)
	}

	for i := 0; i < rand.Intn(5); i++ {
		drawLoadScene(randomFunLoadMessage(), 20, time.Second/5)
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

func randomFunLoadMessage() string {
	messages := []string{
		"Preparing for adventure",
		"Loading dreams and imagination",
		"Embarking on a journey of epic proportions",
		"Brace yourself for a thrilling experience",
		"Unleashing creativity and excitement",
		"Loading pixels and magic",
		"Getting ready to conquer new worlds",
		"Buckle up and get ready for action",
		"Venturing into the unknown",
		"Loading happiness and excitement",
		"Prepare for an epic adventure",
		"Loading dreams and adventures",
		"Unleash your inner hero",
		"Loading pixels and magic",
		"Embark on a journey of a lifetime",
		"Loading... Brace yourself for excitement",
		"Welcome to a world of wonders",
		"Get ready to explore new realms",
		"Loading... Embrace the challenge ahead",
		"Adventure awaits just beyond this screen",
	}
	return messages[rand.Intn(len(messages))]
}

func drawLoadScene(message string, dots int, dur time.Duration) {
	rl.EndDrawing()
	timeOffsetNanos := float64(dur.Nanoseconds()) / float64(dots)
	col := rl.NewColor(
		uint8(rand.Intn(255)),
		uint8(rand.Intn(255)),
		uint8(rand.Intn(255)),
		255,
	)
	for i := 1; i <= dots; i++ {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawText(message+strings.Repeat(".", i), int32(WIDTH)/3, int32(HEIGHT)/2, 50, col)
		rl.EndDrawing()
		time.Sleep(time.Duration(timeOffsetNanos))
	}
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

type CameraUpdateMode int

const (
	FastSmooth CameraUpdateMode = iota
	InstantSmooth
)

func updateCameraWithMode(
	camera *rl.Camera2D,
	pos rl.Vector2,
	delta float32,
	mode CameraUpdateMode) {

	minSpeed := 0.0
	minEffectLength := 10
	fractionSpeed := 0.0

	if mode == FastSmooth {
		minSpeed = 30.0
		fractionSpeed = 5.0
	}

	if mode == InstantSmooth {
		minSpeed = 100.0
		fractionSpeed = 10.0
	}

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
