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

type CollissionCallback func(pos BoxPosition)

type CollissionBox interface {
	ResolveCollission(callback CollissionCallback)
}

type Object interface {
	Draw()
	Update(delta float32)
}
