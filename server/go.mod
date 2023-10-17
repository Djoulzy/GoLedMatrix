module github.com/Djoulzy/GoLedMatrix/server

go 1.21.3

replace clog => ../clog

replace confload => ../confload

replace scenario => ../scenario

require (
	github.com/gorilla/mux v1.8.0
	gopkg.in/gographics/imagick.v3 v3.5.0
)
