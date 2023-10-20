package scenario

import (
	"fmt"
	"github.com/Djoulzy/GoLedMatrix/rgbmatrix"
	"image"
	"time"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/bitmapfont"
)

type StartupParams struct {
	FGColor1 string `json:"fgcolor1"`
	BGColor  string `json:"bgcolor"`
	Message  string `json:"message"`
	FontFace string `json:"fontface"`
	FontSize int    `json:"fontsize"`
}

type Startup struct {
	ctx    *gg.Context
	sprite []*rgbmatrix.Sprite
	// req    StockResponse
}

func (S *Startup) DrawLine(param interface{}) {
	sprite := param.(*rgbmatrix.Sprite)
	S.ctx.SetHexColor(sprite.FgColor)
	S.ctx.DrawString(sprite.Text, float64(sprite.Pos.X), float64(sprite.Pos.Y))
}

func (S *Scenario) Startup(version string) {
	ticker := time.NewTicker(time.Second * time.Duration(S.conf.DefaultConf.StartUpDelay))
	defer func() {
		ticker.Stop()
	}()

	startup := Startup{}

	size := S.tk.Canvas.Bounds().Max
	strHeight := 8
	startup.ctx = gg.NewContext(size.X, size.Y)
	startup.ctx.SetFontFace(bitmapfont.Gothic10r)

	startup.sprite = make([]*rgbmatrix.Sprite, 4)

	startup.sprite[0] = &rgbmatrix.Sprite{
		ScreenSize: size,
		Size:       image.Point{len("GOLedMatrix") * 5, strHeight},
		Pos:        image.Point{5, strHeight},
		Text:       "GOLedMatrix",
		FgColor:    "#FF0000",
		DirX:       -1,
		Draw:       startup.DrawLine,
	}

	startup.sprite[1] = &rgbmatrix.Sprite{
		ScreenSize: size,
		Size:       image.Point{len(version) * 5, strHeight},
		Pos:        image.Point{5, strHeight * 2},
		Text:       version,
		FgColor:    "#f29d0c",
		DirX:       -1,
		Draw:       startup.DrawLine,
	}

	startup.sprite[2] = &rgbmatrix.Sprite{
		ScreenSize: size,
		Size:       image.Point{len("Listen:") * 5, strHeight},
		Pos:        image.Point{5, strHeight * 3},
		Text:       "Listen:",
		FgColor:    "#FF0000",
		DirX:       -1,
		Draw:       startup.DrawLine,
	}

	tmp := fmt.Sprintf("http://%s:%d", S.conf.HTTPserver.Addr, S.conf.HTTPserver.Port)
	startup.sprite[3] = &rgbmatrix.Sprite{
		ScreenSize: size,
		Size:       image.Point{len(tmp) * 5, strHeight},
		Pos:        image.Point{5, strHeight * 4},
		Text:       tmp,
		FgColor:    "#ffe900",
		DirX:       -1,
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
