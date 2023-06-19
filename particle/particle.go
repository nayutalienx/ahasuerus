package particle

import (
	"image/color"
	_ "image/png"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/blizzy78/twodeeparticles"
	"github.com/fogleman/ease"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type ParticleSystemSettings struct {
	MaxParticles             int
	EmissionRate             float64
	EmissionRateVariance     float64
	MoveTime                 float64
	MoveTimeVariance         float64
	FadeOutTime              float64
	StartPositionMaxDistance float64
	StartSpeed               float64
	StartSpeedVariance       float64
	StartScale               float64
	EndScale                 float64
	EndScaleVariance         float64
	MinAlpha                 float64
	Gravity                  twodeeparticles.Vector
}

func DefaultParticleSystemSettings() ParticleSystemSettings {
	return ParticleSystemSettings{
		MaxParticles:             300,
		EmissionRate:             20.0,
		EmissionRateVariance:     10.0,
		MoveTime:                 20.0,
		MoveTimeVariance:         2.0,
		FadeOutTime:              5.0,  // время исчезновения (0 - изчезнет мгновенно)
		StartPositionMaxDistance: 50.0, // оффсет от стартовой позиции
		StartSpeed:               50.0,
		StartSpeedVariance:       10.0,
		StartScale:               0.2,
		EndScale:                 0.65,
		EndScaleVariance:         0.3,
		MinAlpha:                 0.35,
		Gravity:                  twodeeparticles.Vector{0.0, 150},
	}
}

type bubbleData struct {
	speed    float64
	alpha    float64
	endScale float64
}

func (pss *ParticleSystemSettings) Bubbles() *twodeeparticles.ParticleSystem {

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	particleDataPool := &sync.Pool{}
	particleDataPool.New = func() any {
		return &bubbleData{}
	}

	s := twodeeparticles.NewSystem()

	s.MaxParticles = pss.MaxParticles

	s.DataOverLifetime = func(old any, t twodeeparticles.NormalizedDuration, delta time.Duration) any {
		if old != nil {
			return old
		}

		data := particleDataPool.Get().(*bubbleData)
		data.speed = randomValue(pss.StartSpeed-pss.StartSpeedVariance/2.0, pss.StartSpeed+pss.StartSpeedVariance/2.0, rand)
		data.endScale = randomValue(pss.EndScale-pss.EndScaleVariance/2.0, pss.EndScale+pss.EndScaleVariance/2.0, rand)
		data.alpha = randomValue(pss.MinAlpha, 1.0, rand)
		return data
	}

	s.DeathFunc = func(p *twodeeparticles.Particle) {
		particleDataPool.Put(p.Data())
	}

	s.EmissionRateOverTime = func(d time.Duration, delta time.Duration) float64 {
		q := float64(int(d.Seconds())%7)/7.0 - 0.5
		v := pss.EmissionRateVariance * q
		return pss.EmissionRate + v
	}

	s.EmissionPositionOverTime = func(d time.Duration, delta time.Duration) twodeeparticles.Vector {
		a := randomValue(0.0, 360.0, rand)
		dir := angleToDirection(a)
		return dir.Multiply(pss.StartPositionMaxDistance)
	}

	s.LifetimeOverTime = func(d time.Duration, delta time.Duration) time.Duration {
		mt := randomValue(pss.MoveTime-pss.MoveTimeVariance/2.0, pss.MoveTime+pss.MoveTimeVariance/2.0, rand)
		return time.Duration((mt+pss.FadeOutTime)*1000.0) * time.Millisecond
	}

	s.VelocityOverLifetime = func(p *twodeeparticles.Particle, t twodeeparticles.NormalizedDuration, delta time.Duration) twodeeparticles.Vector {
		data := p.Data().(*bubbleData)

		s := t.Duration(p.Lifetime()).Seconds()
		if s == 0 {
			a := randomValue(0.0, 360.0, rand)
			dir := angleToDirection(a)
			return rewind(dir.Multiply(data.speed))
		}

		moveTime := p.Lifetime().Seconds() - pss.FadeOutTime
		if s > moveTime {
			return twodeeparticles.ZeroVector
		}

		dir := p.Velocity().Normalize()
		m := 1.0 - ease.OutSine(s/moveTime)
		return rewind(dir.Multiply(data.speed * m))
	}

	s.ScaleOverLifetime = func(p *twodeeparticles.Particle, t twodeeparticles.NormalizedDuration, delta time.Duration) twodeeparticles.Vector {
		data := p.Data().(*bubbleData)

		s := t.Duration(p.Lifetime()).Seconds()
		if s == 0 {
			return twodeeparticles.Vector{pss.StartScale, pss.StartScale}
		}

		moveTime := p.Lifetime().Seconds() - pss.FadeOutTime
		if s > moveTime {
			sc := (1.0-ease.OutSine((s-moveTime)/pss.FadeOutTime))*(data.endScale-pss.StartScale) + pss.StartScale
			return twodeeparticles.Vector{sc, sc}
		}

		sc := ease.OutSine(s/moveTime)*(data.endScale-pss.StartScale) + pss.StartScale
		return twodeeparticles.Vector{sc, sc}
	}

	s.ColorOverLifetime = func(p *twodeeparticles.Particle, t twodeeparticles.NormalizedDuration, delta time.Duration) color.Color {
		data := p.Data().(*bubbleData)
		s := t.Duration(p.Lifetime()).Seconds()
		moveTime := p.Lifetime().Seconds() - pss.FadeOutTime
		if s <= moveTime {
			return color.RGBA{255, 255, 255, uint8(data.alpha * float64(t) * 255.0)}
		}

		return color.RGBA{255, 255, 255, uint8(data.alpha * (1.0 - ((s - moveTime) / pss.FadeOutTime)) * 255)}
	}

	return s
}

func (pss *ParticleSystemSettings) Fountain() *twodeeparticles.ParticleSystem {

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	s := twodeeparticles.NewSystem()

	s.MaxParticles = 500

	s.EmissionRateOverTime = constant(80.0)
	s.LifetimeOverTime = constantDuration(5 * time.Second)

	s.VelocityOverLifetime = func(p *twodeeparticles.Particle, t twodeeparticles.NormalizedDuration, delta time.Duration) twodeeparticles.Vector {
		var v twodeeparticles.Vector

		if t == 0 {
			a := 2.0 * math.Pi * randomValue(80.0, 100.0, rand) / 360.0
			s := randomValue(315.0-25.0, 315.0+25.0, rand)
			dir := angleToDirection(a)
			v = dir.Multiply(s)
		} else {
			v = p.Velocity()
		}

		return v.Add(pss.Gravity.Multiply(delta.Seconds()))
	}

	s.ScaleOverLifetime = particleConstantVector(twodeeparticles.Vector{0.2, 0.2})

	s.ColorOverLifetime = func(p *twodeeparticles.Particle, t twodeeparticles.NormalizedDuration, delta time.Duration) color.Color {
		if t == 0 {
			return color.RGBA{255, 255, 255, uint8(randomValue(pss.MinAlpha, 1.0, rand) * 255.0)}
		}

		return p.Color()
	}

	s.UpdateFunc = func(p *twodeeparticles.Particle, t twodeeparticles.NormalizedDuration, delta time.Duration) {
		if t < 0.1 || p.Position().Y < 0 {
			return
		}
		p.Kill()
	}

	return s
}

func (pss *ParticleSystemSettings) Vortex() *twodeeparticles.ParticleSystem {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	s := twodeeparticles.NewSystem()

	s.MaxParticles = 150

	s.EmissionRateOverTime = func(d time.Duration, delta time.Duration) float64 {
		if s.NumParticles() >= s.MaxParticles {
			return 0.0
		}
		return 15.0
	}

	s.LifetimeOverTime = constantDuration(24 * time.Hour)

	s.EmissionPositionOverTime = func(d time.Duration, delta time.Duration) twodeeparticles.Vector {
		a := randomValue(0.0, 360.0, rand)
		dir := angleToDirection(a)
		dist := randomValue(140.0, 160.0, rand)
		return dir.Multiply(dist)
	}

	s.VelocityOverLifetime = func(p *twodeeparticles.Particle, t twodeeparticles.NormalizedDuration, delta time.Duration) twodeeparticles.Vector {
		if t == 0 {
			dir := p.Position().Normalize()
			dir = rotate(dir, 2.0*math.Pi*-90.0/360.0)
			return dir.Multiply(200.0)
		}

		v := p.Velocity()
		s := v.Magnitude()
		dir := v.Normalize()
		a := randomValue(105.0, 115.0, rand)
		dir = rotate(dir, 2.0*math.Pi*-a/360.0*delta.Seconds())
		return dir.Multiply(s)
	}

	s.ScaleOverLifetime = func(p *twodeeparticles.Particle, t twodeeparticles.NormalizedDuration, delta time.Duration) twodeeparticles.Vector {
		if t == 0 {
			s := randomValue(0.1, 0.7, rand)
			return twodeeparticles.Vector{s, s}
		}

		return p.Scale()
	}

	s.ColorOverLifetime = func(p *twodeeparticles.Particle, t twodeeparticles.NormalizedDuration, delta time.Duration) color.Color {
		if t == 0 {
			return color.RGBA{255, 255, 255, uint8(randomValue(pss.MinAlpha, 1.0, rand) * 255.0)}
		}

		return p.Color()
	}

	return s
}

func constant(c float64) twodeeparticles.ValueOverTimeFunc {
	return func(d time.Duration, delta time.Duration) float64 {
		return c
	}
}

func constantDuration(d time.Duration) twodeeparticles.DurationOverTimeFunc {
	return func(dt time.Duration, delta time.Duration) time.Duration {
		return d
	}
}

func particleConstantVector(v twodeeparticles.Vector) twodeeparticles.ParticleVectorOverNormalizedTimeFunc {
	return func(p *twodeeparticles.Particle, t twodeeparticles.NormalizedDuration, delta time.Duration) twodeeparticles.Vector {
		return v
	}
}

func randomValue(min float64, max float64, rand *rand.Rand) float64 {
	return min + rand.Float64()*(max-min)
}

func angleToDirection(a float64) twodeeparticles.Vector {
	sin, cos := math.Sincos(a)
	return twodeeparticles.Vector{cos, -sin}
}

func rotate(v twodeeparticles.Vector, a float64) twodeeparticles.Vector {
	// https://matthew-brett.github.io/teaching/rotation_2d.html
	sin, cos := math.Sincos(a)
	return twodeeparticles.Vector{v.X*cos - v.Y*sin, v.X*sin + v.Y*cos}
}

func distance(v1 twodeeparticles.Vector, v2 twodeeparticles.Vector) float64 {
	return v1.Add(v2.Multiply(-1.0)).Magnitude()
}

func rewind(v twodeeparticles.Vector) twodeeparticles.Vector {
	if rl.IsKeyDown(rl.KeyLeftShift) {
		return v.Multiply(5)
	}
	return v
}
