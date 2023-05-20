package models

import (
	"ahasuerus/config"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Orientation int

const (
	Left  Orientation = 1
	Right Orientation = 2
)

var (
	FPS = config.GetFPS()
)

type Animation struct {
	Texture     rl.Texture2D
	Pos         rl.Vector2
	StepInPixel int32
	Orientation     Orientation	

	frame         rl.Rectangle
	currentFrame  int32
	framesCounter int32

	texturePath     string
	steps           int32
	framesPerSecond int32
}

func NewAnimation(texturePath string, steps int32, framesPerSecond int32) *Animation {
	return &Animation{
		texturePath:     texturePath,
		steps:           steps,
		framesPerSecond: framesPerSecond,
	}
}

func (a Animation) Draw() {
	rl.DrawTextureRec(a.Texture, a.frame, a.Pos, rl.White)
}

func (a *Animation) Update(delta float32) {

	if a.Orientation == Left {
		a.frame.Width = (-1) * float32(a.StepInPixel)
	} else {
		a.frame.Width = float32(a.StepInPixel)
	}

	a.framesCounter++

	if a.framesCounter >= FPS/a.framesPerSecond {
		a.framesCounter = 0
		a.currentFrame++
		if a.currentFrame > a.steps {
			a.currentFrame = 0
		}
		a.frame.X = float32(a.currentFrame) * float32(a.StepInPixel)
	}
}

func (a *Animation) Stop() {
	a.currentFrame = 0
	a.frame.X = float32(a.StepInPixel)
}

func (a *Animation) Load() {
	a.Texture = rl.LoadTexture(a.texturePath)
	a.StepInPixel = a.Texture.Width / a.steps
	a.frame = rl.NewRectangle(0, 0, float32(a.StepInPixel), float32(a.Texture.Height))
}

func (a *Animation) Unload() {
	rl.UnloadTexture(a.Texture)
}
