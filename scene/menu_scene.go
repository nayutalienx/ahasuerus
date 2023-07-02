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

type MenuButton int

const (
	StartButton MenuButton = iota
	ExitButton
)

type MenuScene struct {
	menuContainer   *container.ObjectResourceContainer
	menuShouldClose bool
	nextScene       SceneId

	currentButton MenuButton

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
		rl.ClearBackground(rl.Black)

		delta := rl.GetFrameTime()
		m.menuContainer.Update(delta)
		m.menuContainer.Draw()

		m.updateCurrentButton()

		c := models.NewCounter()
		m.drawButton("Start", StartButton, &c)
		m.drawButton("Exit", ExitButton, &c)

		m.processMenuEnter()
		rl.EndDrawing()
	}

	m.menuShouldClose = false

	m.pause()

	if m.nextScene == Close {
		return nil
	}

	return GetScene(m.nextScene)
}

func (m *MenuScene) drawButton(text string, button MenuButton, c *models.Counter) {
	color := rl.White
	if m.currentButton == button {
		color = rl.Orange
	}
	models.DrawSdfText(text, rl.NewVector2(WIDTH/2-200, HEIGHT/10*float32(c.GetAndIncrement())), 100, color)
}

func (m *MenuScene) updateCurrentButton() {

	if rl.IsKeyReleased(rl.KeyDown) {
		m.currentButton++
	}

	if rl.IsKeyReleased(rl.KeyUp) {
		m.currentButton--
	}

	if m.currentButton < StartButton {
		m.currentButton = StartButton
	}

	if m.currentButton > ExitButton {
		m.currentButton = ExitButton
	}

}

func (m *MenuScene) processMenuEnter() {
	m.menuShouldClose = rl.WindowShouldClose()
	if m.menuShouldClose {
		m.nextScene = Close
	}

	if rl.IsKeyReleased(rl.KeyEnter) {
		if m.currentButton == StartButton {
			if lastScene == Menu {
				lastScene = Start
			}
			m.menuShouldClose = true
			m.nextScene = lastScene
		}

		if m.currentButton == ExitButton {
			m.menuShouldClose = true
			m.nextScene = Close
		}
	}
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
