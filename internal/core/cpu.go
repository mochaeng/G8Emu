package core

import (
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
	c8.opcode = (uint16(c8.memory[c8.pc])<<8 | uint16(c8.memory[c8.pc+1]))
	c8.pc += 2
}

func (c8 *Chip8) decodeAndExecute() {
	firstNibble := (c8.opcode & 0xF000) >> 12
	c8.table[firstNibble]()
}

func (c8 *Chip8) Cycle() {
	c8.fetch()
	c8.decodeAndExecute()

	// if c8.delayTimer > 0 {
	// 	c8.delayTimer--
	// }

	// if c8.soundTimer > 0 {
	// 	c8.soundTimer--
	// }
}
