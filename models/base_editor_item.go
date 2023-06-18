package models

import (
	"ahasuerus/collision"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

type BaseEditorItem struct {
	Id       string
	Polygons [2]collision.Polygon

	Rotation float32

	EditSelected           bool `json:"-"`
	ExternalUnselect       bool `json:"-"`
	EditorMoveWithCursor   bool `json:"-"`
	EditorResizeWithCursor bool `json:"-"`
	EditorRotateMode       bool `json:"-"`
}

func NewBaseEditorItem(polygons [2]collision.Polygon) BaseEditorItem {
	return BaseEditorItem{
		Id:       uuid.NewString(),
		Polygons: polygons,
	}
}

func (p *BaseEditorItem) SetEditorMoveWithCursorTrue() {
	p.EditorMoveWithCursor = true
}

func (p *BaseEditorItem) SetEditorResizeWithCursorTrue() {
	p.EditorResizeWithCursor = true
}

func (p *BaseEditorItem) SetEditorRotateModeTrue() {
	p.EditorRotateMode = true
}

func (p *BaseEditorItem) EditorDetectSelection() EditorItemDetectSelectionResult {
	mousePos := rl.GetMousePosition()
	triangle1 := p.Polygons[0].Points
	triangle2 := p.Polygons[1].Points
	RotateTriangleByA(&triangle1[0], &triangle1[1], &triangle1[2], float64(p.Rotation))
	RotateTriangleByA(&triangle2[0], &triangle2[1], &triangle2[2], float64(p.Rotation))
	collission := rl.CheckCollisionPointTriangle(mousePos, triangle1[0], triangle1[1], triangle1[2]) ||
		rl.CheckCollisionPointTriangle(mousePos, triangle2[0], triangle2[1], triangle2[2])
	if collission {
		rl.DrawTriangleLines(triangle1[0], triangle1[1], triangle1[2], rl.Purple)
		rl.DrawTriangleLines(triangle2[0], triangle2[1], triangle2[2], rl.Purple)

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			p.EditSelected = true
		}
	}

	return EditorItemDetectSelectionResult{
		Selected:  p.EditSelected,
		Collision: collission,
	}
}

func (p *BaseEditorItem) ProcessEditorSelection() EditorItemProcessSelectionResult {

	if p.EditorMoveWithCursor {
		mousePos := rl.GetMousePosition()
		rl.DrawCircle(int32(mousePos.X), int32(mousePos.Y), 10, rl.Red)
		offset := 10
		newPosX := mousePos.X - float32(offset)
		newPosY := mousePos.Y - float32(offset)

		p.ChangePosition(rl.NewVector2(newPosX, newPosY))
	}

	if p.EditorResizeWithCursor {
		mousePos := rl.GetMousePosition()
		rl.DrawCircle(int32(mousePos.X), int32(mousePos.Y), 10, rl.Red)
		offset := 10

		newPosX := mousePos.X - float32(offset)
		newPosY := mousePos.Y - float32(offset)

		width := p.Width()
		height := p.Height()

		p.Polygons[0].Points[1].X = newPosX
		p.Polygons[0].Points[1].Y = newPosY - height

		p.Polygons[0].Points[2].X = newPosX
		p.Polygons[0].Points[2].Y = newPosY

		p.Polygons[1].Points[1].X = newPosX - width
		p.Polygons[1].Points[1].Y = newPosY

		p.Polygons[1].Points[2].X = newPosX
		p.Polygons[1].Points[2].Y = newPosY
	}

	if p.EditorRotateMode {
		if rl.IsKeyDown(rl.KeyT) {
			p.Rotation++
		}
		if rl.IsKeyDown(rl.KeyR) {
			p.Rotation--
		}

		if p.Rotation < 0 {
			p.Rotation = 360
		}
		if p.Rotation > 360 {
			p.Rotation = 0
		}
	}

	if (p.EditorMoveWithCursor || p.EditorResizeWithCursor) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		p.EditorMoveWithCursor = false
		p.EditorResizeWithCursor = false
		p.EditSelected = false
		p.EditorRotateMode = false
		return EditorItemProcessSelectionResult{
			Finished: true,
		}
	}

	if p.EditSelected {
		if rl.IsKeyDown(rl.KeyF11) || p.ExternalUnselect {
			p.ExternalUnselect = false
			p.EditorMoveWithCursor = false
			p.EditorResizeWithCursor = false
			p.EditSelected = false
			p.EditorRotateMode = false
			return EditorItemProcessSelectionResult{
				Finished:            true,
				DisableCursor:       true,
				CursorForcePosition: true,
				CursorX:             int(p.Polygons[0].Points[0].X),
				CursorY:             int(p.Polygons[0].Points[0].Y),
			}
		}
	}

	return EditorItemProcessSelectionResult{
		Finished: false,
	}
}

func (p BaseEditorItem) Center() rl.Vector2 {
	tl := p.TopLeft()
	br := p.BottomRight()
	return rl.Vector2Scale(rl.Vector2Add(tl, br), 0.5)
}

func (p BaseEditorItem) TopLeft() rl.Vector2 {
	return p.Polygons[0].Points[0]
}

func (p BaseEditorItem) TopRight() rl.Vector2 {
	return p.Polygons[0].Points[1]
}

func (p BaseEditorItem) BottomRight() rl.Vector2 {
	return p.Polygons[0].Points[2]
}

func (p BaseEditorItem) BottomLeft() rl.Vector2 {
	return p.Polygons[1].Points[1]
}

func (p BaseEditorItem) Width() float32 {
	return p.Polygons[0].Points[2].X - p.Polygons[0].Points[0].X
}

func (p BaseEditorItem) Height() float32 {
	return p.Polygons[0].Points[2].Y - p.Polygons[0].Points[0].Y
}

func (p BaseEditorItem) PolygonsWithRotation() []collision.Polygon {
	polys := p.Polygons[:]
	if p.Rotation != 0 {
		RotateTriangleByA(&polys[0].Points[0], &polys[0].Points[1], &polys[0].Points[2], float64(p.Rotation))
		RotateTriangleByA(&polys[1].Points[0], &polys[1].Points[1], &polys[1].Points[2], float64(p.Rotation))
	}
	return polys
}

func (p *BaseEditorItem) Translate(movement rl.Vector2) {
	newPos := rl.Vector2Add(p.Polygons[0].Points[0], movement)
	p.ChangePosition(newPos)
}

func (p *BaseEditorItem) ChangePosition(newPos rl.Vector2) {
	newPosX := newPos.X
	newPosY := newPos.Y

	width := p.Width()
	height := p.Height()

	p.Polygons[0].Points[0].X = newPosX
	p.Polygons[0].Points[0].Y = newPosY

	p.Polygons[0].Points[1].X = newPosX + width
	p.Polygons[0].Points[1].Y = newPosY

	p.Polygons[0].Points[2].X = newPosX + width
	p.Polygons[0].Points[2].Y = newPosY + height

	p.Polygons[1].Points[0].X = newPosX
	p.Polygons[1].Points[0].Y = newPosY

	p.Polygons[1].Points[1].X = newPosX
	p.Polygons[1].Points[1].Y = newPosY + height

	p.Polygons[1].Points[2].X = newPosX + width
	p.Polygons[1].Points[2].Y = newPosY + height
}

func (p *BaseEditorItem) SetPolygons(polys [2]collision.Polygon) {
	p.Polygons = polys
}

func (p BaseEditorItem) Draw() {
	if p.EditorRotateMode {
		rl.DrawText("Rotate on [R and T]", int32(p.TopLeft().X), int32(p.TopLeft().Y+40), 40, rl.Red)
	}
}
