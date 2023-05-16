package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type textUpdateCallback func(text *Text)

type Text struct {
	x int32
	y int32
	data string
	color rl.Color
	fontSize int32
	updateCallback textUpdateCallback
}

func NewText(x int32, y int32) *Text {
	return &Text{
		x: x,
		y: y,
		color: rl.DarkGray,
		fontSize: 20,
		updateCallback: func(text *Text) {},
	}
}

func (p Text) Draw() {
	rl.DrawText(p.data, p.x, p.y, p.fontSize, p.color)
}

func (p *Text) SetUpdateCallback(callback textUpdateCallback) *Text {
	p.updateCallback = callback
	return p
}

func (p *Text) Update(delta float32) {
	p.updateCallback(p)
}

func (p Text) GetX() int32 {
	return p.x
}

func (p Text) GetY() int32 {
	return p.y
}

func (p *Text) SetX(x int32) {
	p.x = x
}

func (p *Text) SetY(y int32) {
	p.y = y
}

func (p *Text) SetFontSize(size int32) *Text {
	p.fontSize = size
	return p
}

func (p *Text) SetData(data string) *Text {
	p.data = data
	return p
}

func (p *Text) SetColor(col rl.Color) *Text {
	p.color = col
	return p
}