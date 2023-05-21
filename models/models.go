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

type CollisionBoxCallback func(pos BoxPosition)

type CollisionBox interface {
	ResolveCollision(callback CollisionBoxCallback)
}

type CollisionBezierCallback func(bezier *Bezier)

type CollisionBezier interface {
	ResolveCollision(callback CollisionBezierCallback)
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
