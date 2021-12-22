package server

import (
	"encoding/json"
	"html/template"
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
	var err error
	t := template.New("")
	if _, err = t.ParseFiles("./server/templates/home.html"); err != nil {
		clog.Fatal("HTTPServer", "homeHandler", err)
	}
	if err = t.ExecuteTemplate(w, "home", homeVars); err != nil {
		clog.Fatal("HTTPServer", "homeHandler", err)
	}
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
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./server/static"))))

	clog.Output("HTTP Server starting listening on %s", config.HTTPserver.Addr)

	go func() {
		clog.Fatal("server", "HTTP", http.ListenAndServe(config.HTTPserver.Addr, router))
	}()
}
