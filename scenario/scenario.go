package scenario

import (
	"GoLedMatrix/clog"
	"GoLedMatrix/confload"
	"GoLedMatrix/emulator"
	"GoLedMatrix/rgbmatrix"
	"time"
)

type DataParams interface{}

type Scenario struct {
	tk       *rgbmatrix.ToolKit
	conf     *confload.ConfigData
	m        *rgbmatrix.Matrix
	mode     int
	controls DataParams
	quit     chan bool
}

type ControlParams struct {
	Mode         int        `json:"mode"`
	ModuleParams DataParams `json:"params"`
}

func (S *Scenario) Control(params *ControlParams) {
	clog.Test("Scenario", "Control", "%v", params)
	S.mode = params.Mode
	S.controls = params.ModuleParams
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

	S.mode = 1
	S.controls = nil
	S.m = &t
	S.conf = config
	S.tk = rgbmatrix.NewToolKit(t)
	defer S.tk.Close()

	S.quit = make(chan bool, 0)

	for {
		switch S.mode {
		case 1:
			S.OfficeRound()
		case 2:
			S.displayGif()
		case 3:
			S.slideShow()
		case 4:
			S.ScrollText()
		case 5:
			S.FancyClock()
		}
	}
}
