package scenario

import (
	"GoLedMatrix/clog"
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
	"github.com/icza/gox/imagex/colorx"
	"github.com/mitchellh/mapstructure"
)

type ClockParams struct {
	FGColor1 string `json:"fgcolor1"`
	FGColor2 string `json:"fgcolor2"`
	BGColor  string `json:"bgcolor"`
	FontFace string `json:"font"`
	FontSize int    `json:"size"`
}

func validateParams(params DataParams) *ClockParams {
	var clockParams ClockParams
	if params != nil {
		if err := mapstructure.Decode(params, &clockParams); err != nil {
			clog.Fatal("clock", "validateParams", err)
		}
	}

	if clockParams.FontFace == "" {
		clockParams.FontFace = "./media/ttf/digital/TickingTimebomb.ttf"
	}
	if clockParams.FontSize == 0 {
		clockParams.FontSize = 38
	}
	if clockParams.FGColor1 == "" {
		clockParams.FGColor1 = "#FF0000"
	}
	if clockParams.FGColor2 == "" {
		clockParams.FGColor2 = "#FFFFFF"
	}
	if clockParams.BGColor == "" {
		clockParams.BGColor = "#000000"
	}
	return &clockParams
}

func (S *Scenario) SimpleTime(text string) {
	actual := time.Now()
	var test = make([]string, 10)

	test[0] = actual.Format("15:04:05")
	x := rand.Intn(55) + 1
	y := rand.Intn(100) + 1
	// S.tk.DrawText(test, x, y, "./ttf/orange_juice.ttf", 24, 1)
	S.tk.DrawText(test, x, y, "./ttf/Perform.ttf", 12, 1)
}

func (S *Scenario) FancyClock() {
	size := S.tk.Canvas.Bounds().Max
	ctx := gg.NewContext(size.X, size.Y)
	center := image.Point{X: size.X / 2, Y: size.Y / 2}
	ctx.LoadFontFace("./ttf/modern/HappyBomb.ttf", 55)

	for {
		select {
		case <-S.quit:
			return
		default:
			ctx.SetColor(color.Black)
			ctx.Clear()

			actual := time.Now()
			timeHour := actual.Format("15")
			timeMinute := actual.Format("04")
			timeHourWidth, _ := ctx.MeasureString(timeHour)
			timeMinuteWidth, timeMinuteHeight := ctx.MeasureString(timeMinute)

			ctx.SetColor(color.RGBA{255, 131, 0, 255})
			ctx.DrawString(timeHour, float64(center.X)-(timeHourWidth/2), float64(center.Y))
			ctx.SetColor(color.RGBA{123, 224, 222, 255})
			ctx.DrawString(timeMinute, float64(center.X)-(timeMinuteWidth/2), float64(center.Y)+20+timeMinuteHeight)

			S.tk.PlayImage(ctx.Image(), time.Second)
		}
	}
}

func (S *Scenario) OfficeRound() {
	// clog.Test("scenario", "OfficeRound", "%v", clockParams)
	clockParams := validateParams(S.controls)

	size := S.tk.Canvas.Bounds().Max
	ctx := gg.NewContext(size.X, size.Y)
	center := image.Point{X: size.X / 2, Y: size.Y / 2}
	DeuxPi := 2 * math.Pi
	div12 := DeuxPi / 12
	div60 := DeuxPi / 60
	rotate := 90 * math.Pi / 180
	r1 := float64(center.Y) - 8
	r2 := float64(center.Y) - 2

	ctx.LoadFontFace(clockParams.FontFace, float64(clockParams.FontSize))
	bgcol, _ := colorx.ParseHexColor(clockParams.BGColor)
	fgcol1, _ := colorx.ParseHexColor(clockParams.FGColor1)
	fgcol2, _ := colorx.ParseHexColor(clockParams.FGColor2)

	for {
		select {
		case <-S.quit:
			return
		default:
			ctx.SetColor(bgcol)
			ctx.Clear()

			ctx.SetColor(fgcol2)
			var t, x, y float64
			var sec int
			for t = 0; t <= DeuxPi; t += div12 {
				x = float64(center.X) + r1*math.Cos(t)
				y = float64(center.Y) + r1*math.Sin(t)
				ctx.DrawPoint(x, y, 1)
			}
			ctx.Stroke()

			ctx.SetColor(fgcol1)
			actual := time.Now()
			timeString := actual.Format("15:04")
			sec = 0
			for t = 0; t <= DeuxPi; t += div60 {
				x = float64(center.X) + r2*math.Cos(t-rotate)
				y = float64(center.Y) + r2*math.Sin(t-rotate)
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
	}
}
