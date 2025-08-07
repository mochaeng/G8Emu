package emulator

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mochaeng/G8Emu/internal/constants"
)

type Platform struct {
	display    *ebiten.Image
	keymap     map[ebiten.Key]int
	videoScale int
}

func NewPlatform(videoScale int) *Platform {
	p := &Platform{
		display:    ebiten.NewImage(constants.VIDEO_WIDTH, constants.VIDEO_HEIGHT),
		videoScale: videoScale,
	}

	p.keymap = map[ebiten.Key]int{
		ebiten.KeyX: 0x0, ebiten.Key1: 0x1, ebiten.Key2: 0x2, ebiten.Key3: 0x3,
		ebiten.KeyQ: 0x4, ebiten.KeyW: 0x5, ebiten.KeyE: 0x6, ebiten.KeyA: 0x7,
		ebiten.KeyS: 0x8, ebiten.KeyD: 0x9, ebiten.KeyZ: 0xA, ebiten.KeyC: 0xB,
		ebiten.Key4: 0xC, ebiten.KeyR: 0xD, ebiten.KeyF: 0xE, ebiten.KeyV: 0xF,
	}

	return p
}

// Renders scaled display
func (p *Platform) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(p.videoScale), float64(p.videoScale))
	screen.DrawImage(p.display, op)
}

func (p *Platform) Layout(outsideWidth, outsideHeight int) (int, int) {
	return constants.VIDEO_WIDTH * p.videoScale, constants.VIDEO_HEIGHT * p.videoScale
}

func (p *Platform) ProcessInput(keys []bool) bool {
	for key, chipKey := range p.keymap {
		keys[chipKey] = ebiten.IsKeyPressed(key)
	}

	return ebiten.IsKeyPressed(ebiten.KeyEscape)
}

func (p *Platform) UpdateDisplay(videoBuffer []bool) {
	p.display.Clear()
	for y := range constants.VIDEO_HEIGHT {
		for x := range constants.VIDEO_WIDTH {
			if videoBuffer[y*constants.VIDEO_WIDTH+x] {
				p.display.Set(x, y, color.White)
			} else {
				p.display.Set(x, y, color.Black)
			}
		}
	}
}
