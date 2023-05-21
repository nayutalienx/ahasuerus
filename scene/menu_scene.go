package scene

import (
	"ahasuerus/models"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type MenuScene struct {
	menuShouldClose bool
	nextScene models.Scene
}

func NewMenuScene() *MenuScene {
	rg.SetStyle(rg.DEFAULT, rg.TEXT_SIZE, 70)
	return &MenuScene{}
}

func (m *MenuScene) Run() models.Scene {

	for !m.menuShouldClose {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		m.menuShouldClose = rl.WindowShouldClose()

		startButton := rg.Button(rl.NewRectangle(WIDTH/2-200, HEIGHT/6, 500, 200), "START")
		if startButton {
			m.menuShouldClose = true
			m.nextScene = GetScene(Start)
		}

		closeButton := rg.Button(rl.NewRectangle(WIDTH/2-200, HEIGHT/3, 500, 200), "CLOSE")
		if closeButton {
			m.menuShouldClose = true
			m.nextScene = nil
		}

		if m.nextScene == nil {
			rl.DrawText("next scene nil", 100, 100, 50, rl.Blue)
		} else {
			rl.DrawText("next scene not nil", 100, 100, 50, rl.Blue)
		}

		rl.EndDrawing()
	}

	m.menuShouldClose = false

	return m.nextScene
}

func (m MenuScene) Unload() {

}
