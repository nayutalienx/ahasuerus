package scene

import (
	"ahasuerus/container"
	"ahasuerus/models"
	"ahasuerus/repository"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const SCENE_COLLECTION = "start-scene"

type StartScene struct {
	worldContainer       *container.ObjectResourceContainer
	environmentContainer *container.ObjectResourceContainer
	camera               *rl.Camera2D
	player               *models.Player

	paused          bool
	editMode        bool
	cameraEditPos   rl.Vector2
	editCameraSpeed float32
	editLabel       models.Object
}

func NewStartScene() *StartScene {
	startScene := StartScene{
		worldContainer:       container.NewObjectResourceContainer(),
		environmentContainer: container.NewObjectResourceContainer(),
		cameraEditPos:        rl.NewVector2(0, 0),
		editCameraSpeed:      5,
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

	rectangles := repository.GetAllRectangles(SCENE_COLLECTION)

	for i, _ := range rectangles {
		rect := rectangles[i]
		startScene.worldContainer.AddObject(&rect)
		startScene.player.AddCollisionBox(&rect)
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

		if rl.IsKeyDown(rl.KeyF1) && !s.editMode {
			s.editMode = true
			s.cameraEditPos.X = s.player.Pos.X
			s.cameraEditPos.Y = s.player.Pos.Y
			s.player.Pause()

			rl.SetMousePosition(int(s.camera.Target.X), int(HEIGHT)/2)

			s.editLabel = models.NewText(10, 50).
				SetFontSize(40).
				SetColor(rl.Red).
				SetUpdateCallback(func(t *models.Text) {
					t.SetData(fmt.Sprintf("edit mode, camera speed %.1f", s.editCameraSpeed))
				})

			s.environmentContainer.AddObject(
				s.editLabel,
			)

		}
		if rl.IsKeyDown(rl.KeyF2) && s.editMode {
			s.editMode = false
			s.player.Resume()
			s.environmentContainer.RemoveObject(s.editLabel)
		}

		if s.editMode {

			mousePos := rl.GetMousePosition()

			if rl.IsKeyDown(rl.KeyRight) {
				s.cameraEditPos.X += s.editCameraSpeed
				mousePos.X += s.editCameraSpeed
			}

			if rl.IsKeyDown(rl.KeyLeft) {
				s.cameraEditPos.X -= s.editCameraSpeed
				mousePos.X -= s.editCameraSpeed
			}

			rl.SetMousePosition(int(mousePos.X), int(mousePos.Y))

			if rl.IsKeyDown(rl.KeyEqual) {
				s.editCameraSpeed++
			}

			if rl.IsKeyDown(rl.KeyMinus) {
				s.editCameraSpeed--
			}

			updateCameraCenter(s.camera, s.cameraEditPos)
		} else {
			updateCameraSmooth(s.camera, s.player.Pos, delta)
		}

		s.environmentContainer.Update(delta)
		s.environmentContainer.Draw()
		rl.BeginMode2D(*s.camera)
		s.worldContainer.Update(delta)
		s.worldContainer.Draw()
		if s.editMode {
			s.updateEditor()
		}
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

func (s *StartScene) updateEditor() {
	mouse := rl.GetMousePosition()
	rl.DrawCircle(int32(mouse.X), int32(mouse.Y), 10, rl.Red)
	s.worldContainer.ForEachObject(func(obj models.Object) {
		editorItem, ok := obj.(models.EditorItem)
		if ok {
			editorItem.ReactOnCollision()
		}
	})
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
