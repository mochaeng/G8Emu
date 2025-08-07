package emulator

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mochaeng/G8Emu/internal/core"
)

type Game struct {
	platform        *Platform
	chip8           *core.Chip8
	cpuFrequency    int
	lastUpdate      time.Time
	lastTimer       time.Time
	timeAccumulator time.Duration
	cycleTime       time.Duration
}

func NewGame(platform *Platform, chip8 *core.Chip8, cpuFrequency int) *Game {
	return &Game{
		platform:   platform,
		chip8:      chip8,
		lastUpdate: time.Now(),
		cycleTime:  time.Second / time.Duration(cpuFrequency),
	}
}

func (g *Game) Update() error {
	quit := g.platform.ProcessInput(g.chip8.Keypad[:])
	if quit {
		return ebiten.Termination
	}

	currentTime := time.Now()
	elapsed := currentTime.Sub(g.lastUpdate)
	g.lastUpdate = currentTime
	g.timeAccumulator += elapsed

	// cyclesPerSecond := 540
	// cycleTime := time.Second / time.Duration(cyclesPerSecond)

	for g.timeAccumulator >= g.cycleTime {
		g.chip8.Cycle()
		g.timeAccumulator -= g.cycleTime
	}

	if currentTime.Sub(g.lastTimer) >= time.Second/60 {
		if g.chip8.DelayTimer > 0 {
			g.chip8.DelayTimer--
		}
		if g.chip8.SoundTimer > 0 {
			g.chip8.SoundTimer--
		}
		g.lastTimer = time.Now()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// if g.frameCount%2 == 0 {
	// g.platform.UpdateDisplay(g.chip8.Video[:])
	// }
	// g.platform.UpdateDisplay(g.chip8.Video[:])
	g.platform.UpdateDisplay(g.chip8.Video[:])
	g.platform.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.platform.Layout(outsideWidth, outsideHeight)
}
