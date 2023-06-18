package models

import (
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Particle struct {
	LifeRect rl.Rectangle
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

	shader     *rl.Shader `json:"-"`
	shaderLocs []int32    `json:"-"`
}

func NewParticle(
	lifeRect rl.Rectangle,
	texture rl.Texture2D,
	moveDirection rl.Vector2,
	lifetime time.Duration, moveSpeed, scale, rotation, maxOpacity float32,
) *Particle {
	lifetimeInFrames := float64(FPS) * lifetime.Seconds()

	fadeSpeed := maxOpacity / float32(lifetimeInFrames)

	p := Particle{
		LifeRect:      lifeRect,
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

func (p *Particle) WithShader(
	shader *rl.Shader,
	shaderLocs []int32,
) *Particle {
	p.shader = shader
	p.shaderLocs = shaderLocs
	return p
}

func (p *Particle) Update(delta float32) {

	if p.opacity <= 0 {
		p.opacityIncrement = true
		p.pos = p.randomStartPos()
	}

	if uint8(p.opacity) >= uint8(p.maxOpacity) {
		p.opacityIncrement = false
	}

	if p.opacityIncrement {
		p.opacity += p.fadeSpeed
	} else {
		p.opacity -= p.fadeSpeed
	}

	p.pos = rl.Vector2Add(p.pos, rl.Vector2Scale(p.moveDirection, p.moveSpeed))

	if p.shader != nil {
		rl.SetShaderValueTexture(*p.shader, p.shaderLocs[0], p.texture)
		rl.SetShaderValue(*p.shader, p.shaderLocs[1], []float32{float32(p.opacity/p.maxOpacity)}, rl.ShaderUniformFloat)
		rewind := 0.0
		if rl.IsKeyDown(rl.KeyLeftShift) {
			rewind = 1.0
		}
		rl.SetShaderValue(*p.shader, p.shaderLocs[2], []float32{float32(rewind)}, rl.ShaderUniformFloat)
	}

}

func (p *Particle) Draw() {
	if p.shader != nil {
		rl.BeginShaderMode(*p.shader)
		rl.DrawTextureEx(p.texture, p.pos, p.rotation, p.scale, rl.White)
		rl.EndShaderMode()
	} else {
		rl.DrawTextureEx(p.texture, p.pos, p.rotation, p.scale, rl.White)
	}
}

func (p *Particle) randomStartPos() rl.Vector2 {

	minX := int(p.LifeRect.X)
	maxX := int(p.LifeRect.X + p.LifeRect.Width)

	minY := int(p.LifeRect.Y)
	maxY := int(p.LifeRect.Y + p.LifeRect.Height)

	x := rand.Intn(maxX-minX+1) + minX
	y := rand.Intn(maxY-minY+1) + minY

	return rl.NewVector2(float32(x), float32(y))
}
