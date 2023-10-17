package scenario

import (
	"clog"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
)

type ImageParams struct {
	Image string `json:"image"`
	Serie string `json:"serie"`
}

func (S *Scenario) slideShow() {
	var imageParams ImageParams
	mapstructure.Decode(S.controls, &imageParams)

	var d time.Duration = time.Second * 3

	mediaDir := fmt.Sprintf("%simg/", S.conf.DefaultConf.MediaDir)

	files, err := ioutil.ReadDir(mediaDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, finfo := range files {
		select {
		case <-S.quit:
			return
		default:
			f, err := os.Open(mediaDir + finfo.Name())
			if err != nil {
				clog.Fatal("scenario", "slideShow", err)
			}
			img, _, err := image.Decode(f)

			err = S.tk.PlayImage(img, d)
			if err != nil {
				clog.Fatal("scenario", "slideShow", err)
			}
		}
	}
}

func (S *Scenario) stillIlage(photo string) {
	var d time.Duration = time.Second * 3

	f, err := os.Open(photo)
	if err != nil {
		clog.Fatal("scenario", "slideShow", err)
	}
	img, _, err := image.Decode(f)

	err = S.tk.PlayImage(img, d)
	if err != nil {
		clog.Fatal("scenario", "slideShow", err)
	}

}
