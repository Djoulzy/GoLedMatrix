package main

import (
	"fmt"

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
	fmt.Printf("%v", config)

	m, err := rgbmatrix.NewRGBLedMatrix(&config.HardwareConfig, &config.RuntimeOptions)
	fatal(err)

	// BouncingBall(&m)
	displayGif(&m)
}

func fatal(err error) {
	if err != nil {
		clog.Fatal("GoLedMatrix", "main", err)
		panic(err)
	}
}
