module github.com/Djoulzy/GoLedMatrix/scenario

go 1.21.3

replace clog => ../clog

replace rgbmatrix => ../rgbmatrix

replace confload => ../confload

replace emulator => ../emulator

require (
	clog v0.0.0-00010101000000-000000000000
	confload v0.0.0-00010101000000-000000000000
	emulator v0.0.0-00010101000000-000000000000
	github.com/fogleman/gg v1.3.0
	github.com/hajimehoshi/bitmapfont v1.3.1
	github.com/icza/gox v0.0.0-20230924165045-adcb03233bb5
	github.com/mitchellh/mapstructure v1.5.0
	rgbmatrix v0.0.0-00010101000000-000000000000
)

require (
	gioui.org v0.3.1 // indirect
	gioui.org/cpu v0.0.0-20210817075930-8d6a761490d2 // indirect
	gioui.org/shader v1.0.8 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-text/typesetting v0.0.0-20230803102845-24e03d8b5372 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20221012211006-4de253d81b95 // indirect
	golang.org/x/exp/shiny v0.0.0-20220827204233-334a2380cb91 // indirect
	golang.org/x/image v0.13.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
