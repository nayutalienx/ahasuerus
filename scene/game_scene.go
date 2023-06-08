package scene

import (
	"ahasuerus/container"
	"ahasuerus/controls"
	"ahasuerus/models"
	"ahasuerus/repository"
	"ahasuerus/resources"
	"fmt"
	"strings"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

const (
	editorStartMenuPosY    = 110
	editorMenuButtonWidth  = 200
	editorMenuButtonHeight = 50
	envContainer           = "env"
	worldContainer         = "world"
)

type GameScene struct {
	worldContainer       *container.ObjectResourceContainer
	environmentContainer *container.ObjectResourceContainer
	camera               *rl.Camera2D
	player               *models.Player

	sceneName               string
	paused                  bool
	editMode                bool
	editModeShowMenu        bool
	cameraEditPos           rl.Vector2
	editCameraSpeed         float32
	editLabel               models.Object
	selectedGameObjectsItem []models.EditorSelectedItem
	selectedBackgroundItem  []models.EditorSelectedItem

	editMenuBgImageDropMode   bool
	editMenuGameImageDropMode bool
	editHideGameObjectsMode   bool
	editBgImageEditorMode     bool
}

func NewGameScene(sceneName string) *GameScene {
	scene := GameScene{
		sceneName:               sceneName,
		worldContainer:          container.NewObjectResourceContainer(),
		environmentContainer:    container.NewObjectResourceContainer(),
		cameraEditPos:           rl.NewVector2(0, 0),
		editCameraSpeed:         5,
		selectedGameObjectsItem: make([]models.EditorSelectedItem, 0),
	}
	
	worldImages := repository.GetAllImages(scene.sceneName, worldContainer)
	for i, _ := range worldImages {
		img := worldImages[i]
		//img.
			//WithShader(resources.TextureLightShader).
			// AddLightPoint(lightPoint1).
			// AddLightPoint(lightPoint2)
		scene.worldContainer.AddObjectResource(&img)
	}

	beziers := repository.GetAllBeziers(sceneName)

	lines := repository.GetAllLines(sceneName)

	scene.player = models.NewPlayer(100, 100)
	//WithShader(resources.TextureLightShader)

	for i, _ := range beziers {
		bz := beziers[i]
		scene.worldContainer.AddObject(&bz)
		scene.player.AddCollisionBezier(&bz)
	}

	for i, _ := range lines {
		l := lines[i]
		scene.worldContainer.AddObject(&l)
		scene.player.AddCollisionLine(&l)
	}

	rectangles := repository.GetAllRectangles(sceneName)

	for i, _ := range rectangles {
		rect := rectangles[i]
		scene.worldContainer.AddObject(&rect)
		scene.player.AddCollisionBox(&rect)
	}

	scene.worldContainer.AddObjectResource(scene.player)

	// envImages := repository.GetAllImages(scene.sceneName, envContainer)
	// for i, _ := range envImages {
	// 	img := envImages[i]
	// 	scene.environmentContainer.AddObjectResource(&img)
	// }

	// lightPoint1 := models.NewLightPoint(rl.NewVector2(200, 200)).Dynamic(rl.NewVector2(200, 200), rl.NewVector2(7000, 200), 10)
	// scene.worldContainer.AddObject(lightPoint1)

	// lightPoint2 := models.NewLightPoint(rl.NewVector2(3000, 200)).Dynamic(rl.NewVector2(200, 200), rl.NewVector2(7000, 200), 10)
	// scene.worldContainer.AddObject(lightPoint2)

	// scene.player.AddLightPoint(lightPoint1)
	// scene.player.AddLightPoint(lightPoint2)

	// scene.environmentContainer.AddObjectResource(
	// 	models.NewMusicStream("resources/music/theme.mp3").SetVolume(0.2))

	//scene.environmentContainer.AddObjectResource(models.NewMusicStream("resources/music/menu_theme.mp3"))

	scene.environmentContainer.AddObject(
		models.NewText(10, 10).
			SetFontSize(40).
			SetColor(rl.White).
			SetUpdateCallback(func(t *models.Text) {
				t.SetData(fmt.Sprintf("fps: %d [movement(arrow keys), jump(space), edit mode(F1)]", rl.GetFPS()))
			}))

	scene.environmentContainer.Load()
	scene.worldContainer.Load()

	camera := rl.NewCamera2D(
		rl.NewVector2(WIDTH/2, HEIGHT-500),
		rl.NewVector2(0, 0),
		0, 1)
	scene.camera = &camera

	return &scene
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
		if rl.IsKeyDown(rl.KeyF2) && s.editMode && !s.editModeShowMenu {
			s.disableEditMode()
		}

		if s.editMode {
			s.processEditorMode()
		} else {
			updateCameraSmooth(s.camera, s.player.Pos, delta)
		}

		rl.BeginMode2D(*s.camera)
		if !s.editHideGameObjectsMode {
			s.worldContainer.Update(delta)
			s.worldContainer.Draw()
			if s.editMode && !s.editModeShowMenu {
				s.resolveEditorGameObjectsSelection()
			}
			if s.editMode && s.editModeShowMenu {
				s.processEditorGameObjectSelection()
			}
		}
		rl.EndMode2D()

		if s.editModeShowMenu {
			s.processEditorMenuMode()
		}

		if s.editHideGameObjectsMode && s.editBgImageEditorMode {
			hasAnySelection, _ := s.hasAnySelectedBackgroundEditorItem()
			if !hasAnySelection {
				s.resolveEditorBackgroundImageSelection()
			}
			s.processEditorBackgroundSelection()
		}

		s.environmentContainer.Update(delta)
		s.environmentContainer.Draw()

		rl.EndDrawing()
	}

	s.pause()

	return GetScene(Menu)
}

func (m *GameScene) Unload() {
	m.environmentContainer.Unload()
	m.worldContainer.Unload()
}

func (s *GameScene) resolveEditorBackgroundImageSelection() {
	selectedItem := make([]models.EditorSelectedItem, 0)

	s.environmentContainer.ForEachObjectReverseWithPredicate(func(obj models.Object) bool {
		editorItem, ok := obj.(models.EditorItem)
		if ok {
			resolveResult := editorItem.EditorResolveSelect()
			if resolveResult.Selected && resolveResult.Collision {
				selectedItem = append(selectedItem, models.EditorSelectedItem{
					Selected: resolveResult.Selected,
					Item:     editorItem,
				})
			}
			if resolveResult.Collision {
				return true
			}
		}
		return false
	})
	s.selectedBackgroundItem = selectedItem
}

func (s *GameScene) resolveEditorGameObjectsSelection() {
	mouse := rl.GetMousePosition()
	rl.DrawText(fmt.Sprintf("%v %v", mouse.X, mouse.Y), int32(mouse.X), int32(mouse.Y), 40, rl.Red)
	rl.DrawCircle(int32(mouse.X), int32(mouse.Y), 10, rl.Red)

	selectedItem := make([]models.EditorSelectedItem, 0)

	s.worldContainer.ForEachObjectReverseWithPredicate(func(obj models.Object) bool {
		editorItem, ok := obj.(models.EditorItem)
		if ok {
			resolveResult := editorItem.EditorResolveSelect()
			if resolveResult.Selected && resolveResult.Collision {
				selectedItem = append(selectedItem, models.EditorSelectedItem{
					Selected: resolveResult.Selected,
					Item:     editorItem,
				})
			}
			if resolveResult.Collision {
				return true
			}
		}
		return false
	})

	s.selectedGameObjectsItem = selectedItem
}

func (s *GameScene) processEditorGameObjectSelection() {
	for i, _ := range s.selectedGameObjectsItem {
		ei := s.selectedGameObjectsItem[i]
		if ei.Selected {
			processResult := ei.Item.ProcessEditorSelection()
			if processResult.Finished {
				s.editModeShowMenu = false
				s.selectedGameObjectsItem[i].Selected = false
				if processResult.DisableCursor {
					controls.DisableCursor(246)
				}
				if processResult.CursorForcePosition {
					controls.SetMousePosition(processResult.CursorX, processResult.CursorY, 249)
				}
			}
		}
	}
}

func (s *GameScene) processEditorBackgroundSelection() {
	for i, _ := range s.selectedBackgroundItem {
		ei := s.selectedBackgroundItem[i]
		if ei.Selected {
			processResult := ei.Item.ProcessEditorSelection()
			if processResult.Finished {
				s.selectedBackgroundItem[i].Selected = false
			}
		}
	}
}

func (s GameScene) hasAnySelectedGameObjectEditorItem() (bool, models.EditorItem) {
	for i, _ := range s.selectedGameObjectsItem {
		ei := s.selectedGameObjectsItem[i]
		if ei.Selected {
			return true, ei.Item
		}
	}
	return false, nil
}

func (s GameScene) hasAnySelectedBackgroundEditorItem() (bool, models.EditorItem) {
	for i, _ := range s.selectedBackgroundItem {
		ei := s.selectedBackgroundItem[i]
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

			editorItem, ok := obj.(models.EditorItem)
			if ok {
				image, ok := editorItem.(*models.Image)
				if ok {
					repository.SaveImage(s.sceneName, worldContainer, image)
				}
			}
		}
	})
	s.environmentContainer.ForEachObject(func(obj models.Object) {
		editorItem, ok := obj.(models.EditorItem)
		if ok {
			image, ok := editorItem.(*models.Image)
			if ok {
				repository.SaveImage(s.sceneName, envContainer, image)
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
	s.player.Pos.X = 100
	s.player.Pos.Y = 100
	s.editMode = true
	s.cameraEditPos.X = s.player.Pos.X
	s.cameraEditPos.Y = s.player.Pos.Y
	s.player.Pause()

	controls.SetMousePosition(int(s.camera.Target.X), int(HEIGHT)/2, 333)

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
	hasAnySelected, editorItem := s.hasAnySelectedGameObjectEditorItem()
	if hasAnySelected {
		s.reactOnGameObjectEditorSelect(editorItem)
	} else {
		s.drawNonGameFocusedMenu()
	}
}

func (s *GameScene) drawNonGameFocusedMenu() {
	buttonCounter := models.NewCounter()

	newRectangle := false
	newLine := false
	newBezier := false
	newGameImage := false
	newBgImage := false

	toggleHideGameObjectsText := "HIDE GAME OBJECTS"
	if s.editHideGameObjectsMode {
		toggleHideGameObjectsText = "SHOW GAME OBJECTS"
		toggleBgImageEditorText := "ENABLE BG IMAGE EDITOR [PRESS B]"
		if s.editBgImageEditorMode {
			toggleBgImageEditorText = "DISABLE BG IMAGE EDITOR [PRESS V]"
		}

		if rl.IsKeyDown(rl.KeyB) {
			s.editBgImageEditorMode = true
		}
		if rl.IsKeyDown(rl.KeyV) {
			s.editBgImageEditorMode = false
		}

		rl.DrawText(toggleBgImageEditorText, 10, int32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), 30, rl.Red)

		if !s.editBgImageEditorMode {
			newBgImage = rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "NEW BG IMAGE")
		}
	} else {
		newRectangle = rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "NEW RECTANGLE")
		newLine = rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "NEW LINE")
		newBezier = rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "NEW BEZIER")
		newGameImage = rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "NEW IMAGE")
	}
	toggleHideGameObjects := false
	if !s.editBgImageEditorMode {
		toggleHideGameObjects = rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth*2), float32(editorMenuButtonHeight)), toggleHideGameObjectsText)
	}

	toggleModelsDrawText := "HIDE COLLISSION"
	if !models.DRAW_MODELS {
		toggleModelsDrawText = "SHOW COLLISSION"
	}
	toggleCollissionDrawButton := rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth*2), float32(editorMenuButtonHeight)), toggleModelsDrawText)

	if toggleCollissionDrawButton {
		models.DRAW_MODELS = !models.DRAW_MODELS
	}


	if s.editMenuBgImageDropMode || s.editMenuGameImageDropMode {

		rl.DrawText("DROP IMAGE or BACKSPACE TO LEAVE", int32(WIDTH)/2, int32(HEIGHT)/2, 60, rl.Red)

		if rl.IsKeyDown(rl.KeyBackspace) {
			s.editMenuBgImageDropMode = false
			s.editMenuGameImageDropMode = false
		}

		if rl.IsFileDropped() {
			files := rl.LoadDroppedFiles()

			path := "resources" + strings.Split(files[0], "resources")[1]

			image := models.NewImage(s.environmentContainer.Size(), uuid.NewString(), resources.GameTexture(path), 0, 0, 0, 0, 0).
				AfterLoadPreset(func(i *models.Image) {
					if s.editMenuBgImageDropMode {
						i.Pos.X = WIDTH / 2
						i.Pos.Y = HEIGHT / 2
					}
					if s.editMenuGameImageDropMode {
						i.Pos.X = s.camera.Target.X
						i.Pos.Y = s.camera.Target.Y
					}
				})

			image.Load()

			if s.editMenuBgImageDropMode {
				s.environmentContainer.AddObjectResource(
					image,
				)
			} else if s.editMenuGameImageDropMode {
				s.worldContainer.AddObjectResource(
					image,
				)
			}

			s.editMenuBgImageDropMode = false
			s.editMenuGameImageDropMode = false
		}

	}

	if newRectangle {
		blue := rl.Blue
		blue.A = 80
		rect := models.NewRectangle(uuid.NewString(), s.camera.Target.X, s.camera.Target.Y, 200, 100, blue)
		s.worldContainer.AddObject(rect)
		s.player.AddCollisionBox(rect)
	}

	if newLine {
		gold := rl.Gold
		gold.A = 90
		line := models.NewLine(uuid.NewString(), rl.NewVector2(s.camera.Target.X, s.camera.Target.Y), rl.NewVector2(s.camera.Target.X+100, s.camera.Target.Y+100), 10, gold)
		s.worldContainer.AddObject(line)
		s.player.AddCollisionLine(line)
	}

	if newBezier {
		gold := rl.Gold
		gold.A = 90
		bez := models.NewBezier(uuid.NewString(), rl.NewVector2(s.camera.Target.X, s.camera.Target.Y), rl.NewVector2(s.camera.Target.X+100, s.camera.Target.Y+100), 10, gold)
		s.worldContainer.AddObject(bez)
		s.player.AddCollisionBezier(bez)
	}

	if newGameImage {
		s.editMenuGameImageDropMode = true
	}

	if newBgImage {
		s.editMenuBgImageDropMode = true
	}

	if toggleHideGameObjects {
		s.editHideGameObjectsMode = !s.editHideGameObjectsMode
	}

	if s.editHideGameObjectsMode && s.editBgImageEditorMode {
		hasAnySelectedBackgroundItem, backgroundSelectedItem := s.hasAnySelectedBackgroundEditorItem()
		if hasAnySelectedBackgroundItem {

			bgImage, isImage := backgroundSelectedItem.(*models.Image)
			if isImage {
				s.reactOnImageEditorSelection(s.environmentContainer, bgImage, buttonCounter)
			}

		}
	}
}

func (s *GameScene) reactOnImageEditorSelection(container *container.ObjectResourceContainer, image *models.Image, buttonCounter *models.Counter) {

	changeBgPosButton := rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "CHANGE BG POS")
	resizeBgButton := rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "RESIZE IMG")

	moveUpper := rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "MOVE UPPER")
	moveDown := rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "MOVE DOWN")

	replicate := rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "REPLICATE")
	rotateMode := rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "ROTATE MODE")
	
	deleteImage := rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "DELETE")

	shouldDisableCursor := container == s.worldContainer

	if changeBgPosButton {
		image.SetEditorMoveWithCursorTrue()
		if shouldDisableCursor {
			controls.DisableCursor(498)
		}
		controls.SetMousePosition(int(image.Pos.X), int(image.Pos.Y), 500)
	}

	if resizeBgButton {
		image.SetEditorResizeWithCursorTrue()
		if shouldDisableCursor {
			controls.DisableCursor(506)
		}
		controls.SetMousePosition(int(image.Pos.X+image.Box.X), int(image.Pos.Y+image.Box.Y), 508)
	}

	if rotateMode {
		image.SetEditorRotateModeTrue()
	}

	if moveUpper {
		drawIndex := container.MoveUp(image)
		image.DrawIndex = drawIndex
		s.syncDrawIndex(container)
	}

	if moveDown {
		drawIndex := container.MoveDown(image)
		image.DrawIndex = drawIndex
		s.syncDrawIndex(container)
	}

	if replicate {
		imageReplica := image.Replicate(uuid.NewString(), image.Pos.X-100, image.Pos.Y-100)
		imageReplica.Load()
		container.AddObjectResource(imageReplica)
	}

	if deleteImage {
		container.RemoveObject(image)
		repository.DeleteImage(s.sceneName, worldContainer, image)
	}

}

func (s *GameScene) reactOnGameObjectEditorSelect(editorItem models.EditorItem) {

	buttonCounter := models.NewCounter()

	bezier, isBezier := editorItem.(*models.Bezier)

	if isBezier {
		changeStart := rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "CHANGE START")
		changeEnd := rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "CHANGE END")
		delete := rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "DELETE")
		if changeStart || changeEnd {
			if changeStart {
				bezier.SetStartModeTrue()
				controls.DisableCursor(543)
				controls.SetMousePosition(int(bezier.Start.X-20), int(bezier.Start.Y-20), 544)
			}

			if changeEnd {
				bezier.SetEndModeTrue()
				controls.DisableCursor(549)
				controls.SetMousePosition(int(bezier.End.X+20), int(bezier.End.Y+20), 550)
			}
		}
		if delete {
			s.worldContainer.RemoveObject(bezier)
			repository.DeleteBezier(s.sceneName, bezier)
		}
	}

	line, isLine := editorItem.(*models.Line)

	if isLine {
		changeStart := rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "CHANGE START")
		changeEnd := rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "CHANGE END")
		delete := rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "DELETE")
		if changeStart || changeEnd {
			if changeStart {
				line.SetStartModeTrue()
				controls.DisableCursor(563)
				controls.SetMousePosition(int(line.Start.X-20), int(line.Start.Y-20), 564)
			}

			if changeEnd {
				line.SetEndModeTrue()
				controls.DisableCursor(569)
				controls.SetMousePosition(int(line.End.X+20), int(line.End.Y+20), 570)
			}
		}
		if delete {
			s.worldContainer.RemoveObject(line)
			repository.DeleteLine(s.sceneName, line)
		}
	}

	rect, isRect := editorItem.(*models.Rectangle)
	if isRect {
		changePosition := rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "CHANGE POSITION")
		changeSize := rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "CHANGE SIZE")
		delete := rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "DELETE")

		if changePosition {
			rect.SetEditorMoveModeTrue()
			controls.DisableCursor(582)
			controls.SetMousePosition(int(rect.GetPos().X), int(rect.GetPos().Y), 583)
		}

		if changeSize {
			rect.SetEditorSizeModeTrue()
			controls.DisableCursor(588)
			controls.SetMousePosition(int(rect.GetPos().X+rect.GetBox().X), int(rect.GetPos().Y+rect.GetBox().Y), 589)
		}

		if delete {
			s.worldContainer.RemoveObject(rect)
			repository.DeleteRectangle(s.sceneName, rect)
		}
	}

	img, isImg := editorItem.(*models.Image)
	if isImg {
		s.reactOnImageEditorSelection(s.worldContainer, img, buttonCounter)
	}

}

func (s *GameScene) processEditorMode() {

	mousePos := rl.GetMousePosition()

	updateMouse := false

	if rl.IsKeyDown(rl.KeyRight) {
		s.cameraEditPos.X += s.editCameraSpeed
		if !s.editModeShowMenu {
			mousePos.X += s.editCameraSpeed
			updateMouse = true
		}
	}

	if rl.IsKeyDown(rl.KeyLeft) {
		s.cameraEditPos.X -= s.editCameraSpeed
		if !s.editModeShowMenu {
			mousePos.X -= s.editCameraSpeed
			updateMouse = true
		}
	}

	if updateMouse {
		controls.SetMousePosition(int(mousePos.X), int(mousePos.Y), 618)
	}

	if rl.IsKeyDown(rl.KeyEqual) {
		s.editCameraSpeed++
	}

	if rl.IsKeyDown(rl.KeyMinus) {
		s.editCameraSpeed--
	}

	if rl.IsKeyDown(rl.KeyP) {
		s.saveEditor()
		s.environmentContainer.AddObject(
			models.NewText(int32(WIDTH)/2, int32(HEIGHT)/4).
				SetData("DATA SAVED").
				SetFontSize(60).
				SetColor(rl.Red).
				WithExpire(3, func(t *models.Text) {
					s.environmentContainer.RemoveObject(t)
				}),
		)
	}

	hasAnySelected, _ := s.hasAnySelectedGameObjectEditorItem()

	if hasAnySelected {

		if !s.editModeShowMenu {
			s.editModeShowMenu = true
			controls.EnableCursor(653)
		}

	} else {

		if (rl.IsKeyDown(rl.KeyM)) && !s.editModeShowMenu {
			s.editModeShowMenu = true
			controls.EnableCursor(660)
			controls.SetMousePosition(int(WIDTH)/2, int(HEIGHT)/2, 655)
		}

		if rl.IsKeyDown(rl.KeyN) && s.editModeShowMenu {
			s.editModeShowMenu = false
			controls.DisableCursor(666)
			controls.SetMousePosition(int(s.cameraEditPos.X), int(s.cameraEditPos.Y), 661)
		}

	}

	updateCameraCenter(s.camera, s.cameraEditPos)
}

func (s *GameScene) syncDrawIndex(container *container.ObjectResourceContainer) {
	index := 0
	container.ForEachObject(func(obj models.Object) {
		image, ok := obj.(*models.Image)
		if ok {
			image.DrawIndex = index
		}
		index++
	})
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
