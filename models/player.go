package models

import (
	"ahasuerus/resources"
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	JUMP_SPEED = 350
	GRAVITY    = 400
)

type Player struct {
	Pos   rl.Vector2
	speed float32

	fallSpeed int32

	Box rl.Vector2

	collisionBoxes     []CollisionBox
	collisionBoxChecks []CollisionBoxCheck

	collisionLines      []CollisionLine
	collisionLineChecks []CollisionLineCheck

	collisionBeziers      []CollisionBezier
	collisionBezierChecks []CollisionBezierCheck

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
		speed: 5,

		fallSpeed: 0,

		Box: rl.NewVector2(100, 200),

		collisionBoxes:     make([]CollisionBox, 0),
		collisionBoxChecks: make([]CollisionBoxCheck, 0),

		collisionBeziers:      make([]CollisionBezier, 0),
		collisionBezierChecks: make([]CollisionBezierCheck, 0),
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

	hasCurveCollision, collisionedCurve := p.hasCurveCollision()
	hasLineCollision, collisionedLine := p.hasLineCollision()

	spacePressed := rl.IsKeyDown(rl.KeySpace)

	if spacePressed && p.fallSpeed == 0 {
		p.fallSpeed = -JUMP_SPEED
	}

	if rl.IsKeyDown(rl.KeyLeft) && p.canMoveLeft() && !p.paused {
		p.currentAnimation = p.runAnimation
		if hasCurveCollision {
			prev, _ := CalculatePreviousNextPointsOfBezier(collisionedCurve.Point, collisionedCurve.Curve.Start, collisionedCurve.Curve.End)
			diff := rl.Vector2Subtract(prev, rl.NewVector2(p.Pos.X+p.Box.X, p.Pos.Y+p.Box.Y))
			movement := rl.Vector2Scale(rl.Vector2Normalize(diff), p.speed)
			p.Pos = rl.Vector2Add(p.Pos, movement)
			//rl.DrawCircle(int32(p.Pos.X), int32(p.Pos.Y), 4, rl.Pink)
		} else if hasLineCollision {
			prev, _ := CalculatePreviousAndNextPointOfLine(collisionedLine.Point, collisionedLine.Line.Start, collisionedLine.Line.End)
			diff := rl.Vector2Subtract(prev, rl.NewVector2(p.Pos.X+p.Box.X, p.Pos.Y+p.Box.Y+(2*p.speed)))
			movement := rl.Vector2Scale(rl.Vector2Normalize(diff), p.speed)
			p.Pos = rl.Vector2Add(p.Pos, movement)
		} else {
			p.Pos.X -= p.speed
		}
		p.orientation = Left
	}

	if rl.IsKeyDown(rl.KeyRight) && p.canMoveRight() && !p.paused {
		p.currentAnimation = p.runAnimation
		if hasCurveCollision {
			_, next := CalculatePreviousNextPointsOfBezier(collisionedCurve.Point, collisionedCurve.Curve.Start, collisionedCurve.Curve.End)
			diff := rl.Vector2Subtract(next, rl.NewVector2(p.Pos.X, p.Pos.Y+p.Box.Y))
			movement := rl.Vector2Scale(rl.Vector2Normalize(diff), p.speed)
			p.Pos = rl.Vector2Add(p.Pos, movement)
			//rl.DrawCircle(int32(p.Pos.X), int32(p.Pos.Y), 4, rl.Pink)
		} else if hasLineCollision {
			_, next := CalculatePreviousAndNextPointOfLine(collisionedLine.Point, collisionedLine.Line.Start, collisionedLine.Line.End)
			diff := rl.Vector2Subtract(next, rl.NewVector2(p.Pos.X, p.Pos.Y+p.Box.Y+(2*p.speed)))
			movement := rl.Vector2Scale(rl.Vector2Normalize(diff), p.speed)
			p.Pos = rl.Vector2Add(p.Pos, movement)
		} else {
			p.Pos.X += p.speed
		}
		p.orientation = Right
	}

	p.fallSpeed += int32(GRAVITY * delta)

	shouldUpdateY := true

	if hasCurveCollision || hasLineCollision {
		shouldUpdateY = false
	}

	if (hasCurveCollision || hasLineCollision) && spacePressed {
		shouldUpdateY = true
	}

	if shouldUpdateY {
		p.Pos.Y += float32(p.fallSpeed) * delta
	}

	p.updateCollisions()

	if ok, pos, box := p.hasBottomBoxCollision(); ok {

		if p.fallSpeed < 0 {
			p.fallSpeed *= -1
		}

		x1 := pos.X - box.X
		x2 := pos.X + box.X
		if p.Pos.X > x1 && p.Pos.X < x2 && p.Pos.Y > pos.Y {
			p.Pos.Y = pos.Y + box.Y
		}
	}

	if hasCurveCollision || hasLineCollision {
		p.fallSpeed = 0
	}

	if ok, pos, box := p.hasTopBoxCollision(); ok {
		p.fallSpeed = 0

		inaccurracy := float32(math.Min(float64(box.Y), float64(p.Box.Y)/3))

		yRangeBegin := pos.Y - inaccurracy
		yRangeEnd := pos.Y + inaccurracy

		playerBottomLine := p.Pos.Y + p.Box.Y

		if playerBottomLine >= yRangeBegin && playerBottomLine <= yRangeEnd {
			p.Pos.Y = pos.Y - p.Box.Y
		}
	}

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

func (p *Player) AddCollisionBox(cb CollisionBox) *Player {
	p.collisionBoxes = append(p.collisionBoxes, cb)
	p.collisionBoxChecks = append(p.collisionBoxChecks, CollisionBoxCheck{})
	return p
}

func (p *Player) AddCollisionLine(line *Line) *Player {
	p.collisionLines = append(p.collisionLines, line)
	p.collisionLineChecks = append(p.collisionLineChecks, CollisionLineCheck{
		Line: line,
	})
	return p
}

func (p *Player) AddCollisionBezier(bz *Bezier) *Player {
	p.collisionBeziers = append(p.collisionBeziers, bz)
	p.collisionBezierChecks = append(p.collisionBezierChecks, CollisionBezierCheck{
		Curve: bz,
	})
	return p
}

func (p *Player) GetPos() *rl.Vector2 {
	return &p.Pos
}

func (p *Player) GetBox() *rl.Vector2 {
	return &p.Box
}

func (p Player) canMoveRight() bool {
	for _, c := range p.collisionBoxChecks {
		if c.Left && c.Bottom {
			return false
		}
	}
	return true
}

func (p Player) canMoveLeft() bool {
	for _, c := range p.collisionBoxChecks {
		if c.Right && c.Bottom {
			return false
		}
	}
	return true
}

func (p *Player) hasTopBoxCollision() (bool, *rl.Vector2, *rl.Vector2) {
	for _, c := range p.collisionBoxChecks {
		if c.Top {
			return true, c.Y.GetPos(), c.Y.GetBox()
		}
	}
	return false, nil, nil
}

func (p *Player) hasBottomBoxCollision() (bool, *rl.Vector2, *rl.Vector2) {
	for _, c := range p.collisionBoxChecks {
		if c.Bottom {
			return true, c.Y.GetPos(), c.Y.GetBox()
		}
	}
	return false, nil, nil
}

func (p Player) hasCurveCollision() (bool, *CollisionBezierCheck) {
	for i, _ := range p.collisionBezierChecks {
		c := p.collisionBezierChecks[i]
		if c.Colliding {
			return true, &c
		}
	}
	return false, nil
}

func (p Player) hasLineCollision() (bool, *CollisionLineCheck) {
	for i, _ := range p.collisionLineChecks {
		c := p.collisionLineChecks[i]
		if c.Colliding {
			return true, &c
		}
	}
	return false, nil
}

func (p *Player) AddLightPoint(lp *LightPoint) *Player {
	p.LightPoints = append(p.LightPoints, lp)
	return p
}

func (p *Player) updateCollisions() {
	for i, _ := range p.collisionBoxes {
		cb := p.collisionBoxes[i]
		cb.ResolveCollision(func(bp BoxPosition) {
			res := DetectBoxCollision(p, bp)
			p.collisionBoxChecks[i].Intersected = res.Intersected
			p.collisionBoxChecks[i].Top = res.Top
			p.collisionBoxChecks[i].Bottom = res.Bottom
			p.collisionBoxChecks[i].Right = res.Right
			p.collisionBoxChecks[i].Left = res.Left
			p.collisionBoxChecks[i].X = res.X
			p.collisionBoxChecks[i].Y = res.Y
		})
	}

	for i, _ := range p.collisionLines {
		l := p.collisionLines[i]
		l.ResolveCollision(func(line *Line) {
			startLine := rl.NewVector2(p.Pos.X, p.Pos.Y+p.Box.Y)
			endLine := rl.NewVector2(p.Pos.X+p.Box.X, p.Pos.Y+p.Box.Y)
			pointCollision := rl.NewVector2(0, 0)
			isCollision := rl.CheckCollisionLines(startLine, endLine, line.Start, line.End, &pointCollision)

			p.collisionLineChecks[i].Colliding = isCollision
			p.collisionLineChecks[i].Point.X = pointCollision.X
			p.collisionLineChecks[i].Point.Y = pointCollision.Y
		})
	}

	for i, _ := range p.collisionBeziers {
		bz := p.collisionBeziers[i]
		bz.ResolveCollision(func(bezier *Bezier) {

			startLine := rl.NewVector2(p.Pos.X, p.Pos.Y+p.Box.Y)
			endLine := rl.NewVector2(p.Pos.X+p.Box.X, p.Pos.Y+p.Box.Y)

			colPoint := CheckCollisionLineBezier(
				startLine,
				endLine,
				bezier.Start,
				bezier.End,
				float64(bezier.Thick))
			p.collisionBezierChecks[i].Colliding = colPoint.Colliding
			p.collisionBezierChecks[i].Point.X = colPoint.Point.X
			p.collisionBezierChecks[i].Point.Y = colPoint.Point.Y
		})
	}
}

func (p *Player) WithShader(gs resources.GameShader) *Player {
	p.ImageShader = gs
	return p
}

func (p *Player) normalDebug() func(t *Text) {
	return func(t *Text) {

		offset := 0

		t.SetX(int32(p.Pos.X) - int32(offset))
		t.SetY(int32(p.Pos.Y))

		collisions := ""
		for i, _ := range p.collisionBezierChecks {
			cp := p.collisionBezierChecks[i]
			collisions += fmt.Sprintf("%d: %v %.1f %.1f \n", i, cp.Colliding, cp.Point.X, cp.Point.Y)
			if cp.Colliding {
				prev, next := CalculatePreviousNextPointsOfBezier(cp.Point, cp.Curve.Start, cp.Curve.End)
				collisions += fmt.Sprintf("prev {%.1f : %.1f} next {%.1f : %.1f} \n", prev.X, prev.Y, next.X, next.Y)
			}
		}

		t.SetData(fmt.Sprintf("x: %.1f y: %.1f fs: %d; \n%s", p.Pos.X, p.Pos.Y, p.fallSpeed, collisions))
	}
}
