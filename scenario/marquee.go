package scenario

import (
	"clog"
	"image"
	"image/color"
	"time"

	"github.com/fogleman/gg"
	"github.com/icza/gox/imagex/colorx"
	"github.com/mitchellh/mapstructure"
)

type ScrollParams struct {
	Text     string `json:"text"`
	FGColor  string `json:"fgcolor"`
	BGColor  string `json:"bgcolor"`
	FontFace string `json:"font"`
	FontSize int    `json:"size"`
}

type TextAnim struct {
	ctx       *gg.Context
	position  image.Point
	dir       image.Point
	message   string
	txtWidth  float64
	txtHeight float64
	col       color.RGBA
	Quit      chan bool
}

type MarqueeParams struct {
	FGColor1 string `json:"fgcolor1"`
	BGColor  string `json:"bgcolor"`
	Message  string `json:"message"`
	FontFace string `json:"fontface"`
	FontSize int    `json:"fontsize"`
}

func validateMarqueeParams(params DataParams, defParams *MarqueeParams) *MarqueeParams {
	var marqueeParams MarqueeParams
	if params != nil {
		if err := mapstructure.Decode(params, &marqueeParams); err != nil {
			clog.Fatal("Marquee", "validateMarqueeParams", err)
		}
	}

	if marqueeParams.FontFace == "" {
		marqueeParams.FontFace = defParams.FontFace
	}
	if marqueeParams.FontSize == 0 {
		marqueeParams.FontSize = defParams.FontSize
	}
	if marqueeParams.FGColor1 == "" {
		marqueeParams.FGColor1 = defParams.FGColor1
	}
	if marqueeParams.Message == "" {
		marqueeParams.Message = defParams.Message
	}
	if marqueeParams.BGColor == "" {
		marqueeParams.BGColor = defParams.BGColor
	}
	return &marqueeParams
}

func (S *Scenario) ScrollText() {
	defaultParams := MarqueeParams{
		FontFace: "marquee/Bullpen3D.ttf",
		FontSize: 40,
		FGColor1: "#FF8337",
		Message:  "GoLedMatrix",
		BGColor:  "#000000",
	}
	scrollParams := validateMarqueeParams(S.controls, &defaultParams)

	size := S.tk.Canvas.Bounds().Max
	center := image.Point{X: size.X / 2, Y: size.Y / 2}

	anim := &TextAnim{
		ctx:      gg.NewContext(size.X, size.Y),
		dir:      image.Point{-1, 0},
		position: image.Point{128, 64},
		message:  scrollParams.Message,
		Quit:     S.quit,
	}

	anim.ctx.LoadFontFace(S.conf.DefaultConf.FontDir+scrollParams.FontFace, float64(scrollParams.FontSize))
	anim.txtWidth, anim.txtHeight = anim.ctx.MeasureString(scrollParams.Message)
	anim.position = image.Point{size.X, center.Y + int(anim.txtHeight/2)}
	anim.col, _ = colorx.ParseHexColor(scrollParams.FGColor1)

	S.tk.PlayAnimation(anim)
}

func (t *TextAnim) Init() chan bool {
	return t.Quit
}

func (t *TextAnim) Next() (image.Image, <-chan time.Time, error) {
	defer t.updatePosition()

	t.ctx.SetColor(color.Black)
	t.ctx.Clear()

	t.ctx.SetColor(t.col)
	t.ctx.DrawString(t.message, float64(t.position.X), float64(t.position.Y))
	t.ctx.Fill()
	return t.ctx.Image(), time.After(time.Millisecond * 20), nil
}

func (t *TextAnim) updatePosition() {
	t.position.X += 1 * t.dir.X

	if t.position.X < int(0-t.txtWidth) {
		t.position.X = t.ctx.Width()
	}
}
