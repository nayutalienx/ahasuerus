package models

import (
	"ahasuerus/audio"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type MusicStream struct {
	directResourcePath  string
	reverseResourcePath string

	directAudioPanel  *audio.AudioPanel
	reverseAudioPanel *audio.AudioPanel
	currentAudioPanel *audio.AudioPanel

	directAudioSpeed  float64
	reverseAudioSpeed float64

	rewindStarted bool
	rewindSpeed   int32 `json:"-"`
}

func NewMusicStream(directResourcePath, reverseResourcePath string) *MusicStream {
	return &MusicStream{
		directResourcePath:  directResourcePath,
		reverseResourcePath: reverseResourcePath,
		directAudioSpeed:    1.0,
		reverseAudioSpeed:   1.0,
		rewindSpeed:         1,
	}
}

func (p *MusicStream) GetDrawIndex() int {
	return 0
}
func (p *MusicStream) GetId() string {
	return "music-theme"
}

func (p *MusicStream) Draw() {
}

func (p *MusicStream) Update(delta float32) {

	rewindEnabled := rl.IsKeyDown(rl.KeyLeftShift)
	if rewindEnabled {

		if p.rewindStarted {
			p.rewindSpeed = 1
			p.rewindStarted = false
		}

		p.updateRewindSpeed()

		if p.rewindSpeed > 0 {
			if p.currentAudioPanel != p.reverseAudioPanel {
				p.currentAudioPanel = p.reverseAudioPanel
				p.directAudioPanel.Pause()
				mirrorPos := p.directAudioPanel.Length() - p.directAudioPanel.Position()
				p.reverseAudioPanel.SetPosition(mirrorPos)
				p.reverseAudioPanel.Unpause()
			}
		}

		if p.rewindSpeed < 0 {

			if p.currentAudioPanel != p.directAudioPanel {
				p.currentAudioPanel = p.directAudioPanel
				p.reverseAudioPanel.Pause()
				mirrorPos := p.reverseAudioPanel.Length() - p.reverseAudioPanel.Position()
				p.directAudioPanel.SetPosition(mirrorPos)
				p.directAudioPanel.Unpause()
			}

		}

		if p.rewindSpeed != 0 {

			if p.currentAudioPanel.IsPaused() {
				p.currentAudioPanel.Unpause()
			}

			p.currentAudioPanel.SetSpeed(math.Abs(float64(p.rewindSpeed)))
		} else {
			p.currentAudioPanel.SetSpeed(0.1)
		}

	} else {

		p.rewindStarted = true

		if p.currentAudioPanel != p.directAudioPanel {
			p.currentAudioPanel = p.directAudioPanel
			p.reverseAudioPanel.Pause()
			mirrorPos := p.reverseAudioPanel.Length() - p.reverseAudioPanel.Position()
			p.directAudioPanel.SetPosition(mirrorPos)
		}

		if p.directAudioPanel.IsPaused() {
			p.directAudioPanel.Unpause()
			p.directAudioPanel.SetSpeed(1)
		}

	}

}

func (p *MusicStream) Load() {
	p.directAudioPanel = audio.NewAudioPanel(p.directResourcePath)
	p.directAudioPanel.SetVolume(-3.0)

	p.reverseAudioPanel = audio.NewAudioPanel(p.reverseResourcePath)
	p.reverseAudioPanel.SetVolume(-3.0)

	p.directAudioPanel.Play()
	p.reverseAudioPanel.Play()

	p.reverseAudioPanel.Pause()

	p.currentAudioPanel = p.directAudioPanel
}

func (p *MusicStream) Unload() {

	err := p.directAudioPanel.Close()
	if err != nil {
		panic(err)
	}

	err = p.reverseAudioPanel.Close()
	if err != nil {
		panic(err)
	}

}

func (p *MusicStream) Resume() {
	p.directAudioPanel.Unpause()
}

func (p *MusicStream) Pause() {
	p.directAudioPanel.Pause()
	p.reverseAudioPanel.Pause()
}

func (p *MusicStream) updateRewindSpeed() {
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
