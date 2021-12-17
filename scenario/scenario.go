package scenario

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

type Scenario struct {
	tk *rgbmatrix.ToolKit
}

func (S *Scenario) slideShow() {
	var d time.Duration = 1000000000

	files, err := ioutil.ReadDir("./img")
	if err != nil {
		log.Fatal(err)
	}

	for _, finfo := range files {
		clog.Test("main", "displayImage", "Render")
		f, err := os.Open("./img/" + finfo.Name())
		if err != nil {
			clog.Fatal("scenario", "slideShow", err)
		}
		img, _, err := image.Decode(f)

		err = S.tk.PlayImage(img, d)
		if err != nil {
			clog.Fatal("scenario", "slideShow", err)
		}
	}
}

func Setup(m interface{}) {
	modeLoop := Scenario{}

	switch m.(type) {
	case emulator.Emulator:
		duration := time.Second
		time.Sleep(duration)
	case rgbmatrix.Matrix:
	}
	t := m.(rgbmatrix.Matrix)

	modeLoop.tk = rgbmatrix.NewToolKit(t)
	defer modeLoop.tk.Close()

	clog.Trace("scenario", "Setup", "Starting default mode")
	for {
		modeLoop.slideShow()
	}
}
