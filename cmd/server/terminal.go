package main

import (
	"image"
	"time"

	"github.com/Djoulzy/GoLedMatrix/clog"
	"github.com/Djoulzy/GoLedMatrix/rgbmatrix"

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

type Terminal struct {
	ctx        *gg.Context
	tk         *rgbmatrix.ToolKit
	size       image.Point
	lines      []*rgbmatrix.Sprite
	hide       chan bool
	charHeight int
	charWidth  int
	maxLines   int
	nextLine   int
	// req    StockResponse
}

func InitTerminal(m *rgbmatrix.Matrix) *Terminal {
	term := Terminal{
		tk:         rgbmatrix.NewToolKit(*m),
		charHeight: 8,
		charWidth:  5,
	}

	term.size = term.tk.Canvas.Bounds().Max
	term.ctx = gg.NewContext(term.size.X, term.size.Y)
	term.ctx.SetFontFace(bitmapfont.Gothic10r)
	term.maxLines = term.size.Y / term.charHeight
	term.lines = make([]*rgbmatrix.Sprite, term.maxLines)
	term.hide = make(chan bool)

	clog.Warn("Terminal", "InitTerminal", "Nblines = %d", term.maxLines)

	for index := range term.lines {
		term.lines[index] = &rgbmatrix.Sprite{
			ScreenSize: term.size,
			Size:       image.Point{0, term.charHeight},
			Pos:        image.Point{term.charWidth, term.charHeight * (index + 1)},
			Text:       "",
			FgColor:    "#000000",
			DirX:       -1,
			Draw:       term.DrawLine,
		}
	}
	return &term
}

func (T *Terminal) DrawLine(param interface{}) {
	sprite := param.(*rgbmatrix.Sprite)
	T.ctx.SetHexColor(sprite.FgColor)
	T.ctx.DrawString(sprite.Text, float64(sprite.Pos.X), float64(sprite.Pos.Y))
}

func (T *Terminal) ScrollUp() {
	out := T.lines[0]
	for i := 0; i < (T.maxLines - 1); i++ {
		T.lines[i] = T.lines[i+1]
		T.lines[i].Pos = image.Point{T.charWidth, T.charHeight * (i + 1)}
	}
	T.lines[T.maxLines-1] = out
	T.nextLine--
}

func (T *Terminal) AddLine(mess string, color string) {
	if T.nextLine > T.maxLines-1 {
		T.ScrollUp()
	}

	T.lines[T.nextLine].Size = image.Point{len(mess) * T.charWidth, T.charHeight}
	T.lines[T.nextLine].Pos = image.Point{T.charWidth, T.charHeight * (T.nextLine + 1)}
	T.lines[T.nextLine].Text = mess
	T.lines[T.nextLine].FgColor = color
	T.lines[T.nextLine].Draw = T.DrawLine

	T.nextLine++
}

func (T *Terminal) Run() {
	for {
		select {
		case <-T.hide:
			return
		default:
			T.ctx.SetHexColor("#000000")
			T.ctx.Clear()
			for _, line := range T.lines {
				line.Move()
			}
			T.tk.PlayImage(T.ctx.Image(), time.Millisecond*50)
		}
	}
}

func (T *Terminal) Show() {
	go T.Run()
}

func (T *Terminal) Hide() {
	T.hide <- true
}
