package scenario

import (
	"image"
	"image/color"
	"time"

	"github.com/fogleman/gg"
	"github.com/icza/gox/imagex/colorx"
)

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

func (S *Scenario) ScrollText() {
	size := S.tk.Canvas.Bounds().Max
	center := image.Point{X: size.X / 2, Y: size.Y / 2}

	anim := &TextAnim{
		ctx:      gg.NewContext(size.X, size.Y),
		dir:      image.Point{-1, 0},
		position: image.Point{128, 64},
		message:  S.controls.Text,
		Quit:     S.quit,
	}

	anim.ctx.LoadFontFace("./ttf/marquee/Bullpen3D.ttf", 40)
	anim.txtWidth, anim.txtHeight = anim.ctx.MeasureString(S.controls.Text)
	anim.position = image.Point{size.X, center.Y + int(anim.txtHeight/2)}
	anim.col, _ = colorx.ParseHexColor(S.controls.Color)

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
