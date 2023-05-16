package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Bezier struct {
	start rl.Vector2
	end   rl.Vector2
	thick float32
	color rl.Color
}

func NewBezier(start, end rl.Vector2, thick float32) *Bezier {
	return &Bezier{
		start: start,
		end:   end,
		thick: thick,
		color: rl.Gold,
	}
}

func (p *Bezier) Draw() {
	rl.DrawLineBezier(p.start, p.end, p.thick, p.color)
}

func (p *Bezier) Update(delta float32) {

}

func (p Bezier) GetColor() rl.Color {
	return p.color
}

func (p *Bezier) SetColor(col rl.Color) *Bezier {
	p.color = col
	return p
}
