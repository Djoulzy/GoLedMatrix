package scenario

import (
	"fmt"
	"time"

	"github.com/fogleman/gg"
)

func (S *Scenario) Startup() {
	size := S.tk.Canvas.Bounds().Max

	ctx := gg.NewContext(size.X, size.Y)

	ctx.SetHexColor("#000000")
	ctx.Clear()
	//ctx.LoadFontFace(S.conf.DefaultConf.FontDir+"fixed/November.ttf", 10)

	message := "GoLedMatrix"

	_, strHeight := ctx.MeasureString(message)
	ctx.SetHexColor("#FF0000")
	ctx.DrawString(message, 0, strHeight)

	host := fmt.Sprintf("%s:%d", S.conf.HTTPserver.Addr, S.conf.HTTPserver.Port)
	ctx.SetHexColor("#FF0000")
	ctx.DrawString(host, 2, (strHeight*2)+2)

	S.tk.PlayImage(ctx.Image(), time.Second*5)
}
