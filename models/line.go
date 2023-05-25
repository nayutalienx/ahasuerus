package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Line struct {
	Id string
	Start rl.Vector2
	End rl.Vector2

	Thick float32
	color rl.Color

	editSelected  bool
	editStartMode bool
	editEndMode   bool
}

func NewLine(id string, start, end rl.Vector2, thick float32, color rl.Color) *Line {
	return &Line{
		Id: id,
		Start:    start,
		End:    end,
		Thick: thick,
		color: color,
	}
}

func (p *Line) Draw() {
	rl.DrawLineEx(p.Start, p.End, p.Thick, p.color)
}

func (p *Line) Update(delta float32) {

}

func (p *Line) ResolveCollision(callback CollisionLineCallback) {
	callback(p)
}

func (p Line) GetColor() rl.Color {
	return p.color
}

func (p *Line) SetColor(col rl.Color) *Line {
	p.color = col
	return p
}

func (p *Line) SetStartModeTrue() {
	p.editStartMode = true
}

func (p *Line) SetEndModeTrue() {
	p.editEndMode = true
}

func (p *Line) ProcessEditorSelection() EditorItemProcessSelectionResult {
	if p.editStartMode {
		mousePos := rl.GetMousePosition()
		p.Start.X = mousePos.X + 20
		p.Start.Y = mousePos.Y + 20
	}

	if p.editEndMode {
		mousePos := rl.GetMousePosition()
		p.End.X = mousePos.X - 20
		p.End.Y = mousePos.Y - 20
	}

	if (p.editStartMode || p.editEndMode) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		p.editSelected = false
		p.editStartMode = false
		p.editEndMode = false
		return EditorItemProcessSelectionResult{
			Finished:      true,
		}
	}

	if p.editSelected {
		if rl.IsKeyDown(rl.KeyBackspace) {
			p.editSelected = false
			p.editStartMode = false
			p.editEndMode = false
			return EditorItemProcessSelectionResult{
				Finished:      true,
				DisableCursor: true,
				CursorForcePosition: true,
				CursorX: int(p.Start.X),
				CursorY: int(p.Start.Y),
			}
		}
	}

	return EditorItemProcessSelectionResult{
		Finished:      false,
	}
}

func (p *Line) EditorResolveSelect() (EditorItemResolveSelectionResult) {
	mousePos := rl.GetMousePosition()
	isCollision := rl.CheckCollisionPointLine(mousePos, p.Start, p.End, int32(p.Thick))
	if isCollision {
		p.color = rl.Red

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) && !p.editSelected {
			p.editSelected = true
		}

	} else {
		p.color = rl.Gold
	}
	return EditorItemResolveSelectionResult{
		Selected:  p.editSelected,
		Collision: isCollision,
	}
}