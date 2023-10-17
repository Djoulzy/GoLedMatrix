module github.com/Djoulzy/GoLedMatrix

go 1.21.3

replace clog => ./clog

replace confload => ./confload

replace rgbmatrix => ./rgbmatrix

replace scenario => ./scenario

replace server => ./server

replace emulator => ./emulator

require (
	clog v0.0.0-00010101000000-000000000000
	confload v0.0.0-00010101000000-000000000000
	rgbmatrix v0.0.0-00010101000000-000000000000
	scenario v0.0.0-00010101000000-000000000000
	server v0.0.0-00010101000000-000000000000
)

require (
	emulator v0.0.0-00010101000000-000000000000 // indirect
	gioui.org v0.3.1 // indirect
	gioui.org/cpu v0.0.0-20210817075930-8d6a761490d2 // indirect
	gioui.org/shader v1.0.8 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-text/typesetting v0.0.0-20230803102845-24e03d8b5372 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/hajimehoshi/bitmapfont v1.3.1 // indirect
	github.com/icza/gox v0.0.0-20230924165045-adcb03233bb5 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	golang.org/x/exp v0.0.0-20221012211006-4de253d81b95 // indirect
	golang.org/x/exp/shiny v0.0.0-20220827204233-334a2380cb91 // indirect
	golang.org/x/image v0.13.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	gopkg.in/gographics/imagick.v3 v3.5.0 // indirect
)
