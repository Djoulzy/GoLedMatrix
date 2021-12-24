package scenario

import (
	"GoLedMatrix/clog"
	"GoLedMatrix/confload"
	"GoLedMatrix/emulator"
	"GoLedMatrix/rgbmatrix"
	"time"
)

type Scenario struct {
	tk       *rgbmatrix.ToolKit
	conf     *confload.ConfigData
	m        *rgbmatrix.Matrix
	mode     int
	controls *ControlParams
	quit     chan bool
}

type ControlParams struct {
	Mode  int    `json:"mode"`
	Text  string `json:"text"`
	Image string `json:"image"`
	Serie string `json:"serie"`
}

func (S *Scenario) Control(params *ControlParams) {
	clog.Test("Scenario", "Control", "Starting mode: %d", params.Mode)
	S.mode = params.Mode
	S.controls = params
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

	S.controls = &ControlParams{
		Mode:  4,
		Text:  "Joyeux Noël ...",
		Serie: "christmas",
	}
	S.quit = make(chan bool, 0)

	for {
		switch S.controls.Mode {
		case 1:
			S.slideShow()
		case 2:
			S.displayGif()
		case 3:
			S.HorloLed()
		case 4:
			S.ScrollText()
		case 5:
			S.FancyClock()
		}
	}
}
