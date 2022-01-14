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
	Serie string `json:"serie"`
}

func validateGifParams(params DataParams, defParams *GifParams) *GifParams {
	var gifParams GifParams
	if params != nil {
		if err := mapstructure.Decode(params, &gifParams); err != nil {
			clog.Fatal("gif_anims", "validateParams", err)
		}
	}

	if gifParams.Serie == "" {
		gifParams.Serie = defParams.Serie
	}
	return &gifParams
}

func (S *Scenario) displayGif() {
	defaultParams := GifParams{
		Serie: "fun",
	}
	gifParams := validateGifParams(S.controls, &defaultParams)

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
