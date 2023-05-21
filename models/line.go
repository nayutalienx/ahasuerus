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
