package resources

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type FontTtf string

const (
	Literata FontTtf = "resources/font/literata.ttf"
)

type GuiStyle string

const (
	Lavanda GuiStyle = "resources/styles/lavanda.rgs"
)

type GameTexture string

const (
	PlayerStayTexture       GameTexture = "resources/heroes/tim_stay.png"
	PlayerRunTexture        GameTexture = "resources/heroes/tim_run.png"
	PlayerDirectUpTexture   GameTexture = "resources/heroes/tim_direct_up.png"
	PlayerDirectDownTexture GameTexture = "resources/heroes/tim_direct_down.png"
	PlayerSideUpTexture     GameTexture = "resources/heroes/tim_side_up.png"
	PlayerSideDownTexture   GameTexture = "resources/heroes/tim_side_down.png"
)

type GameShader string

const (
	UndefinedShader GameShader = ""
	BloomShader     GameShader = "resources/shader/bloom.fs"
	BlurShader      GameShader = "resources/shader/blur.fs"
	TextureShader   GameShader = "resources/shader/texture.fs"
	PlayerShader    GameShader = "resources/shader/player.fs"
	NpcShader       GameShader = "resources/shader/npc.fs"
	SdfShader       GameShader = "resources/shader/sdf.fs"
)

var (
	textureCache = make(map[GameTexture]rl.Texture2D)
	fontsCache   = make(map[FontTtf]rl.Font)
	shaderCache  = make(map[GameShader]rl.Shader)
)

func LoadFont(f FontTtf) rl.Font {
	font, ok := fontsCache[f]
	if ok {
		return font
	}

	font = rl.LoadFont(string(f))
	fontsCache[f] = font
	return font
}

func UnloadFont(f FontTtf) {
	font, ok := fontsCache[f]
	if ok {
		rl.UnloadFont(font)
	}
}

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

func LoadShaderCache(gameShader GameShader) rl.Shader {
	sh, ok := shaderCache[GameShader(gameShader)]
	if ok {
		return sh
	}
	shaderCache[GameShader(gameShader)] = rl.LoadShader("", string(gameShader))
	return shaderCache[GameShader(gameShader)]
}

func UnloadShader(shader rl.Shader) {
	rl.UnloadShader(shader)
}

func UnloadShaderCache(shader GameShader) {
	sh, ok := shaderCache[GameShader(shader)]
	if ok {
		delete(shaderCache, shader)
		rl.UnloadShader(sh)
	}
}
