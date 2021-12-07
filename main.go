package main

import (
	"flag"
	"image"
	"image/color"
	"time"

	"GoLedMatrix/confload"
	"GoLedMatrix/rgbmatrix"

	"github.com/fogleman/gg"
)

type ConfigData struct {
	rgbmatrix.HardwareConfig
	rgbmatrix.RuntimeOptions
}

var config = &ConfigData{}

func main() {

	confload.Load("config.ini", config)

	m, err := rgbmatrix.NewRGBLedMatrix(&config.HardwareConfig)
	fatal(err)

	tk := rgbmatrix.NewToolKit(m)
	defer tk.Close()

	tk.PlayAnimation(NewAnimation(image.Point{128, 64}))
}

func init() {
	flag.Parse()
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}

type Animation struct {
	ctx      *gg.Context
	position image.Point
	dir      image.Point
	stroke   int
}

func NewAnimation(sz image.Point) *Animation {
	return &Animation{
		ctx:    gg.NewContext(sz.X, sz.Y),
		dir:    image.Point{1, 1},
		stroke: 5,
	}
}

func (a *Animation) Next() (image.Image, <-chan time.Time, error) {
	defer a.updatePosition()

	a.ctx.SetColor(color.Black)
	a.ctx.Clear()

	a.ctx.DrawCircle(float64(a.position.X), float64(a.position.Y), float64(a.stroke))
	a.ctx.SetColor(color.RGBA{255, 0, 0, 255})
	a.ctx.Fill()
	return a.ctx.Image(), time.After(time.Millisecond * 50), nil
}

func (a *Animation) updatePosition() {
	a.position.X += 1 * a.dir.X
	a.position.Y += 1 * a.dir.Y

	if a.position.Y+a.stroke > a.ctx.Height() {
		a.dir.Y = -1
	} else if a.position.Y-a.stroke < 0 {
		a.dir.Y = 1
	}

	if a.position.X+a.stroke > a.ctx.Width() {
		a.dir.X = -1
	} else if a.position.X-a.stroke < 0 {
		a.dir.X = 1
	}
}
