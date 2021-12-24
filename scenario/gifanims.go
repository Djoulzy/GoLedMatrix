package scenario

import (
	"GoLedMatrix/clog"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func (S *Scenario) displayGif(serie string) {
	ticker := time.NewTicker(time.Second * 10)
	defer func() {
		ticker.Stop()
	}()

	files, err := ioutil.ReadDir("./media/anim/" + serie)
	if err != nil {
		log.Fatal(err)
	}

	for _, finfo := range files {
		f, err := os.Open("./media/anim/" + serie + "/" + finfo.Name())
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
