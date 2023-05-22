package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Rectangle struct {
	Id string
	pos   rl.Vector2
	box   rl.Vector2
	color rl.Color

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

	if p.editorMoveWithCursor {
		rl.DrawText("ATTACHED TO CURSOR", int32(p.pos.X), int32(p.pos.Y) + 30, 20, rl.Red)
	}

	if p.editorEditSizeWithCursor {
		rl.DrawText("SIZE ATTACHED TO CURSOR", int32(p.pos.X), int32(p.pos.Y) + 30, 20, rl.Red)
	}
}

func (p *Rectangle) Update(delta float32) {
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

func (p *Rectangle) ProcessEditorSelection() bool {

	return true
}

func (p *Rectangle) EditorResolveSelect() bool {
	rec := rl.NewRectangle(p.pos.X, p.pos.Y, p.box.X, p.box.Y)
	mousePos := rl.GetMousePosition()
	collission := rl.CheckCollisionPointRec(mousePos, rec)
	if collission {
		sizeEditorSquareSize := 20

		rl.DrawRectangleLinesEx(rec, 3.0, rl.Red)

		rl.DrawRectangle(
			int32(p.pos.X) + int32(p.box.X-float32(sizeEditorSquareSize)), 
			int32(p.pos.Y)+int32(p.box.Y-float32(sizeEditorSquareSize)), 
			int32(sizeEditorSquareSize), 
			int32(sizeEditorSquareSize), 
			rl.Red)

		collissionForEditSize := rl.CheckCollisionPointRec(
			mousePos, 
			rl.NewRectangle(
				float32(int32(p.pos.X) + int32(p.box.X-float32(sizeEditorSquareSize))), 
				float32(int32(p.pos.Y)+int32(p.box.Y-float32(sizeEditorSquareSize))),
				float32(sizeEditorSquareSize),
				float32(sizeEditorSquareSize),
			),
		)

		if rl.IsMouseButtonPressed(rl.MouseLeftButton){
			if collissionForEditSize {
				p.editorEditSizeWithCursor = !p.editorEditSizeWithCursor
			} else {
				p.editorMoveWithCursor = !p.editorMoveWithCursor
			}
		}		
	}
	return false
}
