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

type EditScene struct {
	sceneName            string
	worldContainer       *container.ObjectResourceContainer
	environmentContainer *container.ObjectResourceContainer
	camera               *rl.Camera2D
	sourceScene          SceneId

	editModeShowMenu        bool
	cameraEditPos           rl.Vector2
	editCameraSpeed         float32
	selectedGameObjectsItem []models.EditorSelectedItem
	selectedBackgroundItem  []models.EditorSelectedItem

	editMenuBgImageDropMode   bool
	editMenuGameImageDropMode bool
	editHideGameObjectsMode   bool
	editBgImageEditorMode     bool
}

func NewEditScene(
	sceneName string,
	sourceScene SceneId,
) *EditScene {

	scene := &EditScene{
		sourceScene:             sourceScene,
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
		scene.worldContainer.AddObjectResource(&img)
	}

	camera := rl.NewCamera2D(
		rl.NewVector2(WIDTH/2, HEIGHT-500),
		rl.NewVector2(0, 0),
		0, 1)
	scene.camera = &camera

	controls.SetMousePosition(int(scene.cameraEditPos.X), int(scene.cameraEditPos.Y), 661)

	scene.environmentContainer.Load()
	scene.worldContainer.Load()

	return scene
}

func (s EditScene) Run() models.Scene {
	rg.SetStyle(rg.DEFAULT, rg.TEXT_SIZE, 20)

	editLabel := models.NewText(10, 50).
		SetFontSize(40).
		SetColor(rl.Red).
		SetUpdateCallback(func(t *models.Text) {
			t.SetData(fmt.Sprintf("edit mode[movement(arrow keys), cam.speed(+,-,%.1f), save(P), menu(M), off menu(N), exit(F2)]", s.editCameraSpeed))
		})

	s.environmentContainer.AddObject(
		editLabel,
	)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		delta := rl.GetFrameTime()
		s.camera.Zoom += rl.GetMouseWheelMove() * 0.05

		if rl.IsKeyDown(rl.KeyF2) && !s.editModeShowMenu {
			break
		}

		s.processEditorMode()

		rl.BeginMode2D(*s.camera)
		if !s.editHideGameObjectsMode {
			s.worldContainer.Update(delta)
			s.worldContainer.Draw()
			if !s.editModeShowMenu {
				s.resolveEditorGameObjectsSelection()
			}
			if s.editModeShowMenu {
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

	s.Unload()
	return GetScene(s.sourceScene)
}
func (e EditScene) Unload() {
	e.environmentContainer.Unload()
	e.worldContainer.Unload()
}

func (s *EditScene) resolveEditorBackgroundImageSelection() {
	selectedItem := make([]models.EditorSelectedItem, 0)

	s.environmentContainer.ForEachObjectReverseWithPredicate(func(obj models.Object) bool {
		editorItem, ok := obj.(models.EditorItem)
		if ok {
			resolveResult := editorItem.EditorDetectSelection()
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

func (s *EditScene) resolveEditorGameObjectsSelection() {
	mouse := rl.GetMousePosition()
	rl.DrawText(fmt.Sprintf("%v %v", mouse.X, mouse.Y), int32(mouse.X), int32(mouse.Y), 40, rl.Red)
	rl.DrawCircle(int32(mouse.X), int32(mouse.Y), 10, rl.Red)

	selectedItem := make([]models.EditorSelectedItem, 0)

	s.worldContainer.ForEachObjectReverseWithPredicate(func(obj models.Object) bool {
		editorItem, ok := obj.(models.EditorItem)
		if ok {
			resolveResult := editorItem.EditorDetectSelection()
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

func (s *EditScene) processEditorGameObjectSelection() {
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

func (s *EditScene) processEditorBackgroundSelection() {
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

func (s EditScene) hasAnySelectedGameObjectEditorItem() (bool, models.EditorItem) {
	for i, _ := range s.selectedGameObjectsItem {
		ei := s.selectedGameObjectsItem[i]
		if ei.Selected {
			return true, ei.Item
		}
	}
	return false, nil
}

func (s EditScene) hasAnySelectedBackgroundEditorItem() (bool, models.EditorItem) {
	for i, _ := range s.selectedBackgroundItem {
		ei := s.selectedBackgroundItem[i]
		if ei.Selected {
			return true, ei.Item
		}
	}
	return false, nil
}

func (s *EditScene) saveEditor() {
	s.worldContainer.ForEachObject(func(obj models.Object) {
		editorItem, ok := obj.(models.EditorItem)
		if ok {
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

func (s *EditScene) processEditorMenuMode() {
	hasAnySelected, editorItem := s.hasAnySelectedGameObjectEditorItem()
	if hasAnySelected {
		s.reactOnGameObjectEditorSelect(editorItem)
	} else {
		s.drawNonGameFocusedMenu()
	}
}

func (s *EditScene) drawNonGameFocusedMenu() {
	buttonCounter := models.NewCounter()

	newGameImage := false
	newBgImage := false
	//newCollisionBox = false

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
		newGameImage = rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "NEW IMAGE")
		//newCollisionBox = rg.Button(rl.NewRectangle(10, float32(editorStartMenuPosY+editorMenuButtonHeight*buttonCounter.GetAndIncrement()), float32(editorMenuButtonWidth), float32(editorMenuButtonHeight)), "NEW COLLISION BOX")
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

func (s *EditScene) reactOnImageEditorSelection(container *container.ObjectResourceContainer, image *models.Image, buttonCounter *models.Counter) {

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

func (s *EditScene) reactOnGameObjectEditorSelect(editorItem models.EditorItem) {

	buttonCounter := models.NewCounter()

	img, isImg := editorItem.(*models.Image)
	if isImg {
		s.reactOnImageEditorSelection(s.worldContainer, img, buttonCounter)
	}

}

func (s *EditScene) processEditorMode() {

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

func (s *EditScene) syncDrawIndex(container *container.ObjectResourceContainer) {
	index := 0
	container.ForEachObject(func(obj models.Object) {
		image, ok := obj.(*models.Image)
		if ok {
			image.DrawIndex = index
		}
		index++
	})
}
