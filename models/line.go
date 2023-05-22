package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Line struct {
	Start rl.Vector2
	End rl.Vector2

	Thick float32
	color rl.Color
}

func NewLine(start, end rl.Vector2, thick float32) *Line {
	return &Line{
		Start:    start,
		End:    end,
		Thick: thick,
		color: rl.Gold,
	}
}

func (p *Line) Draw() {
	rl.DrawLineEx(p.Start, p.End, p.Thick, p.color)
}

func (p *Line) Update(delta float32) {

}

func (p *Line) ResolveCollision(callback CollisionLineCallback) {
	callback(p)
}

func (p Line) GetColor() rl.Color {
	return p.color
}

func (p *Line) SetColor(col rl.Color) *Line {
	p.color = col
	return p
}

func (p *Line) ProcessEditorSelection() bool {

	return true
}

func (p *Line) EditorResolveSelect() bool {
	mousePos := rl.GetMousePosition()
	isCollision := rl.CheckCollisionPointLine(mousePos, p.Start, p.End, int32(p.Thick))
	if isCollision {
		p.color = rl.Red
	} else {
		p.color = rl.Gold
	}
	return false
}