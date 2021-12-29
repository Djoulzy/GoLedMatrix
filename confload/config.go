package confload

import "GoLedMatrix/rgbmatrix"

type HTTPserver struct {
	Addr    string
	Port    int
	Enabled bool
}

type DefaultConf struct {
	Mode     int
	MediaDir string
	TmpDir   string
	FontDir  string
}

type ConfigData struct {
	rgbmatrix.HardwareConfig
	rgbmatrix.RuntimeOptions
	HTTPserver
	DefaultConf
}
