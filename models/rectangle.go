package models

import rl "github.com/gen2brain/raylib-go/raylib"

type Rectangle struct {
	pos   rl.Vector2
	box   rl.Vector2
	color rl.Color

	editorMoveWithCursor bool
}

func NewRectangle(x, y float32) *Rectangle {
	return &Rectangle{
		pos:   rl.NewVector2(x, y),
		box:   rl.NewVector2(20, 10),
		color: rl.NewColor(0, 121, 241, 100),
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
}

func (p *Rectangle) Update(delta float32) {
	if p.editorMoveWithCursor {
		mousePos := rl.GetMousePosition()
		offset := 10
		p.pos.X = mousePos.X-float32(offset)
		p.pos.Y = mousePos.Y-float32(offset)
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

func (p *Rectangle) ReactOnCollision() {
	rec := rl.NewRectangle(p.pos.X, p.pos.Y, p.box.X, p.box.Y)
	mousePos := rl.GetMousePosition()
	collission := rl.CheckCollisionPointRec(mousePos, rec)
	if collission {
		rl.DrawRectangleLinesEx(rec, 3.0, rl.Red)

		if rl.IsMouseButtonPressed(rl.MouseLeftButton){
			p.editorMoveWithCursor = !p.editorMoveWithCursor
		}		
	}
}
