package scenario

import (
	"crypto/tls"
	"fmt"
	"github.com/Djoulzy/GoLedMatrix/clog"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type AuthType int

const (
	XAPIKEY AuthType = iota
	BEARER
	PARAM
)

type ApiToken struct {
	Value string
	Auth  AuthType
}

func APICall(url string, key ApiToken, method string, action string) ([]byte, error) {
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

	if key.Auth == PARAM {
		action += "&token=" + key.Value
	}

	clog.Test("APICall", method, "%s/%s", url, action)
	req, err := http.NewRequest(method, url+"/"+action, nil)
	if err != nil {
		clog.Error("APICall", method, "%s", err)
		return nil, err
	}
	req.Header.Set("user-agent", "GoLedMatrix Agent")
	req.Header.Add("Content-Type", "application/json")
	switch key.Auth {
	case XAPIKEY:
		req.Header.Add("X-Api-Key", key.Value)
	case BEARER:
		bearer := fmt.Sprintf("Bearer %s", key.Value)
		req.Header.Add("Authorization", bearer)
	}
	clog.Test("APICall", "Header", "%v", req.Header)

	resp, err := client.Do(req)
	if err != nil {
		clog.Error("APICall", method, "%s", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		clog.Warn("APICall", method, "HTTP Response Status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		if err != nil {
			clog.Error("APICall", method, "%s", err)
		}
		clog.Warn("APICall", method, "%s", string(body))
		return nil, err
	}
	return body, nil
	// return []byte(temp), nil
}
