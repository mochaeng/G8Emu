package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	START_ADDRESS = 0x200

	FONTSET_SIZE          = 80
	FONTSET_START_ADDRESS = 0x50

	VIDEO_WIDTH  = 64
	VIDEO_HEIGHT = 32
)

var fontset = [FONTSET_SIZE]uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type Chip8 struct {
	pc         uint16
	sp         uint8
	index      uint16
	delayTimer uint8
	soundTimer uint8
	opcode     uint16

	registers [16]uint8
	memory    [4096]uint8
	stack     [16]uint16
	keypad    [16]uint8
	video     [VIDEO_WIDTH * VIDEO_HEIGHT]bool

	rng *rand.Rand
}

func NewChip8() *Chip8 {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	chip8 := Chip8{
		pc:  START_ADDRESS,
		rng: rng,
	}

	for i := range FONTSET_SIZE {
		chip8.memory[FONTSET_START_ADDRESS+i] = fontset[i]
	}

	return &chip8
}

func (c8 *Chip8) randByte() uint8 {
	return uint8(c8.rng.Intn(256))
}

// Clears the screen
//
// [usage]: CLS
func (c8 *Chip8) Op00E0() {
	for i := range len(c8.video) {
		c8.video[i] = 0
	}
}

// Return from a subroutine
//
// [usage]: RET
func (c8 *Chip8) Op00EE() {
	c8.sp--
	c8.pc = c8.stack[c8.sp]
}

// Jump to location NNN
//
// [usage]: JP addrr
func (c8 *Chip8) Op1NNN() {
	c8.pc = c8.opcode & 0x0FFF
}

// Call a subroutine at NNN
//
// [usage]: CALL addrr
func (c8 *Chip8) Op2NNN() {
	address := c8.opcode & 0x0FFF
	c8.stack[c8.sp] = c8.pc
	c8.sp++
	c8.pc = address
}

// Skips the next instruction if the value in register [Vx]
// is equal to [NN]
//
// [usage]: SE Vx, byte
func (c8 *Chip8) Op3XNN() {
	vx := (c8.opcode & 0x0F00) >> 8
	byte := c8.opcode & 0x00FF
	if c8.registers[vx] == uint8(byte) {
		c8.pc += 2
	}
}

// Skips the next instruction if the value in register [Vx]
// is different of [NN]
//
// [usage]: SNE Vx, byte
func (c8 *Chip8) Op4XNN() {
	vx := (c8.opcode & 0x0F00) >> 8
	byte := c8.opcode & 0x00FF
	if c8.registers[vx] != uint8(byte) {
		c8.pc += 2
	}
}

// Skips the next instruction if the register [Vx] is equal to [Vy]
//
// [usage]: SE Vx, Vy
func (c8 *Chip8) Op5XY0() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4
	if c8.registers[vx] == c8.registers[vy] {
		c8.pc += 2
	}
}

// Set register Vx to NN
//
// [usage]: LD Vx, byte
func (c8 *Chip8) Op6XNN() {
	vx := (c8.opcode & 0x0F00) >> 8
	byte := c8.opcode & 0x00FF
	c8.registers[vx] = uint8(byte)
}

// Set register Vx to [Vx + NN]
//
// [usage]: ADD Vx, byte
func (c8 *Chip8) Op7XNN() {
	vx := (c8.opcode & 0x0F00) >> 8
	byte := c8.opcode & 0x00FF
	c8.registers[vx] += uint8(byte)
}

// Copies the value from register Vy to Vx
//
// [usage]: LD Vx, Vy
func (c8 *Chip8) Op8XY0() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4
	c8.registers[vx] = c8.registers[vy]
}

// Performs a bitwise OR between register Vx and Vy
//
// usage: OR Vx, Vy
func (c8 *Chip8) Op8XY1() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4
	c8.registers[vx] |= c8.registers[vy]
}

// Performs a bitwise AND between register Vx and Vy
//
// [usage]: AND Vx, Vy
func (c8 *Chip8) Op8XY2() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4
	c8.registers[vx] &= c8.registers[vy]
}

// Performs a bitwise XOR between register Vx and Vy
//
// [usage]: XOR Vx, Vy
func (c8 *Chip8) Op8XY3() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4
	c8.registers[vx] ^= c8.registers[vy]
}

// Sums the two registers Vx and Vy. Also set VF = carry
//
// [usage]: ADD Vx, Vy
//
// [details]: If the sum is greater than 8 bits (>255), register VF
// is set to 1. Also, only the 8 bits of the result are kept
// and stored in Vx
func (c8 *Chip8) Op8XY4() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4
	sum := vx + vy
	if sum > 255 {
		c8.registers[0xF] = 1
	} else {
		c8.registers[0xF] = 0
	}
	c8.registers[vx] = uint8(sum & 0xFF)
}

// Subtracts the two registers Vx and Vy.
//
// [usage]: SUB Vx, Vy
//
// [details]: Also set VF = not borrow If Vx > Vy, then VF is
// set to 1, otherwise 0
func (c8 *Chip8) Op8XY5() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4
	if c8.registers[vx] > c8.registers[vy] {
		c8.registers[0xF] = 1
	} else {
		c8.registers[0xF] = 0
	}
	c8.registers[vx] -= c8.registers[vy]
}

// Shifts right a bit from register Vx.
//
// [usage]: = SHR Vx
//
// [details]: If the least-significant
// bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is
// divided by 2
func (c8 *Chip8) Op8XY6() {
	vx := (c8.opcode & 0x0F00) >> 8
	c8.registers[0xF] = uint8(c8.registers[vx] & 1)
	c8.registers[vx] >>= 1
}

// Set Vx = Vy - Vx, set VF = not borrow.
//
// [usage]: SUBN Vx, Vy
//
// [details]: If Vy > Vx, then VF is set
// to 1, otherwise 0.
func (c8 *Chip8) Op8XY7() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4
	if c8.registers[vx] > c8.registers[vy] {
		c8.registers[0xF] = 1
	} else {
		c8.registers[0xF] = 0
	}
	c8.registers[vx] = c8.registers[vy] - c8.registers[vx]
}

// Set Vx = Vx SHL 1
//
// [usage]: SHL Vx
//
// [details]: If the most significant bit of Vx is 1, then VF
// is set to 1, otherwise to 0. Then Vx is multiplied by 2
func (c8 *Chip8) Op8XYE() {
	vx := (c8.opcode & 0x0F00) >> 8
	c8.registers[0xF] = uint8((c8.opcode & 0x80) >> 7)
	c8.registers[vx] <<= 1
}

// Skip the next instruction if Vx != Vy
//
// [usage]: SNE Vx, Vy
func (c8 *Chip8) Op9XY0() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4
	if c8.registers[vx] != c8.registers[vy] {
		c8.pc += 2
	}
}

// Set I = NNN. The value of register I is set to NNN
//
// [usage]: LD I, addr
func (c8 *Chip8) OpANNN() {
	addr := c8.opcode & 0x0FFF
	c8.index = addr
}

// Jump to location NNN + V0. The PC is set to NNN + V0
//
// [usage]: JP V0, addr
func (c8 *Chip8) OpBNNN() {
	addr := c8.opcode & 0x0FFF
	c8.pc = addr + uint16(c8.registers[0])
}

// Set Vx = random byte AND NN. Generates a random number from 0 to 255
//
// [usage]: RND Vx, byte
func (c8 *Chip8) OpCXNN() {
	vx := (c8.opcode & 0x0F00) >> 8
	byte := c8.opcode & 0x00FF
	c8.registers[vx] = c8.randByte() & uint8(byte)
}

// Display n-byte sprite starting at memory location I at (Vx, Vy),
// set VF = collision
func (c8 *Chip8) OpDXYN() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4

	x := c8.registers[vx] % VIDEO_WIDTH
	y := c8.registers[vy] % VIDEO_HEIGHT
	c8.registers[0xF] = 0

	n := c8.opcode & 0x000F
	for range n {
		spriteRow := c8.memory[c8.index+n]

		var mask uint8
		for mask = 0x80; mask != 0; mask >>= 1 {
			currentSpritePixel := spriteRow & mask

			point := y*VIDEO_WIDTH + x
			currentScreenPixel := c8.video[point]
			if currentSpritePixel == 1 && currentScreenPixel {
				c8.video[point] = false
				c8.registers[0xF] = 1
			} else if currentSpritePixel == 1 && !currentScreenPixel {
				c8.video[point] = true
			}

			if x > VIDEO_WIDTH {
				break
			}

			x += 1
		}

		if y > VIDEO_HEIGHT {
			break
		}

		y += 1
	}
}

func (c8 *Chip8) LoadRom(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open ROM file: %w", err)
	}
	defer file.Close()

	buffer, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read ROM file: %w", err)
	}

	if len(buffer) > len(c8.memory)-START_ADDRESS {
		return fmt.Errorf("ROM is too large to fit in memory")
	}

	for i := range len(buffer) {
		c8.memory[START_ADDRESS+i] = buffer[i]
	}

	return nil
}

func main() {

}
