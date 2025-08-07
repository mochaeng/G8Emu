# G8Emu

An emulator reads the original machine code instructions that were assembled for the target machine, interprets them, and then replicates the functionality of the target machine on the host machine. The ROM files contain the instructions, the mulator reads those instructions, and then does work to mimic the original machine.

# Description

## Registers

CHIP-8 has sixteen 8-bit registers, labeled **V0** to **VF**. Each register is able to hold values from _0x00_ to _0xFF_. **VF** is used as a flag to hold information about the result of operations.

## Memory

There are 4096 bytes of memory, meaning the address space is from _0x00_ to _0xFFF_.

- `0x000 - 0x1FF`: Originally reserved for the CHIP-8 interpreter.
- `0x050 - 0x0A0`: Storage space for the 16 built-in characters (0 throught F).
- `0x200 - 0xFFF`: Instructions from the ROM will be stored from here.

## Index Register

A 16-bit special register to store memory addresses for use in operations.

## Program Counter

A 16-bit special register that holds the address of the next instruction to execute.

## Level Stack

A way for a CPU to keep track of the order of executation when it calls into functions. A instruction like `CALL` will cause the CPU to begin executing instructions in a different region of the program. When the program reaches another instruction `RET`, it must be able to go back to where it was when it hit the `CALL` instruction. The stack holds the PC value when the `CALL` instruction was executed, and the `RETURN` statement pull that address from the stack and puts it back into the PC so the CPU will execute it on the next cycle.

## Stack Pointer

Tells us where in the 16-levels of stack our most recent value was placed (top).

- With each `CALL`, the PC is placed where the SP was pointing, and the SP is incremented.
- With each `RET`, the stack pointer is decremented by one and the address that it's pointing to is put into the PC for execution.

## Monochrome Display Memory

Additional memory buffer used for storing the graphics to display (64x32). Each pixel is either on or off.

The draw instruction iterates over each pixel in a sprite and XORs the sprite pixel with the display pixel.

## Running

```sh
go run ./cmd 10 1 ROMs/test_opcode.ch8
```

## Platform

Handles hardware abstraction and I/O

- display rendering (`*ebiten.Image`) as the buffer
- scales the native 64x32 to a larger size `videoScale`
- converts the chip8 boolean video to pixels (white=true,black=false)
- implements `Draw()`

- maps physical keyboard keys to chip8's hex keypad using `keymap`
- `ProcessInput()` check keyboard state and updates chip8's keypad array

- handles scre layout/scaling through `Layout()`
- connection between chip8's internal display and ebitens window system

## Engine

Implements ebiten's game loop and timing control

- uses `cycleDelay` (milliseconds) to control emulation speed
- tracks `lastCycleTime` to maintain consistent timing
- only advances chip8 state when delay time is reached

- `Update()` called every frame
- `Draw()` updates visual output
- `Layout()` delegates to platform to screen sizing

- bridges between ebiten callbacks and chip8 operations
- manages communication between platform and chip8 instances

## Zed

GOOS=js GOARCH=wasm zed .
