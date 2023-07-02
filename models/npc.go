package models

import (
	"ahasuerus/config"
	"ahasuerus/resources"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

var (
	WIDTH, HEIGHT = config.GetResolution()
	INACCURACY    = float32(10)
)

type Npc struct {
	CollisionHitbox

	Dialogues NpcDialog

	BgImagePath  string
	BgImageScale float32

	bgImage        *Image
	drawBgImage    bool        `json:"-"`
	bgImageHidePos rl.Vector2  `json:"-"`
	screenChan     chan Object `json:"-"`
	screenScale    float32     `json:"-"`

	Rewind               [REWIND_BUFFER_SIZE]HitboxRewindData `json:"-"`
	rewindLastIndex      int32                                `json:"-"`
	rewindSpeed          int32                                `json:"-"`
	rewindModeStartIndex int32                                `json:"-"`
	rewindModeStarted    bool                                 `json:"-"`
}

type HitboxRewindData struct {
	Dialogues NpcDialog
}

func (p *Npc) ScreenChan(c chan Object) *Npc {
	p.screenChan = c
	return p
}

func (p *Npc) ScreenScale(scale float32) *Npc {
	p.screenScale = scale
	p.Dialogues.ScreenScale(scale)
	return p
}

func (p *Npc) Load() {
	if p.BgImagePath != "" {
		p.bgImage = NewImage(uuid.NewString(), resources.GameTexture(p.BgImagePath), 0, 0, 0).
			WithShader(resources.NpcShader)
		p.bgImage.Scale = p.screenScale
		p.bgImage.Load()

		if p.BgImageScale > 0 {
			p.bgImage.Scale = p.BgImageScale
		}

	}
}

func (p *Npc) Unload() {
	if p.bgImage != nil {
		p.bgImage.Unload()
	}
}

func (p *Npc) Draw() {
	if DRAW_MODELS {
		polys := p.getDynamicHitbox().Polygons
		for i, _ := range polys {
			rl.DrawTriangleLines(
				polys[i].Points[0],
				polys[i].Points[1],
				polys[i].Points[2],
				rl.Pink,
			)
		}
	}

	if p.hasCollision || p.EditSelected {
		p.screenChan <- &p.Dialogues
	}

	if p.drawBgImage {
		if p.bgImage != nil {
			p.screenChan <- p.bgImage
		}
	}

	if !p.hasCollision && p.drawBgImage {
		if p.bgImage.TopLeft().X > p.bgImageHidePos.X-INACCURACY && p.bgImage.TopLeft().Y > p.bgImageHidePos.Y-INACCURACY {
			p.drawBgImage = false
		}
	}

	p.BaseEditorItem.Draw()
}

func (p *Npc) Update(delta float32) {

	rewindEnabled := rl.IsKeyDown(rl.KeyLeftShift)
	if rewindEnabled {
		p.updateRewindSpeed()
		p.rewindNpc()
		p.rewindModeStarted = true
	} else {
		p.saveNpcToRewind()
		p.rewindModeStarted = false
	}

	detectedCollision, _ := p.CollisionProcessor.Detect(p.getDynamicHitbox())

	if !p.hasCollision && detectedCollision {
		p.enterCollision()
	}

	if p.hasCollision && !detectedCollision {
		p.exitCollision()
	}

	p.hasCollision = detectedCollision
}

func (p *Npc) updateRewindSpeed() {
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

func (p *Npc) enterCollision() {
	if p.bgImage != nil {
		start := rl.NewVector2(WIDTH, 0)
		end := rl.NewVector2(WIDTH-(p.bgImage.Width()*p.bgImage.Scale), 0)
		p.bgImage.StartMove(start, end, 10)
		p.drawBgImage = true
	}
}

func (p *Npc) exitCollision() {
	if p.bgImage != nil {
		p.bgImageHidePos = rl.NewVector2(WIDTH+INACCURACY, 0)
		p.bgImage.StartMove(p.bgImage.TopLeft(), p.bgImageHidePos, 10)
	}
}

func (p *Npc) saveNpcToRewind() {
	if int(p.rewindLastIndex) == len(p.Rewind)-1 {
		p.rewindLastIndex = 0
	}

	p.Rewind[p.rewindLastIndex] = HitboxRewindData{
		Dialogues: p.Dialogues,
	}
	p.rewindLastIndex++
}

func (p *Npc) rewindNpc() {
	if !p.rewindModeStarted {
		p.rewindModeStartIndex = p.rewindLastIndex
		p.rewindSpeed = 1
	}

	rewind := p.Rewind[p.rewindLastIndex]

	if p.rewindLastIndex > p.rewindSpeed && p.rewindLastIndex < p.rewindModeStartIndex+p.rewindSpeed {
		rewind = p.Rewind[p.rewindLastIndex-p.rewindSpeed]
		p.rewindLastIndex -= p.rewindSpeed
	}

	p.Dialogues = rewind.Dialogues
}
