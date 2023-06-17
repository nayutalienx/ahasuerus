package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Particle struct {
	startPos rl.Vector2
	pos      rl.Vector2

	texture    rl.Texture2D
	scale      float32
	rotation   float32
	maxOpacity float32

	fadeSpeed float32
	moveSpeed float32

	opacity          float32
	opacityIncrement bool

	moveDirection rl.Vector2
}

func NewParticle(
	startPos rl.Vector2,
	texture rl.Texture2D,
	moveDirection rl.Vector2,
	fadeSpeed, moveSpeed, scale, rotation, maxOpacity float32,
) *Particle {
	p := Particle{
		startPos:      startPos,
		pos:           startPos,
		texture:       texture,
		maxOpacity:    maxOpacity,
		moveDirection: moveDirection,
		fadeSpeed:     fadeSpeed,
		moveSpeed:     moveSpeed,
		scale:         scale,
		rotation:      rotation,
	}

	return &p
}

func (p *Particle) Update(delta float32) {

	if p.opacity <= 0 {
		p.opacityIncrement = true
		p.pos = p.startPos
	}

	if uint8(p.opacity) >= uint8(p.maxOpacity) {
		p.opacityIncrement = false
	}

	if p.opacityIncrement {
		p.opacity += p.fadeSpeed * delta
	} else {
		p.opacity -= p.fadeSpeed * delta
	}

	p.pos = rl.Vector2Add(p.pos, rl.Vector2Scale(p.moveDirection, p.moveSpeed))

}

func (p *Particle) Draw() {
	rl.DrawTextureEx(p.texture, p.pos, p.rotation, p.scale, rl.White)
}
