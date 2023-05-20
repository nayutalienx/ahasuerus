package scene

import (
	"ahasuerus/container"
	"ahasuerus/models"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type StartScene struct {
}

func NewStartScene() StartScene {
	return StartScene{}
}

func (StartScene) Run() models.Scene {

	worldContainer := container.NewObjectResourceContainer()
	beziers := []models.Bezier{
		*models.NewBezier(rl.NewVector2(0, 150), rl.NewVector2(300, 400), 20.0),
		*models.NewBezier(rl.NewVector2(WIDTH, 600), rl.NewVector2(WIDTH+100, 800), 20.0),
		*models.NewBezier(rl.NewVector2(WIDTH+400, 800), rl.NewVector2(2*WIDTH+100, 600), 20.0),
	}

	player := models.NewPlayer(100, 100)

	for i, _ := range beziers {
		bz := beziers[i]
		worldContainer.AddObject(&bz)
		player.AddCollisionBezier(&bz)
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
		worldContainer.AddObject(&rect)
		player.AddCollisionBox(rect)
	}

	worldContainer.AddObjectResource(player)

	environmentContainer := container.NewObjectResourceContainer()

	environmentContainer.AddObjectResource(
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

	environmentContainer.AddObject(
		models.NewText(10, 10).
			SetFontSize(40).
			SetColor(rl.White).
			SetUpdateCallback(func(t *models.Text) {
				t.SetData(fmt.Sprintf("fps: %d ", rl.GetFPS()))
			}))

	environmentContainer.Load()
	worldContainer.Load()

	camera := rl.NewCamera2D(
		rl.NewVector2(WIDTH/2, HEIGHT/2),
		rl.NewVector2(0, 0),
		0, 1)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		delta := rl.GetFrameTime()
		camera.Zoom += rl.GetMouseWheelMove() * 0.05

		updateCameraSmooth(&camera, player, delta)

		environmentContainer.Update(delta)
		environmentContainer.Draw()

		rl.BeginMode2D(camera)
		worldContainer.Update(delta)
		worldContainer.Draw()
		rl.EndMode2D()

		rl.EndDrawing()
	}

	environmentContainer.Unload()
	worldContainer.Unload()

	return NewMenuScene()
}
