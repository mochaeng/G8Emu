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
	video     [64 * 32]uint32

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
//	(OP-00E0): CLS
func (c8 *Chip8) CLS() {
	for i := range len(c8.video) {
		c8.video[i] = 0
	}
}

// Return from a subroutine
//
// (OP-00EE): RET
func (c8 *Chip8) RET() {
	c8.sp--
	c8.pc = c8.stack[c8.sp]
}

// Jump to location NNN
//
// (OP-1NNN): JP addr
func (c8 *Chip8) JP() {
	c8.pc = c8.opcode & 0x0FFF
}

// [OP-2NNN] = CALL addr
//
// Call a subroutine at NNN
func (c8 *Chip8) CALL() {
	address := c8.opcode & 0x0FFF
	c8.stack[c8.sp] = c8.pc
	c8.sp++
	c8.pc = address
}

// Skips the next instruction if the value in register [Vx]
// is equal to [NN]
//
// (OP-3XNN): SE Vx, byte
func (c8 *Chip8) SE() {
	vx := (c8.opcode & 0x0F00) >> 8
	byte := c8.opcode & 0x00FF
	if c8.registers[vx] == uint8(byte) {
		c8.pc += 2
	}
}

// Skips the next instruction if the value in register [Vx]
// is different of [NN]
//
// (OP-4XNN): SNE Vx, byte
func (c8 *Chip8) SNE_BYTE() {
	vx := (c8.opcode & 0x0F00) >> 8
	byte := c8.opcode & 0x00FF
	if c8.registers[vx] != uint8(byte) {
		c8.pc += 2
	}
}

// Skips the next instruction if the register [Vx] is equal to [Vy]
//
// (5XY0): SE Vx, Vy
func (c8 *Chip8) SE_REGISTERS() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4
	if c8.registers[vx] == c8.registers[vy] {
		c8.pc += 2
	}
}

// Set register Vx to NN
//
// (6XNN): LD Vx, byte
func (c8 *Chip8) LD_BYTE() {
	vx := (c8.opcode & 0x0F00) >> 8
	byte := c8.opcode & 0x00FF
	c8.registers[vx] = uint8(byte)
}

// Set register Vx to [Vx + NN]
//
// (7XNN): ADD Vx, byte
func (c8 *Chip8) AddToRegister() {
	vx := (c8.opcode & 0x0F00) >> 8
	byte := c8.opcode & 0x00FF
	c8.registers[vx] += uint8(byte)
}

// Copies the value from register Vy to Vx
//
// (8XY0): LD Vx, Vy
func (c8 *Chip8) CopyRegister() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4
	c8.registers[vx] = uint8(vy)
}

// Performs a bitwise OR between register Vx and Vy
//
// (8XY1): OR Vx, Vy
func (c8 *Chip8) OrRegisters() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4
	c8.registers[vx] |= uint8(vy)
}

// Performs a bitwise AND between register Vx and Vy
//
// (8XY2): AND Vx, Vy
func (c8 *Chip8) AND() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4
	c8.registers[vx] &= uint8(vy)
}

// Performs a bitwise XOR between register Vx and Vy
//
// (8XY3): XOR Vx, Vy
func (c8 *Chip8) XorRegisters() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4
	c8.registers[vx] ^= uint8(vy)
}

// Sums the two registers Vx and Vy. Also set VF = carry
//
// (8XY4): ADD Vx, Vy
//
// If the sum is greater than 8 bits (>255), register VF
// is set to 1. Also, only the 8 bits of the result are kept
// and stored in Vx
func (c8 *Chip8) AddRegisters() {
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

// [OP-8XY5]: SUB Vx, Vy
//
// Subtracts the two registers Vx and Vy. Also set VF = not borrow
// If Vx > Vy, then VF is set to 1, otherwise 0
func (c8 *Chip8) SubRegisters() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4
	if c8.registers[vx] > c8.registers[vy] {
		c8.registers[0xF] = 1
	} else {
		c8.registers[0xF] = 0
	}
	c8.registers[vx] -= c8.registers[vy]
}

// [OP-8XY6] = SHR Vx
//
// Shifts right a bit from register Vx. If the least-significant
// bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is
// divided by 2
func (c8 *Chip8) ShiftRightRegister() {
	vx := (c8.opcode & 0x0F00) >> 8
	c8.registers[0xF] = uint8(c8.registers[vx] & 1)
	c8.registers[vx] >>= 1
}

// [OP-8XY7] = SUBN Vx, Vy
//
// Set Vx = Vy - Vx, set VF = not borrow. If Vy > Vx, then VF is set
// to 1, otherwise 0.
func (c8 *Chip8) SUBN() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4
	if c8.registers[vx] > c8.registers[vy] {
		c8.registers[0xF] = 1
	} else {
		c8.registers[0xF] = 0
	}
	c8.registers[vx] = c8.registers[vy] - c8.registers[vx]
}

// [8XYE] = SHL Vx
//
// Set Vx = Vx SHL 1. If the most significant bit of Vx is 1, then VF
// is set to 1, otherwise to 0. Then Vx is multiplied by 2
func (c8 *Chip8) SHL() {
	vx := (c8.opcode & 0x0F00) >> 8
	c8.registers[0xF] = uint8((c8.opcode & 0x80) >> 7)
	c8.registers[vx] <<= 1
}

// [9XY0] = SNE Vx, Vy
//
// Skip the next instruction if Vx != Vy
func (c8 *Chip8) SNE_REGISTERS() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4
	if c8.registers[vx] != c8.registers[vy] {
		c8.pc += 2
	}
}

// [ANNN] = LD I, addr
//
// Set I = NNN. The value of register I is set to NNN
func (c8 *Chip8) LD_I() {
	addr := c8.opcode & 0x0FFF
	c8.index = addr
}

// [BNNN] = JP V0, addr
//
// Jump to location NNN + V0. The PC is set to NNN + V0
func (c8 *Chip8) JP_V0() {
	addr := c8.opcode & 0x0FFF
	c8.pc = addr + uint16(c8.registers[0])
}

// [CXNN] = RND Vx, byte
//
// Set Vx = random byte AND NN. Generates a random number from 0 to 255
func (c8 *Chip8) RND() {
	vx := (c8.opcode & 0x0F00) >> 8
	byte := c8.opcode & 0x00FF
	c8.registers[vx] = c8.randByte() & uint8(byte)
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
