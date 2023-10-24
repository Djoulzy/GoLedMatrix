package main

import (
	"encoding/gob"
	"errors"
	"image"
	"image/draw"

	"github.com/Djoulzy/GoLedMatrix/rgbmatrix"
	"github.com/fogleman/gg"
)

type LayerType int

const (
	CONTEXT LayerType = iota
	IMAGE
)

type Layer struct {
	Type LayerType
	CTX  *gg.Context
	img  *image.Image
}

func (L *Layer) GetContent() (image.Image, error) {
	if L.Type == CONTEXT {
		return L.CTX.Image(), nil
	} else {
		if L.img != nil {
			return *L.img, nil
		}
		return nil, errors.New("no image")
	}
}

func (L *Layer) SetImage(img *image.Image) {
	L.img = img
}

type Display struct {
	TK     *rgbmatrix.ToolKit
	CTX    *gg.Context
	Size   image.Point
	Layers []*Layer
}

func NewDisplay(m *rgbmatrix.Matrix) *Display {
	disp := Display{
		TK: rgbmatrix.NewToolKit(*m),
	}

	disp.Size = disp.TK.Canvas.Bounds().Max
	disp.CTX = gg.NewContext(disp.Size.X, disp.Size.Y)
	disp.Layers = make([]*Layer, 0)

	tmpReg := image.NewNRGBA(disp.TK.Canvas.Bounds())
	gob.Register(disp.CTX.Image())
	gob.Register(tmpReg)
	return &disp
}

func (D *Display) GetLayer(t LayerType) *Layer {
	l := Layer{
		Type: t,
	}
	if t == CONTEXT {
		l.CTX = D.CTX
	} else {
		l.img = nil
	}
	D.Layers = append(D.Layers, &l)
	return &l
}

func (D *Display) Render() {
	var img image.Image
	var err error

	for _, l := range D.Layers {
		if img, err = l.GetContent(); err == nil {
			draw.Draw(D.TK.Canvas, D.TK.Canvas.Bounds(), img, image.Point{}, draw.Over)
		}
	}
	D.TK.Canvas.Render()
}
