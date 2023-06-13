package models

import (
	"ahasuerus/collision"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type HitboxType int

const (
	Collision HitboxType = iota
	Light
	Npc
)

type Hitbox struct {
	BaseEditorItem
	Type               HitboxType
	CollisionProcessor collision.CollisionDetector

	hasCollision bool
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
			rl.DrawCircleLines(int32(center.X), int32(center.Y), rl.Vector2Distance(p.TopLeft(), p.TopRight())/6, rl.Gold)
			rl.DrawCircle(int32(center.X), int32(center.Y), 10, rl.Gold)
		}

		if p.Type == Npc {
			polys := p.getDynamicHitbox().Polygons
			for i, _ := range polys {
				rl.DrawTriangleLines(
					polys[i].Points[0],
					polys[i].Points[1],
					polys[i].Points[2],
					rl.Pink,
				)
			}
		}

		p.BaseEditorItem.Draw()
	}

	if p.Type == Npc {
		if p.hasCollision || p.EditSelected {
			pos := p.TopRight()

			offsetX := int32(p.PropertyFloat("blockOffsetX"))
			offsetY := int32(p.PropertyFloat("blockOffsetY"))
			width := int32(p.PropertyFloat("blockWidth"))
			height := int32(p.PropertyFloat("blockHeight"))

			outline := p.PropertyFloat("outlineThick")
			fontSize := int32(p.PropertyFloat("fontSize"))

			textOffsetX := p.PropertyFloat("textOffsetX")
			textOffsetY := p.PropertyFloat("textOffsetY")
			text := p.PropertyString("text")

			rl.DrawRectangle(int32(pos.X)+offsetX, int32(pos.Y)+offsetY, width, height, rl.Black)
			rl.DrawRectangleLinesEx(rl.NewRectangle((pos.X)+float32(offsetX), (pos.Y)+float32(offsetY), float32(width), float32(height)), outline, rl.White)
			rl.DrawText(text, int32(pos.X)+offsetX+int32(textOffsetX), int32(pos.Y)+offsetY+int32(textOffsetY), fontSize, rl.White)
		}
	}

}

func (p *Hitbox) Update(delta float32) {

	if p.Type == Npc {
		p.hasCollision, _ = p.CollisionProcessor.Detect(p.getDynamicHitbox())
	}

}

func (p Hitbox) getDynamicHitbox() collision.Hitbox {
	topLeft := p.TopLeft()
	bottomRight := p.BottomRight()
	width := bottomRight.X - topLeft.X
	height := bottomRight.Y - topLeft.Y
	hb := GetDynamicHitboxFromMap(GetDynamicHitboxMap(topLeft, width, height))
	return hb
}
