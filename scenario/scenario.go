package scenario

import (
	"io/fs"
	"log"
	"os"
	"time"

	"github.com/Djoulzy/GoLedMatrix/clog"
	"github.com/Djoulzy/GoLedMatrix/confload"
	"github.com/Djoulzy/GoLedMatrix/emulator"
	"github.com/Djoulzy/GoLedMatrix/rgbmatrix"
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

func (S *Scenario) GetDirList(src map[string]string) []fs.DirEntry {
	var path string

	switch src["type"] {
	case "img":
	case "anim":
		path = S.conf.DefaultConf.MediaDir + "anim"
	case "ttf":
		path = S.conf.DefaultConf.FontDir
	}
	if serie, ok := src["serie"]; ok {
		path += serie
	}
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func (S *Scenario) Control(params *ControlParams) {
	clog.Test("Scenario", "Control", "%v", params)
	S.mode = params.Mode
	S.controls = params.ModuleParams
	S.quit <- true
}

func (S *Scenario) Init(m interface{}, config *confload.ConfigData, version string) {
	switch m.(type) {
	case emulator.Emulator:
		duration := time.Second * 1
		time.Sleep(duration)
	case rgbmatrix.Matrix:
	}
	t := m.(rgbmatrix.Matrix)

	S.mode = config.DefaultConf.Mode
	S.controls = nil
	S.m = &t
	S.conf = config
	S.quit = make(chan bool)
	S.tk = rgbmatrix.NewToolKit(t)

	S.Startup(version)
}

func (S *Scenario) Run(m interface{}, config *confload.ConfigData, version string) {

	S.Init(m, config, version)
	defer S.tk.Close()

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
		case 6:
			S.Business()
		}
	}
}
