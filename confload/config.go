package confload

import "GoLedMatrix/rgbmatrix"

type HTTPserver struct {
	Addr    string
	Port    int
	Enabled bool
}

type DefaultConf struct {
	StartUpDelay int
	InstallDir   string
	MediaDir     string
	TmpDir       string
	FontDir      string
	Mode         int
}

type API struct {
	QuoteURL      string
	QuoteKey      string
	QuoteSymbols  string
	QuoteInterval int
}

type ConfigData struct {
	rgbmatrix.HardwareConfig
	rgbmatrix.RuntimeOptions
	HTTPserver
	DefaultConf
	API
}
