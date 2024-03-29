package scene

import (
	"ahasuerus/collision"
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
	editorStartMenuPosY     = 110
	editorControlRectWidth  = float32(300)
	editorControlRectHeight = float32(60)
	editorControlMarginLeft = 50
	maxTextSize             = 200
)

type EditScene struct {
	worldContainer *container.ObjectResourceContainer
	camera         *rl.Camera2D
	sourceScene    SceneId

	editorHubEnabled        bool
	editCameraSpeed         float32
	selectedGameObjectsItem []models.EditorSelectedItem

	editMenuGameImageDropMode bool

	onScreenQueue chan models.Object
	screenScale float32
	level repository.Level
}

func NewEditScene(
	sceneName string,
	sourceScene SceneId,
) *EditScene {

	rg.LoadStyleDefault()

	level := repository.GetLevel(sceneName)
	
	levelSize := level.Size()

	screenScale := HEIGHT/levelSize.Y

	camera := rl.NewCamera2D(
		rl.NewVector2(WIDTH/2, HEIGHT/2),
		rl.NewVector2(WIDTH/2, levelSize.Y/2),
		0, screenScale)

	scene := &EditScene{
		level:                   level,
		camera:                  &camera,
		sourceScene:             sourceScene,
		worldContainer:          container.NewObjectResourceContainer(),
		editCameraSpeed:         5,
		selectedGameObjectsItem: make([]models.EditorSelectedItem, 0),
		onScreenQueue:  make(chan models.Object, 2),
		screenScale: screenScale,
	}

	worldImages := scene.level.Images
	for i, _ := range worldImages {
		img := worldImages[i]
		img.Camera(scene.camera)
		scene.worldContainer.AddObjectResource(&img)
	}

	collisionHitboxes := scene.level.CollissionHitboxes
	for i, _ := range collisionHitboxes {
		hb := collisionHitboxes[i]
		scene.worldContainer.AddObjectResource(&hb)
	}

	lights := scene.level.Lights
	for i, _ := range lights {
		light := lights[i]
		scene.worldContainer.AddObject(&light)
	}

	characters := scene.level.Characters
	for i, _ := range characters {
		npc := characters[i].ScreenChan(scene.onScreenQueue).ScreenScale(scene.screenScale)
		scene.worldContainer.AddObjectResource(npc)
	}

	particles := scene.level.ParticleSources
	for i, _ := range particles {
		particle := particles[i]
		scene.worldContainer.AddObjectResource(&particle)
	}

	controls.SetMousePosition(int(scene.camera.Target.X), int(scene.camera.Target.Y), 661)

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

		for len(s.onScreenQueue) > 0 {
			onScreenObject := <-s.onScreenQueue
			onScreenObject.Draw()
			onScreenObject.Update(delta)
		}

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

		_, isImage := obj.(*models.Image)
		if models.DRAW_MODELS && isImage {
			return false // skip image selection when draw game objects
		}

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
	newLevel := s.level

	newLevel.Characters = []models.Npc{}
	newLevel.Lights = []models.Light{}
	newLevel.CollissionHitboxes = []models.CollisionHitbox{}
	newLevel.Images = []models.Image{}
	newLevel.ParticleSources = []models.ParticleSource{}

	s.worldContainer.ForEachObject(func(obj models.Object) {
		editorItem, ok := obj.(models.EditorItem)
		if ok {

			image, ok := editorItem.(*models.Image)
			if ok {
				newLevel.Images = append(newLevel.Images, *image)
			}

			hitbox, ok := editorItem.(*models.CollisionHitbox)
			if ok {
				newLevel.CollissionHitboxes = append(newLevel.CollissionHitboxes, *hitbox)
			}

			light, ok := editorItem.(*models.Light)
			if ok {
				newLevel.Lights = append(newLevel.Lights, *light)
			}

			npc, ok := editorItem.(*models.Npc)
			if ok {
				newLevel.Characters = append(newLevel.Characters, *npc)
			}

			particleSource, ok := editorItem.(*models.ParticleSource)
			if ok {
				newLevel.ParticleSources = append(newLevel.ParticleSources, *particleSource)
			}

		}
	})
	newLevel.SaveLevel()
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

	newGameImage := rg.Button(s.controlRect(&bc), "NEW IMAGE")
	newCollisionBox := rg.Button(s.controlRect(&bc), "NEW COLLISIONBOX")
	newLightBox := rg.Button(s.controlRect(&bc), "NEW LIGHTBOX")
	newNpc := rg.Button(s.controlRect(&bc), "NEW NPC")
	newParticleSource := rg.Button(s.controlRect(&bc), "PARTICLES")

	toggleModelsDrawText := "HIDE COLLISSION"
	if !models.DRAW_MODELS {
		toggleModelsDrawText = "SHOW COLLISSION"
	}
	toggleCollissionDrawButton := rg.Button(s.controlRect(&bc), toggleModelsDrawText)

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
				uuid.NewString(),
				resources.GameTexture(path),
				0,
				0,
				0)

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

	if newCollisionBox || newLightBox || newNpc || newParticleSource {
		height := float32(100)
		width := float32(100)

		topLeft := s.camera.Target
		bottomLeft := rl.Vector2{topLeft.X, topLeft.Y + height}

		topRight := rl.Vector2{topLeft.X + width, topLeft.Y}
		bottomRight := rl.Vector2{topLeft.X + width, topLeft.Y + height}

		var newObject models.Object

		baseEditorItem := models.NewBaseEditorItem([2]collision.Polygon{
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
		})

		if newCollisionBox {
			newObject = &models.CollisionHitbox{
				BaseEditorItem: baseEditorItem,
			}
		}

		if newLightBox {
			newObject = &models.Light{
				BaseEditorItem: baseEditorItem,
			}
		}

		if newParticleSource {
			ps := models.NewParticleSource(baseEditorItem)
			ps.Load()
			newObject = ps
		}

		if newNpc {
			npc := &models.Npc{
				CollisionHitbox: models.CollisionHitbox{
					BaseEditorItem: baseEditorItem,
				},
			}
			npc.ScreenChan(s.onScreenQueue)
			newObject = npc
		}

		s.worldContainer.AddObject(newObject)
	}
}

func (s *EditScene) reactOnEditorItemSelection(container *container.ObjectResourceContainer, item *models.BaseEditorItem, bc *models.Counter) {

	changePos := rg.Button(s.controlRect(bc), "CHANGE POS")
	resize := rg.Button(s.controlRect(bc), "RESIZE")
	rotate := rg.Button(s.controlRect(bc), "ROTATE")
	deleteItem := rg.Button(s.controlRect(bc), "DELETE")
	unselect := rg.Button(s.controlRect(bc), "UNSELECT(F11)")

	moveUpper := rg.Button(s.controlRect(bc), "MOVE UPPER")
	moveDown := rg.Button(s.controlRect(bc), "MOVE DOWN")

	if moveUpper {
		container.MoveUp(item)
		item.DrawIndex--
	}

	if moveDown {
		container.MoveDown(item)
		item.DrawIndex++
	}

	if unselect {
		item.ExternalUnselect = true
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

	if deleteItem {
		s.worldContainer.RemoveObject(item)
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

	replicate := rg.Button(s.controlRect(bc), "REPLICATE")

	if replicate {
		topLeft := image.TopLeft()
		imageReplica := image.Replicate(uuid.NewString(), topLeft.X-100, topLeft.Y-100)
		imageReplica.Load()
		container.AddObjectResource(imageReplica)
	}

}

func (s *EditScene) drawHubForItem(editorItem models.EditorItem) {

	buttonCounter := models.NewCounter()

	img, isImg := editorItem.(*models.Image)
	if isImg {
		s.reactOnEditorItemSelection(s.worldContainer, &img.BaseEditorItem, &buttonCounter)
		s.reactOnImageEditorSelection(s.worldContainer, img, &buttonCounter)
	}

	collisionHitbox, isHitbox := editorItem.(*models.CollisionHitbox)
	if isHitbox {
		s.reactOnEditorItemSelection(s.worldContainer, &collisionHitbox.BaseEditorItem, &buttonCounter)
	}

	light, isLight := editorItem.(*models.Light)
	if isLight {
		s.reactOnEditorItemSelection(s.worldContainer, &light.BaseEditorItem, &buttonCounter)
	}

	npc, isNpc := editorItem.(*models.Npc)
	if isNpc {
		s.reactOnEditorItemSelection(s.worldContainer, &npc.BaseEditorItem, &buttonCounter)
	}

	particleSource, isParticleSource := editorItem.(*models.ParticleSource)
	if isParticleSource {
		s.reactOnEditorItemSelection(s.worldContainer, &particleSource.BaseEditorItem, &buttonCounter)
	}

}

func (s *EditScene) processInputs() {

	mousePos := rl.GetMousePosition()

	updateMouse := false

	if rl.IsKeyDown(rl.KeyRight) {
		s.camera.Target.X += s.editCameraSpeed
		if !s.editorHubEnabled {
			mousePos.X += s.editCameraSpeed
			updateMouse = true
		}
	}

	if rl.IsKeyDown(rl.KeyLeft) {
		s.camera.Target.X -= s.editCameraSpeed
		if !s.editorHubEnabled {
			mousePos.X -= s.editCameraSpeed
			updateMouse = true
		}
	}

	if rl.IsKeyDown(rl.KeySpace) {
		mouseDelta := rl.Vector2Negate(rl.GetMouseDelta())
		s.camera.Target = rl.Vector2Add(s.camera.Target, mouseDelta)
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
			controls.SetMousePosition(int(s.camera.Target.X), int(s.camera.Target.Y), 661)
		}

	}
}
