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

type BoxPosition interface {
	Box
	Position
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

type EditorItem interface {
	ReactOnCollision()
}

// Collision interfaces

type CollisionBoxCallback func(pos BoxPosition)
type CollisionBezierCallback func(bezier *Bezier)
type CollisionLineCallback func(line *Line)

type Collision interface {}

type CollisionBox interface {
	Collision
	ResolveCollision(callback CollisionBoxCallback)
}
type CollisionBezier interface {
	Collision
	ResolveCollision(callback CollisionBezierCallback)
}
type CollisionLine interface {
	Collision
	ResolveCollision(callback CollisionLineCallback)
}
