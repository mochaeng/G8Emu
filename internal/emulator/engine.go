package emulator

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mochaeng/G8Emu/internal/core"
)

type Engine struct {
	platform        *Platform
	chip8           *core.Chip8
	cpuFrequency    int
	lastUpdate      time.Time
	lastTimer       time.Time
	timeAccumulator time.Duration
	cycleTime       time.Duration

	pausedKeyPressed bool
}

func NewGame(platform *Platform, chip8 *core.Chip8, cpuFrequency int) *Engine {
	return &Engine{
		platform:   platform,
		chip8:      chip8,
		lastUpdate: time.Now(),
		cycleTime:  time.Second / time.Duration(cpuFrequency),
	}
}

func (e *Engine) Update() error {
	e.platform.ProcessInput(e.chip8.Keypad[:])

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		// return ebiten.Termination
	}

	isPause := ebiten.IsKeyPressed(ebiten.KeyP)
	if isPause && !e.pausedKeyPressed {
		e.chip8.TogglePause()
		e.pausedKeyPressed = true
	} else if !isPause {
		e.pausedKeyPressed = false
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		e.Reset()
		return nil
	}

	currentTime := time.Now()
	elapsed := currentTime.Sub(e.lastUpdate)
	e.lastUpdate = currentTime
	e.timeAccumulator += elapsed

	for e.timeAccumulator >= e.cycleTime {
		e.chip8.Cycle()
		e.timeAccumulator -= e.cycleTime
	}

	if !e.chip8.IsPaused() && currentTime.Sub(e.lastTimer) >= time.Second/60 {
		if e.chip8.DelayTimer > 0 {
			e.chip8.DelayTimer--
		}
		if e.chip8.SoundTimer > 0 {
			e.chip8.SoundTimer--
		}
		e.lastTimer = time.Now()
	}

	return nil
}

func (e *Engine) Draw(screen *ebiten.Image) {
	e.platform.UpdateDisplay(e.chip8.Video[:])
	e.platform.Draw(screen)
}

func (e *Engine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return e.platform.Layout(outsideWidth, outsideHeight)
}

func (e *Engine) Reset() {
	e.chip8.Reset()

	e.lastUpdate = time.Now()
	e.lastTimer = time.Now()
	e.timeAccumulator = 0

	e.pausedKeyPressed = false
}

func (e *Engine) Pause() {
	e.chip8.Pause()
}

func (e *Engine) Resume() {
	e.chip8.Resume()
}

func (e *Engine) IsPaused() bool {
	return e.chip8.IsPaused()
}
