package core

import (
	"fmt"
	"math/rand"
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
	DelayTimer uint8
	SoundTimer uint8
	opcode     uint16

	registers [16]uint8
	memory    [4096]uint8
	stack     [16]uint16
	Keypad    [16]bool
	Video     [constants.VIDEO_WIDTH * constants.VIDEO_HEIGHT]bool

	rng *rand.Rand

	paused bool

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

	chip8.loadFontset()
	chip8.initTables()

	for i := range len(chip8.memory) {
		chip8.memory[i] = 0x00
	}

	return &chip8
}

func (c8 *Chip8) loadFontset() {
	for i := range FONTSET_SIZE {
		c8.memory[FONTSET_START_ADDRESS+i] = fontset[i]
	}
}

func (c8 *Chip8) randByte() uint8 {
	return uint8(c8.rng.Intn(256))
}

func (c8 *Chip8) fetch() {
	addr1 := c8.pc & 0xFFF
	addr2 := (c8.pc + 1) & 0xFFF
	c8.opcode = (uint16(c8.memory[addr1])<<8 | uint16(c8.memory[addr2]))
	c8.pc = (c8.pc + 2) & 0xFFF
}

func (c8 *Chip8) decodeAndExecute() {
	firstNibble := (c8.opcode & 0xF000) >> 12
	c8.table[firstNibble]()
}

func (c8 *Chip8) Cycle() {
	if c8.paused {
		return
	}

	c8.fetch()

	if c8.opcode == 0x0000 {
		return
	}

	c8.decodeAndExecute()
}

func (c8 *Chip8) memRead(addr uint16) uint8 {
	return c8.memory[addr%4096]
}

func (c8 *Chip8) memWrite(addr uint16, value uint8) {
	c8.memory[addr%4096] = value
}

func (c8 *Chip8) DumpMemory(start, end uint16) {
	for i := start; i <= end; i += 2 {
		opcode := uint16(c8.memory[i])<<8 | uint16(c8.memory[i+1])
		println(fmt.Sprintf("%04X : %04X", i, opcode))
	}
}

func (c8 *Chip8) Reset() {
	c8.pc = START_ADDRESS
	c8.sp = 0
	c8.index = 0
	c8.DelayTimer = 0
	c8.SoundTimer = 0
	c8.opcode = 0
	c8.paused = false
	c8.rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range len(c8.registers) {
		c8.registers[i] = 0
	}

	for i := range len(c8.memory) {
		c8.memory[i] = 0
	}

	// c8.loadFontset()

	for i := range len(c8.stack) {
		c8.stack[i] = 0
	}

	for i := range len(c8.Video) {
		c8.Video[i] = false
	}

}

func (c8 *Chip8) Pause() {
	c8.paused = true
}

func (c8 *Chip8) Resume() {
	c8.paused = false
}

func (c8 *Chip8) TogglePause() {
	c8.paused = !c8.paused
}

func (c8 *Chip8) IsPaused() bool {
	return c8.paused
}
