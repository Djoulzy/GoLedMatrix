package main

import (
	"image"
)

type Image struct {
	Disp  *Display
	Layer *Layer
	img   image.Image
	hide  chan bool
}

func InitImage(d *Display) *Image {
	img := Image{
		Disp:  d,
		Layer: d.GetLayer(IMAGE),
	}

	img.hide = make(chan bool)

	return &img
}

func (I *Image) SetImage(img image.Image) {
	I.img = img
	I.Layer.SetImage(&I.img)
}