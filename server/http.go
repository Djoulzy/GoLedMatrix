package server

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jdeng/goheif"

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

func (h *HTTP) modulesHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	params := mux.Vars(r)
	clog.Test("HTTPServer", "modulesHandler", "%v", params)
	t := template.New("")
	if _, err = t.ParseFiles("./server/templates/" + params["module"] + ".html"); err != nil {
		clog.Fatal("HTTPServer", "modulesHandler", err)
	}
	if err = t.ExecuteTemplate(w, params["module"], homeVars); err != nil {
		clog.Fatal("HTTPServer", "modulesHandler", err)
	}
}

func (h *HTTP) testFunc(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	req := &scenario.ControlParams{}
	json.Unmarshal(body, &req)

	clog.Test("HTTPServer", "testFunc", "%v", req)

	h.scen.Control(req)
}

func cropAndResize(src, dest string) {
	f, err := os.Open("./img/" + src)
	if err != nil {
		clog.Fatal("scenario", "slideShow", err)
	}

	ext := strings.Split(dest, ".")[1]
	switch ext {
	case "jpg":
	case "jpeg":
	case "heic":
		goheif.Decode(f)
	}
}

func (h *HTTP) uploadMedia(w http.ResponseWriter, r *http.Request) {
	clog.Test("HTTPServer", "uploadMedia", "receiving file")
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("imgFile")
	if err != nil {
		clog.Error("HTTPServer", "uploadMedia", "%s", err)
		return
	}
	defer file.Close()

	tempFile, err := ioutil.TempFile("./media/tmp", "tmp-*")
	if err != nil {
		clog.Error("HTTPServer", "uploadMedia", "%s", err)
	}
	defer func() {
		tempFile.Close()
		clog.Trace("HTTPServer", "uploadMedia", "%s %s", "./media/img/"+tempFile.Name(), "./media/img/"+handler.Filename)
		os.Rename(tempFile.Name(), "./media/img/"+handler.Filename)
	}()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		clog.Error("HTTPServer", "uploadMedia", "%s", err)
	}
	tempFile.Write(fileBytes)
}

// StartHTTP : Lance le server HTTP
func (h *HTTP) StartHTTP(config *confload.ConfigData, S *scenario.Scenario) {
	h.scen = S
	router := mux.NewRouter()
	homeVars = templateVars{}

	router.HandleFunc("/", h.homeHandler).Methods("GET")
	router.HandleFunc("/modules/{module}", h.modulesHandler).Methods("GET")
	router.HandleFunc("/upload", h.uploadMedia).Methods("POST")
	router.HandleFunc("/test", h.testFunc).Methods("POST")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./server/static"))))

	clog.Output("HTTP Server starting listening on %s", config.HTTPserver.Addr)

	go func() {
		clog.Fatal("server", "HTTP", http.ListenAndServe(config.HTTPserver.Addr, router))
	}()
}
