package core

import "github.com/mochaeng/G8Emu/internal/constants"

// NULL operation for invalid opcodes
func (c8 *Chip8) OpNULL() {}

// Clears the screen
//
// [instruction]: CLS
func (c8 *Chip8) Op00E0() {
	for i := range len(c8.Video) {
		c8.Video[i] = false
	}
}

// Return from a subroutine
//
// [instruction]: RET
func (c8 *Chip8) Op00EE() {
	c8.sp--
	c8.pc = c8.stack[c8.sp]
}

// Jump to location NNN
//
// [instruction]: JP addrr
func (c8 *Chip8) Op1NNN() {
	c8.pc = c8.opcode & 0x0FFF
}

// Call a subroutine at NNN
//
// [instruction]: CALL addrr
func (c8 *Chip8) Op2NNN() {
	address := c8.opcode & 0x0FFF

	c8.stack[c8.sp] = c8.pc
	c8.sp++
	c8.pc = address
}

// Skips the next instruction if the value in register [Vx]
// is equal to [NN]
//
// [instruction]: SE Vx, byte
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
// [instruction]: SNE Vx, byte
func (c8 *Chip8) Op4XNN() {
	vx := (c8.opcode & 0x0F00) >> 8
	byte := c8.opcode & 0x00FF

	if c8.registers[vx] != uint8(byte) {
		c8.pc += 2
	}
}

// Skips the next instruction if the register [Vx] is equal to [Vy]
//
// [instruction]: SE Vx, Vy
func (c8 *Chip8) Op5XY0() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4

	if c8.registers[vx] == c8.registers[vy] {
		c8.pc += 2
	}
}

// Set register Vx to NN
//
// [instruction]: LD Vx, byte
func (c8 *Chip8) Op6XNN() {
	vx := (c8.opcode & 0x0F00) >> 8
	byte := c8.opcode & 0x00FF

	c8.registers[vx] = uint8(byte)
}

// Set register Vx to [Vx + NN]
//
// [instruction]: ADD Vx, byte
func (c8 *Chip8) Op7XNN() {
	vx := (c8.opcode & 0x0F00) >> 8
	byte := c8.opcode & 0x00FF

	c8.registers[vx] += uint8(byte)
}

// Copies the value from register Vy to Vx
//
// [instruction]: LD Vx, Vy
func (c8 *Chip8) Op8XY0() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4

	c8.registers[vx] = c8.registers[vy]
}

// Performs a bitwise OR between register Vx and Vy
//
// instruction: OR Vx, Vy
func (c8 *Chip8) Op8XY1() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4

	c8.registers[vx] |= c8.registers[vy]
}

// Performs a bitwise AND between register Vx and Vy
//
// [instruction]: AND Vx, Vy
func (c8 *Chip8) Op8XY2() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4

	c8.registers[vx] &= c8.registers[vy]
}

// Performs a bitwise XOR between register Vx and Vy
//
// [instruction]: XOR Vx, Vy
func (c8 *Chip8) Op8XY3() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4

	c8.registers[vx] ^= c8.registers[vy]
}

// Sums the two registers Vx and Vy. Also set VF = carry
//
// [instruction]: ADD Vx, Vy
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
// [instruction]: SUB Vx, Vy
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
// [instruction]: = SHR Vx
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
// [instruction]: SUBN Vx, Vy
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
// [instruction]: SHL Vx
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
// [instruction]: SNE Vx, Vy
func (c8 *Chip8) Op9XY0() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4

	if c8.registers[vx] != c8.registers[vy] {
		c8.pc += 2
	}
}

// Set I = NNN. The value of register I is set to NNN
//
// [instruction]: LD I, addr
func (c8 *Chip8) OpANNN() {
	addr := c8.opcode & 0x0FFF
	c8.index = addr
}

// Jump to location NNN + V0. The PC is set to NNN + V0
//
// [instruction]: JP V0, addr
func (c8 *Chip8) OpBNNN() {
	addr := c8.opcode & 0x0FFF
	c8.pc = uint16(c8.registers[0]) + addr
}

// Set Vx = random byte AND NN. Generates a random number from 0 to 255
//
// [instruction]: RND Vx, byte
func (c8 *Chip8) OpCXNN() {
	vx := (c8.opcode & 0x0F00) >> 8
	byte := c8.opcode & 0x00FF

	c8.registers[vx] = c8.randByte() & uint8(byte)
}

// Display n-byte sprite starting at memory location I at (Vx, Vy)
//
// [instruction]: DRW Vx, Vy, nibble
func (c8 *Chip8) OpDXYN() {
	vx := (c8.opcode & 0x0F00) >> 8
	vy := (c8.opcode & 0x00F0) >> 4

	x := c8.registers[vx]
	y := c8.registers[vy]
	height := c8.opcode & 0x000F
	c8.registers[0xF] = 0

	for row := range height {
		spriteRowData := c8.memory[c8.index+row]

		screenY := y + uint8(row)
		if screenY >= constants.VIDEO_HEIGHT {
			break
		}

		for col := range uint8(8) {
			screenX := x + col
			if screenX >= constants.VIDEO_WIDTH {
				break
			}

			isSpritePixelOn := (spriteRowData & (0x80 >> col)) != 0
			pixelPosition := uint16(screenY)*constants.VIDEO_WIDTH + uint16(screenX)

			if pixelPosition >= constants.VIDEO_WIDTH*constants.VIDEO_HEIGHT {
				continue
			}

			isScreenPixelOn := c8.Video[pixelPosition]

			if isSpritePixelOn && isScreenPixelOn {
				c8.registers[0xF] = 1
			}

			if isSpritePixelOn {
				c8.Video[pixelPosition] = !isScreenPixelOn
			}
		}
	}
}

// Skip next instruction if key with the value of Vx is pressed
//
// [instruction]: SKP Vx
func (c8 *Chip8) OpEX9E() {
	vx := (c8.opcode & 0x0F00) >> 8
	key := c8.registers[vx]

	if c8.Keypad[key] {
		c8.pc += 2
	}
}

// Skip next instruction if key with the value of Vx is not pressed
//
// [instruction]: SKNP Vx
func (c8 *Chip8) OpEXA1() {
	vx := (c8.opcode & 0x0F00) >> 8
	key := c8.registers[vx]

	if !c8.Keypad[key] {
		c8.pc += 2
	}
}

// Set Vx = delay timer value
//
// [instruction]: LD Vx, DT
func (c8 *Chip8) OpFX07() {
	vx := (c8.opcode & 0x0F00) >> 8
	c8.registers[vx] = c8.DelayTimer
}

// Wait for a key press, store the value of the key in Vx
//
// [instruction]: LD Vx, K
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
	} else if c8.Keypad[4] {
		c8.registers[vx] = 4
	} else if c8.Keypad[5] {
		c8.registers[vx] = 5
	} else if c8.Keypad[6] {
		c8.registers[vx] = 6
	} else if c8.Keypad[7] {
		c8.registers[vx] = 7
	} else if c8.Keypad[8] {
		c8.registers[vx] = 8
	} else if c8.Keypad[9] {
		c8.registers[vx] = 9
	} else if c8.Keypad[10] {
		c8.registers[vx] = 10
	} else if c8.Keypad[11] {
		c8.registers[vx] = 11
	} else if c8.Keypad[12] {
		c8.registers[vx] = 12
	} else if c8.Keypad[13] {
		c8.registers[vx] = 13
	} else if c8.Keypad[14] {
		c8.registers[vx] = 14
	} else if c8.Keypad[15] {
		c8.registers[vx] = 15
	} else {
		c8.pc -= 2
	}
}

// Set delay timer = Vx.
//
// [instruction]: LD DT, Vx
func (c8 *Chip8) OpFX15() {
	vx := (c8.opcode & 0x0F00) >> 8
	c8.DelayTimer = c8.registers[vx]
}

// Set sound timer = Vx.
//
// [instruction]: LD ST, Vx
func (c8 *Chip8) OpFX18() {
	vx := (c8.opcode & 0x0F00) >> 8
	c8.SoundTimer = c8.registers[vx]
}

// Set I = I + Vx.
//
// [instruction]:  ADD I, Vx
func (c8 *Chip8) OpFX1E() {
	vx := (c8.opcode & 0x0F00) >> 8
	c8.index += uint16(c8.registers[vx])
}

// Set I = location of sprite for digit Vx
//
// [instruction]: LD F, Vx
func (c8 *Chip8) OpFX29() {
	vx := (c8.opcode & 0x0F00) >> 8
	digit := c8.registers[vx]

	c8.index = FONTSET_START_ADDRESS + CHAR_FONT_SIZE*uint16(digit)
}

// Store BCD representation of Vx in memory locations I, I+1, and I+2.
//
// [instruction]: LD B, Vx
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
// [instruction]: LD [I], Vx
func (c8 *Chip8) OpFX55() {
	vx := (c8.opcode & 0x0F00) >> 8

	for i := range vx + 1 {
		c8.memory[c8.index+i] = c8.registers[i]
	}
}

// Read registers V0 through Vx from memory starting at location I.
//
// [instruction]: LD Vx, [I]
func (c8 *Chip8) OpFX65() {
	vx := (c8.opcode & 0x0F00) >> 8

	for i := range vx + 1 {
		c8.registers[i] = c8.memory[c8.index+i]
	}
}
