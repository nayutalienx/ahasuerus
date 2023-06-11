package scene

import (
	"ahasuerus/container"
	"ahasuerus/controls"
	"ahasuerus/models"
	_ "ahasuerus/repository"
	"ahasuerus/resources"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type MenuScene struct {
	menuContainer   *container.ObjectResourceContainer
	menuShouldClose bool
	nextScene       models.Scene

	paused bool
}

func NewMenuScene() *MenuScene {
	menuScene := &MenuScene{
		menuContainer: container.NewObjectResourceContainer(),
	}

	menuScene.menuContainer.Load()
	return menuScene
}

func (m *MenuScene) Run() models.Scene {

	rl.DisableCursor()
	rg.LoadStyle(string(resources.Lavanda))
	rg.SetStyle(rg.DEFAULT, rg.TEXT_SIZE, 70)

	controls.SetMousePosition(int(WIDTH)/2, int(HEIGHT)/2, 43)

	if m.paused {
		m.resume()
	}

	for !m.menuShouldClose {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Blank)

		delta := rl.GetFrameTime()
		m.menuContainer.Update(delta)
		m.menuContainer.Draw()

		m.menuShouldClose = rl.WindowShouldClose()

		startButton := rg.Button(rl.NewRectangle(WIDTH/2-200, HEIGHT/6, 500, 200), "New game")
		closeButton := rg.Button(rl.NewRectangle(WIDTH/2-200, HEIGHT/3, 500, 200), "See you next time")
		mouse := rl.GetMousePosition()
		rl.DrawCircle(int32(mouse.X), int32(mouse.Y), 10, rl.Green)

		if startButton {
			m.menuShouldClose = true
			m.nextScene = GetScene(Start)
		}

		if closeButton {
			m.menuShouldClose = true
			m.nextScene = nil
		}

		rl.EndDrawing()
	}

	m.menuShouldClose = false

	m.pause()

	return m.nextScene
}

func (m *MenuScene) Unload() {
	m.menuContainer.Unload()
}

func (s *MenuScene) pause() {
	s.menuContainer.Pause()
	s.paused = true
}

func (s *MenuScene) resume() {
	s.menuContainer.Resume()
	s.paused = false
}
