package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Hitbox struct {
	BaseEditorItem
	Id string
}

func (p *Hitbox) Draw() {
	if DRAW_MODELS {
		for i, _ := range p.Polygons {
			rl.DrawTriangleLines(
				p.Polygons[i].Points[0],
				p.Polygons[i].Points[1],
				p.Polygons[i].Points[2],
				rl.Blue,
			)
		}
	}
}

func (p *Hitbox) Update(delta float32) {

}