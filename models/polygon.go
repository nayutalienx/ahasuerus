package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Polygon struct {
	Points   [3]rl.Vector2
	Color    rl.Color
}

func (p *Polygon) Draw() {
	if DRAW_MODELS {
		rl.DrawTriangleLines(p.Points[0], p.Points[1], p.Points[2], p.Color)
	}
}

func (p *Polygon) Update(delta float32) {

}
