package scenario

import (
	"clog"
	"image"
	"math"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
	"github.com/mitchellh/mapstructure"
)

type ClockParams struct {
	FGColor1 string `json:"fgcolor1"`
	FGColor2 string `json:"fgcolor2"`
	BGColor  string `json:"bgcolor"`
	FontFace string `json:"fontface"`
	FontSize int    `json:"fontsize"`
}

func validateClockParams(params DataParams, defParams *ClockParams) *ClockParams {
	var clockParams ClockParams
	if params != nil {
		if err := mapstructure.Decode(params, &clockParams); err != nil {
			clog.Fatal("clock", "validateParams", err)
		}
	}

	if clockParams.FontFace == "" {
		clockParams.FontFace = defParams.FontFace
	}
	if clockParams.FontSize == 0 {
		clockParams.FontSize = defParams.FontSize
	}
	if clockParams.FGColor1 == "" {
		clockParams.FGColor1 = defParams.FGColor1
	}
	if clockParams.FGColor2 == "" {
		clockParams.FGColor2 = defParams.FGColor2
	}
	if clockParams.BGColor == "" {
		clockParams.BGColor = defParams.BGColor
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
	defaultParams := ClockParams{
		FontFace: "modern/HappyBomb.ttf",
		FontSize: 55,
		FGColor1: "#FF8337",
		FGColor2: "#7be0de",
		BGColor:  "#000000",
	}
	clockParams := validateClockParams(S.controls, &defaultParams)

	size := S.tk.Canvas.Bounds().Max
	ctx := gg.NewContext(size.X, size.Y)
	center := image.Point{X: size.X / 2, Y: size.Y / 2}
	ctx.LoadFontFace(S.conf.DefaultConf.FontDir+clockParams.FontFace, float64(clockParams.FontSize))

	for {
		select {
		case <-S.quit:
			return
		default:
			ctx.SetHexColor(clockParams.BGColor)
			ctx.Clear()

			actual := time.Now()
			timeHour := actual.Format("15")
			timeMinute := actual.Format("04")
			timeHourWidth, _ := ctx.MeasureString(timeHour)
			timeMinuteWidth, timeMinuteHeight := ctx.MeasureString(timeMinute)

			ctx.SetHexColor(clockParams.FGColor1)
			ctx.DrawString(timeHour, float64(center.X)-(timeHourWidth/2), float64(center.Y))
			ctx.SetHexColor(clockParams.FGColor2)
			ctx.DrawString(timeMinute, float64(center.X)-(timeMinuteWidth/2), float64(center.Y)+20+timeMinuteHeight)

			S.tk.PlayImage(ctx.Image(), time.Second)
		}
	}
}

func (S *Scenario) OfficeRound() {

	defaultParams := ClockParams{
		FontFace: "digital/TickingTimebomb.ttf",
		FontSize: 38,
		FGColor1: "#FF0000",
		FGColor2: "#FFFFFF",
		BGColor:  "#000000",
	}
	clockParams := validateClockParams(S.controls, &defaultParams)

	size := S.tk.Canvas.Bounds().Max
	ctx := gg.NewContext(size.X, size.Y)
	center := image.Point{X: size.X / 2, Y: size.Y / 2}
	DeuxPi := 2 * math.Pi
	div12 := DeuxPi / 12
	div60 := DeuxPi / 60
	rotate := 90 * math.Pi / 180
	r1 := float64(center.Y) - 8
	r2 := float64(center.Y) - 2

	var dotSize float64
	if size.X > 64 {
		dotSize = 1
	} else {
		dotSize = 0.7
	}

	ctx.LoadFontFace(S.conf.DefaultConf.FontDir+clockParams.FontFace, float64(clockParams.FontSize))

	for {
		select {
		case <-S.quit:
			return
		default:
			ctx.SetHexColor(clockParams.BGColor)
			ctx.Clear()

			ctx.SetHexColor(clockParams.FGColor2)
			var t, x, y float64
			var sec int
			for t = 0; t <= DeuxPi; t += div12 {
				x = float64(center.X) + r1*math.Cos(t)
				y = float64(center.Y) + r1*math.Sin(t)
				ctx.DrawPoint(x, y, dotSize)
			}
			ctx.Stroke()

			ctx.SetHexColor(clockParams.FGColor1)
			actual := time.Now()
			timeString := actual.Format("15:04")
			sec = 0
			for t = 0; t <= DeuxPi; t += div60 {
				x = float64(center.X) + r2*math.Cos(t-rotate)
				y = float64(center.Y) + r2*math.Sin(t-rotate)
				ctx.DrawPoint(x, y, dotSize)
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
