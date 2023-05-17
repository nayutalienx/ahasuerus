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
	WIDTH       = 1280
	HEIGHT      = 720
	APPLICATION = "ahasuerus"
)

func Start() {
	rl.InitWindow(WIDTH, HEIGHT, APPLICATION)
	defer rl.CloseWindow()
	rl.SetTargetFPS(FPS)

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

	camera := rl.NewCamera2D(
		rl.NewVector2(WIDTH/2, HEIGHT/2),
		rl.NewVector2(0, 0),
		0, 1)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		delta := rl.GetFrameTime()
		camera.Zoom += rl.GetMouseWheelMove() * 0.05

		hudObjectContainer.Update(delta)
		hudObjectContainer.Draw()

		updateCameraSmooth(&camera, player, delta)

		rl.BeginMode2D(camera)
		worldObjectContainer.Update(delta)
		worldObjectContainer.Draw()
		rl.EndMode2D()

		rl.EndDrawing()
	}
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