package scenario

import (
	"GoLedMatrix/clog"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
)

type GifParams struct {
	Image string `json:"image"`
	Serie string `json:"serie"`
}

func (S *Scenario) displayGif() {
	var gifParams GifParams
	mapstructure.Decode(S.controls, &gifParams)

	ticker := time.NewTicker(time.Second * 10)
	defer func() {
		ticker.Stop()
	}()

	mediaDir := fmt.Sprintf("%sanim/", S.conf.DefaultConf.MediaDir)

	files, err := ioutil.ReadDir(mediaDir + gifParams.Serie)
	if err != nil {
		log.Fatal(err)
	}

	for _, finfo := range files {
		f, err := os.Open(mediaDir + gifParams.Serie + "/" + finfo.Name())
		if err != nil {
			clog.Fatal("scenario", "slideShow", err)
		}
		close, err := S.tk.PlayGIF(f)
		select {
		case <-ticker.C:
			close <- true
			break
		case <-S.quit:
			close <- true
			return
		}
	}
}
