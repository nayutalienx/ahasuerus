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
		p.BaseEditorItem.Draw()
	}
}

func (p *Hitbox) Update(delta float32) {

	if p.Rotation != 0 {
		RotateTriangleByA(&p.Polygons[0].Points[0], &p.Polygons[0].Points[1], &p.Polygons[0].Points[2], float64(p.Rotation))
		RotateTriangleByA(&p.Polygons[1].Points[0], &p.Polygons[1].Points[1], &p.Polygons[1].Points[2], float64(p.Rotation))
		p.Rotation = 0
	}

}
