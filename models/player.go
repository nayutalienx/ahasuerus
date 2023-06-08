package models

import (
	"ahasuerus/resources"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	JUMP_SPEED     = 350
	GRAVITY        = 400
	MOVE_SPEED     = 280
	MAX_FALL_SPEED = JUMP_SPEED
)

type Player struct {
	Pos   rl.Vector2
	speed float32

	fallSpeed int32

	Box rl.Vector2

	currentAnimation *Animation
	runAnimation     *Animation
	stayAnimation    *Animation
	orientation      Orientation

	Shader      rl.Shader
	ImageShader resources.GameShader
	LightPoints []*LightPoint
	shaderLocs  []int32

	paused bool

	debugText Text
}

func NewPlayer(x float32, y float32) *Player {
	p := &Player{
		Pos:   rl.NewVector2(x, y),
		speed: MOVE_SPEED,

		fallSpeed: 0,

		Box: rl.NewVector2(100, 200),
	}

	p.debugText = *NewText(int32(x), int32(y)).
		SetFontSize(20).
		SetColor(rl.Lime).
		SetUpdateCallback(p.normalDebug())

	return p
}

func (p *Player) Load() {
	p.runAnimation = NewAnimation(resources.PlayerRunTexture, 27, 24)
	p.runAnimation.Load()

	p.stayAnimation = NewAnimation(resources.PlayerStayTexture, 22, 7)
	p.stayAnimation.Load()

	p.Box.X = float32(p.stayAnimation.StepInPixel)
	p.Box.Y = float32(p.stayAnimation.Texture.Height)

	if p.ImageShader != resources.UndefinedShader {
		p.Shader = resources.LoadShader(p.ImageShader)
		if p.ImageShader == resources.TextureLightShader {
			p.shaderLocs = []int32{
				rl.GetShaderLocation(p.Shader, "objectPosBottomLeft"),
				rl.GetShaderLocation(p.Shader, "objectSize"),
				rl.GetShaderLocation(p.Shader, "lightPos"),
				rl.GetShaderLocation(p.Shader, "lightPosSize"),
				rl.GetShaderLocation(p.Shader, "texture0"),
			}
		}
	}
}

func (p *Player) Unload() {
	p.runAnimation.Unload()
	p.stayAnimation.Unload()
	if p.ImageShader != resources.UndefinedShader {
		resources.UnloadShader(p.Shader)
	}
}

func (p *Player) Resume() {
	p.fallSpeed = 0
	p.paused = false
}

func (p *Player) Pause() {
	p.paused = true
}

func (p Player) Draw() {

	if p.ImageShader != resources.UndefinedShader {
		rl.BeginShaderMode(p.Shader)
		p.currentAnimation.Draw()
		rl.EndShaderMode()
	} else {
		p.currentAnimation.Draw()
	}

	//p.debugText.Draw()
	// for _, colPoint := range p.collisionBezierChecks {
	// 	if colPoint.Colliding {
	// 		rl.DrawCircle(int32(colPoint.Point.X), int32(colPoint.Point.Y), 4, rl.Orange)
	// 	}
	// }
}

func (p *Player) Update(delta float32) {

	p.currentAnimation = p.stayAnimation

	spacePressed := rl.IsKeyDown(rl.KeySpace)

	if spacePressed && p.fallSpeed == 0 {
		p.fallSpeed = -JUMP_SPEED
	}

	if p.fallSpeed > MAX_FALL_SPEED {
		p.fallSpeed = MAX_FALL_SPEED
	}

	calculatedSpeed := p.speed * delta

	if rl.IsKeyDown(rl.KeyLeft) && !p.paused {
		p.currentAnimation = p.runAnimation
		p.Pos.X -= calculatedSpeed
		p.orientation = Left
	}

	if rl.IsKeyDown(rl.KeyRight) && !p.paused {
		p.currentAnimation = p.runAnimation
		p.Pos.X += calculatedSpeed
		p.orientation = Right
	}

	p.fallSpeed += int32(GRAVITY * delta)

	shouldUpdateY := true

	if shouldUpdateY {
		p.Pos.Y += float32(p.fallSpeed) * delta
	}

	p.updateCollisions()


	p.currentAnimation.Pos.X = p.Pos.X
	p.currentAnimation.Pos.Y = p.Pos.Y
	p.currentAnimation.Orientation = p.orientation
	p.currentAnimation.Update(delta)

	if p.ImageShader == resources.TextureLightShader {
		lightPoints := make([]float32, 0)
		for i, _ := range p.LightPoints {
			lp := p.LightPoints[i]
			lightPoints = append(lightPoints, float32(lp.Pos.X), float32(lp.Pos.Y))
		}
		rl.SetShaderValue(p.Shader, p.shaderLocs[0], []float32{p.Pos.X, p.Pos.Y + p.Box.Y}, rl.ShaderUniformVec2)
		rl.SetShaderValue(p.Shader, p.shaderLocs[1], []float32{p.Box.X, p.Box.Y}, rl.ShaderUniformVec2)
		rl.SetShaderValueV(p.Shader, p.shaderLocs[2], lightPoints, rl.ShaderUniformVec2, int32(len(p.LightPoints)))
		rl.SetShaderValue(p.Shader, p.shaderLocs[3], []float32{float32(len(p.LightPoints))}, rl.ShaderUniformFloat)

		rl.SetShaderValueTexture(p.Shader, p.shaderLocs[4], p.currentAnimation.Texture)
	}

	//p.debugText.Update(delta)
}


func (p *Player) GetPos() *rl.Vector2 {
	return &p.Pos
}

func (p *Player) GetBox() *rl.Vector2 {
	return &p.Box
}

func (p *Player) AddLightPoint(lp *LightPoint) *Player {
	p.LightPoints = append(p.LightPoints, lp)
	return p
}

func (p *Player) updateCollisions() {

}

func (p *Player) WithShader(gs resources.GameShader) *Player {
	p.ImageShader = gs
	return p
}

func (p *Player) normalDebug() func(t *Text) {
	return func(t *Text) {

	}
}
