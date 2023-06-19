package models

import (
	"ahasuerus/particle"
	"ahasuerus/resources"
	"time"

	"github.com/blizzy78/twodeeparticles"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ParticleSourceType string

const (
	Bubble   ParticleSourceType = "bubble"
	Fountain ParticleSourceType = "fountain"
	Vortex   ParticleSourceType = "vortex"
)

type ParticleSource struct {
	BaseEditorItem
	ParticleTexture resources.GameTexture
	ParticleShader  resources.GameShader
	Type            ParticleSourceType

	texture        rl.Texture2D `json:"-"`
	SystemSettings particle.ParticleSystemSettings
	system         *twodeeparticles.ParticleSystem `json:"-"`

	shader     *rl.Shader `json:"-"`
	shaderLocs []int32    `json:"-"`
}

func NewParticleSource(
	bei BaseEditorItem) *ParticleSource {
	p := ParticleSource{
		BaseEditorItem: bei,
		ParticleShader: resources.ParticleShader,
		Type:           Bubble,
		SystemSettings: particle.DefaultParticleSystemSettings(),
	}

	return &p
}

func (p *ParticleSource) Load() {
	if p.ParticleTexture != "" {
		p.texture = resources.LoadTexture(p.ParticleTexture)
	}

	if p.ParticleShader != resources.UndefinedShader {
		sh := resources.LoadShader(p.ParticleShader)
		p.shader = &sh
		p.shaderLocs = []int32{
			rl.GetShaderLocation(*p.shader, "texture0"),
			rl.GetShaderLocation(*p.shader, "opacity"),
			rl.GetShaderLocation(*p.shader, "rewind"),
		}
	}
	if p.Type == Bubble {
		p.system = p.SystemSettings.Bubbles()
	}
	if p.Type == Fountain {
		p.system = p.SystemSettings.Fountain()
	}
	if p.Type == Vortex {
		p.system = p.SystemSettings.Vortex()
	}
}

func (p *ParticleSource) Unload() {
	if p.ParticleTexture != "" {
		resources.UnloadTexture(p.ParticleTexture)
	}

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
				rl.Violet,
			)
		}

		p.BaseEditorItem.Draw()
	}

	p.system.ForEachParticle(func(particle *twodeeparticles.Particle, t twodeeparticles.NormalizedDuration, delta time.Duration) {

		particlePos := particle.Position()
		translatedPos := rl.Vector2Add(rl.NewVector2(float32(particlePos.X), float32(particlePos.Y)), p.Center())

		_, _, _, a := particle.Color().RGBA()

		col := rl.Orange
		if rl.IsKeyDown(rl.KeyLeftShift) {
			col = rl.Gray
		}

		col.A = uint8(a)

		rl.DrawCircle(int32(translatedPos.X), int32(translatedPos.Y), 10, col)

	}, time.Now())

}

func (p *ParticleSource) Update(delta float32) {
	p.system.Update(time.Now())
}
