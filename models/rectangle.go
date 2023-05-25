package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Rectangle struct {
	Id string
	pos   rl.Vector2
	box   rl.Vector2
	color rl.Color

	editSelected  bool
	editorMoveWithCursor bool
	editorEditSizeWithCursor bool
}

func NewRectangle(id string, x, y, width, height float32, color rl.Color) *Rectangle {
	return &Rectangle{
		Id: id,
		pos:   rl.NewVector2(x, y),
		box:   rl.NewVector2(width, height),
		color: color,
	}
}

func (p Rectangle) ResolveCollision(callback CollisionBoxCallback) {
	callback(&p)
}

func (p *Rectangle) Draw() {
	rl.DrawRectangle(int32(p.pos.X), int32(p.pos.Y), int32(p.box.X), int32(p.box.Y), p.color)
}

func (p *Rectangle) Update(delta float32) {
}

func (p Rectangle) GetColor() rl.Color {
	return p.color
}

func (p *Rectangle) SetColor(col rl.Color) *Rectangle {
	p.color = col
	return p
}

func (p *Rectangle) SetWidth(width float32) *Rectangle {
	p.box.X = width
	return p
}

func (p *Rectangle) SetHeight(height float32) *Rectangle {
	p.box.Y = height
	return p
}

func (p *Rectangle) GetPos() *rl.Vector2 {
	return &p.pos
}

func (p *Rectangle) GetBox() *rl.Vector2 {
	return &p.box
}

func (p *Rectangle) SetEditorSizeModeTrue() {
	p.editorEditSizeWithCursor = true
}

func (p *Rectangle) SetEditorMoveModeTrue() {
	p.editorMoveWithCursor = true
}

func (p *Rectangle) ProcessEditorSelection() EditorItemProcessSelectionResult {

	if p.editorMoveWithCursor {
		mousePos := rl.GetMousePosition()
		offset := 10
		p.pos.X = mousePos.X-float32(offset)
		p.pos.Y = mousePos.Y-float32(offset)
	}

	if p.editorEditSizeWithCursor {
		mousePos := rl.GetMousePosition()
		offset := 10
		p.box.X = mousePos.X+float32(offset) - p.pos.X 
		p.box.Y = mousePos.Y+float32(offset) - p.pos.Y
	}

	if (p.editorMoveWithCursor || p.editorEditSizeWithCursor) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		p.editorEditSizeWithCursor = false
		p.editorMoveWithCursor = false
		p.editSelected = false
		return EditorItemProcessSelectionResult{
			Finished:      true,
		}
	}

	if rl.IsKeyDown(rl.KeyBackspace) {
		p.editorEditSizeWithCursor = false
		p.editorMoveWithCursor = false
		p.editSelected = false
		return EditorItemProcessSelectionResult{
			Finished:      true,
			DisableCursor: true,
			CursorForcePosition: true,
			CursorX: int(p.pos.X),
			CursorY: int(p.pos.Y),
		}
	}

	return EditorItemProcessSelectionResult{
		Finished:      false,
	}
}

func (p *Rectangle) EditorResolveSelect() (EditorItemResolveSelectionResult) {
	rec := rl.NewRectangle(p.pos.X, p.pos.Y, p.box.X, p.box.Y)
	mousePos := rl.GetMousePosition()
	collission := rl.CheckCollisionPointRec(mousePos, rec)
	if collission {
		rl.DrawRectangleLinesEx(rec, 3.0, rl.Red)

		if rl.IsMouseButtonPressed(rl.MouseLeftButton){
			p.editSelected = true
		}		
	}
	return EditorItemResolveSelectionResult{
		Selected:  p.editSelected,
		Collision: collission,
	}
}
