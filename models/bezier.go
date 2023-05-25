package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Bezier struct {
	Id    string
	Start rl.Vector2
	End   rl.Vector2

	Thick float32
	color rl.Color

	editSelected  bool
	editStartMode bool
	editEndMode   bool
}

func NewBezier(id string, start, end rl.Vector2, thick float32, color rl.Color) *Bezier {
	return &Bezier{
		Id:    id,
		Start: start,
		End:   end,
		Thick: thick,
		color: color,
	}
}

func (p *Bezier) Draw() {
	if p.editSelected {
		rl.DrawLineBezier(p.Start, p.End, p.Thick, rl.Red)
	} else {
		rl.DrawLineBezier(p.Start, p.End, p.Thick, p.color)
	}
}

func (p *Bezier) Update(delta float32) {

}

func (p *Bezier) ResolveCollision(callback CollisionBezierCallback) {
	callback(p)
}

func (p Bezier) GetColor() rl.Color {
	return p.color
}

func (p *Bezier) SetColor(col rl.Color) *Bezier {
	p.color = col
	return p
}

func (p *Bezier) SetStartModeTrue() {
	p.editStartMode = true
}

func (p *Bezier) SetEndModeTrue() {
	p.editEndMode = true
}

func (p *Bezier) ProcessEditorSelection() EditorItemProcessSelectionResult {
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

func (p *Bezier) EditorResolveSelect() (EditorItemResolveSelectionResult) {
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
