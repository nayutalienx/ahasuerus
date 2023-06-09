package models

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
	CursorX             int
	CursorY             int
}

type EditorItemDetectSelectionResult struct {
	Selected  bool
	Collision bool
}

type EditorItem interface {
	EditorDetectSelection() EditorItemDetectSelectionResult
	ProcessEditorSelection() EditorItemProcessSelectionResult
}

type EditorSelectedItem struct {
	Selected bool
	Item     EditorItem
}
