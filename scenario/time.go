package scenario

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
)

func (S *Scenario) SimpleTime(text string) {
	actual := time.Now()
	var test = make([]string, 10)

	test[0] = actual.Format("15:04:05")
	x := rand.Intn(55) + 1
	y := rand.Intn(100) + 1
	// S.tk.DrawText(test, x, y, "./ttf/orange_juice.ttf", 24, 1)
	S.tk.DrawText(test, x, y, "./ttf/Perform.ttf", 12, 1)
}

func (S *Scenario) HorloLed() {
	size := S.tk.Canvas.Bounds().Max
	ctx := gg.NewContext(size.X, size.Y)
	center := image.Point{X: size.X / 2, Y: size.Y / 2}

	ctx.LoadFontFace("./ttf/digital/frozencrystal.ttf", 35)
	ctx.SetColor(color.Black)
	ctx.Clear()

	ctx.SetColor(color.White)
	var t, x, y float64
	var sec int
	r := float64(center.Y) - 8
	for t = 0; t <= 2*math.Pi; t += (2 * math.Pi) / 12 {
		x = float64(center.X) + r*math.Cos(t)
		y = float64(center.Y) + r*math.Sin(t)
		ctx.DrawPoint(x, y, 1)
	}
	ctx.Stroke()

	ctx.SetColor(color.RGBA{255, 0, 0, 255})
	actual := time.Now()
	timeString := actual.Format("15:04")
	sec = 0
	r = float64(center.Y) - 2
	for t = 0; t <= 2*math.Pi; t += (2 * math.Pi) / 60 {
		x = float64(center.X) + r*math.Cos(t-90*math.Pi/180)
		y = float64(center.Y) + r*math.Sin(t-90*math.Pi/180)
		ctx.DrawPoint(x, y, 1)
		sec++
		if sec > actual.Second() {
			break
		}
	}
	ctx.Stroke()

	timeWidth, timeHeight := ctx.MeasureString(timeString)

	ctx.DrawString(timeString, float64(center.X)-(timeWidth/2), float64(center.Y)+(timeHeight/2))

	S.tk.PlayImage(ctx.Image(), time.Second)
}
