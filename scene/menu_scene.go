package scene

import (
	"ahasuerus/models"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type MenuScene struct {
}

func NewMenuScene() models.Scene {
	return MenuScene{}
}

func (m MenuScene) Run() models.Scene {
	
	menuShouldClose := false
	var nextScene models.Scene

	rg.SetStyle(rg.DEFAULT, rg.TEXT_SIZE, 70)

	for !menuShouldClose {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		menuShouldClose = rl.WindowShouldClose()

		buttonClicked := rg.Button(rl.NewRectangle(WIDTH/2-200, HEIGHT/2, 500, 200), "START")
		if buttonClicked {
			menuShouldClose = true
			nextScene = NewStartScene()
		}

		rl.EndDrawing()
	}

	return nextScene
}
