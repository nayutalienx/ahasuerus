package scene

import (
	"ahasuerus/collision"
	"ahasuerus/container"
	"ahasuerus/models"
	"ahasuerus/repository"
	"ahasuerus/resources"
	"fmt"
	"math"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameScene struct {
	worldContainer *container.ObjectResourceContainer
	camera         *rl.Camera2D
	player         *models.Player

	level repository.Level

	onScreenQueue chan models.Object

	paused bool
}

func NewGameScene(sceneName string) *GameScene {
	scene := GameScene{
		worldContainer: container.NewObjectResourceContainer(),
		onScreenQueue:  make(chan models.Object, 1),
	}

	scene.level = repository.GetLevel(sceneName)

	camera := rl.NewCamera2D(
		rl.NewVector2(WIDTH/2, HEIGHT-250),
		scene.level.CameraPos,
		0, 1.0)
	scene.camera = &camera

	worldImages := scene.level.Images
	for i, _ := range worldImages {
		img := worldImages[i]
		img.Camera(&camera)
		scene.worldContainer.AddObjectResource(&img)
	}

	particles := scene.level.ParticleSources
	for i, _ := range particles {
		particle := particles[i]
		scene.worldContainer.AddObjectResource(&particle)
	}

	scene.player = models.NewPlayer(float32(scene.level.PlayerPos.X), float32(scene.level.PlayerPos.Y)).WithShader(resources.PlayerShader)

	scene.worldContainer.AddObjectResource(scene.player)

	collisionHitboxes := scene.level.CollissionHitboxes
	for i, _ := range collisionHitboxes {
		hb := collisionHitboxes[i]
		scene.worldContainer.AddObjectResource(&hb)
		scene.player.CollisionProcessor.AddHitbox(&collision.Hitbox{
			Polygons: hb.PolygonsWithRotation(),
		})

	}

	lights := scene.level.Lights
	for i, _ := range lights {
		light := lights[i]
		scene.worldContainer.AddObject(&light)
		scene.player.AddLightbox(light)
	}

	characters := scene.level.Characters
	for i, _ := range characters {
		npc := characters[i]
		npc.CollisionProcessor.AddHitbox(scene.player.GetHitbox())
		scene.worldContainer.AddObjectResource(npc.ScreenChan(scene.onScreenQueue))
	}

	scene.worldContainer.Load()

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
			break
		}

		if rl.IsKeyReleased(rl.KeyF2) { // toggle draw collision box
			models.DRAW_MODELS = !models.DRAW_MODELS
		}

		s.updateCamera(delta)

		rl.BeginMode2D(*s.camera)
		s.worldContainer.Update(delta)
		s.worldContainer.Draw()
		rl.EndMode2D()

		if len(s.onScreenQueue) > 0 {
			onScreenObject := <-s.onScreenQueue
			onScreenObject.Draw()
			onScreenObject.Update(delta)
		}

		if models.DRAW_MODELS {
			models.NewText(10, 10).
				SetFontSize(40).
				SetColor(rl.White).
				SetData(fmt.Sprintf("fps: %d [movement(arrow keys), jump(space), edit mode(F1)]", rl.GetFPS())).
				Draw()
		}

		rl.EndDrawing()
	}

	s.pause()

	return GetScene(nextScene)
}

func (s *GameScene) updateCamera(delta float32) {
	if s.player.Pos.X > s.level.CameraStartEndMove.X {
		if s.player.Pos.X < s.level.CameraStartEndMove.Y {

			cameraNewPos := s.player.Pos
			cameraNewPos.Y = s.camera.Target.Y
			distanceToCamera := math.Abs(float64(s.player.Pos.X - s.camera.Target.X))
			if distanceToCamera > models.PLAYER_MOVE_SPEED+models.PLAYER_MOVE_SPEED/3 {
				updateCameraWithMode(s.camera, cameraNewPos, delta, FastSmooth)
			} else {
				updateCameraWithMode(s.camera, cameraNewPos, delta, InstantSmooth)
			}

		} else {
			updateCameraWithMode(s.camera, rl.NewVector2(float32(s.level.CameraStartEndMove.Y)+s.camera.Offset.X, s.camera.Target.Y), delta, FastSmooth)
		}
	} else {
		updateCameraWithMode(s.camera, rl.NewVector2(0, s.camera.Target.Y), delta, FastSmooth)
	}
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
