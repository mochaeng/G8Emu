package main

import (
	"github.com/mochaeng/G8Emu/internal/core"
	"github.com/mochaeng/G8Emu/internal/emulator"
)

var (
	game     *emulator.Game
	platform *emulator.Platform
	chip8    *core.Chip8
)

func main() {
	const scale = 10
	const frequency = 540

	platform = emulator.NewPlatform(scale)
	chip8 = core.NewChip8()
	game = emulator.NewGame(platform, chip8, frequency)

	js
}
