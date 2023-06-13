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
	properties     map[SceneProp]interface{}

	sceneName string
	paused    bool
}

func NewGameScene(sceneName string) *GameScene {
	scene := GameScene{
		sceneName:      sceneName,
		worldContainer: container.NewObjectResourceContainer(),
	}

	scene.properties = map[SceneProp]interface{}{}
	for k, v := range repository.GetSceneProperties(scene.sceneName) {
		scene.properties[SceneProp(k)] = v
	}

	worldImages := repository.GetAllImages(scene.sceneName)
	for i, _ := range worldImages {
		img := worldImages[i]
		scene.worldContainer.AddObjectResource(&img)
	}

	scene.player = models.NewPlayer(float32(scene.properties[PlayerStartX].(float64)), float32(scene.properties[PlayerStartY].(float64)))
	shaderAsInterface, hasShader := scene.properties[PlayerShader]
	if hasShader {
		shader := resources.GameShader(shaderAsInterface.(string))
		scene.player.WithShader(shader)
	}

	scene.worldContainer.AddObjectResource(scene.player)

	hitboxes := repository.GetAllHitboxes(scene.sceneName)
	for i, _ := range hitboxes {
		hb := hitboxes[i]
		scene.worldContainer.AddObject(&hb)
		if hb.Type == models.Collision {
			scene.player.CollisionProcessor.AddHitbox(&collision.Hitbox{
				Polygons: hb.Polygons(),
			})
		}
		if hb.Type == models.Light {
			scene.player.AddLightbox(hb)
		}
		if hb.Type == models.Npc {
			hb.CollisionProcessor.AddHitbox(scene.player.GetHitbox())
		}

	}

	scene.worldContainer.Load()

	camera := rl.NewCamera2D(
		rl.NewVector2(WIDTH/2, HEIGHT-250),
		rl.NewVector2(0, 0),
		0, 1.0)
	camera.Target.Y = 250
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
	if float64(s.player.Pos.X) > s.properties[StartCameraFollowPos].(float64) {
		if float64(s.player.Pos.X) < s.properties[EndCameraFollowPos].(float64) {

			cameraNewPos := s.player.Pos
			cameraNewPos.Y = s.camera.Target.Y
			distanceToCamera := math.Abs(float64(s.player.Pos.X - s.camera.Target.X))
			if distanceToCamera > models.PLAYER_MOVE_SPEED+models.PLAYER_MOVE_SPEED/3 {
				updateCameraWithMode(s.camera, cameraNewPos, delta, FastSmooth)
			} else {
				updateCameraWithMode(s.camera, cameraNewPos, delta, InstantSmooth)
			}

		} else {
			updateCameraWithMode(s.camera, rl.NewVector2(float32(s.properties[EndCameraFollowPos].(float64))+s.camera.Offset.X, s.camera.Target.Y), delta, FastSmooth)
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
