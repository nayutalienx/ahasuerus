package models

import (
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

	box rl.Vector2

	collisionBoxes     []CollisionBox
	collisionBoxChecks []CollisionBoxCheck

	collisionBeziers      []CollisionBezier
	collisionBezierChecks []CollisionBezierCheck

	debugText Text
}

func NewPlayer(x float32, y float32) *Player {
	p := &Player{
		Pos:   rl.NewVector2(x, y),
		speed: 5,

		fallSpeed: 0,

		box: rl.NewVector2(100, 200),

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

func (p Player) Draw() {
	rl.DrawRectangle(int32(p.Pos.X), int32(p.Pos.Y), int32(p.box.X), int32(p.box.Y), rl.DarkPurple)
	p.debugText.Draw()
}

func (p *Player) Update(delta float32) {

	if rl.IsKeyDown(rl.KeySpace) && p.fallSpeed == 0 {
		p.fallSpeed = -JUMP_SPEED
	}

	if rl.IsKeyDown(rl.KeyLeft) && p.canMoveLeft() {
		p.Pos.X -= p.speed
	}

	if rl.IsKeyDown(rl.KeyRight) && p.canMoveRight() {
		p.Pos.X += p.speed
	}

	p.fallSpeed += int32(GRAVITY * delta)
	p.Pos.Y += float32(p.fallSpeed) * delta

	p.updateCollisions()

	if ok, pos, box := p.hasBottomCollision(); ok {

		if p.fallSpeed < 0 {
			p.fallSpeed *= -1
		}

		x1 := pos.X - box.X
		x2 := pos.X + box.X
		if p.Pos.X > x1 && p.Pos.X < x2 && p.Pos.Y > pos.Y {
			p.Pos.Y = pos.Y + box.Y
		}
	}

	if ok, pos, box := p.hasTopCollision(); ok {
		p.fallSpeed = 0

		inaccurracy := float32(math.Min(float64(box.Y), float64(p.box.Y)/3))

		yRangeBegin := pos.Y - inaccurracy
		yRangeEnd := pos.Y + inaccurracy

		playerBottomLine := p.Pos.Y + p.box.Y

		if playerBottomLine >= yRangeBegin && playerBottomLine <= yRangeEnd {
			p.Pos.Y = pos.Y - p.box.Y
		}
	}

	p.debugText.Update(delta)
}

func (p *Player) AddCollisionBox(cb CollisionBox) *Player {
	p.collisionBoxes = append(p.collisionBoxes, cb)
	p.collisionBoxChecks = append(p.collisionBoxChecks, CollisionBoxCheck{})
	return p
}

func (p *Player) AddCollisionBezier(bz *Bezier) *Player {
	p.collisionBeziers = append(p.collisionBeziers, bz)
	p.collisionBezierChecks = append(p.collisionBezierChecks, CollisionBezierCheck{
		Curve: bz,
	})
	return p
}

func (p Player) GetPos() *rl.Vector2 {
	return &p.Pos
}

func (p Player) GetBox() *rl.Vector2 {
	return &p.box
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

func (p *Player) hasTopCollision() (bool, *rl.Vector2, *rl.Vector2) {
	for _, c := range p.collisionBoxChecks {
		if c.Top {
			return true, c.Y.GetPos(), c.Y.GetBox()
		}
	}
	return false, nil, nil
}

func (p *Player) hasBottomCollision() (bool, *rl.Vector2, *rl.Vector2) {
	for _, c := range p.collisionBoxChecks {
		if c.Bottom {
			return true, c.Y.GetPos(), c.Y.GetBox()
		}
	}
	return false, nil, nil
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

	for i, _ := range p.collisionBeziers {
		bz := p.collisionBeziers[i]
		bz.ResolveCollision(func(bezier *Bezier) {

			startLine := rl.NewVector2(p.Pos.X, p.Pos.Y+p.box.Y)
			endLine := rl.NewVector2(p.Pos.X+p.box.X, p.Pos.Y+p.box.Y)

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

func (p *Player) normalDebug() func(t *Text) {
	return func(t *Text) {

		offset := 0

		t.SetX(int32(p.Pos.X) - int32(offset))
		t.SetY(int32(p.Pos.Y))

		collisions := ""
		for i, _ := range p.collisionBezierChecks {
			cp := p.collisionBezierChecks[i]
			collisions += fmt.Sprintf("%d: %v %.1f %.1f \n", i, cp.Colliding, cp.Point.X, cp.Point.Y)
		}
		// for i, c := range p.collisionsCheck {

		// 	pos1 := c.X.GetPos()
		// 	box1 := c.X.GetBox()

		// 	pos2 := c.Y.GetPos()
		// 	box2 := c.Y.GetBox()

		// 	collisions += fmt.Sprintf("%d t: %v b: %v r: %v l: %v [{%.1f:%.1f %1.f:%1.f}, {%.1f:%.1f %1.f:%1.f}]\n",
		// 		i, c.Top, c.Bottom, c.Right, c.Left, pos1.X, pos1.Y, box1.X, box1.Y, pos2.X, pos2.Y, box2.X, box2.Y)
		// }

		t.SetData(fmt.Sprintf("x: %.1f y: %.1f fs: %d; \n%s", p.Pos.X, p.Pos.Y, p.fallSpeed, collisions))
	}
}
