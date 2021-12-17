package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"GoLedMatrix/clog"
	"GoLedMatrix/confload"
	"GoLedMatrix/scenario"
)

type HTTP struct {
	scen *scenario.Scenario
}

type templateVars struct {
}

var homeVars templateVars

func (h *HTTP) homeHandler(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, homeVars)
}

func (h *HTTP) testFunc(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	req := &scenario.ControlParams{}
	json.Unmarshal(body, &req)

	clog.Test("HTTPServer", "testFunc", "%v", req)

	h.scen.Control(req)
}

// StartHTTP : Lance le server HTTP
func (h *HTTP) StartHTTP(config *confload.ConfigData, S *scenario.Scenario) {
	h.scen = S
	router := mux.NewRouter()
	homeVars = templateVars{}
	router.HandleFunc("/", h.homeHandler).Methods("GET")
	router.HandleFunc("/test", h.testFunc).Methods("POST")
	clog.Output("HTTP Server starting listening on %s", config.HTTPserver.Addr)

	go func() {
		clog.Fatal("server", "HTTP", http.ListenAndServe(config.HTTPserver.Addr, router))
	}()
}
