package models

import (
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type NpcDialog struct {
	CurrentInteraction uint
	Interactions       []NpcInteraction
}

type NpcInteraction struct {
	Text          string
	Options       []string
	Routes        []uint
	CurrentOption uint
}

func (p *NpcDialog) Draw() {

	p.drawDialog()

}

func (p *NpcDialog) Update(delta float32) {

	if rl.IsKeyReleased(rl.KeyEnter) {
		currentInteraction := p.Interactions[p.CurrentInteraction]
		p.CurrentInteraction = currentInteraction.Routes[currentInteraction.CurrentOption]
	}

	if rl.IsKeyReleased(rl.KeyDown) {
		p.Interactions[p.CurrentInteraction].CurrentOption++
	}

	if rl.IsKeyReleased(rl.KeyUp) {
		p.Interactions[p.CurrentInteraction].CurrentOption--
	}

	maxOption := uint(len(p.Interactions[p.CurrentInteraction].Options))
	if maxOption > 0 {
		maxOption--
	}

	if p.Interactions[p.CurrentInteraction].CurrentOption > maxOption {
		p.Interactions[p.CurrentInteraction].CurrentOption = 0
	}

}

func (p *NpcDialog) GetDrawIndex() int {
	return -999
}

func (p *NpcDialog) GetId() string {
	return ""
}

func (p *NpcDialog) drawDialog() {

	interaction := p.Interactions[p.CurrentInteraction]

	npcTextRows := p.splitNpcTextAndAddLine(interaction.Text, 80)

	allRows := append(npcTextRows, interaction.Options...)

	dialogRectangle, positions := p.getRectangleForRows(allRows)

	rl.DrawRectangleRounded(dialogRectangle, 0, 0, rl.NewColor(0, 0, 0, 150))

	for i, _ := range positions {
		color := rl.White
		if i == len(npcTextRows)+int(interaction.CurrentOption) {
			color = rl.Orange
		}
		DrawSdfText(allRows[i], positions[i], 60, color)
	}

}

func (p *NpcDialog) getRectangleForRows(rows []string) (rl.Rectangle, []rl.Vector2) {

	width := WIDTH - (WIDTH / 2)
	height := HEIGHT / 5

	posx := WIDTH/2 - width/2
	posy := (HEIGHT/2 + HEIGHT/4)

	positions := []rl.Vector2{}
	for i, _ := range rows {
		rowPosx := posx + 20
		rowPosy := posy + 10 + float32(i*50)

		positions = append(positions, rl.Vector2{
			X: rowPosx,
			Y: rowPosy,
		})
	}

	return rl.NewRectangle(posx, posy, width, height), positions
}

func (p *NpcDialog) splitNpcTextAndAddLine(text string, maxCharsOnRow int) []string {
	result := []string{}
	for i := 0; i <= len(text)/maxCharsOnRow; i++ {
		positionStart := i * maxCharsOnRow
		positionEnd := (i + 1) * maxCharsOnRow
		if positionEnd > len(text) {
			positionEnd = len(text)
		}
		result = append(result, text[positionStart:positionEnd])
	}
	result = append(result, strings.Repeat("-", maxCharsOnRow))
	return result
}
