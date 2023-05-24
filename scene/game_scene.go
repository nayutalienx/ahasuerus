package scene

import (
	"ahasuerus/container"
	"ahasuerus/models"
	"ahasuerus/repository"
	"fmt"
	"strings"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

type GameScene struct {
	worldContainer       *container.ObjectResourceContainer
	environmentContainer *container.ObjectResourceContainer
	camera               *rl.Camera2D
	player               *models.Player

	sceneName        string
	paused           bool
	editMode         bool
	editModeShowMenu bool
	cameraEditPos    rl.Vector2
	editCameraSpeed  float32
	editLabel        models.Object
	selectedItem     []models.EditorSelectedItem

	editMenuBgImageDropMode bool
	editHideGameObjectsMode        bool
}

func NewGameScene(sceneName string) *GameScene {
	startScene := GameScene{
		sceneName:            sceneName,
		worldContainer:       container.NewObjectResourceContainer(),
		environmentContainer: container.NewObjectResourceContainer(),
		cameraEditPos:        rl.NewVector2(0, 0),
		editCameraSpeed:      5,
		selectedItem:         make([]models.EditorSelectedItem, 0),
	}

	beziers := repository.GetAllBeziers(sceneName)

	lines := repository.GetAllLines(sceneName)

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

	rectangles := repository.GetAllRectangles(sceneName)

	for i, _ := range rectangles {
		rect := rectangles[i]
		startScene.worldContainer.AddObject(&rect)
		startScene.player.AddCollisionBox(&rect)
	}

	startScene.worldContainer.AddObjectResource(startScene.player)

	// startScene.environmentContainer.AddObjectResource(
	// 	models.NewImage("resources/bg/1.jpg", 0, 0).AfterLoadPreset(func(i *models.Image) {
	// 		i.Texture.Width = int32(WIDTH)
	// 		i.Texture.Height = int32(HEIGHT)
	// 	}),
	// 	models.NewImage("resources/heroes/girl1.png", 0, 0).
	// 		Scale(1.3).
	// 		AfterLoadPreset(func(girl *models.Image) {
	// 			girl.Pos.X = WIDTH - WIDTH/12 - float32(girl.Texture.Width)
	// 			girl.Pos.Y = HEIGHT - float32(girl.Texture.Height)
	// 		}),
	// 	models.NewMusicStream("resources/music/theme.mp3").SetVolume(0.2))

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

func (s *GameScene) Run() models.Scene {

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
			s.enableEditMode()
		}
		if rl.IsKeyDown(rl.KeyF2) && s.editMode {
			s.disableEditMode()
		}

		if s.editMode {
			s.processEditorMode()
		} else {
			updateCameraSmooth(s.camera, s.player.Pos, delta)
		}

		s.environmentContainer.Update(delta)
		s.environmentContainer.Draw()

		if s.editModeShowMenu {
			s.processEditorMenuMode()
		}

		rl.BeginMode2D(*s.camera)
		if !s.editHideGameObjectsMode {
			s.worldContainer.Update(delta)
			s.worldContainer.Draw()
			if s.editMode && !s.editModeShowMenu {
				s.resolveEditorSelection()
			}
			if s.editMode && s.editModeShowMenu {
				s.processEditorSelection()
			}
		}
		rl.EndMode2D()

		rl.EndDrawing()
	}

	s.pause()

	return GetScene(Menu)
}

func (m *GameScene) Unload() {
	m.environmentContainer.Unload()
	m.worldContainer.Unload()
}

func (s *GameScene) resolveEditorSelection() {
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

func (s *GameScene) processEditorSelection() {
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

func (s GameScene) hasAnySelectedEditorItem() (bool, models.EditorItem) {
	for i, _ := range s.selectedItem {
		ei := s.selectedItem[i]
		if ei.Selected {
			return true, ei.Item
		}
	}
	return false, nil
}

func (s *GameScene) saveEditor() {
	s.worldContainer.ForEachObject(func(obj models.Object) {
		editorItem, ok := obj.(models.EditorItem)
		if ok {
			rect, ok := editorItem.(*models.Rectangle)
			if ok {
				repository.SaveRectangle(s.sceneName, rect)
			}

			bez, ok := editorItem.(*models.Bezier)
			if ok {
				repository.SaveBezier(s.sceneName, bez)
			}

			line, ok := editorItem.(*models.Line)
			if ok {
				repository.SaveLine(s.sceneName, line)
			}
		}
	})
}

func (s *GameScene) disableEditMode() {
	s.editMode = false
	s.player.Resume()
	s.environmentContainer.RemoveObject(s.editLabel)
}

func (s *GameScene) enableEditMode() {
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

func (s *GameScene) processEditorMenuMode() {
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

		rect, isRect := editorItem.(*models.Rectangle)
		if isRect {
			changePosition := rg.Button(rl.NewRectangle(10, 110, 200, 100), "CHANGE POSITION")
			changeSize := rg.Button(rl.NewRectangle(10, 220, 200, 100), "CHANGE SIZE")

			if changePosition {
				rect.SetEditorMoveModeTrue()
				rl.DisableCursor()
				rl.SetMousePosition(int(rect.GetPos().X), int(rect.GetPos().Y))
			}

			if changeSize {
				rect.SetEditorSizeModeTrue()
				rl.DisableCursor()
				rl.SetMousePosition(int(rect.GetPos().X+rect.GetBox().X), int(rect.GetPos().Y+rect.GetBox().Y))
			}
		}

	} else {
		
		buttonWidth := 200
		buttonHeight := 50
		startMenuPosY := 110

		newRectangle := rg.Button(rl.NewRectangle(10, float32(startMenuPosY+buttonHeight), float32(buttonWidth), float32(buttonHeight)), "NEW RECTANGLE")
		newLine := rg.Button(rl.NewRectangle(10, float32(startMenuPosY+buttonHeight*2), float32(buttonWidth), float32(buttonHeight)), "NEW LINE")
		newBezier := rg.Button(rl.NewRectangle(10, float32(startMenuPosY+buttonHeight*3), float32(buttonWidth), float32(buttonHeight)), "NEW BEZIER")
		newBgImage := rg.Button(rl.NewRectangle(10, float32(startMenuPosY+buttonHeight*4), float32(buttonWidth), float32(buttonHeight)), "NEW BG IMAGE")
		
		toggleHideGameObjectsText := "HIDE GAME OBJECTS"
		if s.editHideGameObjectsMode {
			toggleHideGameObjectsText = "SHOW GAME OBJECTS"
		}
		toggleHideGameObjects := rg.Button(rl.NewRectangle(10, float32(startMenuPosY+buttonHeight*5), float32(buttonWidth*2), float32(buttonHeight)), toggleHideGameObjectsText)

		if s.editMenuBgImageDropMode {

			rl.DrawText("DROP IMAGE or BACKSPACE TO LEAVE", int32(WIDTH)/2, int32(HEIGHT)/2, 60, rl.Red)

			if rl.IsKeyDown(rl.KeyBackspace) {
				s.editMenuBgImageDropMode = false
			}

			if rl.IsFileDropped() {
				files := rl.LoadDroppedFiles()

				path := "resources" + strings.Split(files[0], "resources")[1]

				image := models.NewImage(path, 0, 0).
					AfterLoadPreset(func(girl *models.Image) {
						girl.Pos.X = WIDTH / 2
						girl.Pos.Y = HEIGHT / 2
					})

				image.Load()

				s.environmentContainer.AddObjectResource(
					image,
				)

				s.editMenuBgImageDropMode = false
			}

		}

		if newRectangle {
			rect := models.NewRectangle(uuid.NewString(), s.camera.Target.X, s.camera.Target.Y, 200, 100, rl.Blue)
			s.worldContainer.AddObject(rect)
			s.player.AddCollisionBox(rect)
		}

		if newLine {
			line := models.NewLine(uuid.NewString(), rl.NewVector2(s.camera.Target.X, s.camera.Target.Y), rl.NewVector2(s.camera.Target.X+100, s.camera.Target.Y+100), 10, rl.Gold)
			s.worldContainer.AddObject(line)
			s.player.AddCollisionLine(line)
		}

		if newBezier {
			bez := models.NewBezier(uuid.NewString(), rl.NewVector2(s.camera.Target.X, s.camera.Target.Y), rl.NewVector2(s.camera.Target.X+100, s.camera.Target.Y+100), 10, rl.Gold)
			s.worldContainer.AddObject(bez)
			s.player.AddCollisionBezier(bez)
		}

		if newBgImage {
			s.editMenuBgImageDropMode = true
		}

		if toggleHideGameObjects {
			s.editHideGameObjectsMode = !s.editHideGameObjectsMode
		}
	}
}

func (s *GameScene) processEditorMode() {

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
}

func (s *GameScene) pause() {
	s.worldContainer.Pause()
	s.environmentContainer.Pause()
	s.paused = true
}

func (s *GameScene) resume() {
	s.worldContainer.Resume()
	s.environmentContainer.Resume()
	s.paused = false
}
