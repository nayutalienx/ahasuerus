package scene

import (
	"ahasuerus/collision"
	"ahasuerus/container"
	"ahasuerus/models"
	"ahasuerus/repository"
	"fmt"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	editorStartMenuPosY    = 110
	editorMenuButtonWidth  = 200
	editorMenuButtonHeight = 50
	envContainer           = "env"
	worldContainer         = "world"
)

type GameScene struct {
	worldContainer *container.ObjectResourceContainer
	camera         *rl.Camera2D
	player         *models.Player

	sceneName string
	paused    bool
}

func NewGameScene(sceneName string) *GameScene {
	scene := GameScene{
		sceneName:      sceneName,
		worldContainer: container.NewObjectResourceContainer(),
	}

	worldImages := repository.GetAllImages(scene.sceneName, worldContainer)
	for i, _ := range worldImages {
		img := worldImages[i]
		scene.worldContainer.AddObjectResource(&img)
	}

	scene.player = models.NewPlayer(100, 100)

	scene.worldContainer.AddObjectResource(scene.player)

	polygonFirst := &models.Polygon{
		Points: [3]rl.Vector2{
			{0, 400}, {2000, 400}, {2000, 600},
		},
		Color: rl.Blue,
	}
	polygonSecond := &models.Polygon{
		Points: [3]rl.Vector2{
			{0, 400}, {0, 600}, {2000, 600},
		},
		Color: rl.Blue,
	}

	scene.worldContainer.AddObject(polygonFirst)
	scene.worldContainer.AddObject(polygonSecond)


	// lightPoint1 := models.NewLightPoint(rl.NewVector2(200, 200)).Dynamic(rl.NewVector2(200, 200), rl.NewVector2(7000, 200), 10)
	// scene.worldContainer.AddObject(lightPoint1)

	// lightPoint2 := models.NewLightPoint(rl.NewVector2(3000, 200)).Dynamic(rl.NewVector2(200, 200), rl.NewVector2(7000, 200), 10)
	// scene.worldContainer.AddObject(lightPoint2)

	// scene.player.AddLightPoint(lightPoint1)
	// scene.player.AddLightPoint(lightPoint2)

	scene.worldContainer.Load()

	scene.player.CollisionProcessor.AddHitbox(collision.Hitbox{
		Polygons: []collision.Polygon{
			{
				Points: polygonFirst.Points,
			},
			{
				Points: polygonSecond.Points,
			},
		},
	})

	camera := rl.NewCamera2D(
		rl.NewVector2(WIDTH/2, HEIGHT-500),
		rl.NewVector2(0, 0),
		0, 1)
	scene.camera = &camera

	return &scene
}

func (s *GameScene) Run() models.Scene {

	rg.SetStyle(rg.DEFAULT, rg.TEXT_SIZE, 20)

	if s.paused {
		s.resume()
	}

	nextScene := Menu

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		delta := rl.GetFrameTime()
		s.camera.Zoom += rl.GetMouseWheelMove() * 0.05

		if rl.IsKeyDown(rl.KeyF1) { // jump to editor scene
			nextScene = Editor
			s.Unload()
			break
		}

		updateCameraSmooth(s.camera, s.player.Pos, delta)

		rl.BeginMode2D(*s.camera)
		s.worldContainer.Update(delta)
		s.worldContainer.Draw()
		rl.EndMode2D()

		models.NewText(10, 10).
			SetFontSize(40).
			SetColor(rl.White).
			SetData(fmt.Sprintf("fps: %d [movement(arrow keys), jump(space), edit mode(F1)]", rl.GetFPS())).
			Draw()

		rl.EndDrawing()
	}

	s.pause()

	return GetScene(nextScene)
}

func (m *GameScene) Unload() {
	m.worldContainer.Unload()
}

func (s *GameScene) pause() {
	s.worldContainer.Pause()
	s.paused = true
}

func (s *GameScene) resume() {
	s.worldContainer.Resume()
	s.paused = false
}
