package main

import (
	"GoLedMatrix/clog"
	"GoLedMatrix/confload"
	"GoLedMatrix/rgbmatrix"
	"GoLedMatrix/scenario"
	"GoLedMatrix/server"
)

var config = &confload.ConfigData{}

func main() {

	confload.Load("config.ini", config)
	clog.LogLevel = 5
	clog.StartLogging = true

	if config.HTTPserver.Enabled {
		server.StartHTTP(config)
	}

	m, err := rgbmatrix.NewRGBLedMatrix(&config.HardwareConfig, &config.RuntimeOptions)
	clog.Fatal("GoLedMatrix", "main", err)

	// go BouncingBall(&m)
	// go displayGif(&m)
	go scenario.Setup(m)
	m.Start()
}
