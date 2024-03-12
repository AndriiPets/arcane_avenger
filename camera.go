package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/math/f64"
)

type Camera struct {
	ViewPort f64.Vec2
	Position f64.Vec2
	Rotation int
	Zoom     int
	cursorX  float64
	cursorY  float64
}

func (c *Camera) ViewPortCenter() f64.Vec2 {

	return f64.Vec2{
		c.ViewPort[0] * 0.5,
		c.ViewPort[1] * 0.5,
	}

}

func (c *Camera) WorldMatrix() ebiten.GeoM {

	m := ebiten.GeoM{}
	m.Translate(-c.Position[0], -c.Position[1])
	m.Translate(-c.ViewPortCenter()[0], -c.ViewPortCenter()[1])

	m.Scale(
		math.Pow(1.01, float64(c.Zoom)),
		math.Pow(1.01, float64(c.Zoom)),
	)

	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)
	m.Translate(c.ViewPortCenter()[0], c.ViewPortCenter()[1])

	return m
}

func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {

	invMatrix := c.WorldMatrix()
	if invMatrix.IsInvertible() {

		invMatrix.Invert()
		return invMatrix.Apply(float64(posX), float64(posY))

	} else {

		return math.NaN(), math.NaN()

	}
}

func (c *Camera) String() string {
	return fmt.Sprintf(
		"T: %.1f, R: %d, S: %d",
		c.Position, c.Rotation, c.Zoom,
	)
}

func (c *Camera) Render(world, screen *ebiten.Image) {
	screen.DrawImage(world, &ebiten.DrawImageOptions{
		GeoM: c.WorldMatrix(),
	})

	//draw aim circle
	mouseX, mouseY := ebiten.CursorPosition()
	mx, my := float32(mouseX), float32(mouseY)

	aimColor := color.RGBA{0, 225, 0, 225}

	vector.StrokeCircle(screen, mx, my, 12, 2, aimColor, false)
	vector.DrawFilledCircle(screen, mx, my, 2, aimColor, false)
}

func (c *Camera) Reset() {

	c.Rotation = 0
	c.Zoom = 0
}

func (c *Camera) Update(px, py float64) {

	c.Position[0] = px
	c.Position[1] = py

	mouseX, mouseY := c.ScreenToWorld(ebiten.CursorPosition())
	c.cursorX = mouseX
	c.cursorY = mouseY

}
