package models

import (
	"ahasuerus/collision"
	"ahasuerus/resources"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	JUMP_SPEED = 350
	GRAVITY    = 400
	MOVE_SPEED = 5
)

type Player struct {
	Pos                rl.Vector2
	CollisionProcessor collision.CollisionDetector
	velocity           rl.Vector2

	width, height    float32
	currentAnimation *Animation
	runAnimation     *Animation
	stayAnimation    *Animation
	orientation      Orientation

	Shader      rl.Shader
	ImageShader resources.GameShader
	LightPoints []*LightPoint
	shaderLocs  []int32

	paused bool
}

func NewPlayer(x float32, y float32) *Player {
	p := &Player{
		Pos: rl.NewVector2(x, y),
	}
	return p
}

func (p *Player) Load() {
	p.runAnimation = NewAnimation(resources.PlayerRunTexture, 27, 24)
	p.runAnimation.Load()

	p.stayAnimation = NewAnimation(resources.PlayerStayTexture, 22, 7)
	p.stayAnimation.Load()

	p.width = float32(p.stayAnimation.StepInPixel)
	p.height = float32(p.stayAnimation.Texture.Height)


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
}

func (p *Player) Update(delta float32) {
	p.currentAnimation = p.stayAnimation
	
	p.velocity.X = 0
	p.velocity.Y = GRAVITY * delta
	
	p.processMoveXInput()

	futurePos := rl.Vector2Add(p.Pos, p.velocity)

	hitbox := p.getHitboxFromPosition(futurePos)
	hasCollision, _ := p.CollisionProcessor.Detect(hitbox)
	if hasCollision {
		p.velocity.Y = 0
		futurePos = rl.Vector2Add(p.Pos, p.velocity)
	}

	p.Pos = futurePos

	p.updateAnimation(delta)

	if p.ImageShader == resources.TextureLightShader {
		lightPoints := make([]float32, 0)
		for i, _ := range p.LightPoints {
			lp := p.LightPoints[i]
			lightPoints = append(lightPoints, float32(lp.Pos.X), float32(lp.Pos.Y))
		}
		rl.SetShaderValue(p.Shader, p.shaderLocs[0], []float32{p.Pos.X, p.Pos.Y + p.height}, rl.ShaderUniformVec2)
		rl.SetShaderValue(p.Shader, p.shaderLocs[1], []float32{p.width, p.height}, rl.ShaderUniformVec2)
		rl.SetShaderValueV(p.Shader, p.shaderLocs[2], lightPoints, rl.ShaderUniformVec2, int32(len(p.LightPoints)))
		rl.SetShaderValue(p.Shader, p.shaderLocs[3], []float32{float32(len(p.LightPoints))}, rl.ShaderUniformFloat)

		rl.SetShaderValueTexture(p.Shader, p.shaderLocs[4], p.currentAnimation.Texture)
	}
}

func (p *Player) updateAnimation(delta float32) {
	p.currentAnimation.Pos.X = p.Pos.X
	p.currentAnimation.Pos.Y = p.Pos.Y
	p.currentAnimation.Orientation = p.orientation
	p.currentAnimation.Update(delta)
}

func (p *Player) processMoveYInput() {
	//spacePressed := rl.IsKeyDown(rl.KeySpace)
}

func (p *Player) processMoveXInput() {
	if rl.IsKeyDown(rl.KeyLeft) && !p.paused {
		p.currentAnimation = p.runAnimation
		p.velocity.X = (-1) * MOVE_SPEED
		p.orientation = Left
	}

	if rl.IsKeyDown(rl.KeyRight) && !p.paused {
		p.currentAnimation = p.runAnimation
		p.velocity.X = MOVE_SPEED
		p.orientation = Right
	}
}

func (p *Player) AddLightPoint(lp *LightPoint) *Player {
	p.LightPoints = append(p.LightPoints, lp)
	return p
}

func (p *Player) WithShader(gs resources.GameShader) *Player {
	p.ImageShader = gs
	return p
}

func (p *Player) getHitboxFromPosition(pos rl.Vector2) collision.Hitbox {
	topLeft := pos
	bottomLeft := rl.Vector2{pos.X, pos.Y + p.height}

	topRight := rl.Vector2{pos.X + p.width, pos.Y}
	bottomRight := rl.Vector2{pos.X + p.width, pos.Y + p.height}

	return collision.Hitbox{
		Polygons: []collision.Polygon{
			{
				Points: [3]rl.Vector2{
					topLeft, topRight, bottomRight,
				},
			},
			{
				Points: [3]rl.Vector2{
					topLeft, bottomLeft, bottomRight,
				},
			},
		},
	}
}
