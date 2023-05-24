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
	preset       func(i *Image)

	editSelected           bool
	editorMoveWithCursor   bool
	editorResizeWithCursor bool
}

func NewImage(drawIndex int, id string, path string, x, y, width, height float32) *Image {
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
	}
	return img
}

func (p *Image) Draw() {
	rl.DrawTexture(p.Texture, int32(p.Pos.X), int32(p.Pos.Y), rl.White)
	if p.editSelected {
		rl.DrawText(fmt.Sprintf("DrawIndex: %d", p.DrawIndex), int32(p.Pos.X), int32(p.Pos.Y), 40, rl.Red)
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

func (p *Image) EditorResolveSelect() (bool, bool) {
	rec := rl.NewRectangle(p.Pos.X, p.Pos.Y, float32(p.Texture.Width), float32(p.Texture.Height))
	mousePos := rl.GetMousePosition()
	collission := rl.CheckCollisionPointRec(mousePos, rec)
	if collission {
		rl.DrawRectangleLinesEx(rec, 3.0, rl.Red)

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			p.editSelected = true
		}
	}

	if p.editSelected {
		if rl.IsKeyDown(rl.KeyBackspace) {
			p.editSelected = false
		}
	}

	return p.editSelected, collission
}

func (p *Image) ProcessEditorSelection() bool {

	if p.editorMoveWithCursor {
		mousePos := rl.GetMousePosition()
		offset := 10
		p.Pos.X = mousePos.X - float32(offset)
		p.Pos.Y = mousePos.Y - float32(offset)
	}

	if p.editorResizeWithCursor {
		mousePos := rl.GetMousePosition()
		offset := 10
		p.Box.X = mousePos.X + float32(offset) - p.Pos.X
		p.Box.Y = mousePos.Y + float32(offset) - p.Pos.Y
		p.syncBoxWithTexture()
	}

	if (p.editorMoveWithCursor || p.editorResizeWithCursor) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		p.editorMoveWithCursor = false
		p.editorResizeWithCursor = false
		p.editSelected = false
		return true
	}

	return false
}

func (p *Image) Unload() {
	rl.UnloadTexture(p.Texture) // clear VRAM
}

func (p *Image) AfterLoadPreset(preset func(i *Image)) *Image {
	p.preset = preset
	return p
}

func (p *Image) syncBoxWithTexture() {
	p.Texture.Width = int32(p.Box.X)	
	p.Texture.Height = int32(p.Box.Y)	
}