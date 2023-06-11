package models

import (
	"ahasuerus/collision"
	"ahasuerus/resources"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Image struct {
	BaseEditorItem
	Pos          rl.Vector2
	WidthHeight  rl.Vector2
	DrawIndex    int
	Id           string
	Texture      rl.Texture2D
	Shader       rl.Shader
	ImageTexture resources.GameTexture
	ImageShader  resources.GameShader
	preset       func(i *Image)

	Lightboxes []Hitbox
	shaderLocs []int32
}

func NewImage(drawIndex int, id string, imageTexture resources.GameTexture, x, y, width, height, rotation float32) *Image {
	img := &Image{
		BaseEditorItem: BaseEditorItem{
			Rotation: rotation,
		},
		DrawIndex:    drawIndex,
		Id:           id,
		ImageTexture: imageTexture,
		Pos:          rl.NewVector2(x, y),
		WidthHeight:  rl.NewVector2(width, height),
		Lightboxes:   make([]Hitbox, 0),
		shaderLocs:   make([]int32, 0),
	}
	if width != 0 && height != 0 {
		img.initEditorItem()
	}
	return img
}

func (p *Image) Draw() {
	if p.ImageShader != resources.UndefinedShader {
		rl.BeginShaderMode(p.Shader)
		rl.DrawTextureEx(p.Texture, p.Pos, p.Rotation, 1, rl.White)
		rl.EndShaderMode()
	} else {
		rl.DrawTextureEx(p.Texture, p.Pos, p.Rotation, 1, rl.White)
	}
	if p.EditSelected {
		rl.DrawText(fmt.Sprintf("DrawIndex: %d", p.DrawIndex), int32(p.Pos.X), int32(p.Pos.Y), 40, rl.Red)
	}
	p.BaseEditorItem.Draw()
}

func (p *Image) Update(delta float32) {
	p.Pos = p.TopLeft()
	p.WidthHeight = rl.NewVector2(p.Width(), p.Height())
	if p.ImageShader == resources.TextureLightShader {
		lightPoints := make([]float32, 0)
		for i, _ := range p.Lightboxes {
			lp := p.Lightboxes[i]
			lightPoints = append(lightPoints, float32(lp.Center().X), float32(lp.Center().Y))
		}
		rl.SetShaderValue(p.Shader, p.shaderLocs[0], []float32{p.Pos.X, p.Pos.Y + p.WidthHeight.Y}, rl.ShaderUniformVec2)
		rl.SetShaderValue(p.Shader, p.shaderLocs[1], []float32{p.WidthHeight.X, p.WidthHeight.Y}, rl.ShaderUniformVec2)
		rl.SetShaderValueV(p.Shader, p.shaderLocs[2], lightPoints, rl.ShaderUniformVec2, int32(len(p.Lightboxes)))
		rl.SetShaderValue(p.Shader, p.shaderLocs[3], []float32{float32(len(p.Lightboxes))}, rl.ShaderUniformFloat)
	}
	p.syncBoxWithTexture()
}

func (p *Image) Load() {

	shouldInitEditorItem := p.WidthHeight.X == 0 && p.WidthHeight.Y == 0

	p.Texture = resources.LoadTexture(p.ImageTexture)

	if p.preset != nil {
		p.preset(p)
	}

	if p.WidthHeight.X > 0 && p.WidthHeight.Y > 0 { // scale image
		p.Texture.Width = int32(p.WidthHeight.X)
		p.Texture.Height = int32(p.WidthHeight.Y)
	} else {
		p.WidthHeight.X = float32(p.Texture.Width)
		p.WidthHeight.Y = float32(p.Texture.Height)
	}

	if shouldInitEditorItem {
		p.initEditorItem()
	}

	if p.ImageShader != resources.UndefinedShader {
		p.Shader = resources.LoadShader(p.ImageShader)
		textureLoc := rl.GetShaderLocation(p.Shader, "texture0")
		rl.SetShaderValueTexture(p.Shader, textureLoc, p.Texture)

		if p.ImageShader == resources.TextureLightShader {
			p.shaderLocs = []int32{
				rl.GetShaderLocation(p.Shader, "objectPosBottomLeft"),
				rl.GetShaderLocation(p.Shader, "objectSize"),
				rl.GetShaderLocation(p.Shader, "lightPos"),
				rl.GetShaderLocation(p.Shader, "lightPosSize"),
			}
		}

	}
}

func (p *Image) Resume() {

}

func (p *Image) Pause() {

}

func (p *Image) Unload() {
	resources.UnloadTexture(p.ImageTexture)
	if p.ImageShader != resources.UndefinedShader {
		resources.UnloadShader(p.Shader)
	}
}

func (p *Image) WithShader(gs resources.GameShader) *Image {
	p.ImageShader = gs
	return p
}

func (p *Image) AddLightbox(lp Hitbox) *Image {
	p.Lightboxes = append(p.Lightboxes, lp)
	return p
}

func (p *Image) AfterLoadPreset(preset func(i *Image)) *Image {
	p.preset = preset
	return p
}

func (p Image) Replicate(id string, x, y float32) *Image {
	return NewImage(p.DrawIndex, id, p.ImageTexture, x, y, p.WidthHeight.X, p.WidthHeight.Y, p.Rotation)
}

func (p *Image) syncBoxWithTexture() {
	p.Texture.Width = int32(p.WidthHeight.X)
	p.Texture.Height = int32(p.WidthHeight.Y)
}

func (img *Image) initEditorItem() {
	img.BaseEditorItem.SetPolygons([2]collision.Polygon{
		{
			Points: [3]rl.Vector2{
				img.Pos, {img.Pos.X + img.WidthHeight.X, img.Pos.Y}, {img.Pos.X + img.WidthHeight.X, img.Pos.Y + img.WidthHeight.Y},
			},
		},
		{
			Points: [3]rl.Vector2{
				img.Pos, {img.Pos.X, img.Pos.Y + img.WidthHeight.Y}, {img.Pos.X + img.WidthHeight.X, img.Pos.Y + img.WidthHeight.Y},
			},
		},
	})
}
