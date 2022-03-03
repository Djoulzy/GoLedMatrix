package scenario

import (
	"GoLedMatrix/clog"
	"GoLedMatrix/confload"
	"crypto/tls"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"time"
)

func SendToPointTaxi(url string, method string, action string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest(method, url+"/"+action, nil)
	// clog.Test("LePointTaxi", "SendToPointTaxi", "%s", content)
	if err != nil {
		clog.Error("LePointTaxi", "SendToPointTaxi", "%s", err)
		return nil, err
	}
	req.Header.Set("user-agent", "golang ATA Locator")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", url)

	resp, err := client.Do(req)
	if err != nil {
		clog.Error("LePointTaxi", "SendToPointTaxi", "%s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		clog.Warn("LePointTaxi", "SendToPointTaxi", "HTTP Response Status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		if err != nil {
			clog.Error("LePointTaxi", "SendToPointTaxi", "%s", err)
		}
		clog.Warn("LePointTaxi", "SendToPointTaxi", "%s", string(body))
		return nil, err
	}
	return body, nil
}

func (S *Scenario) business(config *confload.ConfigData, text string) {
	actual := time.Now()
	var test = make([]string, 10)

	test[0] = actual.Format("15:04:05")
	x := rand.Intn(55) + 1
	y := rand.Intn(100) + 1
	// S.tk.DrawText(test, x, y, "./ttf/orange_juice.ttf", 24, 1)
	S.tk.DrawText(test, x, y, "./ttf/Perform.ttf", 12, 1)
}
