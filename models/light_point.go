package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type LightPoint struct {
	Pos      rl.Vector2
	move     bool
	moveBack bool
	StartPos rl.Vector2
	EndPos   rl.Vector2
	Speed    float32
}

func NewLightPoint(pos rl.Vector2) *LightPoint {
	return &LightPoint{
		Pos: pos,
	}
}

func (p *LightPoint) Draw() {
	if DRAW_MODELS {
		rl.DrawCircle(int32(p.Pos.X), int32(p.Pos.Y), 10, rl.Yellow)
	}
}

func (p *LightPoint) Update(delta float32) {
	if p.move {
		if !p.moveBack {
			p.Pos.X += p.Speed
			if p.Pos.X >= p.EndPos.X {
				p.moveBack = true
			}
		} else {
			p.Pos.X -= p.Speed
			if p.Pos.X <= p.StartPos.X {
				p.moveBack = false
			}
		}
	}
}

func (p *LightPoint) Dynamic(start rl.Vector2, end rl.Vector2, speed float32) *LightPoint {
	p.move = true
	p.StartPos = start
	p.EndPos = end
	p.Speed = speed
	return p
}
