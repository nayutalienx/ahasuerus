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
		polys := p.Polygons()
		for i, _ := range polys {
			rl.DrawTriangleLines(
				polys[i].Points[0],
				polys[i].Points[1],
				polys[i].Points[2],
				rl.Blue,
			)
		}
		p.BaseEditorItem.Draw()
	}
}

func (p *Hitbox) Update(delta float32) {

}
