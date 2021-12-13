package main

import (
	"GoLedMatrix/rgbmatrix"
	"image"
	"os"
	"time"
)

func displayImlage(m *rgbmatrix.Matrix) {
	var d time.Duration = 100000000000

	tk := rgbmatrix.NewToolKit(*m)
	defer tk.Close()

	f, err := os.Open("./img/IMG_5427.jpg")
	fatal(err)
	img, _, err := image.Decode(f)

	err = tk.PlayImage(img, d)
	fatal(err)
}
