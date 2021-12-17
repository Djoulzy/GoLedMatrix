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

	m, err := rgbmatrix.NewRGBLedMatrix(&config.HardwareConfig, &config.RuntimeOptions)
	if err != nil {
		clog.Fatal("GoLedMatrix", "main", err)
	}

	var scen scenario.Scenario
	var http server.HTTP

	if config.HTTPserver.Enabled {
		http.StartHTTP(config, &scen)
	}

	go scen.Run(m, config)
	m.Start()
}
