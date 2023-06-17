package models

import (
	"ahasuerus/collision"
	"ahasuerus/resources"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Image struct {
	BaseEditorItem
	Pos            rl.Vector2
	WidthHeight    rl.Vector2
	DrawIndex      int
	Texture        rl.Texture2D
	Shader         rl.Shader
	ImageTexture   resources.GameTexture
	ImageShader    resources.GameShader
	Parallax       float32
	parallaxOffset float32
	preset         func(i *Image)

	camera        *rl.Camera2D
	cameraLastPos rl.Vector2

	Scale float32

	isMoveMode   bool
	startMovePos rl.Vector2
	endMovePos   rl.Vector2
	moveSpeed    float32

	Lightboxes []Hitbox
	shaderLocs []int32
}

func NewImage(drawIndex int, id string, imageTexture resources.GameTexture, x, y, width, height, rotation float32) *Image {
	img := &Image{
		BaseEditorItem: BaseEditorItem{
			Id:       id,
			Rotation: rotation,
		},
		DrawIndex:    drawIndex,
		ImageTexture: imageTexture,
		Pos:          rl.NewVector2(x, y),
		WidthHeight:  rl.NewVector2(width, height),
		Lightboxes:   make([]Hitbox, 0),
		shaderLocs:   make([]int32, 0),
		Scale:        1,
	}
	if width != 0 && height != 0 {
		img.initEditorItem()
	}
	return img
}

func (p *Image) Camera(camera *rl.Camera2D) *Image {
	p.camera = camera
	p.cameraLastPos = camera.Target
	return p
}

func (p *Image) StartMove(startMovePos rl.Vector2, endMovePos rl.Vector2, moveSpeed float32) {
	p.isMoveMode = true
	p.startMovePos = startMovePos
	p.endMovePos = endMovePos
	p.moveSpeed = moveSpeed

	p.Pos = p.startMovePos
}

func (p *Image) Draw() {
	if p.ImageShader != resources.UndefinedShader {
		rl.BeginShaderMode(p.Shader)
		rl.DrawTextureEx(p.Texture, p.Pos, p.Rotation, p.Scale, rl.White)
		rl.EndShaderMode()
	} else {
		rl.DrawTextureEx(p.Texture, p.Pos, p.Rotation, p.Scale, rl.White)
	}
	if p.EditSelected {
		rl.DrawText(fmt.Sprintf("DrawIndex: %d", p.DrawIndex), int32(p.Pos.X), int32(p.Pos.Y), 40, rl.Red)
	}
	p.BaseEditorItem.Draw()
}

func (p *Image) Update(delta float32) {
	if p.isMoveMode {
		p.Pos = rl.Vector2Lerp(p.Pos, p.endMovePos, p.moveSpeed*delta)
		if p.Pos == p.endMovePos {
			p.isMoveMode = false
		}
	} else {
		p.Pos = p.TopLeft()
	}
	p.WidthHeight = rl.NewVector2(p.Width(), p.Height())
	p.syncBoxWithTexture()

	if p.ImageShader != resources.UndefinedShader {
		rl.SetShaderValueTexture(p.Shader, p.shaderLocs[0], p.Texture)
		rewind := 0.0
		if rl.IsKeyDown(rl.KeyLeftShift) {
			rewind = 1.0
		}
		rl.SetShaderValue(p.Shader, p.shaderLocs[1], []float32{float32(rewind)}, rl.ShaderUniformFloat)
	}
	if p.Parallax > 0 {
		delta := p.camera.Target.X - p.cameraLastPos.X
		p.parallaxOffset -= delta * p.Parallax
		p.cameraLastPos = p.camera.Target

		p.Pos.X += p.parallaxOffset
	}
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
		p.shaderLocs = []int32{
			rl.GetShaderLocation(p.Shader, "texture0"),
			rl.GetShaderLocation(p.Shader, "rewind"),
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
