package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
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

func (h *HTTP) homeHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	t := template.New("")
	if _, err = t.ParseFiles("./server/templates/home.html", "./server/templates/headers.html"); err != nil {
		clog.Fatal("HTTPServer", "homeHandler", err)
	}
	if err = t.ExecuteTemplate(w, "home", nil); err != nil {
		clog.Fatal("HTTPServer", "homeHandler", err)
	}
}

func (h *HTTP) modulesHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	params := mux.Vars(r)
	t := template.New("")
	if _, err = t.ParseFiles("./server/templates/"+params["module"]+".html", "./server/templates/headers.html"); err != nil {
		clog.Fatal("HTTPServer", "modulesHandler", err)
	}
	if err = t.ExecuteTemplate(w, params["module"], nil); err != nil {
		clog.Fatal("HTTPServer", "modulesHandler", err)
	}
}

func (h *HTTP) getFontSize(w http.ResponseWriter, r *http.Request) {
	var err error

	var homeVars struct {
		DefNum  int
		NumList map[int]int
	}

	params := mux.Vars(r)
	debut, _ := strconv.Atoi(params["start"])
	fin, _ := strconv.Atoi(params["end"])

	homeVars.DefNum, _ = strconv.Atoi(r.URL.Query()["default"][0])
	homeVars.NumList = make(map[int]int)
	for i := debut; i <= fin; i++ {
		homeVars.NumList[i] = i
	}
	t := template.New("")
	if _, err = t.ParseFiles("./server/templates/sizelist.html"); err != nil {
		clog.Fatal("HTTPServer", "getFontSize", err)
	}
	if err = t.ExecuteTemplate(w, "sizelist", homeVars); err != nil {
		clog.Fatal("HTTPServer", "getFontSize", err)
	}
}

func (h *HTTP) getDir(w http.ResponseWriter, r *http.Request) {
	var err error
	params := mux.Vars(r)

	var homeVars struct {
		DefVal string
		FList  []fs.FileInfo
	}

	homeVars.DefVal = r.URL.Query()["default"][0]
	homeVars.FList = h.scen.GetDirList(params)
	clog.Warn("HTTPServer", "getDir", "%s", homeVars.DefVal)
	t := template.New("")
	if _, err = t.ParseFiles("./server/templates/dirlist.html"); err != nil {
		clog.Fatal("HTTPServer", "getDir", err)
	}
	if err = t.ExecuteTemplate(w, "dirlist", homeVars); err != nil {
		clog.Fatal("HTTPServer", "getDir", err)
	}
}

func (h *HTTP) setControls(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	req := &scenario.ControlParams{}
	json.Unmarshal(body, &req)

	clog.Test("HTTPServer", "setControls", "%v", req)

	h.scen.Control(req)
	json.NewEncoder(w).Encode("OK")
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

	// clog.Trace("cropAndResize", "SetImageOrientation", "ORIENTATION_LEFT_TOP")
	// mw.SetImageOrientation(imagick.ORIENTATION_LEFT_TOP)

	if ow > oh {
		crop = oh
	} else {
		crop = ow
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

func (h *HTTP) shutdown(w http.ResponseWriter, r *http.Request) {
	out, err := exec.Command("poweroff").Output()
	if err != nil {
		clog.Error("HTTPServer", "shutdown", "%s - %s", out, err)
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode("OK")
	}
}

// StartHTTP : Lance le server HTTP
func (h *HTTP) StartHTTP(config *confload.ConfigData, S *scenario.Scenario) {
	h.scen = S
	router := mux.NewRouter()

	router.HandleFunc("/", h.homeHandler).Methods("GET")
	router.HandleFunc("/modules/{module}", h.modulesHandler).Methods("GET")
	router.HandleFunc("/getDir/{type}", h.getDir).Methods("GET")
	router.HandleFunc("/getDir/{type}/{serie:[a-zA-Z0-9]+}", h.getDir).Methods("GET")
	router.HandleFunc("/getSize/{start}/{end}", h.getFontSize).Methods("GET")
	router.HandleFunc("/upload", h.uploadMedia).Methods("POST")
	router.HandleFunc("/controls", h.setControls).Methods("POST")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./server/static"))))
	router.HandleFunc("/shutdown", h.shutdown).Methods("GET")

	host := fmt.Sprintf("%s:%d", config.HTTPserver.Addr, config.HTTPserver.Port)
	clog.Output("HTTP Server starting listening on %s", host)

	go func() {
		clog.Fatal("server", "HTTP", http.ListenAndServe(host, router))
	}()
}
