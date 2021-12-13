package main

import (
	"GoLedMatrix/rgbmatrix"
	"image"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func displayImage(m *rgbmatrix.Matrix) {
	var d time.Duration = 10000000000

	tk := rgbmatrix.NewToolKit(*m)
	defer tk.Close()

	files, err := ioutil.ReadDir("./img/")
	if err != nil {
		log.Fatal(err)
	}

	for _, finfo := range files {
		f, err := os.Open(finfo.Name())
		fatal(err)
		img, _, err := image.Decode(f)

		err = tk.PlayImage(img, d)
		fatal(err)
	}
}
