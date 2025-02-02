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

	randGen  *rand.Rand
	randByte uint8
}

func NewChip8() *Chip8 {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	chip8 := Chip8{
		pc:       START_ADDRESS,
		randGen:  rng,
		randByte: uint8(rng.Intn(256)),
	}

	for i := 0; i < FONTSET_SIZE; i++ {
		chip8.memory[FONTSET_START_ADDRESS+i] = fontset[i]
	}

	return &chip8
}

// Clears the screen (OP-00E0)
func (chip8 *Chip8) ClearDisplay() {
	for i := 0; i < len(chip8.video); i++ {
		chip8.video[i] = 0
	}
}

// Return from a subroutine (OP-00EE)
func (chip8 *Chip8) Return() {
	chip8.sp--
	chip8.pc = chip8.stack[chip8.sp]
}

// Jump to location NNN (OP-1NNN)
func (chip8 *Chip8) Goto() {
	address := chip8.opcode & 0x0FFF
	chip8.pc = address
}

func (chip8 *Chip8) LoadRom(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open ROM file: %w", err)
	}
	defer file.Close()

	buffer, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read ROM file: %w", err)
	}

	if len(buffer) > len(chip8.memory)-START_ADDRESS {
		return fmt.Errorf("ROM is too large to fit in memory")
	}

	for i := 0; i < len(buffer); i++ {
		chip8.memory[START_ADDRESS+i] = buffer[i]
	}

	return nil
}

func main() {

}
