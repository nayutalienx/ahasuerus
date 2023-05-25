package models

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Image struct {
	DrawIndex    int
	Id           string
	Texture      rl.Texture2D
	Pos          rl.Vector2
	Box          rl.Vector2
	ResourcePath string
	Rotation     float32
	preset       func(i *Image)

	editSelected           bool
	editorMoveWithCursor   bool
	editorResizeWithCursor bool
	editorRotateMode       bool
}

func NewImage(drawIndex int, id string, path string, x, y, width, height, rotation float32) *Image {
	img := &Image{
		DrawIndex:    drawIndex,
		Id:           id,
		ResourcePath: path,
		Pos: rl.Vector2{
			X: x,
			Y: y,
		},
		Box: rl.Vector2{
			X: width,
			Y: height,
		},
		Rotation: rotation,
	}
	return img
}

func (p *Image) Draw() {
	rl.DrawTextureEx(p.Texture, p.Pos, p.Rotation, 1, rl.White)
	if p.editSelected {
		rl.DrawText(fmt.Sprintf("DrawIndex: %d", p.DrawIndex), int32(p.Pos.X), int32(p.Pos.Y), 40, rl.Red)
	}
	if p.editorRotateMode {
		rl.DrawText(fmt.Sprintf("Rotate on [R and T]: %.1f", p.Rotation), int32(p.Pos.X), int32(p.Pos.Y+40), 40, rl.Red)
	}
}

func (p *Image) Update(delta float32) {
}

func (p *Image) Load() {
	img := rl.LoadImage(p.ResourcePath)      // load img to RAM
	p.Texture = rl.LoadTextureFromImage(img) // move img to VRAM
	rl.UnloadImage(img)                      // clear ram
	if p.Box.X > 0 && p.Box.Y > 0 {          // scale image
		p.Texture.Width = int32(p.Box.X)
		p.Texture.Height = int32(p.Box.Y)
	} else {
		p.Box.X = float32(p.Texture.Width)
		p.Box.Y = float32(p.Texture.Height)
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

func (p *Image) EditorResolveSelect() EditorItemResolveSelectionResult {
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

	return EditorItemResolveSelectionResult{
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
	rl.UnloadTexture(p.Texture) // clear VRAM
}

func (p *Image) AfterLoadPreset(preset func(i *Image)) *Image {
	p.preset = preset
	return p
}

func (p Image) Replicate(id string, x, y float32) *Image {
	return NewImage(p.DrawIndex, id, p.ResourcePath, x, y, p.Box.X, p.Box.Y, p.Rotation)
}

func (p *Image) syncBoxWithTexture() {
	p.Texture.Width = int32(p.Box.X)
	p.Texture.Height = int32(p.Box.Y)
}
