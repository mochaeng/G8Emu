package main

import (
	"syscall/js"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mochaeng/G8Emu/internal/core"
	"github.com/mochaeng/G8Emu/internal/emulator"
)

func main() {
	const scale = 10
	const frequency = 540

	chip8 := core.NewChip8()
	platform := emulator.NewPlatform(scale)
	engine := emulator.NewGame(platform, chip8, frequency)

	loadRom := func(this js.Value, args []js.Value) any {
		println("what about this one?")
		if len(args) == 0 || args[0].IsNull() {
			return js.ValueOf("No ROM data provided")
		}

		dataLength := args[0].Get("length").Int()
		romData := make([]byte, dataLength)
		js.CopyBytesToGo(romData, args[0])

		if err := chip8.LoadRomBytes(romData); err != nil {
			js.Global().Call("alert", "ROM load error: "+err.Error())
			return js.ValueOf(err.Error())
		}

		return nil
	}

	resetEmulator := func(this js.Value, args []js.Value) any {
		println("is this being called?")
		engine.Reset()
		return nil
	}

	togglePause := func(this js.Value, args []js.Value) any {
		if engine.IsPaused() {
			engine.Resume()
		} else {
			engine.Pause()
		}
		return nil
	}

	setCpuFrequency := func(this js.Value, args []js.Value) any {
		// freq := args[0].Int()
		// game.setCpuFrequency(freq)
		return nil
	}

	js.Global().Set("loadRom", js.FuncOf(loadRom))
	js.Global().Set("resetEmulator", js.FuncOf(resetEmulator))
	js.Global().Set("togglePause", js.FuncOf(togglePause))
	js.Global().Set("setCpuFrequency", js.FuncOf(setCpuFrequency))

	go func() {
		if err := ebiten.RunGame(engine); err != nil {
			println("Game error: ", err)
		}
	}()

	select {}
}
