package models

import (
	"ahasuerus/collision"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type BaseEditorItem struct {
	Polygons [2]collision.Polygon

	Rotation               float32
	EditSelected           bool
	EditorMoveWithCursor   bool
	EditorResizeWithCursor bool
	EditorRotateMode       bool
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
		if rl.IsKeyDown(rl.KeyBackspace) {
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
