package scene

import (
	"ahasuerus/models"
	"ahasuerus/config"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	WIDTH, HEIGHT = config.GetResolution()
)

func updateCameraSmooth(camera *rl.Camera2D, player *models.Player, delta float32) {
	minSpeed := 60.0
	minEffectLength := 10
	fractionSpeed := 0.8

	diff := rl.Vector2Subtract(player.Pos, camera.Target)
	length := rl.Vector2Length(diff)

	if length > float32(minEffectLength) {
		speed := float32(math.Max(fractionSpeed*float64(length),minSpeed))
		camera.Target = rl.Vector2Add(camera.Target, rl.Vector2Scale(diff, speed*delta/length))
	}
}