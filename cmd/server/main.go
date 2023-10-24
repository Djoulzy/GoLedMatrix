package main

import (
	"fmt"
	"net"
	"runtime"

	"github.com/Djoulzy/GoLedMatrix/clog"
	"github.com/Djoulzy/GoLedMatrix/confload"
	"github.com/Djoulzy/GoLedMatrix/rgbmatrix"
)

var version = "No Version Provided"
var goVersion = runtime.Version()
var config = &confload.ConfigData{}
var display *Display
var terminal *Terminal
var graphic *Image

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
	BuildVersion := fmt.Sprintf("%s (%s)", version, goVersion)

	if config.HTTPserver.Addr == "detect" {
		detectedIp := getIP()
		if detectedIp != "" {
			config.HTTPserver.Addr = detectedIp
		}
	}

	m, err := rgbmatrix.NewRGBLedMatrix(&config.HardwareConfig, &config.RuntimeOptions)
	if err != nil {
		clog.Fatal("GoLedServer", "main", err)
	}

	display = NewDisplay(&m)
	terminal = InitTerminal(display)
	graphic = InitImage(display)

	terminal.AddLine("GOLedServer", "#FF0000")
	terminal.AddLine(BuildVersion, "#f29d0c")
	terminal.AddLine("Listen:", "#FF0000")
	terminal.AddLine(fmt.Sprintf("%s:%d", config.HTTPserver.Addr, config.HTTPserver.Port), "#ffe900")
	display.Render()

	go Listener()
	m.Start()
}
