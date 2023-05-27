package resources

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameTexture string

const (
	PlayerStayTexture GameTexture = "resources/heroes/tim_stay.png"
	PlayerRunTexture  GameTexture = "resources/heroes/tim_run.png"
	Girl1Texture      GameTexture = "resources/heroes/girl1.png"
	Girl2Texture      GameTexture = "resources/heroes/girl2.png"
	Girl3Texture      GameTexture = "resources/heroes/girl3.png"
	Girl4Texture      GameTexture = "resources/heroes/girl4.png"
	Girl5Texture      GameTexture = "resources/heroes/girl5.png"
	GameBgTexture     GameTexture = "resources/bg/1.jpg"
	MenuBgTexture     GameTexture = "resources/bg/menu-bg.png"
	GameRoadTexture   GameTexture = "resources/game/road.png"
)

type GameShader string

const (
	UndefinedShader GameShader = ""
	BloomShader     GameShader = "resources/shader/bloom.fs"
	BlurShader      GameShader = "resources/shader/blur.fs"
	TextureShader   GameShader = "resources/shader/texture.fs"
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
	fmt.Println("INFO: Load shader " + gameShader)
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
	fmt.Println("INFO: Unload shader " + gameShader)
	shader, ok := shaderCache[gameShader]
	if ok {
		rl.UnloadShader(shader)
		delete(shaderCache, gameShader)
	} else {
		fmt.Println("WARN: Shader not found for unload")
	}
}
