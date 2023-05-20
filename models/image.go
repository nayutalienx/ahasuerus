package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Image struct {
	Texture      rl.Texture2D
	Pos          rl.Vector2
	resourcePath string
	scale float32
	preset func (i *Image)
}

func NewImage(path string, x, y float32) *Image {
	return &Image{
		resourcePath: path,
		Pos: rl.Vector2{
			X: x,
			Y: y,
		},
	}
}

func (p *Image) Scale(scale float32) *Image {
	p.scale = scale
	return p
}

func (p *Image) Draw() {
	rl.DrawTexture(p.Texture, int32(p.Pos.X), int32(p.Pos.Y), rl.White)
}

func (p *Image) Update(delta float32) {
}

func (p *Image) Load() {
	img := rl.LoadImage(p.resourcePath)        // load img to RAM
	p.Texture = rl.LoadTextureFromImage(img) // move img to VRAM
	rl.UnloadImage(img)                        // clear ram
	if p.scale > 0 { // scale image
		p.Texture.Width = int32(float32(p.Texture.Width)*p.scale)
		p.Texture.Height = int32(float32(p.Texture.Height)*p.scale)
	}
	if p.preset != nil {
		p.preset(p)
	}
}

func (p *Image) Unload() {
	rl.UnloadTexture(p.Texture) // clear VRAM
}

func (p *Image) AfterLoadPreset(preset func (i *Image)) *Image {
	p.preset = preset
	return p
}
