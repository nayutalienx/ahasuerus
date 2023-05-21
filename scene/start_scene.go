package scene

import (
	"ahasuerus/container"
	"ahasuerus/models"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type StartScene struct {
	worldContainer *container.ObjectResourceContainer
	environmentContainer *container.ObjectResourceContainer
	camera *rl.Camera2D
	player *models.Player

	paused bool
}

func NewStartScene() *StartScene {
	startScene := StartScene{
		worldContainer: container.NewObjectResourceContainer(),
		environmentContainer: container.NewObjectResourceContainer(),
	}

	beziers := []models.Bezier{
		*models.NewBezier(rl.NewVector2(0, 150), rl.NewVector2(300, 400), 20.0),
		*models.NewBezier(rl.NewVector2(WIDTH, 600), rl.NewVector2(WIDTH+100, 800), 20.0),
		*models.NewBezier(rl.NewVector2(WIDTH+400, 800), rl.NewVector2(2*WIDTH+100, 600), 20.0),
	}

	startScene.player = models.NewPlayer(100, 100)

	for i, _ := range beziers {
		bz := beziers[i]
		startScene.worldContainer.AddObject(&bz)
		startScene.player.AddCollisionBezier(&bz)
	}

	rectangles := []models.Rectangle{
		*models.NewRectangle(100, 350).SetWidth(200).SetHeight(20),
		*models.NewRectangle(0, 600).SetWidth(WIDTH).SetHeight(100),
		*models.NewRectangle(800, 450).SetWidth(300).SetHeight(20),
		*models.NewRectangle(WIDTH+100, 800).SetWidth(WIDTH).SetHeight(100),
		*models.NewRectangle(2*WIDTH+100, 600).SetWidth(WIDTH).SetHeight(100),
	}

	for i, _ := range rectangles {
		rect := rectangles[i]
		startScene.worldContainer.AddObject(&rect)
		startScene.player.AddCollisionBox(rect)
	}

	startScene.worldContainer.AddObjectResource(startScene.player)

	startScene.environmentContainer.AddObjectResource(
		models.NewImage("resources/bg/1.jpg", 0, 0).AfterLoadPreset(func(i *models.Image) {
			i.Texture.Width = int32(WIDTH)
			i.Texture.Height = int32(HEIGHT)
		}),
		models.NewImage("resources/heroes/girl1.png", 0, 0).
			Scale(1.3).
			AfterLoadPreset(func(girl *models.Image) {
				girl.Pos.X = WIDTH - WIDTH/12 - float32(girl.Texture.Width)
				girl.Pos.Y = HEIGHT - float32(girl.Texture.Height)
			}),
		models.NewMusicStream("resources/music/theme.mp3").SetVolume(0.2))

	startScene.environmentContainer.AddObject(
		models.NewText(10, 10).
			SetFontSize(40).
			SetColor(rl.White).
			SetUpdateCallback(func(t *models.Text) {
				t.SetData(fmt.Sprintf("fps: %d ", rl.GetFPS()))
			}))

	startScene.environmentContainer.Load()
	startScene.worldContainer.Load()

	camera := rl.NewCamera2D(
		rl.NewVector2(WIDTH/2, HEIGHT/2),
		rl.NewVector2(0, 0),
		0, 1)
	startScene.camera = &camera

	return &startScene
}

func (s *StartScene) Run() models.Scene {

	if s.paused {
		s.resume()
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		delta := rl.GetFrameTime()
		s.camera.Zoom += rl.GetMouseWheelMove() * 0.05

		updateCameraSmooth(s.camera, s.player, delta)

		s.environmentContainer.Update(delta)
		s.environmentContainer.Draw()

		rl.BeginMode2D(*s.camera)
		s.worldContainer.Update(delta)
		s.worldContainer.Draw()
		rl.EndMode2D()

		rl.EndDrawing()
	}

	s.pause()

	return GetScene(Menu)
}

func (m *StartScene) Unload() {
	m.environmentContainer.Unload()
	m.worldContainer.Unload()	
}

func (s *StartScene) pause() {
	s.worldContainer.Pause()
	s.environmentContainer.Pause()
	s.paused = true
}

func (s *StartScene) resume() {
	s.worldContainer.Resume()
	s.environmentContainer.Resume()
	s.paused = false
}