package models

import (
	"ahasuerus/resources"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawSdfText(text string, textPos rl.Vector2, fontSize float32, color rl.Color) {
	rl.BeginShaderMode(resources.LoadShaderCache(resources.SdfShader))
	rl.DrawTextEx(resources.LoadFont(resources.Literata), text, textPos, (fontSize), 2, color)
	rl.EndShaderMode()
}
