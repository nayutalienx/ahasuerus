package models

import (
	"ahasuerus/collision"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type CollisionHitbox struct {
	BaseEditorItem
	CollisionProcessor collision.CollisionDetector
	hasCollision       bool
}

func (p *CollisionHitbox) Load() {

}

func (p *CollisionHitbox) Unload() {

}

func (p *CollisionHitbox) Pause()  {}
func (p *CollisionHitbox) Resume() {}

func (p *CollisionHitbox) Draw() {
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

func (p *CollisionHitbox) Update(delta float32) {

}

func (p CollisionHitbox) getDynamicHitbox() collision.Hitbox {
	topLeft := p.TopLeft()
	bottomRight := p.BottomRight()
	width := bottomRight.X - topLeft.X
	height := bottomRight.Y - topLeft.Y
	hb := GetDynamicHitboxFromMap(GetDynamicHitboxMap(topLeft, width, height))
	return hb
}
