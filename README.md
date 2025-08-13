# G8Emu

A CHIP-8 emulator written in Go with [Ebitengine](https://github.com/hajimehoshi/ebiten) for graphics and platform support. G8Emu provides accurate emulation of the classic CHIP-8 system with support for desktop platforms (Windows, macOS, Linux) and web browsers through WebAssembly. This emulator faithfully recreates the original system with its 4KB of RAM, 16 general-purpose registers, 64x32 monochrome display, and 16-key hexadecimal keypad.

Load your favorite CHIP-8 ROMs and experience retro gaming with cross-platform compatibility.

## Controls

The CHIP-8 keypad is mapped to your keyboard as follows:

```sh
CHIP-8 Keypad    →    QWERTY Keyboard
┌─────────────────────────────────────┐
│  1  2  3  C     →     1  2  3  4   │
│  4  5  6  D     →     Q  W  E  R   │
│  7  8  9  E     →     A  S  D  F   │
│  A  0  B  F     →     Z  X  C  V   │
└─────────────────────────────────────┘
```

Additional Controls:

- P: Pause/Resume emulation
- R: Reset emulator

## Where to Find ROMs

- [dmatlack/chip8](https://github.com/dmatlack/chip8/tree/master/roms/games)

Dowloand the files with `.ch8` extension

## Features

- [x] Complete CHIP-8 instruction set
- [x] Cross-platform desktop
- [x] Web version through WebAssembly
- [ ] Sound output
- [ ] Dynamic CPU frequency
- [ ] Save and load emulator states
- [ ] Additional SUPER-CHIP instruction set

## Getting Started

#### Desktop version

##### Pre-built binaries

Dowloand from [releases page](https://github.com/mochaeng/G8Emu/releases):

- Linux: `g8emu-linux-amd64.zip`
- Windows: `g8emu-windows-amd64.zip`
- macOS: `g8emu-macos-arm64.zip`

##### Usage

```sh
./g8emu <scale> <rom-file>
```

##### Example

For playing a tetris ROM with 640x320 resolution on linux:

```sh
./g8emu 10 tetris.ch8
```

#### Web Version

Visit: []

## Building from Source

#### Desktop

```sh
# Go 1.23+
go build -o g8emu cmd/desktop/main.go
```

#### Web

Compile the project with webassembly support:

```sh
GOOS=js GOARCH=wasm go build -o web-react/public/g8emu.wasm ./cmd/wasm/main.go
@cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js web-react/public/
```

Run the react application:

```sh
cd web-react
pnpm install
pnpm build
```

## Close look

##### Desktop

![a chip8 ROM](docs/imgs/Screenshot%20from%202025-08-13%2010-17-45.png)
![a tetris ROM for chip8](docs/imgs/Screenshot%20from%202025-08-13%2010-16-50.png)

##### Web

![a tetris ROM for chip8 running on browser with webassembly support](docs/imgs/Screenshot%20from%202025-08-13%2010-20-05.png)

## Resources and References

- [Chip-8 Technical Reference v1.0](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#00Cn)
- [Building a CHIP-8 Emulator [C++]](https://austinmorlan.com/posts/chip8_emulator/)
- [How to write an emulator (CHIP-8 interpreter)](https://multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/)
- [Guide to making a CHIP-8 emulator](https://tobiasvl.github.io/blog/write-a-chip-8-emulator/)
