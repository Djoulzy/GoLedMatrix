package main

import (
	"GoLedMatrix/clog"
	"GoLedMatrix/emulator"
	"GoLedMatrix/rgbmatrix"
	"image"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func displayImage(m interface{}) {

	switch m.(type) {
	case emulator.Emulator:
		duration := time.Second
		time.Sleep(duration)
	case rgbmatrix.Matrix:
	}
	t := m.(rgbmatrix.Matrix)

	var d time.Duration = 10000000000

	tk := rgbmatrix.NewToolKit(t)
	defer tk.Close()

	files, err := ioutil.ReadDir("./img")
	if err != nil {
		log.Fatal(err)
	}

	for _, finfo := range files {
		clog.Test("main", "displayImage", "Render")
		f, err := os.Open("./img/" + finfo.Name())
		fatal(err)
		img, _, err := image.Decode(f)

		err = tk.PlayImage(img, d)
		fatal(err)
	}
}
