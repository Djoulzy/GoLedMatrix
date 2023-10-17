package confload

import "rgbmatrix"

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

type QuoteAPI struct {
	QuoteURL      string
	QuoteKey      string
	QuoteSymbols  string
	QuoteInterval int
}

type WeatherAPI struct {
	WeatherURL      string
	WeatherKey      string
	WeatherRoute    string
	WeatherINSEE    string
	WeatherInterval int
}

type ConfigData struct {
	rgbmatrix.HardwareConfig
	rgbmatrix.RuntimeOptions
	HTTPserver
	DefaultConf
	QuoteAPI
	WeatherAPI
}
