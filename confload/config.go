package confload

import "GoLedMatrix/rgbmatrix"

type HTTPserver struct {
	Addr    string
	Enabled bool
}

type ConfigData struct {
	rgbmatrix.HardwareConfig
	rgbmatrix.RuntimeOptions
	HTTPserver
}
