package scenario

import (
	"fmt"
	"time"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/bitmapfont"
)

type displayLine struct {
	message string
	width   float64
	posx    float64
	posy    float64
	dir     int
	color   string
}

func (S *Scenario) Startup() {
	size := S.tk.Canvas.Bounds().Max
	strHeight := float64(8)
	ctx := gg.NewContext(size.X, size.Y)
	ctx.SetFontFace(bitmapfont.Gothic10r)

	lines := make([]*displayLine, 4)

	message := "GOLedMatrix"
	lines[0] = &displayLine{
		message: message,
		width:   float64(len(message) * 5),
		posx:    5,
		posy:    strHeight,
		color:   "#FF0000",
		dir:     -1,
	}

	message = "v0.99 Build 2021-12-30"
	lines[1] = &displayLine{
		message: message,
		width:   float64(len(message) * 5),
		posx:    5,
		posy:    strHeight * 2,
		color:   "#f29d0c",
		dir:     -1,
	}

	message = "Listen:"
	lines[2] = &displayLine{
		message: message,
		width:   float64(len(message) * 5),
		posx:    5,
		posy:    strHeight * 3,
		color:   "#FF0000",
		dir:     -1,
	}

	message = fmt.Sprintf("http://%s:%d", S.conf.HTTPserver.Addr, S.conf.HTTPserver.Port)
	lines[3] = &displayLine{
		message: message,
		width:   float64(len(message) * 5),
		posx:    5,
		posy:    strHeight * 4,
		color:   "#ffe900",
		dir:     -1,
	}

	for {
		ctx.SetHexColor("#000000")
		ctx.Clear()
		for _, line := range lines {
			ctx.SetHexColor(line.color)
			ctx.DrawString(line.message, line.posx, line.posy)
			if line.width > float64(size.X) {
				line.posx += float64(line.dir)
				if line.posx+line.width < float64(size.X) {
					line.dir = 1
				}
				if line.posx == 5 {
					line.dir = -1
				}
			}
		}
		S.tk.PlayImage(ctx.Image(), time.Millisecond*50)
	}
}
