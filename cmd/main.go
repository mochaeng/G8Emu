package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mochaeng/G8Emu/internal/constants"
	"github.com/mochaeng/G8Emu/internal/emulator"
	"github.com/mochaeng/G8Emu/internal/engine"
	"github.com/mochaeng/G8Emu/internal/platform"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <Scale> <Delay> <ROM>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "   Scale: Integer scale factor (e.g., 10)\n")
		fmt.Fprintf(os.Stderr, "   Delay: Milliseconds between cycles (e.g., 1)\n")
		fmt.Fprintf(os.Stderr, "   ROM: Path to ROM file\n")
		os.Exit(1)
	}

	videoScale, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("invalid scale factor: %v", err)
	}

	cycleDelay, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("invalid cycle delay: %v", err)
	}

	romFilename := os.Args[3]

	platform := platform.NewPlatform(videoScale)
	chip8 := emulator.NewChip8()

	if err := chip8.LoadRom(romFilename); err != nil {
		log.Fatalf("failed to load ROM: %v", err)
	}

	ebiten.SetWindowSize(constants.VIDEO_WIDTH*constants.SCALE_FACTOR, constants.VIDEO_HEIGHT*constants.SCALE_FACTOR)
	ebiten.SetWindowTitle("G8Emu")

	game := engine.NewGame(platform, chip8, cycleDelay)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
