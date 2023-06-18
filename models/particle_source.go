package models

import (
	"ahasuerus/resources"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ParticleSource struct {
	BaseEditorItem

	particleTexture resources.GameTexture `json:"-"`
	texture         rl.Texture2D          `json:"-"`

	particles []Particle `json:"-"`
}

func NewParticleSource(
	bei BaseEditorItem,
	particleTexture resources.GameTexture,
	particlesSize int) *ParticleSource {
	p := ParticleSource{
		BaseEditorItem:  bei,
		particleTexture: particleTexture,
		particles:       make([]Particle, 0),
	}
	p.texture = resources.LoadTexture(p.particleTexture)

	startPos := bei.Center()

	for i := 0; i < particlesSize; i++ {
		p.particles = append(p.particles, *NewParticle(
			startPos,
			p.texture,
			rl.NewVector2(0, -1),
			100.0,
			2.5,
			0.5,
			0.0,
			100.0,
		))
	}

	return &p
}

func (p *ParticleSource) Load() {

}

func (p *ParticleSource) Unload() {
	resources.UnloadTexture(p.particleTexture)
}

func (p *ParticleSource) Draw() {
	if DRAW_MODELS {
		polys := p.PolygonsWithRotation()

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

	for i, _ := range p.particles {
		p.particles[i].Draw()
	}
}

func (p *ParticleSource) Update(delta float32) {
	for i, _ := range p.particles {
		p.particles[i].Update(delta)
	}
}
