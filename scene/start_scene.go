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
	
	musicTheme := rl.LoadMusicStream("resources/music/theme.mp3")
	rl.PlayMusicStream(musicTheme)
	worldObjectContainer := container.NewObjectContainer()
	beziers := []models.Bezier{
		*models.NewBezier(rl.NewVector2(0,150), rl.NewVector2(300,400), 20.0),
		*models.NewBezier(rl.NewVector2(WIDTH,600), rl.NewVector2(WIDTH+100,800), 20.0),
		*models.NewBezier(rl.NewVector2(WIDTH+400,800), rl.NewVector2(2*WIDTH+100,600), 20.0),
	}

	player := models.NewPlayer(100, 100)

	for i, _ := range beziers {
		bz := beziers[i]
		worldObjectContainer.AddObject(&bz)
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
		worldObjectContainer.AddObject(&rect)	
		player.AddCollisionBox(rect)
	}

	worldObjectContainer.AddObject(player)

	hudObjectContainer := container.NewObjectContainer()
	hudObjectContainer.AddObject(models.NewText(10, 10).
		SetFontSize(40).
		SetColor(rl.White).
		SetUpdateCallback(func(t *models.Text) {
			t.SetData(fmt.Sprintf("fps: %d ", rl.GetFPS()))
		}))

	bg := rl.LoadImage("resources/bg/1.jpg")
	girl := rl.LoadImage("resources/heroes/girl1.png")

	bgTexture :=rl.LoadTextureFromImage(bg)
	girlTexture :=rl.LoadTextureFromImage(girl)

	bgTexture.Width = int32(WIDTH)
	bgTexture.Height = int32(HEIGHT)

	girlTexture.Width = int32(float32(girlTexture.Width)*1.3)
	girlTexture.Height = int32(float32(girlTexture.Height)*1.3)

	rl.UnloadImage(bg)
	rl.UnloadImage(girl)

	camera := rl.NewCamera2D(
		rl.NewVector2(WIDTH/2, HEIGHT/2),
		rl.NewVector2(0, 0),
		0, 1)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		delta := rl.GetFrameTime()
		camera.Zoom += rl.GetMouseWheelMove() * 0.05

		rl.DrawTexture(bgTexture, 0, 0, rl.White)
		rl.DrawTexture(girlTexture, int32(WIDTH-WIDTH/12)-girl.Width, int32(HEIGHT)-girlTexture.Height, rl.White)

		hudObjectContainer.Update(delta)
		hudObjectContainer.Draw()

		updateCameraSmooth(&camera, player, delta)
		rl.UpdateMusicStream(musicTheme)

		rl.BeginMode2D(camera)
		worldObjectContainer.Update(delta)
		worldObjectContainer.Draw()
		rl.EndMode2D()

		rl.EndDrawing()
	}
	rl.UnloadTexture(bgTexture)
	rl.UnloadTexture(girlTexture)
	rl.UnloadMusicStream(musicTheme)

	return NewMenuScene()
}
