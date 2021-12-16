package main

import (
	"GoLedMatrix/rgbmatrix"
	"os"
	"time"
)

func displayGif(m *rgbmatrix.Matrix) {
	duration := time.Second
	time.Sleep(duration)

	tk := rgbmatrix.NewToolKit(*m)
	defer tk.Close()

	f, err := os.Open("./anim/muppet.gif")
	fatal(err)

	close, err := tk.PlayGIF(f)
	fatal(err)

	time.Sleep(time.Second * 30)
	close <- true
}
