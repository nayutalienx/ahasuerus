package models

import (
	"ahasuerus/collision"
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

type HitboxType int

const (
	Collision HitboxType = iota
	Light
	Npc
)

type HitboxRewindData struct {
	TextCounter   string
	Choosed       string
	CurrentChoise string
}

type Hitbox struct {
	BaseEditorItem
	Type               HitboxType
	CollisionProcessor collision.CollisionDetector

	bgImage        *Image
	bgImageHidePos rl.Vector2
	screenChan     chan Object

	hasCollision bool
	drawBgImage  bool

	Rewind          []HitboxRewindData
	rewindLastIndex int32
}

func (p *Hitbox) ScreenChan(c chan Object) *Hitbox {
	p.screenChan = c
	return p
}

func (p *Hitbox) Load() {
	bgImage, ok := p.Properties["bgImage"]
	if ok {
		p.bgImage = NewImage(0, uuid.NewString(), resources.GameTexture(bgImage), 0, 0, 0, 0, 0).WithShader(resources.NpbBgImageShader)
		p.bgImage.Load()

		scale := p.PropertyFloat("bgImageScale")
		if scale > 0 {
			p.bgImage.Scale = scale
		}

	}

	if p.Type == Npc {
		p.Rewind = make([]HitboxRewindData, REWIND_BUFFER_SIZE)
	}
}

func (p *Hitbox) Unload() {
	if p.bgImage != nil {
		p.bgImage.Unload()
	}
}

func (p *Hitbox) Pause()  {}
func (p *Hitbox) Resume() {}

func (p *Hitbox) Draw() {
	if DRAW_MODELS {

		if p.Type == Collision {
			polys := p.Polygons()

			for i, _ := range polys {
				rl.DrawTriangleLines(
					polys[i].Points[0],
					polys[i].Points[1],
					polys[i].Points[2],
					rl.Blue,
				)
			}
		}

		if p.Type == Light {
			center := p.Center()
			rl.DrawCircleLines(int32(center.X), int32(center.Y), rl.Vector2Distance(p.TopLeft(), p.TopRight())/2, rl.Gold)
			rl.DrawCircleLines(int32(center.X), int32(center.Y), rl.Vector2Distance(p.TopLeft(), p.TopRight())/6, rl.Gold)
			rl.DrawCircle(int32(center.X), int32(center.Y), 10, rl.Gold)
		}

		if p.Type == Npc {
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

		p.BaseEditorItem.Draw()
	}

	if p.Type == Npc {
		if p.hasCollision || p.EditSelected {
			p.drawDialog()
		}

		if p.drawBgImage {
			if p.bgImage != nil {
				p.screenChan <- p.bgImage
			}
		}

		if !p.hasCollision && p.drawBgImage {
			if p.bgImage.Pos.X > p.bgImageHidePos.X-INACCURACY && p.bgImage.Pos.Y > p.bgImageHidePos.Y-INACCURACY {
				p.drawBgImage = false
			}
		}

	}

}

func (p *Hitbox) Update(delta float32) {

	if p.Type == Npc {

		rewindEnabled := rl.IsKeyDown(rl.KeyLeftShift)
		if rewindEnabled {
			p.rewindNpc()
		} else {
			p.saveNpcToRewind()
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

			textCounter := int32(p.PropertyFloat("textCounter"))
			currentChoice := int(p.PropertyFloat("currentChoice"))

			if rl.IsKeyReleased(rl.KeyEnter) {
				phrases := strings.Split(p.PropertyString("text"), ";")
				if len(phrases) > int(textCounter+1) {
					p.Properties["textCounter"] = fmt.Sprintf("%.1f", float32(textCounter+1))
					choosed := p.PropertyString("choosed")
					p.Properties["choosed"] = fmt.Sprintf("%s;%d", choosed, currentChoice)
					p.Properties["currentChoice"] = fmt.Sprintf("%.1f", float32(0))
				}
			}

			if rl.IsKeyReleased(rl.KeyDown) || rl.IsKeyReleased(rl.KeyUp) {

				choicesByPhrace := strings.Split(p.PropertyString("choice"), ";")
				if len(choicesByPhrace) > int(textCounter) {
					choices := strings.Split(choicesByPhrace[textCounter], ":")

					futureChoice := 0

					if rl.IsKeyReleased(rl.KeyDown) {
						futureChoice = currentChoice + 1
					}

					if rl.IsKeyReleased(rl.KeyUp) {
						futureChoice = currentChoice - 1
					}

					if len(choices) > futureChoice && futureChoice >= 0 {
						p.Properties["currentChoice"] = fmt.Sprintf("%.1f", float32(futureChoice))
					}

				}

			}

		}

	}

}

func (p *Hitbox) enterCollision() {
	if p.bgImage != nil {
		start := rl.NewVector2(WIDTH, 0)
		end := rl.NewVector2(WIDTH-(p.bgImage.WidthHeight.X*p.bgImage.Scale), 0)
		p.bgImage.StartMove(start, end, 10)
		p.drawBgImage = true
	}
}

func (p *Hitbox) exitCollision() {
	if p.bgImage != nil {
		p.bgImageHidePos = rl.NewVector2(WIDTH+INACCURACY, 0)
		p.bgImage.StartMove(p.bgImage.Pos, p.bgImageHidePos, 10)
	}
}

func (p *Hitbox) saveNpcToRewind() {
	if int(p.rewindLastIndex) == len(p.Rewind)-1 {
		p.rewindLastIndex = 0
	}

	p.Rewind[p.rewindLastIndex] = HitboxRewindData{
		TextCounter:   p.Properties["textCounter"],
		Choosed:       p.Properties["choosed"],
		CurrentChoise: p.Properties["currentChoise"],
	}
	p.rewindLastIndex++
}

func (p *Hitbox) rewindNpc() {
	if p.rewindLastIndex > 0 {
		rewind := p.Rewind[p.rewindLastIndex-1]
		p.Properties["textCounter"] = rewind.TextCounter
		p.Properties["choosed"] = rewind.Choosed
		p.Properties["currentChoise"] = rewind.CurrentChoise
		p.rewindLastIndex--
	}
}

func (p Hitbox) drawSelectDialog(dialogRec rl.Rectangle) {

	choice := p.PropertyString("choice")
	if choice == "" {
		return
	}

	fontSize := int32(p.PropertyFloat("fontSize"))
	textOffsetX := p.PropertyFloat("textOffsetX")
	textOffsetY := p.PropertyFloat("textOffsetY")

	textCounter := int32(p.PropertyFloat("textCounter"))

	choicesByPhrace := strings.Split(choice, ";")

	if len(choicesByPhrace) > int(textCounter) {

		rectColor := rl.Black
		rectColor.A = 150
		dialogRec.X += dialogRec.Width / 1.1
		dialogRec.Y += dialogRec.Height / 2

		choices := strings.Split(choicesByPhrace[textCounter], ":")

		dialogRec.Height = float32(fontSize*int32(len(choices))) + textOffsetY
		rl.DrawRectangleRounded(dialogRec, 0.5, 0, rectColor)

		for i, _ := range choices {
			choice := choices[i]

			textPos := rl.NewVector2(
				dialogRec.X+textOffsetX,
				dialogRec.Y+textOffsetY/2+float32(fontSize)*float32(i),
			)

			color := rl.White
			if i == int(p.PropertyFloat("currentChoice")) {
				color = rl.Orange
			}

			rl.DrawTextEx(resources.LoadFont(resources.Literata), choice, textPos, float32(fontSize), 2, color)

		}

	}

}

func (p Hitbox) drawDialog() {
	pos := p.TopRight()

	offsetX := int32(p.PropertyFloat("blockOffsetX"))
	offsetY := int32(p.PropertyFloat("blockOffsetY"))

	fontSize := int32(p.PropertyFloat("fontSize"))

	textOffsetX := p.PropertyFloat("textOffsetX")
	textOffsetY := p.PropertyFloat("textOffsetY")
	phrases := strings.Split(p.PropertyString("text"), ";")
	textCounter := int32(p.PropertyFloat("textCounter"))

	text := "empty phrase"
	if len(phrases) > int(textCounter) {
		text = phrases[textCounter]
	}

	maxXLen := 0
	splittenByNewLine := strings.Split(text, "\n")
	for i, _ := range splittenByNewLine {
		if len(splittenByNewLine[i]) > maxXLen {
			maxXLen = len(splittenByNewLine[i])
		}
	}

	width := int32(maxXLen * int(float64(fontSize)/2.0))
	height := int32(float64(fontSize)+(float64(fontSize)/1.5)) * (1 + (int32(strings.Count(text, "\n"))))

	if width < 400 {
		width = 400
	}

	rectColor := rl.Black
	rectColor.A = 150

	roundedRec := rl.NewRectangle(float32(int32(pos.X)+offsetX), float32(int32(pos.Y)+offsetY), float32(width), float32(height))

	rl.DrawRectangleRounded(roundedRec, 0.5, 0, rectColor)

	textPos := rl.NewVector2(float32(int32(pos.X)+offsetX+int32(textOffsetX)), float32(int32(pos.Y)+offsetY+int32(textOffsetY)))
	rl.DrawTextEx(resources.LoadFont(resources.Literata), text, textPos, float32(fontSize), 2, rl.White)

	p.drawSelectDialog(roundedRec)
}

func (p Hitbox) getDynamicHitbox() collision.Hitbox {
	topLeft := p.TopLeft()
	bottomRight := p.BottomRight()
	width := bottomRight.X - topLeft.X
	height := bottomRight.Y - topLeft.Y
	hb := GetDynamicHitboxFromMap(GetDynamicHitboxMap(topLeft, width, height))
	return hb
}
