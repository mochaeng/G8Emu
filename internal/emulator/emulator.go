package emulator

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/mochaeng/G8Emu/internal/constants"
)

const (
	START_ADDRESS = 0x200

	FONTSET_SIZE          = 80
	FONTSET_START_ADDRESS = 0x50

	CHAR_FONT_SIZE = 5
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
	Keypad    [16]bool
	Video     [constants.VIDEO_WIDTH * constants.VIDEO_HEIGHT]bool

	rng *rand.Rand

	table  [0xF + 1]func()
	table0 [0xE + 1]func()
	table8 [0xE + 1]func()
	tableE [0xE + 1]func()
	tableF [0x65 + 1]func()
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

	chip8.table[0x0] = chip8.Table0
	chip8.table[0x1] = chip8.Op1NNN
	chip8.table[0x2] = chip8.Op2NNN
	chip8.table[0x3] = chip8.Op3XNN
	chip8.table[0x4] = chip8.Op4XNN
	chip8.table[0x5] = chip8.Op5XY0
	chip8.table[0x6] = chip8.Op6XNN
	chip8.table[0x7] = chip8.Op7XNN
	chip8.table[0x8] = chip8.Table8
	chip8.table[0x9] = chip8.Op9XY0
	chip8.table[0xA] = chip8.OpANNN
	chip8.table[0xB] = chip8.OpBNNN
	chip8.table[0xC] = chip8.OpCXNN
	chip8.table[0xD] = chip8.OpDXYN
	chip8.table[0xE] = chip8.TableE
	chip8.table[0xF] = chip8.TableF

	for i := 0; i <= 0xE; i++ {
		chip8.table0[i] = chip8.OpNULL
		chip8.table8[i] = chip8.OpNULL
		chip8.tableE[i] = chip8.OpNULL
	}

	chip8.table0[0x0] = chip8.Op00E0
	chip8.table0[0xE] = chip8.Op00EE

	chip8.table8[0x0] = chip8.Op8XY0
	chip8.table8[0x1] = chip8.Op8XY1
	chip8.table8[0x2] = chip8.Op8XY2
	chip8.table8[0x3] = chip8.Op8XY3
	chip8.table8[0x4] = chip8.Op8XY4
	chip8.table8[0x5] = chip8.Op8XY5
	chip8.table8[0x6] = chip8.Op8XY6
	chip8.table8[0x7] = chip8.Op8XY7
	chip8.table8[0xE] = chip8.Op8XYE

	chip8.tableE[0x1] = chip8.OpEXA1
	chip8.tableE[0xE] = chip8.OpEX9E

	for i := 0; i <= 0x65; i++ {
		chip8.tableF[i] = chip8.OpNULL
	}

	chip8.tableF[0x07] = chip8.OpFX07
	chip8.tableF[0x0A] = chip8.OpFX0A
	chip8.tableF[0x15] = chip8.OpFX15
	chip8.tableF[0x18] = chip8.OpFX18
	chip8.tableF[0x1E] = chip8.OpFX1E
	chip8.tableF[0x29] = chip8.OpFX29
	chip8.tableF[0x33] = chip8.OpFX33
	chip8.tableF[0x55] = chip8.OpFX55
	chip8.tableF[0x65] = chip8.OpFX65

	return &chip8
}

func (c8 *Chip8) Table0() {
	c8.table0[c8.opcode&0x000F]()
}

func (c8 *Chip8) Table8() {
	c8.table8[c8.opcode&0x000F]()
}

func (c8 *Chip8) TableE() {
	c8.tableE[c8.opcode&0x000F]()
}

func (c8 *Chip8) TableF() {
	c8.tableF[c8.opcode&0x00FF]()
}

func (c8 *Chip8) randByte() uint8 {
	return uint8(c8.rng.Intn(256))
}

func (c8 *Chip8) Cycle() {
	c8.opcode = (uint16(c8.memory[c8.pc])<<8 | uint16(c8.memory[c8.pc+1]))

	c8.pc += 2

	firstNibble := (c8.opcode & 0xF000) >> 12
	c8.table[firstNibble]()

	if c8.delayTimer > 0 {
		c8.delayTimer--
	}

	if c8.soundTimer > 0 {
		c8.soundTimer--
	}
}

// NULL operation for invalid opcodes
func (c8 *Chip8) OpNULL() {}

// Clears the screen
//
// [usage]: CLS
func (c8 *Chip8) Op00E0() {
	for i := range len(c8.Video) {
		c8.Video[i] = false
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

	sum := uint16(c8.registers[vx]) + uint16(c8.registers[vy])

	if sum > 0xFF {
		c8.registers[0xF] = 1
	} else {
		c8.registers[0xF] = 0
	}

	// c8.registers[vx] = uint8(sum & 0xFF)
	c8.registers[vx] = uint8(sum)
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

	if c8.registers[vy] > c8.registers[vx] {
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
	// c8.registers[0xF] = uint8((c8.opcode & 0x80) >> 7)
	c8.registers[0xF] = (c8.registers[vx] & 0x80) >> 7
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
	c8.pc = uint16(c8.registers[0]) + addr
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
//
// [usage]:  DRW Vx, Vy, nibble
// func (c8 *Chip8) OpDXYN() {
// 	vx := (c8.opcode & 0x0F00) >> 8
// 	vy := (c8.opcode & 0x00F0) >> 4

// 	x := c8.registers[vx] % constants.VIDEO_WIDTH
// 	y := c8.registers[vy] % constants.VIDEO_HEIGHT
// 	c8.registers[0xF] = 0

// 	n := c8.opcode & 0x000F
// 	for row := range n {
// 		spriteRow := c8.memory[c8.index+row]

// 		for col := range 8 {
// 			spritePixel := spriteRow & (0x80 >> col)

// 			// point := y*constants.VIDEO_WIDTH + x
// 			point := (y+uint8(row))*constants.VIDEO_WIDTH + (x + uint8(col))
// 			screenPixel := c8.Video[point]

// 			if spritePixel != 0 {
// 				if screenPixel {
// 					c8.registers[0xF] = 1
// 					c8.Video[point] = false
// 				}
// 				if !screenPixel {
// 					c8.Video[point] = true
// 				}
// 				// c8.Video[point] = !screenPixel
// 			}
// 		}
// 	}
// }

// Fixed OpDXYN - Display n-byte sprite starting at memory location I at (Vx, Vy)
func (c8 *Chip8) OpDXYN() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4

	x := c8.registers[vx]
	y := c8.registers[vy]
	height := c8.opcode & 0x000F
	c8.registers[0xF] = 0

	for row := uint16(0); row < height; row++ {
		spriteRow := c8.memory[c8.index+row]

		// Check if we're going off screen vertically
		yCoord := y + uint8(row)
		if yCoord >= constants.VIDEO_HEIGHT {
			break // Clip vertically - don't draw rows that go off screen
		}

		for col := uint8(0); col < 8; col++ {
			// Check if we're going off screen horizontally
			xCoord := x + col
			if xCoord >= constants.VIDEO_WIDTH {
				break // Clip horizontally - don't draw pixels that go off screen
			}

			// Get the current bit from the sprite row
			spritePixel := (spriteRow & (0x80 >> col)) != 0

			// Calculate pixel index
			pixelIndex := uint16(yCoord)*constants.VIDEO_WIDTH + uint16(xCoord)

			// Safety check (shouldn't be needed with proper clipping, but good practice)
			if pixelIndex >= constants.VIDEO_WIDTH*constants.VIDEO_HEIGHT {
				continue
			}

			// Get current screen state
			currentPixel := c8.Video[pixelIndex]

			// Check for collision (both pixels are on)
			if spritePixel && currentPixel {
				c8.registers[0xF] = 1
			}

			// XOR the sprite pixel with screen pixel
			if spritePixel {
				c8.Video[pixelIndex] = !currentPixel
			}
		}
	}
}

// Skip next instruction if key with the value of Vx is pressed
//
// [usage]: SKP Vx
func (c8 *Chip8) OpEX9E() {
	vx := (c8.opcode & 0x0F00) >> 8
	key := c8.registers[vx]

	if c8.Keypad[key] {
		c8.pc += 2
	}
}

// Skip next instruction if key with the value of Vx is not pressed
//
// [usage]: SKNP Vx
func (c8 *Chip8) OpEXA1() {
	vx := (c8.opcode & 0x0F00) >> 8
	key := c8.registers[vx]

	if !c8.Keypad[key] {
		c8.pc += 2
	}
}

// Set Vx = delay timer value
//
// [usage]: LD Vx, DT
func (c8 *Chip8) OpFX07() {
	vx := (c8.opcode & 0x0F00) >> 8
	c8.registers[vx] = c8.delayTimer
}

// Wait for a key press, store the value of the key in Vx
//
// [usage]: LD Vx, K
func (c8 *Chip8) OpFX0A() {
	vx := (c8.opcode & 0x0F00) >> 8

	if c8.Keypad[0] {
		c8.registers[vx] = 0
	} else if c8.Keypad[1] {
		c8.registers[vx] = 1
	} else if c8.Keypad[2] {
		c8.registers[vx] = 2
	} else if c8.Keypad[3] {
		c8.registers[vx] = 3
	} else {
		c8.pc -= 2
	}
}

// Set delay timer = Vx.
//
// [usage]: LD DT, Vx
func (c8 *Chip8) OpFX15() {
	vx := (c8.opcode & 0x0F00) >> 8
	c8.delayTimer = c8.registers[vx]
}

// Set sound timer = Vx.
//
// [usage]: LD ST, Vx
func (c8 *Chip8) OpFX18() {
	vx := (c8.opcode & 0x0F00) >> 8
	c8.soundTimer = c8.registers[vx]
}

// Set I = I + Vx.
//
// [usage]:  ADD I, Vx
func (c8 *Chip8) OpFX1E() {
	vx := (c8.opcode & 0x0F00) >> 8
	c8.index += uint16(c8.registers[vx])
}

// Set I = location of sprite for digit Vx
//
// [usage]: LD F, Vx
func (c8 *Chip8) OpFX29() {
	vx := (c8.opcode & 0x0F00) >> 8
	digit := c8.registers[vx]

	c8.index = FONTSET_START_ADDRESS + CHAR_FONT_SIZE*uint16(digit)
}

// Store BCD representation of Vx in memory locations I, I+1, and I+2.
//
// [usage]: LD B, Vx
func (c8 *Chip8) OpFX33() {
	vx := (c8.opcode & 0x0F00) >> 8
	value := c8.registers[vx]

	c8.memory[c8.index+2] = value % 10
	value /= 10

	c8.memory[c8.index+1] = value % 10
	value /= 10

	c8.memory[c8.index] = value % 10
}

// Store registers V0 through Vx in memory starting at location I.
//
// [usage]: LD [I], Vx
func (c8 *Chip8) OpFX55() {
	vx := (c8.opcode & 0x0F00) >> 8

	for i := range vx + 1 {
		c8.memory[c8.index+i] = c8.registers[i]
	}
}

// Read registers V0 through Vx from memory starting at location I.
//
// [usage]: LD Vx, [I]
func (c8 *Chip8) OpFX65() {
	vx := (c8.opcode & 0x0F00) >> 8

	for i := range vx + 1 {
		c8.registers[i] = c8.memory[c8.index+i]
	}
}

func (c8 *Chip8) LoadRom(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read ROM file: %v", err)
	}

	if len(data) > len(c8.memory)-START_ADDRESS {
		return fmt.Errorf("ROM too large to fit in memory: %d bytes (max %d)", len(data), len(c8.memory)-START_ADDRESS)
	}

	copy(c8.memory[START_ADDRESS:], data)

	return nil
}
