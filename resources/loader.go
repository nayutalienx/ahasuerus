package resources

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

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

type Ellipse struct {
	XMLName xml.Name `xml:"ellipse"`
	Color   string   `xml:"fill,attr"`
	X       string   `xml:"cx,attr"`
	Y       string   `xml:"cy,attr"`
	Radius  string   `xml:"rx,attr"`
}

type Geometry struct {
	XMLName   xml.Name  `xml:"g"`
	Transform string    `xml:"transform,attr"`
	Ellipses  []Ellipse `xml:"ellipse"`
}

type Particles struct {
	XMLName  xml.Name `xml:"svg"`
	Geometry Geometry `xml:"g"`
}

type ParticlesDto struct {
	Pos    rl.Vector2
	Radius int32
	Color  rl.Color
}

func LoadParticles(particleJson string) []ParticlesDto {
	data, err := ioutil.ReadFile(particleJson)
	if err != nil {
		panic(err)
	}

	var particles Particles

	err = xml.Unmarshal(data, &particles)
	if err != nil {
		panic(err)
	}

	result := make([]ParticlesDto, 0)

	scaleString := strings.Split(particles.Geometry.Transform, " ")[0]
	scaleString = strings.Replace(scaleString, "scale(", "", -1)
	scaleString = strings.Replace(scaleString, ")", "", -1)
	scaleFloat, err := strconv.ParseFloat(scaleString, 32)
	if err != nil {
		panic(err)
	}

	for i, _ := range particles.Geometry.Ellipses {
		e := particles.Geometry.Ellipses[i]

		radius, err := strconv.ParseFloat(e.Radius, 32)
		if err != nil {
			panic(err)
		}

		if int32(radius*scaleFloat) >= 400 {
			continue
		}

		x, err := strconv.ParseFloat(e.X, 32)
		if err != nil {
			panic(err)
		}

		y, err := strconv.ParseFloat(e.Y, 32)
		if err != nil {
			panic(err)
		}

		r, g, b, err := HexToRGB(e.Color)
		if err != nil {
			panic(err)
		}

		result = append(result, ParticlesDto{
			Pos:    rl.NewVector2(float32(x*scaleFloat), float32(y*scaleFloat)),
			Radius: int32(radius*scaleFloat),
			Color: rl.NewColor(
				uint8(r),
				uint8(g),
				uint8(b),
				uint8(100),
			),
		})
	}

	return result
}

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

func HexToRGB(hex string) (int, int, int, error) {
	// Remove the "#" prefix if present
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}

	// Parse the hex values
	rgbValue, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return 0, 0, 0, err
	}

	// Extract the RGB components
	red := (rgbValue >> 16) & 0xFF
	green := (rgbValue >> 8) & 0xFF
	blue := rgbValue & 0xFF

	return int(red), int(green), int(blue), nil
}
