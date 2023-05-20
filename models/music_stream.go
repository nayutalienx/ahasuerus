package models

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type MusicStream struct {
	resourcePath string
	music *rl.Music
}

func NewMusicStream(resourcePath string) *MusicStream {
	return &MusicStream{
		resourcePath: resourcePath,
	}
}

func (p *MusicStream) Draw() {
}

func (p *MusicStream) Update(delta float32) {
	rl.UpdateMusicStream(*p.music)
}

func (p *MusicStream) Load() {
	musicTheme := rl.LoadMusicStream(p.resourcePath)
	rl.PlayMusicStream(musicTheme)
	p.music = &musicTheme
}

func (p *MusicStream) Unload() {
	rl.UnloadMusicStream(*p.music)
}