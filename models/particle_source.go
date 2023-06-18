package models

import (
	"ahasuerus/resources"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ParticleSource struct {
	BaseEditorItem
	ParticleTexture    resources.GameTexture
	ParticleShader     resources.GameShader
	Amount             int
	BackgroundParticle bool
	texture            rl.Texture2D `json:"-"`

	particles []Particle `json:"-"`

	shader     *rl.Shader `json:"-"`
	shaderLocs []int32    `json:"-"`
}

func NewParticleSource(
	bei BaseEditorItem,
	particleTexture resources.GameTexture,
	particlesSize int) *ParticleSource {
	p := ParticleSource{
		BaseEditorItem:  bei,
		ParticleTexture: particleTexture,
		Amount:          particlesSize,
	}

	return &p
}

func (p *ParticleSource) Load() {
	p.texture = resources.LoadTexture(p.ParticleTexture)

	if p.ParticleShader != resources.UndefinedShader {
		sh := resources.LoadShader(p.ParticleShader)
		p.shader = &sh
		p.shaderLocs = []int32{
			rl.GetShaderLocation(*p.shader, "texture0"),
			rl.GetShaderLocation(*p.shader, "opacity"),
			rl.GetShaderLocation(*p.shader, "rewind"),
		}
	}

	p.loadParticleBuffer()
}

func (p *ParticleSource) Unload() {
	resources.UnloadTexture(p.ParticleTexture)

	if p.ParticleShader != resources.UndefinedShader {
		resources.UnloadShader(*p.shader)
	}
}

func (p *ParticleSource) Pause() {}

func (p *ParticleSource) Resume() {

}

func (p *ParticleSource) Draw() {
	if DRAW_MODELS {
		polys := p.PolygonsWithRotation()

		for i, _ := range polys {
			rl.DrawTriangleLines(
				polys[i].Points[0],
				polys[i].Points[1],
				polys[i].Points[2],
				rl.Brown,
			)
		}

		p.BaseEditorItem.Draw()
	}

	for i, _ := range p.particles {
		p.particles[i].Draw()
	}
}

func (p *ParticleSource) Update(delta float32) {
	for i, _ := range p.particles {
		if p.EditorMoveWithCursor {
			p.particles[i].LifeRect = rl.NewRectangle(p.TopLeft().X, p.TopLeft().Y, p.Width(), p.Height())
		}

		p.particles[i].Update(delta)
	}
}

func (p *ParticleSource) loadParticleBuffer() {
	p.particles = make([]Particle, 0)

	lifeRect := rl.NewRectangle(p.TopLeft().X, p.TopLeft().Y, p.Width(), p.Height())

	for i := 0; i < p.Amount; i++ {
		p.particles = append(p.particles, *NewParticle(
			lifeRect,
			p.texture,
			rl.NewVector2(0, -1),
			time.Second*5,
			0.5,
			0.5,
			0.0,
			100.0,
		).WithShader(p.shader, p.shaderLocs))
	}

}
