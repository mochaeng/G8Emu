package core

func (c8 *Chip8) initTables() {
	c8.table[0x0] = c8.Table0
	c8.table[0x1] = c8.Op1NNN
	c8.table[0x2] = c8.Op2NNN
	c8.table[0x3] = c8.Op3XNN
	c8.table[0x4] = c8.Op4XNN
	c8.table[0x5] = c8.Op5XY0
	c8.table[0x6] = c8.Op6XNN
	c8.table[0x7] = c8.Op7XNN
	c8.table[0x8] = c8.Table8
	c8.table[0x9] = c8.Op9XY0
	c8.table[0xA] = c8.OpANNN
	c8.table[0xB] = c8.OpBNNN
	c8.table[0xC] = c8.OpCXNN
	c8.table[0xD] = c8.OpDXYN
	c8.table[0xE] = c8.TableE
	c8.table[0xF] = c8.TableF

	for i := 0; i <= 0xE; i++ {
		c8.table0[i] = c8.OpNULL
		c8.table8[i] = c8.OpNULL
		c8.tableE[i] = c8.OpNULL
	}

	c8.table0[0x0] = c8.Op00E0
	c8.table0[0xE] = c8.Op00EE

	c8.table8[0x0] = c8.Op8XY0
	c8.table8[0x1] = c8.Op8XY1
	c8.table8[0x2] = c8.Op8XY2
	c8.table8[0x3] = c8.Op8XY3
	c8.table8[0x4] = c8.Op8XY4
	c8.table8[0x5] = c8.Op8XY5
	c8.table8[0x6] = c8.Op8XY6
	c8.table8[0x7] = c8.Op8XY7
	c8.table8[0xE] = c8.Op8XYE

	c8.tableE[0x1] = c8.OpEXA1
	c8.tableE[0xE] = c8.OpEX9E

	for i := 0; i <= 0x65; i++ {
		c8.tableF[i] = c8.OpNULL
	}

	c8.tableF[0x07] = c8.OpFX07
	c8.tableF[0x0A] = c8.OpFX0A
	c8.tableF[0x15] = c8.OpFX15
	c8.tableF[0x18] = c8.OpFX18
	c8.tableF[0x1E] = c8.OpFX1E
	c8.tableF[0x29] = c8.OpFX29
	c8.tableF[0x33] = c8.OpFX33
	c8.tableF[0x55] = c8.OpFX55
	c8.tableF[0x65] = c8.OpFX65
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
