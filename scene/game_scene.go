package scene

import (
	"ahasuerus/collision"
	"ahasuerus/container"
	"ahasuerus/models"
	"ahasuerus/repository"
	"ahasuerus/resources"
	"fmt"

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
	size rl.Vector2

	screenScale float32
}

func NewGameScene(sceneName string) *GameScene {
	scene := GameScene{
		worldContainer: container.NewObjectResourceContainer(),
		onScreenQueue:  make(chan models.Object, 2),
	}

	scene.level = repository.GetLevel(sceneName)

	scene.size = scene.level.Size()
	scene.screenScale = HEIGHT/scene.size.Y

	camera := rl.NewCamera2D(
		rl.NewVector2(WIDTH/2, HEIGHT/2),
		rl.NewVector2(WIDTH/2, scene.size.Y/2),
		0, scene.screenScale)

	scene.camera = &camera

	scene.player = models.NewPlayer(float32(scene.level.PlayerPos.X), float32(scene.level.PlayerPos.Y)).WithShader(resources.GameShader(scene.level.PlayerShader))

	if scene.level.MusicTheme != "" {
		scene.worldContainer.AddObjectResource(models.NewMusicStream(scene.level.MusicTheme, scene.level.MusicThemeReverse).SetRewindCollisionCheck(scene.player.IsCollisionRewind))
	}

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

	scene.worldContainer.AddObjectResource(scene.player)

	collisionHitboxes := scene.level.CollissionHitboxes
	for i, _ := range collisionHitboxes {
		hb := collisionHitboxes[i]
		scene.worldContainer.AddObjectResource(&hb)
		scene.player.CollisionProcessor.AddHitbox(&collision.Hitbox{
			Polygons: hb.PolygonsWithRotation(),
			Rotation: hb.Rotation,
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
		scene.worldContainer.AddObjectResource(npc.ScreenChan(scene.onScreenQueue).ScreenScale(scene.screenScale))
	}

	scene.worldContainer.Sort()
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

		for len(s.onScreenQueue) > 0 {
			onScreenObject := <-s.onScreenQueue
			onScreenObject.Draw()
			onScreenObject.Update(delta)
		}
		isWannaChangeScene, sc := models.IsWannaChangeScene()
		if isWannaChangeScene {
			nextScene = SceneId(sc)
			break
		}

		if models.DRAW_MODELS {
			models.NewText(10, 10).
				SetFontSize(40).
				SetColor(rl.White).
				SetData(fmt.Sprintf("fps: %d [movement(arrow keys), jump(space), edit mode(F1)] camera: %.1f %.1f %.1f", rl.GetFPS(), s.camera.Target.X, s.camera.Target.Y, s.camera.Zoom)).
				Draw()
		}

		rl.EndDrawing()
	}

	s.pause()

	return GetScene(nextScene)
}

func (s *GameScene) updateCamera(delta float32) {
	cameraNewPos := s.player.Pos
	cameraNewPos.Y = s.camera.Target.Y

	leftCameraLimit := WIDTH/2 + (WIDTH - (WIDTH*s.camera.Zoom))
	if cameraNewPos.X < leftCameraLimit {
		cameraNewPos.X = leftCameraLimit
	}

	updateCameraWithMode(s.camera, cameraNewPos, delta)
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
