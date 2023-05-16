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

type CollisionCallback func(pos BoxPosition)

type CollisionBox interface {
	ResolveCollision(callback CollisionCallback)
}

type Object interface {
	Draw()
	Update(delta float32)
}
