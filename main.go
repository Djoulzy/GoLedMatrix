package main

import (
	"GoLedMatrix/clog"
	"GoLedMatrix/confload"
	"GoLedMatrix/rgbmatrix"
)

type ConfigData struct {
	rgbmatrix.HardwareConfig
	rgbmatrix.RuntimeOptions
}

var config = &ConfigData{}

func main() {

	confload.Load("config.ini", config)

	m, err := rgbmatrix.NewRGBLedMatrix(&config.HardwareConfig, &config.RuntimeOptions)
	fatal(err)

	BouncingBall(&m)
	// displayGif(&m)
	// displayImlage(&m)
}

func fatal(err error) {
	if err != nil {
		clog.Fatal("GoLedMatrix", "main", err)
		panic(err)
	}
}
