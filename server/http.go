package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"GoLedMatrix/clog"
	"GoLedMatrix/confload"
)

type templateVars struct {
}

var homeVars templateVars

func homeHandler(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, homeVars)
}

// StartHTTP : Lance le server HTTP
func StartHTTP(config *confload.ConfigData) {
	router := mux.NewRouter()
	homeVars = templateVars{}
	router.HandleFunc("/", homeHandler).Methods("GET")
	clog.Output("HTTP Server starting listening on %s", config.HTTPserver.Addr)

	go func() {
		clog.Fatal("server", "HTTP", http.ListenAndServe(config.HTTPserver.Addr, router))
	}()
}
