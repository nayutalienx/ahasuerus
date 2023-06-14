package scene

import (
	"ahasuerus/collision"
	"ahasuerus/container"
	"ahasuerus/controls"
	"ahasuerus/models"
	"ahasuerus/repository"
	"ahasuerus/resources"
	"fmt"
	"sort"

	"strings"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

const (
	editorStartMenuPosY     = 110
	editorControlRectWidth  = float32(300)
	editorControlRectHeight = float32(60)
	editorControlMarginLeft = 50
	maxTextSize             = 200
)

type EditScene struct {
	sceneName      string
	worldContainer *container.ObjectResourceContainer
	camera         *rl.Camera2D
	sourceScene    SceneId

	editorHubEnabled        bool
	cameraEditPos           rl.Vector2
	editCameraSpeed         float32
	selectedGameObjectsItem []models.EditorSelectedItem

	editMenuGameImageDropMode bool
}

func NewEditScene(
	sceneName string,
	sourceScene SceneId,
) *EditScene {

	rg.LoadStyleDefault()

	scene := &EditScene{
		sourceScene:             sourceScene,
		sceneName:               sceneName,
		worldContainer:          container.NewObjectResourceContainer(),
		cameraEditPos:           rl.NewVector2(0, 0),
		editCameraSpeed:         5,
		selectedGameObjectsItem: make([]models.EditorSelectedItem, 0),
	}

	worldImages := repository.GetAllImages(scene.sceneName)
	for i, _ := range worldImages {
		img := worldImages[i]
		scene.worldContainer.AddObjectResource(&img)
	}

	hitboxes := repository.GetAllHitboxes(scene.sceneName)
	for i, _ := range hitboxes {
		hb := hitboxes[i]
		scene.worldContainer.AddObjectResource(&hb)
	}

	camera := rl.NewCamera2D(
		rl.NewVector2(WIDTH/2, HEIGHT-500),
		rl.NewVector2(0, 0),
		0, 1)
	scene.camera = &camera

	controls.SetMousePosition(int(scene.cameraEditPos.X), int(scene.cameraEditPos.Y), 661)

	scene.worldContainer.Load()

	models.DRAW_MODELS = true

	return scene
}

func (s EditScene) Run() models.Scene {
	rg.SetStyle(rg.DEFAULT, rg.TEXT_SIZE, 30)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		delta := rl.GetFrameTime()
		s.camera.Zoom += rl.GetMouseWheelMove() * 0.05

		if rl.IsKeyDown(rl.KeyF2) && !s.editorHubEnabled {
			break
		}

		s.processInputs()

		updateCameraCenter(s.camera, s.cameraEditPos, delta)

		rl.BeginMode2D(*s.camera)

		s.worldContainer.Update(delta)
		s.worldContainer.Draw()

		if !s.editorHubEnabled {
			s.resolveEditorGameObjectsSelection()
		}

		if s.editorHubEnabled {
			s.processEditorGameObjectSelection()
		}

		rl.EndMode2D()

		if s.editorHubEnabled {
			s.drawEditorHub()
		}

		models.NewText(10, 50).
			SetFontSize(40).
			SetColor(rl.Red).SetData(fmt.Sprintf("edit mode[movement(arrow keys), cam.speed(+,-,%.1f), save(F10), menu(M), off menu(N), exit(F2)]", s.editCameraSpeed)).
			Draw()

		rl.EndDrawing()
	}

	s.Unload()
	return GetScene(s.sourceScene)
}
func (e EditScene) Unload() {
	e.worldContainer.Unload()
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
				s.editorHubEnabled = false
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

func (s EditScene) hasAnySelectedGameObjectEditorItem() (bool, models.EditorItem) {
	for i, _ := range s.selectedGameObjectsItem {
		ei := s.selectedGameObjectsItem[i]
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

			image, ok := editorItem.(*models.Image)
			if ok {
				repository.SaveImage(s.sceneName, image)
			}

			hitbox, ok := editorItem.(*models.Hitbox)
			if ok {
				repository.SaveHitbox(s.sceneName, hitbox)
			}

		}
	})
}

func (s *EditScene) drawEditorHub() {
	hasAnySelected, editorItem := s.hasAnySelectedGameObjectEditorItem()
	if hasAnySelected {
		s.drawHubForItem(editorItem)
	} else {
		s.drawMainHub()
	}
}

func (s *EditScene) drawMainHub() {
	bc := models.NewCounter()

	newGameImage := rg.Button(s.controlRect(bc), "NEW IMAGE")
	newCollisionBox := rg.Button(s.controlRect(bc), "NEW COLLISIONBOX")
	newLightBox := rg.Button(s.controlRect(bc), "NEW LIGHTBOX")
	newNpc := rg.Button(s.controlRect(bc), "NEW NPC")

	toggleModelsDrawText := "HIDE COLLISSION"
	if !models.DRAW_MODELS {
		toggleModelsDrawText = "SHOW COLLISSION"
	}
	toggleCollissionDrawButton := rg.Button(s.controlRect(bc), toggleModelsDrawText)

	if toggleCollissionDrawButton {
		models.DRAW_MODELS = !models.DRAW_MODELS
	}

	if s.editMenuGameImageDropMode {

		rl.DrawText("DROP IMAGE or F11 TO LEAVE", int32(WIDTH)/2, int32(HEIGHT)/2, 60, rl.Red)

		if rl.IsKeyDown(rl.KeyF11) {
			s.editMenuGameImageDropMode = false
		}

		if rl.IsFileDropped() {
			files := rl.LoadDroppedFiles()
			path := "resources" + strings.Split(files[0], "resources")[1]
			image := models.NewImage(
				s.worldContainer.Size(),
				uuid.NewString(),
				resources.GameTexture(path),
				s.camera.Target.X,
				s.camera.Target.Y, 0, 0, 0)

			image.Load()

			s.worldContainer.AddObjectResource(
				image,
			)

			s.editMenuGameImageDropMode = false
		}

	}

	if newGameImage {
		s.editMenuGameImageDropMode = true
	}

	if newCollisionBox || newLightBox || newNpc {
		height := float32(100)
		width := float32(100)

		topLeft := s.camera.Target
		bottomLeft := rl.Vector2{topLeft.X, topLeft.Y + height}

		topRight := rl.Vector2{topLeft.X + width, topLeft.Y}
		bottomRight := rl.Vector2{topLeft.X + width, topLeft.Y + height}

		hitboxType := models.Collision
		if newLightBox {
			hitboxType = models.Light
		}
		if newNpc {
			hitboxType = models.Npc
		}

		hitbox := models.Hitbox{
			Type: hitboxType,
			BaseEditorItem: models.NewBaseEditorItem([2]collision.Polygon{
				{
					Points: [3]rl.Vector2{
						topLeft, topRight, bottomRight,
					},
				},
				{
					Points: [3]rl.Vector2{
						topLeft, bottomLeft, bottomRight,
					},
				},
			}),
		}

		if newNpc {
			hitbox.Properties = map[string]string{
				"label": "npc",
				"blockOffsetX": "50.0",
				"blockOffsetY": "-150.0",
				"fontSize":     "60.0",
				"textOffsetX":  "20.0",
				"textOffsetY":  "15.0",
				"textCounter":  "0.0",
				"text":         "Hi there!;How are you?;Glad to see you here",
				"choice": "Hi:...;Fine:...;Bye",
				"choosed": "",
				"currentChoice": "0",
			}
		}

		s.worldContainer.AddObject(&hitbox)
	}
}

func (s *EditScene) reactOnEditorItemSelection(container *container.ObjectResourceContainer, item *models.BaseEditorItem, bc *models.Counter) bool {

	changePos := rg.Button(s.controlRect(bc), "CHANGE POS")
	resize := rg.Button(s.controlRect(bc), "RESIZE")
	rotate := rg.Button(s.controlRect(bc), "ROTATE")
	deleteItem := rg.Button(s.controlRect(bc), "DELETE")
	toggleProps := rg.Button(s.controlRect(bc), "PROPERTIES")
	unselect := rg.Button(s.controlRect(bc), "UNSELECT(F11)")

	if unselect {
		item.ExternalUnselect = true
	}

	if toggleProps {
		item.ShowProperties = !item.ShowProperties
	}

	if item.ShowProperties {
		s.drawEditorItemProperties(item, bc)
	}

	if changePos {
		item.SetEditorMoveWithCursorTrue()
		controls.DisableCursor(498)
		controls.SetMousePosition(int(item.TopLeft().X), int(item.TopLeft().Y), 500)
	}

	if resize {
		item.SetEditorResizeWithCursorTrue()
		controls.DisableCursor(506)
		controls.SetMousePosition(int(item.TopLeft().X+item.Width()), int(item.TopLeft().Y+item.Height()), 508)
	}

	if rotate {
		item.SetEditorRotateModeTrue()
	}

	return deleteItem
}

func (s EditScene) drawEditorItemProperties(item *models.BaseEditorItem, bc *models.Counter) {
	propertiesMargin := float32(20)
	propsLen := float32(len(item.Properties)) + 2
	propertiesPanelRect := s.controlRectWithMarginUp(bc, 10)
	propertiesPanelRect.Width *= 2.2
	propertiesPanelRect.Height *= propsLen
	rg.Panel(propertiesPanelRect, "Properties")

	keys := []string{}
	for k, _ := range item.Properties {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for i, _ := range keys {
		k := keys[i]
		propertiesRect := s.controlRectWithMargin(bc, propertiesMargin)

		rg.Label(propertiesRect, k)
		propertiesRect.X += propertiesRect.Width + propertiesMargin
		val := item.Properties[k]
		isActive := false
		if rl.CheckCollisionPointRec(rl.GetMousePosition(), propertiesRect) {
			isActive = true
		}

		shortText := val
		if len(val) >= 20 {
			shortText = val[:20]
			isActive = false
			rg.TextBoxMulti(propertiesRect, &shortText, maxTextSize, isActive)
		} else {
			rg.TextBoxMulti(propertiesRect, &val, maxTextSize, isActive)
		}

		item.Properties[k] = val
	}

}

func (s EditScene) itemPosY(buttonCounter *models.Counter) float32 {
	return float32(editorStartMenuPosY + int(editorControlRectHeight)*buttonCounter.GetAndIncrement())
}

func (s EditScene) controlRect(bc *models.Counter) rl.Rectangle {
	return rl.NewRectangle(editorControlMarginLeft, s.itemPosY(bc), editorControlRectWidth, editorControlRectHeight)
}

func (s EditScene) controlRectWithMarginUp(bc *models.Counter, margin float32) rl.Rectangle {
	return rl.NewRectangle(editorControlMarginLeft, s.itemPosY(bc)+margin, editorControlRectWidth, editorControlRectHeight)
}

func (s EditScene) controlRectWithMargin(bc *models.Counter, margin float32) rl.Rectangle {
	return rl.NewRectangle(editorControlMarginLeft+margin, s.itemPosY(bc)+margin, editorControlRectWidth, editorControlRectHeight)
}

func (s *EditScene) reactOnImageEditorSelection(container *container.ObjectResourceContainer, image *models.Image, bc *models.Counter) {

	moveUpper := rg.Button(s.controlRect(bc), "MOVE UPPER")
	moveDown := rg.Button(s.controlRect(bc), "MOVE DOWN")

	replicate := rg.Button(s.controlRect(bc), "REPLICATE")

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

}

func (s *EditScene) drawHubForItem(editorItem models.EditorItem) {

	buttonCounter := models.NewCounter()

	img, isImg := editorItem.(*models.Image)
	if isImg {
		delete := s.reactOnEditorItemSelection(s.worldContainer, &img.BaseEditorItem, buttonCounter)
		s.reactOnImageEditorSelection(s.worldContainer, img, buttonCounter)
		if delete {
			s.worldContainer.RemoveObject(img)
			repository.DeleteImage(s.sceneName, img)
		}
	}

	hitbox, isHitbox := editorItem.(*models.Hitbox)
	if isHitbox {
		delete := s.reactOnEditorItemSelection(s.worldContainer, &hitbox.BaseEditorItem, buttonCounter)
		if delete {
			s.worldContainer.RemoveObject(hitbox)
			repository.DeleteHitbox(s.sceneName, hitbox.Id)
		}
	}

}

func (s *EditScene) processInputs() {

	mousePos := rl.GetMousePosition()

	updateMouse := false

	if rl.IsKeyDown(rl.KeyRight) {
		s.cameraEditPos.X += s.editCameraSpeed
		if !s.editorHubEnabled {
			mousePos.X += s.editCameraSpeed
			updateMouse = true
		}
	}

	if rl.IsKeyDown(rl.KeyLeft) {
		s.cameraEditPos.X -= s.editCameraSpeed
		if !s.editorHubEnabled {
			mousePos.X -= s.editCameraSpeed
			updateMouse = true
		}
	}

	if rl.IsKeyDown(rl.KeySpace) {
		mouseDelta := rl.Vector2Negate(rl.GetMouseDelta())
		s.cameraEditPos = rl.Vector2Add(s.cameraEditPos, mouseDelta)
		mousePos = rl.Vector2Add(mousePos, mouseDelta)
		updateMouse = true
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

	if rl.IsKeyDown(rl.KeyF10) {
		s.saveEditor()
		s.worldContainer.AddObject(
			models.NewText(int32(s.camera.Target.X-s.camera.Offset.X+WIDTH/2), int32(s.camera.Target.Y-s.camera.Offset.Y+HEIGHT/2)).
				SetData("DATA SAVED").
				SetFontSize(60).
				SetColor(rl.Red).
				WithExpire(3, func(t *models.Text) {
					s.worldContainer.RemoveObject(t)
				}),
		)
	}

	hasAnySelected, _ := s.hasAnySelectedGameObjectEditorItem()

	if hasAnySelected {

		if !s.editorHubEnabled {
			s.editorHubEnabled = true
			controls.EnableCursor(653)
		}

	} else {

		if (rl.IsKeyDown(rl.KeyM)) && !s.editorHubEnabled {
			s.editorHubEnabled = true
			controls.EnableCursor(660)
			controls.SetMousePosition(int(WIDTH)/2, int(HEIGHT)/2, 655)
		}

		if rl.IsKeyDown(rl.KeyN) && s.editorHubEnabled {
			s.editorHubEnabled = false
			controls.DisableCursor(666)
			controls.SetMousePosition(int(s.cameraEditPos.X), int(s.cameraEditPos.Y), 661)
		}

	}
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
