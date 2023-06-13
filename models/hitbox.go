package models

import (
	"ahasuerus/collision"
	"ahasuerus/resources"
	"fmt"
	"strings"

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

			outline := p.PropertyFloat("outlineThick")
			fontSize := int32(p.PropertyFloat("fontSize"))

			textOffsetX := p.PropertyFloat("textOffsetX")
			textOffsetY := p.PropertyFloat("textOffsetY")
			phrases := strings.Split(p.PropertyString("text"), ";")
			textCounter := int32(p.PropertyFloat("textCounter"))

			text := "empty phrase"
			if len(phrases) > int(textCounter) {
				text = phrases[textCounter]
			}

			maxXLen := 0
			splittenByNewLine := strings.Split(text, "\n")
			for i, _ := range splittenByNewLine {
				if len(splittenByNewLine[i]) > maxXLen {
					maxXLen = len(splittenByNewLine[i])
				}
			}

			width := int32(maxXLen * int(float64(fontSize)/1.5))
			height := int32(float64(fontSize)+(float64(fontSize)/1.5)) * (1 + (int32(strings.Count(text, "\n"))))

			rl.DrawRectangle(int32(pos.X)+offsetX, int32(pos.Y)+offsetY, width, height, rl.Black)
			rl.DrawRectangleLinesEx(rl.NewRectangle((pos.X)+float32(offsetX), (pos.Y)+float32(offsetY), float32(width), float32(height)), outline, rl.White)
			// rl.DrawText(text, int32(pos.X)+offsetX+int32(textOffsetX), int32(pos.Y)+offsetY+int32(textOffsetY), fontSize, rl.White)
			rl.DrawTextEx(resources.LoadFont(resources.Literata), text, rl.NewVector2(float32(int32(pos.X)+offsetX+int32(textOffsetX)), float32(int32(pos.Y)+offsetY+int32(textOffsetY))), float32(fontSize), 2, rl.White )
		}
	}

}

func (p *Hitbox) Update(delta float32) {

	if p.Type == Npc {
		p.hasCollision, _ = p.CollisionProcessor.Detect(p.getDynamicHitbox())

		if p.hasCollision {
			if rl.IsKeyReleased(rl.KeyEnter) {
				phrases := strings.Split(p.PropertyString("text"), ";")
				textCounter := int32(p.PropertyFloat("textCounter"))
				if len(phrases) > int(textCounter+1) {
					p.Properties["textCounter"] = fmt.Sprintf("%.1f", float32(textCounter+1))
				}
			}

		}

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
