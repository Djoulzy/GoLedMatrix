module github.com/Djoulzy/GoLedMatrix/rgbmatrix

go 1.21.3

replace emulator => ../emulator

replace clog => ../clog

require (
	clog v0.0.0-00010101000000-000000000000
	emulator v0.0.0-00010101000000-000000000000
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	golang.org/x/image v0.13.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
)

require (
	gioui.org v0.3.1 // indirect
	gioui.org/cpu v0.0.0-20210817075930-8d6a761490d2 // indirect
	gioui.org/shader v1.0.8 // indirect
	github.com/go-text/typesetting v0.0.0-20230803102845-24e03d8b5372 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/kr/text v0.1.0 // indirect
	golang.org/x/exp v0.0.0-20221012211006-4de253d81b95 // indirect
	golang.org/x/exp/shiny v0.0.0-20220827204233-334a2380cb91 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.13.0 // indirect
)
