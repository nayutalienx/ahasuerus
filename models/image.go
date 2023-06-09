package models

import (
	"ahasuerus/resources"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Image struct {
	DrawIndex    int
	Id           string
	Texture      rl.Texture2D
	Shader       rl.Shader
	Pos          rl.Vector2
	Box          rl.Vector2
	ImageTexture resources.GameTexture
	ImageShader  resources.GameShader
	Rotation     float32
	preset       func(i *Image)

	LightPoints []*LightPoint
	shaderLocs  []int32

	editSelected           bool
	editorMoveWithCursor   bool
	editorResizeWithCursor bool
	editorRotateMode       bool
}

func NewImage(drawIndex int, id string, imageTexture resources.GameTexture, x, y, width, height, rotation float32) *Image {
	img := &Image{
		DrawIndex:    drawIndex,
		Id:           id,
		ImageTexture: imageTexture,
		Pos: rl.Vector2{
			X: x,
			Y: y,
		},
		Box: rl.Vector2{
			X: width,
			Y: height,
		},
		Rotation:    rotation,
		LightPoints: make([]*LightPoint, 0),
		shaderLocs:  make([]int32, 0),
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
	if p.editSelected {
		rl.DrawText(fmt.Sprintf("DrawIndex: %d", p.DrawIndex), int32(p.Pos.X), int32(p.Pos.Y), 40, rl.Red)
	}
	if p.editorRotateMode {
		rl.DrawText(fmt.Sprintf("Rotate on [R and T]: %.1f", p.Rotation), int32(p.Pos.X), int32(p.Pos.Y+40), 40, rl.Red)
	}
}

func (p *Image) Update(delta float32) {
	if p.ImageShader == resources.TextureLightShader {
		lightPoints := make([]float32, 0)
		for i, _ := range p.LightPoints {
			lp := p.LightPoints[i]
			lightPoints = append(lightPoints, float32(lp.Pos.X), float32(lp.Pos.Y))
		}
		rl.SetShaderValue(p.Shader, p.shaderLocs[0], []float32{p.Pos.X, p.Pos.Y + p.Box.Y}, rl.ShaderUniformVec2)
		rl.SetShaderValue(p.Shader, p.shaderLocs[1], []float32{p.Box.X, p.Box.Y}, rl.ShaderUniformVec2)
		rl.SetShaderValueV(p.Shader, p.shaderLocs[2], lightPoints, rl.ShaderUniformVec2, int32(len(p.LightPoints)))
		rl.SetShaderValue(p.Shader, p.shaderLocs[3], []float32{float32(len(p.LightPoints))}, rl.ShaderUniformFloat)
	}
}

func (p *Image) Load() {
	p.Texture = resources.LoadTexture(p.ImageTexture)
	if p.Box.X > 0 && p.Box.Y > 0 { // scale image
		p.Texture.Width = int32(p.Box.X)
		p.Texture.Height = int32(p.Box.Y)
	} else {
		p.Box.X = float32(p.Texture.Width)
		p.Box.Y = float32(p.Texture.Height)
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

	if p.preset != nil {
		p.preset(p)
	}
}

func (p *Image) Resume() {

}

func (p *Image) Pause() {

}

func (p *Image) SetEditorMoveWithCursorTrue() {
	p.editorMoveWithCursor = true
}

func (p *Image) SetEditorResizeWithCursorTrue() {
	p.editorResizeWithCursor = true
}

func (p *Image) SetEditorRotateModeTrue() {
	p.editorRotateMode = true
}

func (p *Image) EditorDetectSelection() EditorItemDetectSelectionResult {
	mousePos := rl.GetMousePosition()
	triangle1 := []rl.Vector2{p.Pos, rl.NewVector2(p.Pos.X+p.Box.X, p.Pos.Y), rl.NewVector2(p.Pos.X+p.Box.X, p.Pos.Y+p.Box.Y)}
	triangle2 := []rl.Vector2{p.Pos, rl.NewVector2(p.Pos.X, p.Pos.Y+p.Box.Y), rl.NewVector2(p.Pos.X+p.Box.X, p.Pos.Y+p.Box.Y)}
	RotateTriangleByA(&triangle1[0], &triangle1[1], &triangle1[2], float64(p.Rotation))
	RotateTriangleByA(&triangle2[0], &triangle2[1], &triangle2[2], float64(p.Rotation))
	collission := rl.CheckCollisionPointTriangle(mousePos, triangle1[0], triangle1[1], triangle1[2]) ||
		rl.CheckCollisionPointTriangle(mousePos, triangle2[0], triangle2[1], triangle2[2])
	if collission {
		rl.DrawTriangleLines(triangle1[0], triangle1[1], triangle1[2], rl.Red)
		rl.DrawTriangleLines(triangle2[0], triangle2[1], triangle2[2], rl.Red)

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			p.editSelected = true
		}
	}

	return EditorItemDetectSelectionResult{
		Selected:  p.editSelected,
		Collision: collission,
	}
}

func (p *Image) ProcessEditorSelection() EditorItemProcessSelectionResult {

	if p.editorMoveWithCursor {
		mousePos := rl.GetMousePosition()
		rl.DrawCircle(int32(mousePos.X), int32(mousePos.Y), 10, rl.Red)
		offset := 10
		p.Pos.X = mousePos.X - float32(offset)
		p.Pos.Y = mousePos.Y - float32(offset)
	}

	if p.editorResizeWithCursor {
		mousePos := rl.GetMousePosition()
		rl.DrawCircle(int32(mousePos.X), int32(mousePos.Y), 10, rl.Red)
		offset := 10
		p.Box.X = mousePos.X + float32(offset) - p.Pos.X
		p.Box.Y = mousePos.Y + float32(offset) - p.Pos.Y
		p.syncBoxWithTexture()
	}

	if p.editorRotateMode {
		if rl.IsKeyDown(rl.KeyT) {
			p.Rotation++
		}
		if rl.IsKeyDown(rl.KeyR) {
			p.Rotation--
		}
		if p.Rotation < 0 {
			p.Rotation = 360
		}
		if p.Rotation > 360 {
			p.Rotation = 0
		}
	}

	if (p.editorMoveWithCursor || p.editorResizeWithCursor) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		p.editorMoveWithCursor = false
		p.editorResizeWithCursor = false
		p.editSelected = false
		p.editorRotateMode = false
		return EditorItemProcessSelectionResult{
			Finished: true,
		}
	}

	if p.editSelected {
		if rl.IsKeyDown(rl.KeyBackspace) {
			p.editorMoveWithCursor = false
			p.editorResizeWithCursor = false
			p.editSelected = false
			p.editorRotateMode = false
			return EditorItemProcessSelectionResult{
				Finished:            true,
				DisableCursor:       true,
				CursorForcePosition: true,
				CursorX:             int(p.Pos.X),
				CursorY:             int(p.Pos.Y),
			}
		}
	}

	return EditorItemProcessSelectionResult{
		Finished: false,
	}
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

func (p *Image) AddLightPoint(lp *LightPoint) *Image {
	p.LightPoints = append(p.LightPoints, lp)
	return p
}

func (p *Image) AfterLoadPreset(preset func(i *Image)) *Image {
	p.preset = preset
	return p
}

func (p Image) Replicate(id string, x, y float32) *Image {
	return NewImage(p.DrawIndex, id, p.ImageTexture, x, y, p.Box.X, p.Box.Y, p.Rotation)
}

func (p *Image) syncBoxWithTexture() {
	p.Texture.Width = int32(p.Box.X)
	p.Texture.Height = int32(p.Box.Y)
}
