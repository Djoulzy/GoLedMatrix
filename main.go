package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"time"

	"GoLedMatrix/clog"
	"GoLedMatrix/confload"
	"GoLedMatrix/rgbmatrix"

	"github.com/fogleman/gg"
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

	BouncingBall()
}

func fatal(err error) {
	if err != nil {
		clog.Fatal("GoLedMatrix", "main", err)
		panic(err)
	}
}
