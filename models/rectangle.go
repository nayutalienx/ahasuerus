package models

import rl "github.com/gen2brain/raylib-go/raylib"

type Rectangle struct {
	pos   rl.Vector2
	box   rl.Vector2
	color rl.Color
}

func NewRectangle(x, y float32) *Rectangle {
	return &Rectangle{
		pos:   rl.NewVector2(x, y),
		box:   rl.NewVector2(20, 10),
		color: rl.Blue,
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
