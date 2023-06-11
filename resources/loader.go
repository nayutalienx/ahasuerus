package resources

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GuiStyle string

const (
	Lavanda GuiStyle = "resources\\styles\\lavanda.rgs"
)

type GameTexture string

const (
	PlayerStayTexture GameTexture = "resources/heroes/tim_stay.png"
	PlayerRunTexture  GameTexture = "resources/heroes/tim_run.png"
)

type GameShader string

const (
	UndefinedShader    GameShader = ""
	BloomShader        GameShader = "resources/shader/bloom.fs"
	BlurShader         GameShader = "resources/shader/blur.fs"
	TextureShader      GameShader = "resources/shader/texture.fs"
	TextureLightShader GameShader = "resources/shader/texture_light.fs"
)

var (
	textureCache = make(map[GameTexture]rl.Texture2D)
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
	shader := rl.LoadShader("", string(gameShader))
	return shader
}

func UnloadShader(shader rl.Shader) {
	rl.UnloadShader(shader)
}
