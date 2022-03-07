package confload

import "GoLedMatrix/rgbmatrix"

type HTTPserver struct {
	Addr    string
	Port    int
	Enabled bool
}

type DefaultConf struct {
	InstallDir string
	MediaDir   string
	TmpDir     string
	FontDir    string
}

type API struct {
	QuoteURL     string
	QuoteKey     string
	QuoteSymbols string
}

type ConfigData struct {
	rgbmatrix.HardwareConfig
	rgbmatrix.RuntimeOptions
	HTTPserver
	DefaultConf
	API
}
