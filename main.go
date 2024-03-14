package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/solarlune/resolv"
	"golang.org/x/image/font"
	"golang.org/x/image/math/f64"

	"github.com/AndriiPets/ArcaneAvenger/resources/fonts"
)

const (
	screenWidth  int = 640
	screenHeight int = 360
)

type Game struct {
	Width, Height int
	World         WorldInterface
	Screen        *ebiten.Image
	WorldScreen   *ebiten.Image
	Camera        *Camera
	Time          float64
	Debug         bool
	Hud           *Hud
	FontFace      font.Face
	Running       bool
}

func NewGame() *Game {

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Arcane Avenger")

	g := &Game{
		Width:   screenWidth,
		Height:  screenHeight,
		Running: true,
	}

	g.Camera = &Camera{ViewPort: f64.Vec2{float64(screenWidth), float64(screenHeight)}}
	g.WorldScreen = ebiten.NewImage(screenWidth*2, screenHeight*2)

	g.World = NewWorld(g)

	m := map[int]resolv.Vector{
		1: resolv.NewVector(100, 100),
		2: resolv.NewVector(500, 100),
		3: resolv.NewVector(500, 100),
		4: resolv.NewVector(500, 500),
	}

	//monster spawners
	for i := 0; i <= 4; i++ {
		ms := NewMonsterSpawner(g.World, m[i])
		g.World.AddMonsterSpawner(ms)
	}

	//fonts
	fontData, _ := truetype.Parse(fonts.ExcelFont)

	opts := &truetype.Options{
		Size:    50,
		DPI:     72,
		Hinting: font.HintingFull,
	}

	g.FontFace = truetype.NewFace(fontData, opts)

	g.Hud = NewHud(g.World, g.Width, g.Height, g.FontFace)

	go func() {

		for {

			fmt.Println("FPS: ", ebiten.ActualFPS())
			fmt.Println("Ticks: ", ebiten.ActualTPS())
			time.Sleep(time.Second)

		}

	}()

	return g
}

func (g *Game) Update() error {

	var quit error

	g.check_for_game_end()

	if g.Running {

		if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
			g.Debug = !g.Debug
		}

		world := g.World

		world.Update()
		g.Hud.Update()

		g.Camera.Update(world.GetPlayerPos(g.Width, g.Height))

		g.Time += 1.0 / 60.0

	} else {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.Running = true
			g.World.SetKillCount(0)
		}
	}

	return quit
}

func (g *Game) check_for_game_end() {
	if g.World.IsPlayerDead() {
		g.Running = false
		g.World.Reset()
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.Running {

		g.Screen = screen
		//screen.Fill(color.RGBA{20, 20, 40, 225})
		g.WorldScreen.Fill(color.RGBA{194, 178, 128, 225})
		g.World.Draw(g.WorldScreen)

		g.Camera.Render(g.WorldScreen, screen)
		g.Hud.Draw(screen)

		//camera debug stuff
		worldX, worldY := g.Camera.ScreenToWorld(ebiten.CursorPosition())
		ebitenutil.DebugPrint(
			screen,
			fmt.Sprintf("TPS: %0.2f\nMove (WASD/Arrows)\nZoom (QE)\nRotate (R)\nReset (Space)", ebiten.ActualTPS()),
		)

		ebitenutil.DebugPrintAt(
			screen,
			fmt.Sprintf("%s\nCursor World Pos: %.2f,%.2f",
				g.Camera.String(),
				worldX, worldY),
			0, screenHeight-32,
		)
	} else {
		g.WorldScreen.Fill(color.RGBA{194, 0, 128, 225})
		text.Draw(screen, "Your Score!", g.FontFace, g.Width/2-100, 50, color.RGBA{255, 255, 255, 225})
		text.Draw(screen, strconv.Itoa(g.World.GetKillCount()), g.FontFace, g.Width/2, g.Height/2, color.RGBA{255, 255, 255, 225})
		text.Draw(screen, "Press space to reset", g.FontFace, 100, 300, color.RGBA{255, 255, 255, 225})
	}
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
