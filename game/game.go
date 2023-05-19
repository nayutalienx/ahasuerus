package game

import (
	"ahasuerus/container"
	"ahasuerus/models"
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	FPS         = 60
	WIDTH       = 2560
	HEIGHT      = 1440
	APPLICATION = "ahasuerus"
)

func Start() {
	rl.InitWindow(WIDTH, HEIGHT, APPLICATION)
	defer rl.CloseWindow()
	rl.SetTargetFPS(FPS)

	rl.InitAudioDevice()
	musicTheme := rl.LoadMusicStream("resources/music/theme.mp3")

	rl.PlayMusicStream(musicTheme)

	worldObjectContainer := container.NewObjectContainer()

	rl.SetConfigFlags(rl.FlagMsaa4xHint)

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

	bgTexture.Width = WIDTH
	bgTexture.Height = HEIGHT

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
		rl.DrawTexture(girlTexture, WIDTH-WIDTH/12-girl.Width, HEIGHT-girlTexture.Height, rl.White)

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

	rl.UnloadMusicStream(musicTheme)
	rl.CloseAudioDevice()

	rl.UnloadTexture(bgTexture)
	rl.UnloadTexture(girlTexture)
}

func updateCameraSmooth(camera *rl.Camera2D, player *models.Player, delta float32) {
	minSpeed := 60.0
	minEffectLength := 10
	fractionSpeed := 0.8

	diff := rl.Vector2Subtract(player.Pos, camera.Target)
	length := rl.Vector2Length(diff)

	if length > float32(minEffectLength) {
		speed := float32(math.Max(fractionSpeed*float64(length),minSpeed))
		camera.Target = rl.Vector2Add(camera.Target, rl.Vector2Scale(diff, speed*delta/length))
	}
}