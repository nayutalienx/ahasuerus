package models

import (
	"ahasuerus/config"
	"ahasuerus/resources"
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

var (
	WIDTH, HEIGHT = config.GetResolution()
	INACCURACY    = float32(10)
)

type Npc struct {
	CollisionHitbox

	BgImagePath   string
	BgImageScale  float32
	TextCounter   int
	CurrentChoise int
	Text          string
	Choosed       string
	Choice        string
	FontSize      int32
	TextOffset    rl.Vector2
	BlockOffset   rl.Vector2

	bgImage        *Image
	drawBgImage    bool        `json:"-"`
	bgImageHidePos rl.Vector2  `json:"-"`
	screenChan     chan Object `json:"-"`

	Rewind               [REWIND_BUFFER_SIZE]HitboxRewindData `json:"-"`
	rewindLastIndex      int32              `json:"-"`
	rewindSpeed          int32              `json:"-"`
	rewindModeStartIndex int32              `json:"-"`
	rewindModeStarted    bool               `json:"-"`
}

type HitboxRewindData struct {
	TextCounter   int
	Choosed       string
	CurrentChoise int
}

func (p *Npc) ScreenChan(c chan Object) *Npc {
	p.screenChan = c
	return p
}

func (p *Npc) Load() {
	if p.BgImagePath != "" {
		p.bgImage = NewImage(0, uuid.NewString(), resources.GameTexture(p.BgImagePath), 0, 0, 0).
			WithShader(resources.NpcShader)
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
		p.drawDialog()
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

	if p.hasCollision {

		if rl.IsKeyReleased(rl.KeyEnter) {
			phrases := strings.Split(p.Text, ";")
			if len(phrases) > int(p.TextCounter+1) {
				p.TextCounter++
				p.Choosed = fmt.Sprintf("%s;%d", p.Choosed, p.CurrentChoise)
				p.CurrentChoise = 0
			}
		}

		if rl.IsKeyReleased(rl.KeyDown) || rl.IsKeyReleased(rl.KeyUp) {

			choicesByPhrace := strings.Split(p.Choice, ";")
			if len(choicesByPhrace) > int(p.TextCounter) {
				choices := strings.Split(choicesByPhrace[p.TextCounter], ":")

				futureChoice := 0

				if rl.IsKeyReleased(rl.KeyDown) {
					futureChoice = p.CurrentChoise + 1
				}

				if rl.IsKeyReleased(rl.KeyUp) {
					futureChoice = p.CurrentChoise - 1
				}

				if len(choices) > futureChoice && futureChoice >= 0 {
					p.CurrentChoise = futureChoice
				}

			}

		}

	}
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
		TextCounter:   p.TextCounter,
		Choosed:       p.Choosed,
		CurrentChoise: p.CurrentChoise,
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

	p.TextCounter = rewind.TextCounter
	p.Choosed = rewind.Choosed
	p.CurrentChoise = rewind.CurrentChoise
}

func (p Npc) drawSelectDialog(dialogRec rl.Rectangle) {

	if p.Choice == "" {
		return
	}

	choicesByPhrace := strings.Split(p.Choice, ";")

	if len(choicesByPhrace) > int(p.TextCounter) {

		rectColor := rl.Black
		rectColor.A = 150
		dialogRec.X += dialogRec.Width / 1.1
		dialogRec.Y += dialogRec.Height / 2

		choices := strings.Split(choicesByPhrace[p.TextCounter], ":")

		dialogRec.Height = float32(p.FontSize*int32(len(choices))) + p.TextOffset.Y
		rl.DrawRectangleRounded(dialogRec, 0.5, 0, rectColor)

		for i, _ := range choices {
			choice := choices[i]

			textPos := rl.NewVector2(
				dialogRec.X+p.TextOffset.X,
				dialogRec.Y+p.TextOffset.Y/2+float32(p.FontSize)*float32(i),
			)

			color := rl.White
			if i == int(p.CurrentChoise) {
				color = rl.Orange
			}

			DrawSdfText(choice, textPos, float32(p.FontSize), color)

		}

	}

}

func (p Npc) drawDialog() {
	pos := p.TopRight()

	phrases := strings.Split(p.Text, ";")

	text := "empty phrase"
	if len(phrases) > int(p.TextCounter) {
		text = phrases[p.TextCounter]
	}

	maxXLen := 0
	splittenByNewLine := strings.Split(text, "\n")
	for i, _ := range splittenByNewLine {
		if len(splittenByNewLine[i]) > maxXLen {
			maxXLen = len(splittenByNewLine[i])
		}
	}

	width := int32(maxXLen * int(float64(p.FontSize)/2.0))
	height := int32(float64(p.FontSize)+(float64(p.FontSize)/1.5)) * (1 + (int32(strings.Count(text, "\n"))))

	if width < 400 {
		width = 400
	}

	rectColor := rl.Black
	rectColor.A = 150

	roundedRec := rl.NewRectangle(float32(int32(pos.X)+int32(p.BlockOffset.X)), float32(int32(pos.Y)+int32(p.BlockOffset.Y)), float32(width), float32(height))

	rl.DrawRectangleRounded(roundedRec, 0.5, 0, rectColor)

	textPos := rl.NewVector2(float32(int32(pos.X)+int32(p.BlockOffset.X)+int32(p.TextOffset.X)), float32(int32(pos.Y)+int32(p.BlockOffset.Y)+int32(p.TextOffset.Y)))

	DrawSdfText(text, textPos, float32(p.FontSize), rl.White)

	p.drawSelectDialog(roundedRec)
}
