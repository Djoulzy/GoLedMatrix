package main

import (
	"GoLedMatrix/clog"
	"GoLedMatrix/confload"
	"GoLedMatrix/rgbmatrix"
	"GoLedMatrix/scenario"
	"GoLedMatrix/server"
	"net"
)

var config = &confload.ConfigData{}

func getIP() string {
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.IsPrivate() && (v.IP.DefaultMask() != nil) {
					ip = v.IP
					clog.Trace("main", "getIP", "Found: %s", ip.String())
					return ip.String()
				}
			}
		}
	}
	return ""
}

func main() {

	confload.Load("config.ini", config)
	clog.LogLevel = 5
	clog.StartLogging = true

	m, err := rgbmatrix.NewRGBLedMatrix(&config.HardwareConfig, &config.RuntimeOptions)
	if err != nil {
		clog.Fatal("GoLedMatrix", "main", err)
	}

	if config.HTTPserver.Addr == "detect" {
		detectedIp := getIP()
		if detectedIp != "" {
			config.HTTPserver.Addr = detectedIp
		}
	}

	var scen scenario.Scenario
	var http server.HTTP

	if config.HTTPserver.Enabled {
		http.StartHTTP(config, &scen)
	}

	go scen.Run(m, config)
	m.Start()
}
