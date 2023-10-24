package main

import (
	"image"

	"github.com/Djoulzy/GoLedMatrix/rgbmatrix"

	"github.com/hajimehoshi/bitmapfont"
)

type Terminal struct {
	Disp       *Display
	Layer      *Layer
	lines      []*rgbmatrix.Sprite
	hide       chan bool
	charHeight int
	charWidth  int
	maxLines   int
	nextLine   int
}

func InitTerminal(d *Display) *Terminal {
	term := Terminal{
		Disp:       d,
		Layer:      d.GetLayer(CONTEXT),
		charHeight: 8,
		charWidth:  5,
	}

	term.Layer.CTX.SetFontFace(bitmapfont.Gothic10r)
	term.maxLines = term.Disp.Size.Y / term.charHeight
	term.lines = make([]*rgbmatrix.Sprite, term.maxLines)
	term.hide = make(chan bool)

	for index := range term.lines {
		term.lines[index] = &rgbmatrix.Sprite{
			ScreenSize: term.Disp.Size,
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
	T.Layer.CTX.SetHexColor(sprite.FgColor)
	T.Layer.CTX.DrawString(sprite.Text, float64(sprite.Pos.X), float64(sprite.Pos.Y))
}

func (T *Terminal) DeleteLine() {

}

func (T *Terminal) ScrollUp() {
	T.Layer.CTX.SetHexColor("#000000")
	T.Layer.CTX.Clear()

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
	T.Refresh()
}

func (T *Terminal) Refresh() {
	for _, line := range T.lines {
		line.Move()
	}
}
