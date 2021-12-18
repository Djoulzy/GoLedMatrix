package scenario

import (
	"GoLedMatrix/clog"
	"GoLedMatrix/confload"
	"GoLedMatrix/emulator"
	"GoLedMatrix/rgbmatrix"
	"image"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Scenario struct {
	tk   *rgbmatrix.ToolKit
	conf *confload.ConfigData
	m    *rgbmatrix.Matrix
	mode int
	quit chan bool
}

type ControlParams struct {
	Param1 int    `json:"param1"`
	Test   string `json:"test"`
	Age    int    `json:"age"`
}

func (S *Scenario) slideShow() {
	var d time.Duration = 1000000000

	files, err := ioutil.ReadDir("./img")
	if err != nil {
		log.Fatal(err)
	}

	for _, finfo := range files {
		select {
		case <-S.quit:
			return
		default:
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
}

func (S *Scenario) drawText(text string) {
	actual := time.Now()
	var test = make([]string, 10)

	test[0] = actual.Format("15:04:05")
	S.tk.DrawText(test)
}

func (S *Scenario) displayGif() {

	f, err := os.Open("./anim/muppet.gif")
	if err != nil {
		clog.Fatal("scenario", "displayGif", err)
	}

	close, err := S.tk.PlayGIF(f)
	if err != nil {
		clog.Fatal("scenario", "displayGif", err)
	}

	time.Sleep(time.Second * 30)
	close <- true
}

func (S *Scenario) Control(params *ControlParams) {
	clog.Test("Scenario", "Control", "Params age: %d", params.Age)
	S.mode = 2
	S.quit <- true
}

func (S *Scenario) Run(m interface{}, config *confload.ConfigData) {
	switch m.(type) {
	case emulator.Emulator:
		duration := time.Second * 2
		time.Sleep(duration)
	case rgbmatrix.Matrix:
	}
	t := m.(rgbmatrix.Matrix)

	S.m = &t
	S.conf = config
	S.tk = rgbmatrix.NewToolKit(t)
	defer S.tk.Close()

	S.mode = 3
	S.quit = make(chan bool, 0)

	for {
		switch S.mode {
		case 1:
			S.slideShow()
		case 2:
			S.displayGif()
		case 3:
			S.drawText("test")
		}
	}
}
