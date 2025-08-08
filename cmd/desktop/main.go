package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mochaeng/G8Emu/internal/constants"
	"github.com/mochaeng/G8Emu/internal/core"
	"github.com/mochaeng/G8Emu/internal/emulator"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <Scale> <Delay> <ROM>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "   Scale: Integer scale factor (e.g., 10)\n")
		// fmt.Fprintf(os.Stderr, "   Delay: Milliseconds between cycles (e.g., 1)\n")
		fmt.Fprintf(os.Stderr, "   ROM: Path to ROM file\n")
		os.Exit(1)
	}

	videoScale, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("invalid scale factor: %v", err)
	}

	romFilename := os.Args[2]

	platform := emulator.NewPlatform(videoScale)
	chip8 := core.NewChip8()

	if err := chip8.LoadRomFile(romFilename); err != nil {
		log.Fatalf("failed to load ROM: %v", err)
	}

	ebiten.SetWindowSize(constants.VIDEO_WIDTH*videoScale, constants.VIDEO_HEIGHT*videoScale)
	ebiten.SetWindowTitle("G8Emu")

	cpuFrequency := 540
	game := emulator.NewGame(platform, chip8, cpuFrequency)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
