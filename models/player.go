package models

import (
	"ahasuerus/collision"
	"ahasuerus/resources"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	JUMP_SPEED        = 350
	GRAVITY           = 10
	PLAYER_MOVE_SPEED = 5
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
	Lightboxes  []Hitbox
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
				rl.GetShaderLocation(p.Shader, "lightMaxDistance"),
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

	if DRAW_MODELS {
		p.drawHitbox()
	}
}

func (p *Player) Update(delta float32) {
	p.currentAnimation = p.stayAnimation

	p.velocity.X = 0
	p.velocity.Y += GRAVITY * delta

	moveByXButtonPressed := p.processMoveXInput()

	futurePos := rl.Vector2Add(p.Pos, p.velocity)

	futureHitboxMap := p.getHitboxMap(futurePos)
	hitbox := p.getHitboxFromMap(futureHitboxMap)
	hasCollision, collisionMap := p.CollisionProcessor.Detect(hitbox)
	if hasCollision {
		futurePos = p.resolveCollission(moveByXButtonPressed, collisionMap, delta)
	}

	p.Pos = futurePos

	p.updateAnimation(delta)

	if p.ImageShader == resources.TextureLightShader {
		lightPoints := make([]float32, 0)
		lightPointsRadius := make([]float32, 0)
		for i, _ := range p.Lightboxes {
			lp := p.Lightboxes[i]
			lightPoints = append(lightPoints, float32(lp.Center().X), float32(lp.Center().Y))
			radius := rl.Vector2Distance(lp.TopLeft(), lp.TopRight()) / 2
			lightPointsRadius = append(lightPointsRadius, radius)
		}

		playerHitboxMap := p.getHitboxMap(p.Pos)

		rl.SetShaderValue(p.Shader, p.shaderLocs[0], []float32{playerHitboxMap.center.X, playerHitboxMap.center.Y}, rl.ShaderUniformVec2)
		rl.SetShaderValue(p.Shader, p.shaderLocs[1], []float32{p.width, p.height}, rl.ShaderUniformVec2)
		rl.SetShaderValueV(p.Shader, p.shaderLocs[2], lightPoints, rl.ShaderUniformVec2, int32(len(p.Lightboxes)))
		rl.SetShaderValue(p.Shader, p.shaderLocs[3], []float32{float32(len(p.Lightboxes))}, rl.ShaderUniformFloat)

		rl.SetShaderValueTexture(p.Shader, p.shaderLocs[4], p.currentAnimation.Texture)
		rl.SetShaderValueV(p.Shader, p.shaderLocs[5], lightPointsRadius, rl.ShaderUniformFloat, int32(len(p.Lightboxes)))
	}
}

func (p *Player) updateAnimation(delta float32) {
	p.currentAnimation.Pos.X = p.Pos.X
	p.currentAnimation.Pos.Y = p.Pos.Y
	p.currentAnimation.Orientation = p.orientation
	p.currentAnimation.Update(delta)
}

func (p *Player) processMoveXInput() bool {
	if rl.IsKeyDown(rl.KeyLeft) && !p.paused {
		p.currentAnimation = p.runAnimation
		p.velocity.X = (-1) * PLAYER_MOVE_SPEED
		p.orientation = Left
		return true
	}

	if rl.IsKeyDown(rl.KeyRight) && !p.paused {
		p.currentAnimation = p.runAnimation
		p.velocity.X = PLAYER_MOVE_SPEED
		p.orientation = Right
		return true
	}

	return false
}

func (p *Player) AddLightbox(lp Hitbox) *Player {
	p.Lightboxes = append(p.Lightboxes, lp)
	return p
}

func (p *Player) WithShader(gs resources.GameShader) *Player {
	p.ImageShader = gs
	return p
}

func (p *Player) resolveCollission(moveByXButtonPressed bool, collisionMap map[int]bool, delta float32) rl.Vector2 {
	// _, topLeft := collisionMap[0]
	// _, topRight := collisionMap[1]

	_, rightTop := collisionMap[2]
	_, rightBottom := collisionMap[3]

	_, bottomRight := collisionMap[4]
	_, bottomLeft := collisionMap[5]

	_, leftBottom := collisionMap[6]
	_, leftTop := collisionMap[7]

	if bottomRight || bottomLeft { // fall on ground
		p.velocity.Y = 0
	}

	pushFromWall := false

	if leftTop && leftBottom { // left wall collision (push side)
		p.velocity.X = GRAVITY * 5 * delta
		pushFromWall = true
	}

	if rightTop && rightBottom { // right wall collision (push side)
		p.velocity.X = (-1) * GRAVITY * 5 * delta
		pushFromWall = true
	}

	if (rightBottom && bottomRight || bottomLeft && leftBottom) && moveByXButtonPressed && !pushFromWall { // push hero up when go stairs
		p.velocity.Y = (-1) * GRAVITY * 5 * delta
	}

	// jump
	if bottomRight || bottomLeft {
		spacePressed := rl.IsKeyDown(rl.KeySpace)
		if spacePressed {
			p.velocity.Y = (-1) * (GRAVITY / 1.5)
		}
	}

	return rl.Vector2Add(p.Pos, p.velocity)
}

type playerHitboxMap struct {
	topLeftOne rl.Vector2
	topLeftTwo rl.Vector2

	topMiddle rl.Vector2

	topRightOne rl.Vector2
	topRightTwo rl.Vector2

	leftMiddle  rl.Vector2
	center      rl.Vector2
	rightMiddle rl.Vector2

	bottomLeftOne rl.Vector2
	bottomLeftTwo rl.Vector2

	bottomMiddle rl.Vector2

	bottomRightOne rl.Vector2
	bottomRightTwo rl.Vector2
}

func (p *Player) drawHitbox() {
	hitbox := p.getHitboxFromMap(p.getHitboxMap(p.Pos))
	for i, _ := range hitbox.Polygons {
		poly := hitbox.Polygons[i]
		rl.DrawTriangleLines(poly.Points[0], poly.Points[1], poly.Points[2], rl.Gold)
		rl.DrawText(fmt.Sprintf("%v", p.velocity), int32(p.Pos.X)-100, int32(p.Pos.Y)-100, 50, rl.Red)
	}
}

func (p *Player) getHitboxMap(pos rl.Vector2) playerHitboxMap {
	cornerOffset := float32(10)
	return playerHitboxMap{
		topLeftOne: rl.Vector2{pos.X + cornerOffset, pos.Y},
		topLeftTwo: rl.Vector2{pos.X, pos.Y + cornerOffset},

		topMiddle: rl.Vector2{pos.X + p.width/2, pos.Y},

		topRightOne: rl.Vector2{pos.X + p.width - cornerOffset, pos.Y},
		topRightTwo: rl.Vector2{pos.X + p.width, pos.Y + cornerOffset},

		leftMiddle:  rl.Vector2{pos.X, pos.Y + p.height/2},
		center:      rl.Vector2{pos.X + p.width/2, pos.Y + p.height/2},
		rightMiddle: rl.Vector2{pos.X + p.width, pos.Y + p.height/2},

		bottomRightOne: rl.Vector2{pos.X + p.width, pos.Y + p.height - cornerOffset},
		bottomRightTwo: rl.Vector2{pos.X + p.width - cornerOffset, pos.Y + p.height},

		bottomMiddle: rl.Vector2{pos.X + p.width/2, pos.Y + p.height},

		bottomLeftOne: rl.Vector2{pos.X + cornerOffset, pos.Y + p.height},
		bottomLeftTwo: rl.Vector2{pos.X, pos.Y + p.height - cornerOffset},
	}
}

func (p *Player) getHitboxFromMap(m playerHitboxMap) collision.Hitbox {
	return collision.Hitbox{
		Polygons: []collision.Polygon{
			{
				Points: [3]rl.Vector2{
					m.topLeftOne, m.topMiddle, m.center,
				},
			},
			{
				Points: [3]rl.Vector2{
					m.topMiddle, m.topRightOne, m.center,
				},
			},
			{
				Points: [3]rl.Vector2{
					m.topRightTwo, m.rightMiddle, m.center,
				},
			},
			{
				Points: [3]rl.Vector2{
					m.rightMiddle, m.bottomRightOne, m.center,
				},
			},
			{
				Points: [3]rl.Vector2{
					m.bottomRightTwo, m.bottomMiddle, m.center,
				},
			},
			{
				Points: [3]rl.Vector2{
					m.bottomMiddle, m.bottomLeftOne, m.center,
				},
			},
			{
				Points: [3]rl.Vector2{
					m.bottomLeftTwo, m.leftMiddle, m.center,
				},
			},
			{
				Points: [3]rl.Vector2{
					m.leftMiddle, m.topLeftTwo, m.center,
				},
			},
		},
	}
}
