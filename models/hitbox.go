package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type HitboxType int

const (
	Collision HitboxType = iota
	Light
)

type Hitbox struct {
	BaseEditorItem
	Id   string
	Type HitboxType
}

func (p *Hitbox) Draw() {
	if DRAW_MODELS {
		if p.Type == Collision {
			polys := p.Polygons()
			for i, _ := range polys {
				rl.DrawTriangleLines(
					polys[i].Points[0],
					polys[i].Points[1],
					polys[i].Points[2],
					rl.Blue,
				)
			}
		}
		if p.Type == Light {
			center := p.Center()
			rl.DrawCircleLines(int32(center.X), int32(center.Y), rl.Vector2Distance(p.TopLeft(), p.TopRight())/2, rl.Gold)
		}
		p.BaseEditorItem.Draw()
	}
}

func (p *Hitbox) Update(delta float32) {

}
