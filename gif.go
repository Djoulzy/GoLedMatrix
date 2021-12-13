package main

import (
	"GoLedMatrix/rgbmatrix"
	"os"
	"time"
)

func displayGif(m *rgbmatrix.Matrix) {
	tk := rgbmatrix.NewToolKit(*m)
	defer tk.Close()

	f, err := os.Open("./muppet.gif")
	fatal(err)

	close, err := tk.PlayGIF(f)
	fatal(err)

	time.Sleep(time.Second * 30)
	close <- true
}
