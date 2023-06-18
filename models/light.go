package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Light struct {
	BaseEditorItem
}

func (p *Light) Draw() {
	if DRAW_MODELS {
		center := p.Center()
		rl.DrawCircleLines(int32(center.X), int32(center.Y), rl.Vector2Distance(p.TopLeft(), p.TopRight())/2, rl.Gold)
		rl.DrawCircleLines(int32(center.X), int32(center.Y), rl.Vector2Distance(p.TopLeft(), p.TopRight())/6, rl.Gold)
		rl.DrawCircle(int32(center.X), int32(center.Y), 10, rl.Gold)
	}
	p.BaseEditorItem.Draw()
}

func (p *Light) Update(delta float32) {

}
