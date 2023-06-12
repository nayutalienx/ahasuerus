package models

import (
	"ahasuerus/config"
	"ahasuerus/resources"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Orientation int

const (
	Left  Orientation = 1
	Right Orientation = 2
)

type AnimationType int

const (
	Loop AnimationType = iota
	Temporary
)

var (
	FPS = config.GetFPS()
)

type Animation struct {
	Texture     rl.Texture2D
	Shader      rl.Shader
	Pos         rl.Vector2
	StepInPixel int32
	Orientation Orientation

	animationType AnimationType
	frame         rl.Rectangle
	currentFrame  int32
	framesCounter int32

	GameTexture     resources.GameTexture
	steps           int32
	framesPerSecond int32
	timeInSeconds   float32

	preset func(*Animation)
}

func NewAnimation(gameTexture resources.GameTexture, steps int32, animType AnimationType) *Animation {
	return &Animation{
		GameTexture:   gameTexture,
		steps:         steps,
		animationType: animType,
	}
}

func (a *Animation) Begin() {
	a.currentFrame = 0
	a.framesCounter = 0
}

func (a *Animation) FramesPerSecond(framesPerSecond int32) *Animation {
	a.framesPerSecond = framesPerSecond
	return a
}

func (a *Animation) TimeInSeconds(timeInSeconds float32) *Animation {
	a.timeInSeconds = timeInSeconds
	return a
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

	if a.animationType == Loop {
		if a.framesCounter >= FPS/a.framesPerSecond {
			a.framesCounter = 0
			a.currentFrame++
			if a.currentFrame > a.steps {
				a.currentFrame = 0
			}
		}
	}

	if a.animationType == Temporary {
		if a.framesCounter >= int32(a.timeInSeconds*float32(FPS)) {
			a.framesCounter = 0
			a.currentFrame = 0
		}

		offsetInTime := a.timeInSeconds / float32(a.steps)
		offsetInFrames := int32(offsetInTime * float32(FPS))

		if a.framesCounter%offsetInFrames == 0 {
			a.currentFrame++
		}

	}

	a.frame.X = float32(a.currentFrame) * float32(a.StepInPixel)
}

func (a *Animation) Stop() {
	a.currentFrame = 0
	a.frame.X = float32(a.StepInPixel)
}

func (a *Animation) Load() {
	a.Texture = resources.LoadTexture(a.GameTexture)
	a.StepInPixel = a.Texture.Width / a.steps
	a.frame = rl.NewRectangle(0, 0, float32(a.StepInPixel), float32(a.Texture.Height))
	if a.preset != nil {
		a.preset(a)
	}
}

func (a *Animation) Unload() {
	resources.UnloadTexture(a.GameTexture)
}

func (a *Animation) AfterLoadPreset(cb func(*Animation)) *Animation {
	a.preset = cb
	return a
}
