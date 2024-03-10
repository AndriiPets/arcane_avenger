package main

import (
	//"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/resolv"
)

type Game struct {
	Width, Height int
	World         WorldInterface
	Screen        *ebiten.Image
	Time          float64
	Debug         bool
}

func NewGame() *Game {

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Arcane Avenger")

	g := &Game{
		Width:  640,
		Height: 360,
	}

	g.World = NewWorld(g)

	go func() {

		for {

			//fmt.Println("FPS: ", ebiten.ActualFPS())
			//fmt.Println("Ticks: ", ebiten.ActualTPS())
			time.Sleep(time.Second)

		}

	}()

	return g
}

func (g *Game) Update() error {

	var quit error

	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.Debug = !g.Debug
	}

	world := g.World
	world.Update()

	g.Time += 1.0 / 60.0

	return quit
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Screen = screen
	screen.Fill(color.RGBA{20, 20, 40, 225})
	g.World.Draw(screen)
}

func (g *Game) DebugDraw(screen *ebiten.Image, space *resolv.Space) {

	for y := 0; y < space.Height(); y++ {

		for x := 0; x < space.Width(); x++ {

			cell := space.Cell(x, y)

			cw := float64(space.CellWidth)
			ch := float64(space.CellHeight)
			cx := float64(cell.X) * cw
			cy := float64(cell.Y) * ch

			drawColor := color.RGBA{20, 20, 20, 255}

			if cell.Occupied() {
				drawColor = color.RGBA{255, 255, 0, 255}
			}

			ebitenutil.DrawLine(screen, cx, cy, cx+cw, cy, drawColor)

			ebitenutil.DrawLine(screen, cx+cw, cy, cx+cw, cy+ch, drawColor)

			ebitenutil.DrawLine(screen, cx+cw, cy+ch, cx, cy+ch, drawColor)

			ebitenutil.DrawLine(screen, cx, cy+ch, cx, cy, drawColor)
		}

	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.Width, g.Height
}

func main() {
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
