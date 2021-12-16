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
	clog.LogLevel = 5
	clog.StartLogging = true
	if config.HTTPserver.Enabled {
		server.StartHTTP(config)
	}

	m, err := rgbmatrix.NewRGBLedMatrix(&config.HardwareConfig, &config.RuntimeOptions)
	fatal(err)

	clog.Trace("main", "main", "start")
	go BouncingBall(&m)
	// go displayGif(&m)
	// go displayImage(m)
	m.Start()
}

func fatal(err error) {
	if err != nil {
		clog.Fatal("GoLedMatrix", "main", err)
		panic(err)
	}
}
