package main

import (
	"image"
)

type Image struct {
	Disp *Display
	img  *image.Image
	hide chan bool
}

func InitImage(d *Display) *Image {
	img := Image{
		Disp: d,
	}

	img.hide = make(chan bool)

	return &img
}

func (I *Image) SetImage(img *image.Image) {
	I.img = img
}

func (I *Image) Display() {
	I.Disp.Render()
}
