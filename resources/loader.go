package resources

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameTexture string

const (
	PlayerStay GameTexture = "resources/heroes/tim_stay.png"
	PlayerRun  GameTexture = "resources/heroes/tim_run.png"
	Girl1      GameTexture = "resources/heroes/girl1.png"
	Girl2      GameTexture = "resources/heroes/girl2.png"
	Girl3      GameTexture = "resources/heroes/girl3.png"
	Girl4      GameTexture = "resources/heroes/girl4.png"
	Girl5      GameTexture = "resources/heroes/girl5.png"
	GameBg     GameTexture = "resources/bg/1.jpg"
	MenuBg     GameTexture = "resources/bg/menu-bg.png"
	GameRoad   GameTexture = "resources/game/road.png"
)

type GameShader string

const (
	UndefinedShader   GameShader = ""
	BloomShader GameShader = "resources/shader/bloom.fs"
)

var (
	textureCache = make(map[GameTexture]rl.Texture2D)
	shaderCache  = make(map[GameShader]rl.Shader)
)

func LoadTexture(gameTexture GameTexture) rl.Texture2D {
	loadedTexture, ok := textureCache[gameTexture]
	if ok {
		fmt.Println("WARN: Texture already loaded. Using cache")
		return loadedTexture
	}
	img := rl.LoadImage(string(gameTexture)) // load img to RAM
	texture := rl.LoadTextureFromImage(img)  // move img to VRAM
	rl.UnloadImage(img)                      // clear ram
	textureCache[gameTexture] = texture
	return texture
}

func UnloadTexture(gameTexture GameTexture) {
	texture, ok := textureCache[gameTexture]
	if ok {
		rl.UnloadTexture(texture)
		delete(textureCache, gameTexture)
	} else {
		fmt.Println("WARN: Texture not found for unload")
	}
}

func LoadShader(gameShader GameShader) rl.Shader {
	loadedShader, ok := shaderCache[gameShader]
	if ok {
		fmt.Println("WARN: Shader already loaded. Using cache")
		return loadedShader
	}

	shader := rl.LoadShader("", string(gameShader))
	shaderCache[gameShader] = shader
	return shader
}

func UnloadShader(gameShader GameShader) {
	shader, ok := shaderCache[gameShader]
	if ok {
		rl.UnloadShader(shader)
		delete(shaderCache, gameShader)
	} else {
		fmt.Println("WARN: Shader not found for unload")
	}
}
