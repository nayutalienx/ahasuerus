package scene

import (
	"ahasuerus/container"
	"ahasuerus/models"
	"ahasuerus/repository"
	"fmt"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const SCENE_COLLECTION = "start-scene"

type StartScene struct {
	worldContainer       *container.ObjectResourceContainer
	environmentContainer *container.ObjectResourceContainer
	camera               *rl.Camera2D
	player               *models.Player

	paused           bool
	editMode         bool
	editModeShowMenu bool
	cameraEditPos    rl.Vector2
	editCameraSpeed  float32
	editLabel        models.Object
	selectedItem     []models.EditorSelectedItem
}

func NewStartScene() *StartScene {
	startScene := StartScene{
		worldContainer:       container.NewObjectResourceContainer(),
		environmentContainer: container.NewObjectResourceContainer(),
		cameraEditPos:        rl.NewVector2(0, 0),
		editCameraSpeed:      5,
		selectedItem:         make([]models.EditorSelectedItem, 0),
	}

	beziers := repository.GetAllBeziers(SCENE_COLLECTION)

	lines := []models.Line{
		*models.NewLine(rl.NewVector2(400, 300), rl.NewVector2(800, 350), 10),
		*models.NewLine(rl.NewVector2(900, 400), rl.NewVector2(1400, 150), 10),
	}

	startScene.player = models.NewPlayer(100, 100)

	for i, _ := range beziers {
		bz := beziers[i]
		startScene.worldContainer.AddObject(&bz)
		startScene.player.AddCollisionBezier(&bz)
	}

	for i, _ := range lines {
		l := lines[i]
		startScene.worldContainer.AddObject(&l)
		startScene.player.AddCollisionLine(&l)
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
				t.SetData(fmt.Sprintf("fps: %d [movement(arrow keys), jump(space), edit mode(F1)]", rl.GetFPS()))
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

	rg.SetStyle(rg.DEFAULT, rg.TEXT_SIZE, 20)

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
					t.SetData(fmt.Sprintf("edit mode[movement(arrow keys), cam.speed(+,-,%.1f), save(P), menu(M), off menu(N), exit(F2)]", s.editCameraSpeed))
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
				if !s.editModeShowMenu {
					mousePos.X += s.editCameraSpeed
				}
			}

			if rl.IsKeyDown(rl.KeyLeft) {
				s.cameraEditPos.X -= s.editCameraSpeed
				if !s.editModeShowMenu {
					mousePos.X -= s.editCameraSpeed
				}
			}

			rl.SetMousePosition(int(mousePos.X), int(mousePos.Y))

			if rl.IsKeyDown(rl.KeyEqual) {
				s.editCameraSpeed++
			}

			if rl.IsKeyDown(rl.KeyMinus) {
				s.editCameraSpeed--
			}

			if rl.IsKeyDown(rl.KeyP) {
				s.saveEditor()
			}

			hasAnySelected, _ := s.hasAnySelectedEditorItem()

			if (rl.IsKeyDown(rl.KeyM) || hasAnySelected) && !s.editModeShowMenu {
				s.editModeShowMenu = true
				rl.EnableCursor()
				rl.SetMousePosition(int(WIDTH)/2, int(HEIGHT)/2)
			}

			if rl.IsKeyDown(rl.KeyN) && s.editModeShowMenu {
				s.editModeShowMenu = false
				rl.DisableCursor()
				rl.SetMousePosition(int(s.cameraEditPos.X), int(s.cameraEditPos.Y))
			}

			updateCameraCenter(s.camera, s.cameraEditPos)
		} else {
			updateCameraSmooth(s.camera, s.player.Pos, delta)
		}

		s.environmentContainer.Update(delta)
		s.environmentContainer.Draw()

		if s.editModeShowMenu {
			hasAnySelected, editorItem := s.hasAnySelectedEditorItem()
			if hasAnySelected {
				
				bezier, isBezier := editorItem.(*models.Bezier)

				if isBezier {
					changeStart := rg.Button(rl.NewRectangle(10, 110, 200, 100), "CHANGE START")
					changeEnd := rg.Button(rl.NewRectangle(10, 220, 200, 100), "CHANGE END")
					if changeStart || changeEnd {
						if changeStart {
							bezier.SetStartModeTrue()
							rl.DisableCursor()
							rl.SetMousePosition(int(bezier.Start.X-20), int(bezier.Start.Y-20))
						}

						if changeEnd {
							bezier.SetEndModeTrue()
							rl.DisableCursor()
							rl.SetMousePosition(int(bezier.End.X+20), int(bezier.End.Y+20))
						}
					}
				}

				line, isLine := editorItem.(*models.Line)

				if isLine {
					changeStart := rg.Button(rl.NewRectangle(10, 110, 200, 100), "CHANGE START")
					changeEnd := rg.Button(rl.NewRectangle(10, 220, 200, 100), "CHANGE END")
					if changeStart || changeEnd {
						if changeStart {
							line.SetStartModeTrue()
							rl.DisableCursor()
							rl.SetMousePosition(int(line.Start.X-20), int(line.Start.Y-20))
						}

						if changeEnd {
							line.SetEndModeTrue()
							rl.DisableCursor()
							rl.SetMousePosition(int(line.End.X+20), int(line.End.Y+20))
						}
					}
				}

			} else {
				rg.Button(rl.NewRectangle(50, 100, 200, 100), "NEW RECTANGLE")
			}
		}

		rl.BeginMode2D(*s.camera)
		s.worldContainer.Update(delta)
		s.worldContainer.Draw()
		if s.editMode && !s.editModeShowMenu {
			s.resolveEditorSelection()
		}
		if s.editMode && s.editModeShowMenu {
			s.processEditorSelection()
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

func (s *StartScene) resolveEditorSelection() {
	mouse := rl.GetMousePosition()
	rl.DrawCircle(int32(mouse.X), int32(mouse.Y), 10, rl.Red)

	selectedItem := make([]models.EditorSelectedItem, 0)

	s.worldContainer.ForEachObject(func(obj models.Object) {
		editorItem, ok := obj.(models.EditorItem)
		if ok {
			selected := editorItem.EditorResolveSelect()
			if selected {
				selectedItem = append(selectedItem, models.EditorSelectedItem{
					Selected: selected,
					Item:     editorItem,
				})
			}
		}
	})

	s.selectedItem = selectedItem
}

func (s *StartScene) processEditorSelection() {
	for i, _ := range s.selectedItem {
		ei := s.selectedItem[i]
		if ei.Selected {
			finishedProcessSelection := ei.Item.ProcessEditorSelection()
			if finishedProcessSelection {
				s.editModeShowMenu = false
				s.selectedItem[i].Selected = false
			}
		}
	}	
}

func (s StartScene) hasAnySelectedEditorItem() (bool, models.EditorItem) {
	for i, _ := range s.selectedItem {
		ei := s.selectedItem[i]
		if ei.Selected {
			return true, ei.Item
		}
	}
	return false, nil
}

func (s *StartScene) saveEditor() {
	s.worldContainer.ForEachObject(func(obj models.Object) {
		editorItem, ok := obj.(models.EditorItem)
		if ok {
			rect, ok := editorItem.(*models.Rectangle)
			if ok {
				repository.SaveRectangle(SCENE_COLLECTION, rect)
			}

			bez, ok := editorItem.(*models.Bezier)
			if ok {
				repository.SaveBezier(SCENE_COLLECTION, bez)
			}
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
