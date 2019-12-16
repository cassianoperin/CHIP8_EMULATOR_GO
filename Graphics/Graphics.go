package Graphics

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"Chip8/CPU"
)

var (
	win	*pixelgl.Window
)

const (
	sizeX		= 64
	sizeY		= 32
	screenWidth	= float64(1024)
	screenHeight 	= float64(768)
)

// Print Graphics on Console
func drawGraphicsConsole() {
	newline := 64
	for index := 0; index < 64*32; index++ {
		switch index {
		case newline:
		  fmt.Printf("\n")
			newline += 64
	  }
    if CPU.Graphics[index] == 0 {
			fmt.Printf(" ")
		} else {
			fmt.Printf("#")
		}
	}
	fmt.Printf("\n")
}


func renderGraphics() {
	cfg := pixelgl.WindowConfig{
		Title:  "Chip8",
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  true,
	}
	var err error
	win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
}


func drawGraphics(graphics [64 * 32]byte) {
	// Background color
	win.Clear(colornames.Black)
	imd := imdraw.New(nil)
	imd.Color = pixel.RGB(1, 1, 1)
	screenWidth := win.Bounds().W()
	width := screenWidth/sizeX
	height := screenHeight/sizeY

	for gfxindex := 0 ; gfxindex < len(CPU.Graphics) ; gfxindex++ {
		if (CPU.Graphics[gfxindex] ==1 ){
			if (gfxindex < 64) {
				x := gfxindex
				y := 31
				imd.Push(pixel.V ( width*float64(x), height*float64(y) ) )
				imd.Push(pixel.V ( width*float64(x)+width, height*float64(y)+height ) )
				imd.Rectangle(0)
			} else {
				y := 31
				x := 0
				nro := gfxindex
				for nro >= 64 {
					nro -= 64
					y = y - 1
					x = nro
				}
				imd.Push(pixel.V ( width*float64(x), height*float64(y) ) )
				imd.Push(pixel.V ( width*float64(x)+width, height*float64(y)+height ) )
				imd.Rectangle(0)
			}
		}
	}

	imd.Draw(win)
	win.Update()
}


func Keyboard() {
	for index, key := range CPU.KeyPressed {
		if win.Pressed(key) {
			CPU.Key[index] = 1
		}else {
			CPU.Key[index] = 0
		}
	}
}


func Run() {

	// Set up render system
	renderGraphics()

	// Main Infinite Loop
	for !win.Closed() {

		// Esc to quit program
		if win.Pressed(pixelgl.KeyEscape) {
			break
		}

		// Handle Keys pressed
		Keyboard()

		// Calls CPU Interpreter
		CPU.Interpreter()

		// If necessary, DRAW
		if CPU.DrawFlag {
			drawGraphics(CPU.Graphics)
		}

		// Draw Graphics on Console
		//drawGraphicsConsole()

		// Update Input Events
		win.UpdateInput()
	}

}
