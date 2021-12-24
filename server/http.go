package server

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/gographics/imagick.v3/imagick"

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
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	var crop uint

	_, err := os.Open(src)
	if err != nil {
		clog.Fatal("scenario", "slideShow", err)
	}

	err = mw.ReadImage(src)
	if err != nil {
		panic(err)
	}
	mw.SetFormat("jpg")

	fileInfos := strings.Split(dest, ".")

	ow := mw.GetImageWidth()
	oh := mw.GetImageHeight()

	if ow > oh {
		crop = oh
	} else {
		mw.SetImageOrientation(imagick.ORIENTATION_TOP_LEFT)
		crop = oh
	}
	if err = mw.CropImage(crop, crop, int(ow/2-crop/2), int(oh/2-crop/2)); err != nil {
		clog.Fatal("cropAndResize", "CropImage", err)
	}
	if err = mw.ResizeImage(128, 128, imagick.FILTER_LANCZOS2_SHARP); err != nil {
		clog.Fatal("cropAndResize", "ResizeImage", err)
	}
	if err = mw.WriteImage("./media/img/" + fileInfos[0] + ".jpg"); err != nil {
		clog.Fatal("cropAndResize", "WriteImage", err)
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
		clog.Trace("HTTPServer", "uploadMedia", "%s %s", tempFile.Name(), handler.Filename)
		cropAndResize(tempFile.Name(), handler.Filename)
		os.Remove(tempFile.Name())
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
