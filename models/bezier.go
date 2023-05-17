package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Bezier struct {
	Start rl.Vector2
	End rl.Vector2

	Thick float32
	color rl.Color
}

func NewBezier(start, end rl.Vector2, thick float32) *Bezier {
	return &Bezier{
		Start:    start,
		End:    end,
		Thick: thick,
		color: rl.Gold,
	}
}

func (p *Bezier) Draw() {
	rl.DrawLineBezier(p.Start, p.End, p.Thick, p.color)
}

func (p *Bezier) Update(delta float32) {

}

func (p *Bezier) ResolveCollision(callback CollisionBezierCallback) {
	callback(p)
}

func (p Bezier) GetColor() rl.Color {
	return p.color
}

func (p *Bezier) SetColor(col rl.Color) *Bezier {
	p.color = col
	return p
}
