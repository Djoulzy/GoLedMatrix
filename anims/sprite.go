package anim

import (
	"image"
)

type Sprite struct {
	ScreenSize image.Point
	Size       image.Point
	Pos        image.Point
	Dir        int
	Draw       func(interface{})
	Text       string
	BgColor    string
	FgColor    string
}

func (S *Sprite) Move() {
	S.Draw(S)
	if S.Size.X > S.ScreenSize.X {
		S.Pos.X += S.Dir
		if S.Pos.X+S.Size.X < S.ScreenSize.X {
			S.Dir = 1
		}
		if S.Pos.X == 5 {
			S.Dir = -1
		}
	}
}
