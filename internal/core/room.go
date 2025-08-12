package core

import (
	"fmt"
	"os"
)

func (c8 *Chip8) LoadRomFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read ROM file: %v", err)
	}

	return c8.LoadRomBytes(data)
}

func (c8 *Chip8) LoadRomBytes(data []byte) error {
	if len(data) > len(c8.memory)-START_ADDRESS {
		return fmt.Errorf("ROM too large to fit in memory: %d bytes (max %d)", len(data), len(c8.memory)-START_ADDRESS)
	}

	copy(c8.memory[START_ADDRESS:], data)

	return nil
}
