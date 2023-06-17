package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ParticleType int

const (
	Circle ParticleType = iota
)

type Particle struct {
	Pos    rl.Vector2
	Color  rl.Color
	Type   ParticleType
	Radius float32

	FadeSpeed float32
	FallSpeed float32

	opacity          float32
	opacityIncrement bool
	yOffset          float32
}

func (p *Particle) Update(delta float32) {

	if p.opacity <= 0 {
		p.opacityIncrement = true
		p.yOffset = 0
	}

	if uint8(p.opacity) >= p.Color.A {
		p.opacityIncrement = false
	}

	if p.opacityIncrement {
		p.opacity += p.FadeSpeed * delta
	} else {
		p.opacity -= p.FadeSpeed * delta
	}

	p.yOffset += p.FallSpeed * delta

}

func (p *Particle) Draw() {
	c := rl.NewColor(
		p.Color.R,
		p.Color.G,
		p.Color.B,
		uint8(p.opacity),
	)

	rl.DrawCircle(int32(p.Pos.X), int32(p.Pos.Y)+int32(p.yOffset), p.Radius, c)
}
