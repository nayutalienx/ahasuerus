package models

import (
	"ahasuerus/collision"
	"ahasuerus/resources"
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	JUMP_SPEED        = 350
	FALL_GRAVITY      = 20
	PUSH_FORCE        = 5
	JUMP_FORCE        = 10
	PLAYER_MOVE_SPEED = 5

	MIN_REWIND_SPEED = -4
	MAX_REWIND_SPEED = 4
)

type PlayerRewindData struct {
	Pos              rl.Vector2
	orientation      Orientation
	currentAnimation *Animation
	velocity         rl.Vector2
}

type Player struct {
	Pos                rl.Vector2                  `json:"-"`
	CollisionProcessor collision.CollisionDetector `json:"-"`
	velocity           rl.Vector2                  `json:"-"`
	jumpCounter        uint8                       `json:"-"`

	width, height float32           `json:"-"`
	orientation   Orientation       `json:"-"`
	currentHitbox *collision.Hitbox `json:"-"`

	currentAnimation    *Animation `json:"-"`
	runAnimation        *Animation `json:"-"`
	stayAnimation       *Animation `json:"-"`
	directUpAnimation   *Animation `json:"-"`
	directDownAnimation *Animation `json:"-"`
	sideUpAnimation     *Animation `json:"-"`
	sideDownAnimation   *Animation `json:"-"`

	Shader      rl.Shader            `json:"-"`
	ImageShader resources.GameShader `json:"-"`
	Lightboxes  []Light              `json:"-"`
	shaderLocs  []int32              `json:"-"`

	Rewind               [REWIND_BUFFER_SIZE]PlayerRewindData `json:"-"`
	rewindLastIndex      int32                                `json:"-"`
	rewindSpeed          int32                                `json:"-"`
	rewindModeStartIndex int32                                `json:"-"`
	rewindModeStarted    bool                                 `json:"-"`
	rewindCollision      bool

	paused bool `json:"-"`
}

func NewPlayer(x float32, y float32) *Player {

	p := &Player{
		Pos:         rl.NewVector2(x, y),
		rewindSpeed: 1,
	}
	hb := GetDynamicHitboxFromMap(GetDynamicHitboxMap(p.Pos, p.width, p.height))
	p.currentHitbox = &hb

	return p
}

func (p *Player) GetId() string {
	return "player"
}

func (p *Player) GetDrawIndex() int {
	return -999
}

func (p *Player) IsCollisionRewind() bool {
	return p.rewindCollision
}

func (p *Player) Load() {
	p.runAnimation = NewAnimation(resources.PlayerRunTexture, 27, Loop).FramesPerSecond(30)
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
		if p.ImageShader == resources.PlayerShader {
			p.shaderLocs = []int32{
				rl.GetShaderLocation(p.Shader, "texture0"),
				rl.GetShaderLocation(p.Shader, "objectPosCenter"),
				rl.GetShaderLocation(p.Shader, "lightPosSize"),
				rl.GetShaderLocation(p.Shader, "lightPos"),
				rl.GetShaderLocation(p.Shader, "lightMaxDistance"),
				rl.GetShaderLocation(p.Shader, "playerWidth"),
				rl.GetShaderLocation(p.Shader, "playerHeight"),
				rl.GetShaderLocation(p.Shader, "rewind"),
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

	p.drawRewindSpeed()
}

func (p *Player) Update(delta float32) {

	rewindEnabled := rl.IsKeyDown(rl.KeyLeftShift)

	if !rewindEnabled {
		newVelocity := p.velocity

		newVelocity = p.movementResist(newVelocity, 1, delta)
		newVelocity.Y += FALL_GRAVITY * delta

		newVelocity = p.processMoveXInput(newVelocity)

		futureHitboxMap := GetDynamicHitboxMap(rl.Vector2Add(p.Pos, newVelocity), p.width, p.height)
		hitbox := GetDynamicHitboxFromMap(futureHitboxMap)
		hasCollision, collisionMapArray := p.CollisionProcessor.Detect(hitbox)
		if hasCollision {
			for i, _ := range collisionMapArray {
				newVelocity = p.resolveCollission(newVelocity, collisionMapArray[i], delta)
			}
		}

		p.velocity = newVelocity
		futurePos := rl.Vector2Add(p.Pos, p.velocity)

		posDelta := rl.Vector2Subtract(p.Pos, futurePos)

		p.Pos = futurePos

		p.resolveAndUpdateAnimation(hasCollision, posDelta, delta)

		p.savePlayerToRewind()
		p.rewindModeStarted = false
		p.rewindCollision = false
	} else {
		p.updateRewindSpeed()
		p.rewindPlayer()
		if p.rewindSpeed > 0 {
			p.currentAnimation.Reverse(true)
			p.updateAnimation(delta, uint8(p.rewindSpeed))
			p.currentAnimation.Reverse(false)
		} else if p.rewindSpeed < 0 {
			p.updateAnimation(delta, uint8(math.Abs(float64(p.rewindSpeed))))
		}
		p.rewindModeStarted = true
	}

	// update hitbox for others
	p.updateCurrentHitbox()

	if p.ImageShader == resources.PlayerShader {
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
		rl.SetShaderValue(p.Shader, p.shaderLocs[5], []float32{float32(p.currentAnimation.Texture.Width)}, rl.ShaderUniformFloat)
		rl.SetShaderValue(p.Shader, p.shaderLocs[6], []float32{float32(p.currentAnimation.Texture.Height)}, rl.ShaderUniformFloat)
		rewind := 0.0
		if p.rewindModeStarted {
			rewind = 1.0
		}
		rl.SetShaderValue(p.Shader, p.shaderLocs[7], []float32{float32(rewind)}, rl.ShaderUniformFloat)
	}

	if p.jumpCounter > 0 {
		p.jumpCounter--
	}
}

func (p *Player) savePlayerToRewind() {
	if int(p.rewindLastIndex) == len(p.Rewind)-1 {
		p.rewindLastIndex = 0
	}

	p.Rewind[p.rewindLastIndex] = PlayerRewindData{
		Pos:              p.Pos,
		orientation:      p.orientation,
		currentAnimation: p.currentAnimation,
		velocity:         p.velocity,
	}
	p.rewindLastIndex++
}

func (p *Player) updateRewindSpeed() {
	rewindEnabled := rl.IsKeyDown(rl.KeyLeftShift)
	if rewindEnabled {
		if rl.IsKeyReleased(rl.KeyDown) {
			p.rewindSpeed--
			if p.rewindSpeed < MIN_REWIND_SPEED {
				p.rewindSpeed = MIN_REWIND_SPEED
			}
		}

		if rl.IsKeyReleased(rl.KeyUp) {
			p.rewindSpeed++
			if p.rewindSpeed > MAX_REWIND_SPEED {
				p.rewindSpeed = MAX_REWIND_SPEED
			}
		}
	}
}

func (p *Player) rewindPlayer() {
	if !p.rewindModeStarted {
		p.rewindModeStartIndex = p.rewindLastIndex
		p.rewindSpeed = 1
	}

	rewind := p.Rewind[p.rewindLastIndex]

	if p.rewindLastIndex > p.rewindSpeed && p.rewindLastIndex < p.rewindModeStartIndex+p.rewindSpeed {
		rewind = p.Rewind[p.rewindLastIndex-p.rewindSpeed]
		p.rewindLastIndex -= p.rewindSpeed
		p.rewindCollision = false
	} else {
		p.rewindCollision = true
	}

	p.Pos = rewind.Pos
	p.orientation = rewind.orientation
	p.currentAnimation = rewind.currentAnimation
	p.velocity = rewind.velocity
}

func (p *Player) drawRewindSpeed() {
	rewindEnabled := rl.IsKeyDown(rl.KeyLeftShift)
	if rewindEnabled {
		DrawSdfText(fmt.Sprintf("%dx", p.rewindSpeed), p.Pos, 60, rl.White)
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
	p.updateAnimation(delta, 1)
}

func (p *Player) updateAnimation(delta float32, speed uint8) {
	p.currentAnimation.Pos.X = p.Pos.X
	p.currentAnimation.Pos.Y = p.Pos.Y
	p.currentAnimation.Orientation = p.orientation
	p.currentAnimation.AnimationSpeed(speed)
	p.currentAnimation.Update(delta)
}

func (p *Player) processMoveXInput(velocity rl.Vector2) rl.Vector2 {
	if rl.IsKeyDown(rl.KeyLeft) && !p.paused {
		velocity.X = (-1) * PLAYER_MOVE_SPEED
		p.orientation = Left
	}

	if rl.IsKeyDown(rl.KeyRight) && !p.paused {
		velocity.X = PLAYER_MOVE_SPEED
		p.orientation = Right
	}
	return velocity
}

func (p *Player) AddLightbox(lp Light) *Player {
	p.Lightboxes = append(p.Lightboxes, lp)
	return p
}

func (p *Player) WithShader(gs resources.GameShader) *Player {
	p.ImageShader = gs
	return p
}

func (p *Player) resolveCollission(velocity rl.Vector2, collisionMap map[int]float32, delta float32) rl.Vector2 {
	_, topLeft := collisionMap[0]
	_, topRight := collisionMap[1]

	_, rightTop := collisionMap[2]
	_, rightBottom := collisionMap[3]

	bottomRightRotation, bottomRight := collisionMap[4]
	bottomLeftRotation, bottomLeft := collisionMap[5]

	_, leftBottom := collisionMap[6]
	_, leftTop := collisionMap[7]

	if bottomRight && bottomLeft { // fall on ground
		velocity.Y = 0
		velocity = p.movementResist(velocity, 7, delta)
		velocity = Vec2Rotate(velocity, math.Max(float64(bottomLeftRotation), float64(bottomRightRotation)))
	}

	if leftTop || leftBottom { // left wall collision
		velocity.X = 0
	}

	if rightTop || rightBottom { // right wall collision
		velocity.X = 0
	}

	bottomLeftStuck := bottomLeft && leftBottom && !bottomRight && !leftTop // slide down-side when on border of collision box
	if bottomLeftStuck {
		velocity.X = 3
		velocity.Y = 3
	}
	bottomRightStuck := bottomRight && rightBottom && !bottomLeft && !rightTop // slide down-side when on border of collision box
	if bottomRightStuck {
		velocity.X = -3
		velocity.Y = 3
	}

	if topLeft || topRight { // top collission
		if velocity.Y < 0 {
			velocity.Y = 0
		}
	}

	topLeftStuck := topLeft && leftTop && !leftBottom && !topRight // topLeft border stuck SLIDE DOWN
	if topLeftStuck {
		velocity.X = 3
		velocity.Y = 3
	}

	topRightStuck := topRight && rightTop && !rightBottom && !topLeft // top right border stuck SLIDE DOWN
	if topRightStuck {
		velocity.X = -3
		velocity.Y = 3
	}

	if !leftBottom && leftTop && !topLeft { // LEFT SINGLE TOP COLLISION
		if velocity.Y < 0 {
			velocity.Y = 0
		}
	}

	if !rightBottom && rightTop && !topRight { // RIGHT SINGLE TOP COLLISION
		if velocity.Y < 0 {
			velocity.Y = 0
		}
	}

	if bottomRight && bottomLeft {
		p.jumpCounter = 0
	}

	// jump
	if bottomRight || bottomLeft && (!leftBottom && !rightBottom && !leftTop && !rightTop) {
		spacePressed := rl.IsKeyDown(rl.KeySpace)
		if spacePressed && p.jumpCounter == 0 {
			velocity.Y = (-1) * (JUMP_FORCE)
			p.jumpCounter = uint8(FPS)
		}
	}

	return velocity
}

func (p *Player) movementResist(velocity rl.Vector2, resistScale float32, delta float32) rl.Vector2 {
	if velocity.X > 0 {
		velocity.X += -1 * PLAYER_MOVE_SPEED * resistScale * delta
		if velocity.X < 0 {
			velocity.X = 0
		}
	}
	if velocity.X < 0 {
		velocity.X += PLAYER_MOVE_SPEED * resistScale * delta
		if velocity.X > 0 {
			velocity.X = 0
		}
	}
	return velocity
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
