package engine

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mochaeng/G8Emu/internal/emulator"
	"github.com/mochaeng/G8Emu/internal/platform"
)

type Game struct {
	platform      *platform.Platform
	chip8         *emulator.Chip8
	cycleDelay    time.Duration
	lastCycleTime time.Time
}

func NewGame(platform *platform.Platform, chip8 *emulator.Chip8, cycleDelay int) *Game {
	return &Game{
		platform:      platform,
		chip8:         chip8,
		cycleDelay:    time.Duration(cycleDelay) * time.Millisecond,
		lastCycleTime: time.Now(),
	}
}

func (g *Game) Update() error {
	quit := g.platform.ProcessInput(g.chip8.Keypad[:])
	if quit {
		return ebiten.Termination
	}

	currentTime := time.Now()
	dt := currentTime.Sub(g.lastCycleTime)
	if dt >= g.cycleDelay {
		g.lastCycleTime = currentTime
		g.chip8.Cycle()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.platform.UpdateDisplay(g.chip8.Video[:])
	g.platform.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.platform.Layout(outsideWidth, outsideHeight)
}
