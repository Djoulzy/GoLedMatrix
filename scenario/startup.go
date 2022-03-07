package scenario

import (
	anim "GoLedMatrix/anims"
	"fmt"
	"image"
	"time"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/bitmapfont"
)

type Startup struct {
	ctx    *gg.Context
	sprite []*anim.Sprite
	req    StockResponse
}

func (S *Startup) DrawLine(param interface{}) {
	sprite := param.(*anim.Sprite)
	S.ctx.SetHexColor(sprite.FgColor)
	S.ctx.DrawString(sprite.Text, float64(sprite.Pos.X), float64(sprite.Pos.Y))
}

func (S *Scenario) Startup() {
	ticker := time.NewTicker(time.Second * 10)
	defer func() {
		ticker.Stop()
	}()

	startup := Startup{}

	size := S.tk.Canvas.Bounds().Max
	strHeight := 8
	startup.ctx = gg.NewContext(size.X, size.Y)
	startup.ctx.SetFontFace(bitmapfont.Gothic10r)

	startup.sprite = make([]*anim.Sprite, 4)

	startup.sprite[0] = &anim.Sprite{
		ScreenSize: size,
		Size:       image.Point{len("GOLedMatrix") * 5, strHeight},
		Pos:        image.Point{5, strHeight},
		Text:       "GOLedMatrix",
		FgColor:    "#FF0000",
		Dir:        -1,
		Draw:       startup.DrawLine,
	}

	startup.sprite[1] = &anim.Sprite{
		ScreenSize: size,
		Size:       image.Point{len("v0.99 Build 2021-12-30") * 5, strHeight},
		Pos:        image.Point{5, strHeight * 2},
		Text:       "v0.99 Build 2021-12-30",
		FgColor:    "#f29d0c",
		Dir:        -1,
		Draw:       startup.DrawLine,
	}

	startup.sprite[2] = &anim.Sprite{
		ScreenSize: size,
		Size:       image.Point{len("Listen:") * 5, strHeight},
		Pos:        image.Point{5, strHeight * 3},
		Text:       "Listen:",
		FgColor:    "#FF0000",
		Dir:        -1,
		Draw:       startup.DrawLine,
	}

	tmp := fmt.Sprintf("http://%s:%d", S.conf.HTTPserver.Addr, S.conf.HTTPserver.Port)
	startup.sprite[3] = &anim.Sprite{
		ScreenSize: size,
		Size:       image.Point{len(tmp) * 5, strHeight},
		Pos:        image.Point{5, strHeight * 4},
		Text:       tmp,
		FgColor:    "#ffe900",
		Dir:        -1,
		Draw:       startup.DrawLine,
	}

	for {
		select {
		case <-ticker.C:
			return
		default:
			startup.ctx.SetHexColor("#000000")
			startup.ctx.Clear()
			for _, line := range startup.sprite {
				line.Move()
			}
			S.tk.PlayImage(startup.ctx.Image(), time.Millisecond*50)
		}
	}
}
