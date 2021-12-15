package emulator

import (
	"image"
	"image/color"
	"log"
	"os"
	"sync"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"golang.org/x/exp/shiny/screen"
)

const DefaultPixelPitch = 12
const windowTitle = "RGB led matrix emulator"

type Emulator struct {
	PixelPitch              int
	Gutter                  int
	Width                   int
	Height                  int
	GutterColor             color.Color
	PixelPitchToGutterRatio int
	Margin                  int

	leds []color.Color
	w    *app.Window
	// s    screen.Screen
	wg sync.WaitGroup

	isReady bool
}

func NewEmulator(w, h, pixelPitch int, autoInit bool) *Emulator {
	e := &Emulator{
		Width:                   w,
		Height:                  h,
		GutterColor:             color.Gray{Y: 20},
		PixelPitchToGutterRatio: 2,
		Margin:                  10,
	}
	e.updatePixelPitchForGutter(pixelPitch / e.PixelPitchToGutterRatio)

	if autoInit {
		e.Init()
	}

	return e
}

// Init initialize the emulator, creating a new Window and waiting until is
// painted. If something goes wrong the function panics
func (e *Emulator) Init() {
	e.leds = make([]color.Color, e.Width*e.Height)

	// e.wg.Add(1)
	// go driver.Main(e.mainWindowLoop)
	// e.wg.Wait()

	dims := e.matrixWithMarginsRect()
	go func() {
		w := app.NewWindow(
			app.Size(unit.Dp(float32(dims.Max.X)), unit.Dp(float32(dims.Max.Y))),
			app.Title(windowTitle),
		)
		if err := e.mainWindowLoop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func (e *Emulator) mainWindowLoop(w *app.Window) error {
	// var err error
	e.w = w
	// // Calculate initial window size based on whatever our gutter/pixel pitch currently is.
	// dims := e.matrixWithMarginsRect()
	// e.w, err = s.NewWindow(&screen.NewWindowOptions{
	// 	Title:  windowTitle,
	// 	Width:  dims.Max.X,
	// 	Height: dims.Max.Y,
	// })

	// if err != nil {
	// 	panic(err)
	// }

	// defer e.w.Release()
	var ops op.Ops
	for {
		event := <-w.Events()
		switch evt := event.(type) {
		case system.DestroyEvent:
			return evt.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, evt)
			e.drawContext(gtx, evt)
			evt.Frame(gtx.Ops)
		}
	}
	// var sz size.Event
	// for {
	// 	evn := e.w.NextEvent()
	// 	switch evn := evn.(type) {
	// 	case paint.Event:
	// 		e.drawContext(sz)
	// 		if e.isReady {
	// 			continue
	// 		}

	// 		e.Apply(make([]color.Color, e.Width*e.Height))
	// 		e.wg.Done()
	// 		e.isReady = true
	// 	case size.Event:
	// 		sz = evn

	// 	case error:
	// 		fmt.Fprintln(os.Stderr, e)
	// 	}
	// }
}

func (e *Emulator) drawContext(gtx layout.Context,sz system.FrameEvent) {
	e.updatePixelPitchForGutter(e.calculateGutterForViewableArea(sz.Size))
	// Fill entire background with white.
	paint.Fill(gtx.Ops, color.NRGBA{R: 0xff, G: 0xfe, B: 0xe0, A: 0xff})
	// Fill matrix display rectangle with the gutter color.
	paint.FillShape(gtx.Ops, e.GutterColor, e.matrixWithMarginsRect())
	// Set all LEDs to black.
	e.Apply(make([]color.Color, e.Width*e.Height))
}

// Some formulas that allowed me to better understand the drawable area. I found that the math was
// easiest when put in terms of the Gutter width, hence the addition of PixelPitchToGutterRatio.
//
// PixelPitch = PixelPitchToGutterRatio * Gutter
// DisplayWidth = (PixelPitch * LEDColumns) + (Gutter * (LEDColumns - 1)) + (2 * Margin)
// Gutter = (DisplayWidth - (2 * Margin)) / (PixelPitchToGutterRatio * LEDColumns + LEDColumns - 1)
//
//  MMMMMMMMMMMMMMMM.....MMMM
//  MGGGGGGGGGGGGGGG.....GGGM
//  MGLGLGLGLGLGLGLG.....GLGM
//  MGGGGGGGGGGGGGGG.....GGGM
//  MGLGLGLGLGLGLGLG.....GLGM
//  MGGGGGGGGGGGGGGG.....GGGM
//  .........................
//  MGGGGGGGGGGGGGGG.....GGGM
//  MGLGLGLGLGLGLGLG.....GLGM
//  MGGGGGGGGGGGGGGG.....GGGM
//  MMMMMMMMMMMMMMMM.....MMMM
//
//  where:
//    M = Margin
//    G = Gutter
//    L = LED

// matrixWithMarginsRect Returns a Rectangle that describes entire emulated RGB Matrix, including margins.
func (e *Emulator) matrixWithMarginsRect() clip.Op {
	upperLeftLED := e.ledRect(0, 0)
	lowerRightLED := e.ledRect(e.Width-1, e.Height-1)
	mShape := clip.Rect(image.Rect(upperLeftLED.Min.X-e.Margin, upperLeftLED.Min.Y-e.Margin, lowerRightLED.Max.X+e.Margin, lowerRightLED.Max.Y+e.Margin))
	return mShape.Op()
}

// ledRect Returns a Rectangle for the LED at col and row.
func (e *Emulator) ledRect(col int, row int) image.Rectangle {
	x := (col * (e.PixelPitch + e.Gutter)) + e.Margin
	y := (row * (e.PixelPitch + e.Gutter)) + e.Margin
	return image.Rect(x, y, x+e.PixelPitch, y+e.PixelPitch)
}

// calculateGutterForViewableArea As the name states, calculates the size of the gutter for a given viewable area.
// It's easier to understand the geometry of the matrix on screen when put in terms of the gutter,
// hence the shift toward calculating the gutter size.
func (e *Emulator) calculateGutterForViewableArea(size image.Point) int {
	maxGutterInX := (size.X - 2*e.Margin) / (e.PixelPitchToGutterRatio*e.Width + e.Width - 1)
	maxGutterInY := (size.Y - 2*e.Margin) / (e.PixelPitchToGutterRatio*e.Height + e.Height - 1)
	if maxGutterInX < maxGutterInY {
		return maxGutterInX
	}
	return maxGutterInY
}

func (e *Emulator) updatePixelPitchForGutter(gutterWidth int) {
	e.PixelPitch = e.PixelPitchToGutterRatio * gutterWidth
	e.Gutter = gutterWidth
}

func (e *Emulator) Geometry() (width, height int) {
	return e.Width, e.Height
}

func (e *Emulator) Apply(leds []color.Color) error {
	defer func() { e.leds = make([]color.Color, e.Height*e.Width) }()

	var c color.Color
	for col := 0; col < e.Width; col++ {
		for row := 0; row < e.Height; row++ {
			c = e.At(col + (row * e.Width))
			e.w.Fill(e.ledRect(col, row), c, screen.Over)
		}
	}

	e.w.Publish()
	return nil
}

func (e *Emulator) Render() error {
	return e.Apply(e.leds)
}

func (e *Emulator) At(position int) color.Color {
	if e.leds[position] == nil {
		return color.Black
	}

	return e.leds[position]
}

func (e *Emulator) Set(position int, c color.Color) {
	e.leds[position] = color.RGBAModel.Convert(c)
}

func (e *Emulator) Close() error {
	return nil
}
