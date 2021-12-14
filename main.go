package main

import (
	"GoLedMatrix/clog"
	"GoLedMatrix/confload"
	"GoLedMatrix/rgbmatrix"
	"GoLedMatrix/server"
)

var config = &confload.ConfigData{}

func main() {

	confload.Load("config.ini", config)

	if config.HTTPserver.Enabled {
		server.StartHTTP(config)
	}

	m, err := rgbmatrix.NewRGBLedMatrix(&config.HardwareConfig, &config.RuntimeOptions)
	fatal(err)

	BouncingBall(&m)
	// displayGif(&m)
	// displayImage(&m)
}

func fatal(err error) {
	if err != nil {
		clog.Fatal("GoLedMatrix", "main", err)
		panic(err)
	}
}
