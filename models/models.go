package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Position interface {
	GetPos() *rl.Vector2
}

type Box interface {
	GetBox() *rl.Vector2
}

type Scene interface {
	Run() Scene
	Unload()
}

type Object interface {
	Draw()
	Update(delta float32)
}

type Resource interface {
	Load()
	Pause()
	Resume()
	Unload()
}

type ObjectResource interface {
	Object
	Resource
}

type EditorItemProcessSelectionResult struct {
	Finished      bool
	EnableCursor  bool
	DisableCursor bool
	
	CursorForcePosition bool
	CursorX int
	CursorY int
}

type EditorItemResolveSelectionResult struct {
	Selected  bool
	Collision bool
}

type EditorItem interface {
	EditorResolveSelect() EditorItemResolveSelectionResult
	ProcessEditorSelection() EditorItemProcessSelectionResult
}

type EditorSelectedItem struct {
	Selected bool
	Item     EditorItem
}

// Collision interfaces
type Collision interface{}

