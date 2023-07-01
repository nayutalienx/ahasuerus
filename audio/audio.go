package audio

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var (
	speakerInited = false
	volume     *effects.Volume
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

	if !speakerInited {
		err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/27))
		if err != nil {
			panic(err)
		}	
		speakerInited = true
	}

	return newAudioPanel(format.SampleRate, streamer)
}

func newAudioPanel(sampleRate beep.SampleRate, streamer beep.StreamSeekCloser) *AudioPanel {
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}
	resampler := beep.ResampleRatio(4, 1, ctrl)
	if volume == nil {
		volume = &effects.Volume{Streamer: resampler, Base: 2}
		speaker.Play(volume)
	}
	return &AudioPanel{sampleRate, streamer, ctrl, resampler, volume}
}

func (ap *AudioPanel) IsPaused() bool {
	speaker.Lock()
	defer speaker.Unlock()
	return ap.ctrl.Paused
}

func (ap *AudioPanel) Pause() {
	speaker.Lock()
	defer speaker.Unlock()
	ap.ctrl.Paused = true
}

func (ap *AudioPanel) Unpause() {
	speaker.Lock()
	defer speaker.Unlock()
	ap.ctrl.Paused = false
	volume.Streamer = ap.resampler
}

func (ap *AudioPanel) Position() int {
	speaker.Lock()
	defer speaker.Unlock()
	position := ap.streamer.Position()
	return position
}

func (ap *AudioPanel) Length() int {
	speaker.Lock()
	defer speaker.Unlock()
	length := ap.streamer.Len()
	return length
}

func (ap *AudioPanel) Volume() float64 {
	speaker.Lock()
	defer speaker.Unlock()
	volume := ap.volume.Volume
	return volume
}

func (ap *AudioPanel) Speed() float64 {
	speaker.Lock()
	defer speaker.Unlock()
	speed := ap.resampler.Ratio()
	return speed
}

func (ap *AudioPanel) SetPosition(newPos int) {
	speaker.Lock()
	defer speaker.Unlock()
	if err := ap.streamer.Seek(newPos); err != nil {
		panic(err)
	}
}

func (ap *AudioPanel) SetVolume(newVol float64) {
	speaker.Lock()
	defer speaker.Unlock()
	ap.volume.Volume = newVol
}

func (ap *AudioPanel) SetSpeed(newSpeed float64) {
	speaker.Lock()
	defer speaker.Unlock()
	ap.resampler.SetRatio(newSpeed)
}

func (ap *AudioPanel) Close() error {
	return ap.streamer.Close()
}
