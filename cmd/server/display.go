package main

import (
	"image"
	"image/draw"

	"github.com/Djoulzy/GoLedMatrix/rgbmatrix"
	"github.com/fogleman/gg"
)

const (
	CONTEXT int = iota
	IMAGE
)

type Layer struct {
	Type int
	CTX  *gg.Context
	img  image.Image
}

func (L *Layer) GetContent() image.Image {
	if L.Type == CONTEXT {
		return L.CTX.Image()
	}
	return L.img
}

type Display struct {
	TK   *rgbmatrix.ToolKit
	Size image.Point

	Layers []*Layer
}

func NewDisplay(m *rgbmatrix.Matrix) *Display {
	disp := Display{
		TK: rgbmatrix.NewToolKit(*m),
	}

	disp.Size = disp.TK.Canvas.Bounds().Max
	// disp.CTX = gg.NewContext(disp.Size.X, disp.Size.Y)
	disp.Layers = make([]*Layer, 0)

	return &disp
}

func (D *Display) GetLayer() int {
	l := Layer{}
	D.Layers = append(D.Layers, &l)
	return len(D.Layers) - 1
}

func (D *Display) SetImage(img image.Image) {
	D.img = img
}

func (D *Display) Render() {
	draw.Draw(D.TK.Canvas, D.TK.Canvas.Bounds(), D.CTX.Image(), image.Point{}, draw.Over)
	if D.img != nil {
		draw.Draw(D.TK.Canvas, D.TK.Canvas.Bounds(), D.img, image.Point{}, draw.Over)
	}
	D.TK.Canvas.Render()
}
