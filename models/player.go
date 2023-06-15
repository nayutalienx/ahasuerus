package models

import (
	"ahasuerus/collision"
	"ahasuerus/resources"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	JUMP_SPEED        = 350
	GRAVITY           = 10
	PLAYER_MOVE_SPEED = 5

	rewindBufferSize = 60 * 60 * 30
)

type RewindItem struct {
	Pos              rl.Vector2
	orientation      Orientation
	currentAnimation *Animation
}

type Player struct {
	Pos                rl.Vector2
	CollisionProcessor collision.CollisionDetector
	velocity           rl.Vector2

	width, height float32
	orientation   Orientation
	currentHitbox *collision.Hitbox

	currentAnimation    *Animation
	runAnimation        *Animation
	stayAnimation       *Animation
	directUpAnimation   *Animation
	directDownAnimation *Animation
	sideUpAnimation     *Animation
	sideDownAnimation   *Animation

	Shader      rl.Shader
	ImageShader resources.GameShader
	Lightboxes  []Hitbox
	shaderLocs  []int32

	Rewind          [rewindBufferSize]RewindItem
	rewindLastIndex int32

	paused bool
}

func NewPlayer(x float32, y float32) *Player {

	p := &Player{
		Pos: rl.NewVector2(x, y),
	}
	hb := GetDynamicHitboxFromMap(GetDynamicHitboxMap(p.Pos, p.width, p.height))
	p.currentHitbox = &hb

	return p
}

func (p *Player) Load() {
	p.runAnimation = NewAnimation(resources.PlayerRunTexture, 27, Loop).FramesPerSecond(24)
	p.runAnimation.Load()

	p.stayAnimation = NewAnimation(resources.PlayerStayTexture, 22, Loop).FramesPerSecond(7)
	p.stayAnimation.Load()

	p.directUpAnimation = NewAnimation(resources.PlayerDirectUpTexture, 5, Temporary).TimeInSeconds(1)
	p.directUpAnimation.Load()

	p.directDownAnimation = NewAnimation(resources.PlayerDirectDownTexture, 6, Temporary).TimeInSeconds(1.5)
	p.directDownAnimation.Load()

	p.sideUpAnimation = NewAnimation(resources.PlayerSideUpTexture, 12, Temporary).TimeInSeconds(1)
	p.sideUpAnimation.Load()

	p.sideDownAnimation = NewAnimation(resources.PlayerSideDownTexture, 12, Temporary).TimeInSeconds(1.5)
	p.sideDownAnimation.Load()

	p.currentAnimation = p.stayAnimation

	p.width = float32(p.stayAnimation.StepInPixel)
	p.height = float32(p.stayAnimation.Texture.Height)

	if p.ImageShader != resources.UndefinedShader {
		p.Shader = resources.LoadShader(p.ImageShader)
		if p.ImageShader == resources.TextureLightShader {
			p.shaderLocs = []int32{
				rl.GetShaderLocation(p.Shader, "texture0"),
				rl.GetShaderLocation(p.Shader, "objectPosCenter"),
				rl.GetShaderLocation(p.Shader, "lightPosSize"),
				rl.GetShaderLocation(p.Shader, "lightPos"),
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

	rewindEnabled := rl.IsKeyDown(rl.KeyLeftShift)

	if !rewindEnabled {
		p.velocity.X = 0
		p.velocity.Y += GRAVITY * delta

		moveByXButtonPressed := p.processMoveXInput()

		futurePos := rl.Vector2Add(p.Pos, p.velocity)

		futureHitboxMap := GetDynamicHitboxMap(futurePos, p.width, p.height)
		hitbox := GetDynamicHitboxFromMap(futureHitboxMap)
		hasCollision, collisionMap := p.CollisionProcessor.Detect(hitbox)
		if hasCollision {
			futurePos = p.resolveCollission(moveByXButtonPressed, collisionMap, delta)
		}

		posDelta := rl.Vector2Subtract(p.Pos, futurePos)

		p.Pos = futurePos

		p.resolveAndUpdateAnimation(hasCollision, posDelta, delta)

		p.savePlayerToRewind()
	} else {
		p.rewindPlayer()
		p.currentAnimation.Reverse(true)
		p.updateAnimation(delta)
		p.currentAnimation.Reverse(false)
	}

	// update hitbox for others
	p.updateCurrentHitbox()

	if p.ImageShader == resources.TextureLightShader {
		lightPoints := make([]float32, 0)
		lightPointsRadius := make([]float32, 0)
		for i, _ := range p.Lightboxes {
			lp := p.Lightboxes[i]
			lightPoints = append(lightPoints, float32(lp.Center().X), float32(lp.Center().Y))
			radius := rl.Vector2Distance(lp.TopLeft(), lp.TopRight()) / 2
			lightPointsRadius = append(lightPointsRadius, radius)
		}

		playerHitboxMap := GetDynamicHitboxMap(p.Pos, p.width, p.height)

		rl.SetShaderValueTexture(p.Shader, p.shaderLocs[0], p.currentAnimation.Texture)
		rl.SetShaderValue(p.Shader, p.shaderLocs[1], []float32{playerHitboxMap.center.X, playerHitboxMap.center.Y}, rl.ShaderUniformVec2)
		rl.SetShaderValue(p.Shader, p.shaderLocs[2], []float32{float32(len(p.Lightboxes))}, rl.ShaderUniformFloat)

		rl.SetShaderValueV(p.Shader, p.shaderLocs[3], lightPoints, rl.ShaderUniformVec2, int32(len(p.Lightboxes)))
		rl.SetShaderValueV(p.Shader, p.shaderLocs[4], lightPointsRadius, rl.ShaderUniformFloat, int32(len(p.Lightboxes)))
	}
}

func (p *Player) savePlayerToRewind() {
	p.Rewind[p.rewindLastIndex] = RewindItem{
		Pos:              p.Pos,
		orientation:      p.orientation,
		currentAnimation: p.currentAnimation,
	}
	p.Rewind[p.rewindLastIndex+1] = p.Rewind[p.rewindLastIndex] // fix jump to 0,0 when rewind
	p.rewindLastIndex++
}

func (p *Player) rewindPlayer() {
	if p.rewindLastIndex > 0 {
		rewind := p.Rewind[p.rewindLastIndex-1]
		p.Pos = rewind.Pos
		p.orientation = rewind.orientation
		p.currentAnimation = rewind.currentAnimation
		p.rewindLastIndex--
	}
}

func (p *Player) resolveAndUpdateAnimation(hasCollision bool, posDelta rl.Vector2, delta float32) {

	prevAnimation := p.currentAnimation

	if posDelta.X != 0 && hasCollision {
		p.currentAnimation = p.runAnimation
	}

	if posDelta.X == 0 && posDelta.Y == 0 {
		p.currentAnimation = p.stayAnimation
	}

	if posDelta.Y > 1 {
		p.currentAnimation = p.directUpAnimation
		if posDelta.X != 0 {
			p.currentAnimation = p.sideUpAnimation
		}
	}

	if posDelta.Y < -2 {
		p.currentAnimation = p.directDownAnimation
		if posDelta.X != 0 {
			p.currentAnimation = p.sideDownAnimation
		}
	}

	if p.currentAnimation != prevAnimation {
		p.currentAnimation.Begin()
	}
	p.updateAnimation(delta)
}

func (p *Player) updateAnimation(delta float32) {
	p.currentAnimation.Pos.X = p.Pos.X
	p.currentAnimation.Pos.Y = p.Pos.Y
	p.currentAnimation.Orientation = p.orientation
	p.currentAnimation.Update(delta)
}

func (p *Player) processMoveXInput() bool {
	if rl.IsKeyDown(rl.KeyLeft) && !p.paused {
		p.velocity.X = (-1) * PLAYER_MOVE_SPEED
		p.orientation = Left
		return true
	}

	if rl.IsKeyDown(rl.KeyRight) && !p.paused {
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

func (p *Player) drawHitbox() {
	playerHitboxMap := GetDynamicHitboxMap(p.Pos, p.width, p.height)
	hitbox := GetDynamicHitboxFromMap(playerHitboxMap)
	for i, _ := range hitbox.Polygons {
		poly := hitbox.Polygons[i]
		rl.DrawTriangleLines(poly.Points[0], poly.Points[1], poly.Points[2], rl.Gold)
	}
}

func (p Player) GetHitbox() *collision.Hitbox {
	return p.currentHitbox
}

func (p *Player) updateCurrentHitbox() {
	updatedHb := GetDynamicHitboxFromMap(GetDynamicHitboxMap(p.Pos, p.width, p.height))
	p.currentHitbox.Polygons = updatedHb.Polygons
}
