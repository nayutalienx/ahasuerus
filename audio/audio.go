package audio

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type AudioPanel struct {
	sampleRate beep.SampleRate
	streamer   beep.StreamSeekCloser
	ctrl       *beep.Ctrl
	resampler  *beep.Resampler
	volume     *effects.Volume
}

func NewAudioPanel(path string) *AudioPanel {

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		panic(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/27))

	return newAudioPanel(format.SampleRate, streamer)
}

func newAudioPanel(sampleRate beep.SampleRate, streamer beep.StreamSeekCloser) *AudioPanel {
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}
	resampler := beep.ResampleRatio(4, 1, ctrl)
	volume := &effects.Volume{Streamer: resampler, Base: 2}
	return &AudioPanel{sampleRate, streamer, ctrl, resampler, volume}
}

func (ap *AudioPanel) Play() {
	speaker.Play(ap.volume)
}

func (ap *AudioPanel) IsPaused() bool {
	return ap.ctrl.Paused
}

func (ap *AudioPanel) Pause() {
	ap.ctrl.Paused = true
}

func (ap *AudioPanel) Unpause() {
	ap.ctrl.Paused = false
}

func (ap *AudioPanel) Position() int {
	position := ap.streamer.Position()
	return position
}

func (ap *AudioPanel) Length() int {
	length := ap.streamer.Len()
	return length
}

func (ap *AudioPanel) Volume() float64 {
	volume := ap.volume.Volume
	return volume
}

func (ap *AudioPanel) Speed() float64 {
	speed := ap.resampler.Ratio()
	return speed
}

func (ap *AudioPanel) SetPosition(newPos int) {
	if err := ap.streamer.Seek(newPos); err != nil {
		panic(err)
	}
}

func (ap *AudioPanel) SetVolume(newVol float64) {
	ap.volume.Volume = newVol
}

func (ap *AudioPanel) SetSpeed(newSpeed float64) {
	ap.resampler.SetRatio(newSpeed)
}

func (ap *AudioPanel) Close() error {
	return ap.streamer.Close()
}
