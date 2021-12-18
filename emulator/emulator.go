// +build !arm64

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
)

const DefaultPixelPitch = 12
const windowTitle = "RGB led matrix emulator"

type Emulator struct {
	mu                      sync.RWMutex
	PixelPitch              int
	Gutter                  int
	Width                   int
	Height                  int
	GutterColor             color.NRGBA
	PixelPitchToGutterRatio int
	Margin                  int

	leds []color.Color
	w    *app.Window
	gtx  layout.Context
	wg   sync.WaitGroup

	isReady bool
}

func ConvertColorToNRGBA(col color.Color) color.NRGBA {
	if col == nil {
		return color.NRGBA{0, 0, 0, 0}
	}
	rt, gt, bt, at := col.RGBA()
	nrgba := color.NRGBA{}
	nrgba.R = uint8(rt)
	nrgba.G = uint8(gt)
	nrgba.B = uint8(bt)
	nrgba.A = uint8(at)
	return nrgba
}

func NewEmulator(w, h, pixelPitch int, autoInit bool) *Emulator {
	e := &Emulator{
		Width:                   w,
		Height:                  h,
		GutterColor:             ConvertColorToNRGBA(color.Gray{Y: 20}),
		PixelPitchToGutterRatio: 2,
		Margin:                  10,
		isReady:                 false,
	}
	e.updatePixelPitchForGutter(pixelPitch / e.PixelPitchToGutterRatio)
	e.leds = make([]color.Color, e.Width*e.Height)
	// if autoInit {
	// 	e.Init()
	// }

	return e
}

// Init initialize the emulator, creating a new Window and waiting until is
// painted. If something goes wrong the function panics
func (e *Emulator) Start() {
	go func() {
		dims := e.matrixWithMarginsRect()
		w := app.NewWindow(
			app.Size(unit.Px(float32(dims.Max.X)), unit.Px(float32(dims.Max.Y))),
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
	e.w = w

	var ops op.Ops
	for {
		event := <-w.Events()
		switch evt := event.(type) {
		case system.DestroyEvent:
			return evt.Err
		case system.FrameEvent:
			e.gtx = layout.NewContext(&ops, evt)
			e.drawContext(e.gtx, evt)
			e.Apply(nil)
			evt.Frame(e.gtx.Ops)
			e.isReady = true
		}
	}
}

func (e *Emulator) drawContext(gtx layout.Context, sz system.FrameEvent) {
	e.updatePixelPitchForGutter(e.calculateGutterForViewableArea(sz.Size))
	// Fill entire background with white.
	paint.Fill(gtx.Ops, ConvertColorToNRGBA(color.Black))
	// Fill matrix display rectangle with the gutter color.
	paint.FillShape(gtx.Ops, e.GutterColor, e.matrixWithMarginsRect().Op())
	// Set all LEDs to black.
	// e.Apply(nil)
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
func (e *Emulator) matrixWithMarginsRect() clip.Rect {
	upperLeftLED := e.ledRect(0, 0)
	lowerRightLED := e.ledRect(e.Width-1, e.Height-1)
	mShape := clip.Rect(image.Rect(upperLeftLED.Min.X-e.Margin, upperLeftLED.Min.Y-e.Margin, lowerRightLED.Max.X+e.Margin, lowerRightLED.Max.Y+e.Margin))
	return mShape
}

// ledRect Returns a Rectangle for the LED at col and row.
func (e *Emulator) ledRect(col int, row int) clip.Rect {
	x := (col * (e.PixelPitch + e.Gutter)) + e.Margin
	y := (row * (e.PixelPitch + e.Gutter)) + e.Margin
	return clip.Rect(image.Rect(x, y, x+e.PixelPitch, y+e.PixelPitch))
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
	// defer func() { e.leds = make([]color.Color, e.Height*e.Width) }()

	e.mu.Lock()
	var c color.Color
	for col := 0; col < e.Width; col++ {
		for row := 0; row < e.Height; row++ {
			c = e.At(col + (row * e.Width))
			paint.FillShape(e.gtx.Ops, ConvertColorToNRGBA(c), e.ledRect(col, row).Op())
		}
	}
	e.mu.Unlock()
	return nil
}

func (e *Emulator) Render() error {
	if e.isReady {
		// e.Apply(nil)
		e.w.Invalidate()
	}
	return nil
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
