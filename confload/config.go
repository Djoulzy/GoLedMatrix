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
	FgColor  string
	BgColor  string
	TimeFont string
	TextFont string
}

type ConfigData struct {
	rgbmatrix.HardwareConfig
	rgbmatrix.RuntimeOptions
	HTTPserver
	DefaultConf
}
