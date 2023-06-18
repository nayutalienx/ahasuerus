package models

import (
	"ahasuerus/collision"
	"ahasuerus/resources"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Image struct {
	BaseEditorItem
	Texture        rl.Texture2D `json:"-"`
	Shader         rl.Shader    `json:"-"`
	ImageTexture   resources.GameTexture
	ImageShader    resources.GameShader
	Parallax       float32
	parallaxOffset float32        `json:"-"`
	preset         func(i *Image) `json:"-"`

	camera        *rl.Camera2D `json:"-"`
	cameraLastPos rl.Vector2   `json:"-"`

	Scale float32

	isMoveMode   bool       `json:"-"`
	startMovePos rl.Vector2 `json:"-"`
	endMovePos   rl.Vector2 `json:"-"`
	moveSpeed    float32    `json:"-"`

	Lightboxes []CollisionHitbox `json:"-"`
	shaderLocs []int32           `json:"-"`
}

func NewImage(id string, imageTexture resources.GameTexture, x, y, rotation float32) *Image {
	topLeft := rl.NewVector2(x, y)
	bei := NewBaseEditorItem([2]collision.Polygon{
		{
			Points: [3]rl.Vector2{
				topLeft,
			},
		},
	})
	bei.Id = id
	bei.Rotation = rotation

	img := &Image{
		BaseEditorItem: bei,
		ImageTexture:   imageTexture,
		Lightboxes:     make([]CollisionHitbox, 0),
		ImageShader:    resources.TextureShader,
		shaderLocs:     make([]int32, 0),
		Scale:          1,
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
	p.ChangePosition(p.startMovePos)
}

func (p *Image) Draw() {
	if p.ImageShader != resources.UndefinedShader {
		rl.BeginShaderMode(p.Shader)
		rl.DrawTextureEx(p.Texture, p.TopLeft(), p.Rotation, p.Scale, rl.White)
		rl.EndShaderMode()
	} else {
		rl.DrawTextureEx(p.Texture, p.TopLeft(), p.Rotation, p.Scale, rl.White)
	}
	p.BaseEditorItem.Draw()
}

func (p *Image) Update(delta float32) {
	if p.isMoveMode {
		p.ChangePosition(rl.Vector2Lerp(p.TopLeft(), p.endMovePos, p.moveSpeed*delta))
		if p.TopLeft() == p.endMovePos {
			p.isMoveMode = false
		}
	}

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
		p.Translate(rl.NewVector2(p.parallaxOffset, 0))
	}
}

func (p *Image) Load() {

	p.Texture = resources.LoadTexture(p.ImageTexture)

	if p.preset != nil {
		p.preset(p)
	}

	width := float32(p.Texture.Width)
	height := float32(p.Texture.Height)

	pos := p.TopLeft()

	topLeft := rl.NewVector2(pos.X, pos.Y)
	bottomLeft := rl.Vector2{topLeft.X, topLeft.Y + height}

	topRight := rl.Vector2{topLeft.X + width, topLeft.Y}
	bottomRight := rl.Vector2{topLeft.X + width, topLeft.Y + height}

	p.BaseEditorItem.Polygons = [2]collision.Polygon{
		{
			Points: [3]rl.Vector2{
				topLeft, topRight, bottomRight,
			},
		},
		{
			Points: [3]rl.Vector2{
				topLeft, bottomLeft, bottomRight,
			},
		},
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

func (p *Image) AddLightbox(lp CollisionHitbox) *Image {
	p.Lightboxes = append(p.Lightboxes, lp)
	return p
}

func (p *Image) AfterLoadPreset(preset func(i *Image)) *Image {
	p.preset = preset
	return p
}

func (p Image) Replicate(id string, x, y float32) *Image {
	return NewImage(id, p.ImageTexture, x, y, p.Rotation)
}

func (img *Image) randomNotZero(n int) float32 {
	for {
		x := rand.Intn(n)
		if x > 0 {
			return float32(x)
		}
	}
}

func (i *Image) loadFromProperties() {

}
