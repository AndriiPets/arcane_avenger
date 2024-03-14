package main

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

type Hud struct {
	World     WorldInterface
	ScreenW   int
	ScreenH   int
	ColorInfo map[string]int
	KillCount int
	FontFace  font.Face
}

func NewHud(world WorldInterface, w int, h int, font font.Face) *Hud {

	c_map := make(map[string]int)

	hud := &Hud{
		World:     world,
		ScreenW:   w,
		ScreenH:   h,
		KillCount: 0,
		FontFace:  font,
	}

	hud.ColorInfo = c_map

	return hud
}

func (h *Hud) Update() {

	colors := h.World.GetColorInfo()
	h.ColorInfo = colors

	h.KillCount = h.World.GetKillCount()
}

func (h *Hud) Draw(screen *ebiten.Image) {

	red := h.ColorInfo["red"]
	green := h.ColorInfo["green"]
	blue := h.ColorInfo["blue"]

	//draw color bars

	//red
	vector.DrawFilledRect(screen, 20, float32(h.ScreenH)-40, float32(red)*5, 10, color.RGBA{255, 0, 0, 255}, false)

	//border
	vector.StrokeRect(screen, 20, float32(h.ScreenH)-40, 51*5, 10, 2, color.RGBA{255, 255, 255, 255}, false)

	//green
	vector.DrawFilledRect(screen, 20, float32(h.ScreenH)-30, float32(green)*5, 10, color.RGBA{0, 255, 0, 255}, false)

	//border
	vector.StrokeRect(screen, 20, float32(h.ScreenH)-30, 51*5, 10, 2, color.RGBA{255, 255, 255, 255}, false)

	//blue
	vector.DrawFilledRect(screen, 20, float32(h.ScreenH)-20, float32(blue)*5, 10, color.RGBA{0, 0, 255, 255}, false)

	//border
	vector.StrokeRect(screen, 20, float32(h.ScreenH)-20, 51*5, 10, 2, color.RGBA{255, 255, 255, 255}, false)

	//kill count
	text.Draw(screen, strconv.Itoa(h.KillCount), h.FontFace, h.ScreenW/2, 50, color.RGBA{255, 255, 255, 225})

}
