package scenario

import (
	"GoLedMatrix/clog"
	"GoLedMatrix/rgbmatrix"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/bitmapfont"
)

type StockSymbol struct {
	Symbol              string  `mapstructure:"symbol" json:"symbol"`
	QuoteType           string  `mapstructure:"quoteType" json:"quoteType"`
	FromCurrency        string  `mapstructure:"fromCurrency" json:"fromCurrency"`
	CoinImageUrl        string  `mapstructure:"coinImageUrl" json:"coinImageUrl"`
	RegularMarketPrice  float64 `mapstructure:"regularMarketPrice" json:"regularMarketPrice"`
	RegularMarketChange float64 `mapstructure:"regularMarketChange" json:"regularMarketChange"`
}

type StockResult struct {
	Result []StockSymbol `mapstructure:"result" json:"result"`
}

type StockResponse struct {
	Data StockResult `mapstructure:"quoteResponse" json:"quoteResponse"`
}

type Stock struct {
	ctx     *gg.Context
	sprite  *rgbmatrix.Sprite
	req     StockResponse
	message string
	active  bool
}

var temp string = `{
	"quoteResponse": {
	  "result": [
		{
		  "language": "fr-FR",
		  "region": "FR",
		  "quoteType": "CRYPTOCURRENCY",
		  "quoteSourceName": "CoinMarketCap",
		  "triggerable": false,
		  "customPriceAlertConfidence": "LOW",
		  "currency": "EUR",
		  "marketState": "REGULAR",
		  "exchange": "CCC",
		  "shortName": "Neblio EUR",
		  "messageBoardId": "finmb_NEBL_CCC_lang_fr",
		  "exchangeTimezoneName": "UTC",
		  "exchangeTimezoneShortName": "UTC",
		  "gmtOffSetMilliseconds": 0,
		  "market": "ccc_market",
		  "esgPopulated": false,
		  "firstTradeDateMilliseconds": 1510358400000,
		  "priceHint": 6,
		  "circulatingSupply": 18738626,
		  "lastMarket": "CoinMarketCap",
		  "volume24Hr": 77910,
		  "volumeAllCurrencies": 77910,
		  "fromCurrency": "NEBL",
		  "toCurrency": "EUR=X",
		  "fiftyDayAverage": 0.5523587,
		  "fiftyDayAverageChange": -0.105748236,
		  "fiftyDayAverageChangePercent": -0.1914485,
		  "twoHundredDayAverage": 0.98599607,
		  "twoHundredDayAverageChange": -0.5393856,
		  "twoHundredDayAverageChangePercent": -0.5470464,
		  "marketCap": 8368865,
		  "sourceInterval": 15,
		  "exchangeDataDelayedBy": 0,
		  "tradeable": false,
		  "regularMarketChange": -0.0053860843,
		  "regularMarketChangePercent": -1.1916225,
		  "regularMarketTime": 1646291894,
		  "regularMarketPrice": 0.44661045,
		  "regularMarketDayHigh": 0.4575282,
		  "regularMarketDayRange": "0.44100997 - 0.4575282",
		  "regularMarketDayLow": 0.44100997,
		  "regularMarketVolume": 77910,
		  "regularMarketPreviousClose": 0.4546528,
		  "fullExchangeName": "CCC",
		  "regularMarketOpen": 0.4546528,
		  "averageDailyVolume3Month": 640999,
		  "averageDailyVolume10Day": 103056,
		  "startDate": 1505174400,
		  "coinImageUrl": "https://s.yimg.com/uc/fin/img/reports-thumbnails/1955.png",
		  "fiftyTwoWeekLowChange": 0.082891464,
		  "fiftyTwoWeekLowChangePercent": 0.22789975,
		  "fiftyTwoWeekRange": "0.363719 - 4.452649",
		  "fiftyTwoWeekHighChange": -4.0060387,
		  "fiftyTwoWeekHighChangePercent": -0.89969784,
		  "fiftyTwoWeekLow": 0.363719,
		  "fiftyTwoWeekHigh": 4.452649,
		  "symbol": "NEBL-EUR"
		},
		{
		  "language": "fr-FR",
		  "region": "FR",
		  "quoteType": "CRYPTOCURRENCY",
		  "quoteSourceName": "CoinMarketCap",
		  "triggerable": false,
		  "customPriceAlertConfidence": "LOW",
		  "currency": "EUR",
		  "marketState": "REGULAR",
		  "exchange": "CCC",
		  "shortName": "Algorand EUR",
		  "messageBoardId": "finmb_ALGO_CCC_lang_fr",
		  "exchangeTimezoneName": "UTC",
		  "exchangeTimezoneShortName": "UTC",
		  "gmtOffSetMilliseconds": 0,
		  "market": "ccc_market",
		  "esgPopulated": false,
		  "firstTradeDateMilliseconds": 1561075200000,
		  "priceHint": 6,
		  "circulatingSupply": 6622257152,
		  "lastMarket": "CoinMarketCap",
		  "volume24Hr": 397833152,
		  "volumeAllCurrencies": 397833152,
		  "fromCurrency": "ALGO",
		  "toCurrency": "EUR=X",
		  "fiftyDayAverage": 0.89223146,
		  "fiftyDayAverageChange": -0.1418873,
		  "fiftyDayAverageChangePercent": -0.15902522,
		  "twoHundredDayAverage": 1.2961924,
		  "twoHundredDayAverageChange": -0.54584825,
		  "twoHundredDayAverageChangePercent": -0.42111668,
		  "marketCap": 4968971776,
		  "sourceInterval": 15,
		  "exchangeDataDelayedBy": 0,
		  "tradeable": false,
		  "regularMarketChange": -0.0023061037,
		  "regularMarketChangePercent": -0.30639586,
		  "regularMarketTime": 1646291894,
		  "regularMarketPrice": 0.75034416,
		  "regularMarketDayHigh": 0.76501757,
		  "regularMarketDayRange": "0.7275646 - 0.76501757",
		  "regularMarketDayLow": 0.7275646,
		  "regularMarketVolume": 397833152,
		  "regularMarketPreviousClose": 0.7385995,
		  "fullExchangeName": "CCC",
		  "regularMarketOpen": 0.7385995,
		  "averageDailyVolume3Month": 308882150,
		  "averageDailyVolume10Day": 237881846,
		  "startDate": 1560988800,
		  "coinImageUrl": "https://s.yimg.com/uc/fin/img/reports-thumbnails/4030.png",
		  "fiftyTwoWeekLowChange": 0.19188517,
		  "fiftyTwoWeekLowChangePercent": 0.34359762,
		  "fiftyTwoWeekRange": "0.558459 - 2.497785",
		  "fiftyTwoWeekHighChange": -1.7474409,
		  "fiftyTwoWeekHighChangePercent": -0.69959617,
		  "fiftyTwoWeekLow": 0.558459,
		  "fiftyTwoWeekHigh": 2.497785,
		  "symbol": "ALGO-EUR"
		},
		{
		  "language": "fr-FR",
		  "region": "FR",
		  "quoteType": "CRYPTOCURRENCY",
		  "quoteSourceName": "CoinMarketCap",
		  "triggerable": false,
		  "customPriceAlertConfidence": "LOW",
		  "currency": "EUR",
		  "marketState": "REGULAR",
		  "exchange": "CCC",
		  "shortName": "Dogecoin EUR",
		  "messageBoardId": "finmb_DOGE_CCC_lang_fr",
		  "exchangeTimezoneName": "UTC",
		  "exchangeTimezoneShortName": "UTC",
		  "gmtOffSetMilliseconds": 0,
		  "market": "ccc_market",
		  "esgPopulated": false,
		  "firstTradeDateMilliseconds": 1510358400000,
		  "priceHint": 6,
		  "circulatingSupply": 132670767104,
		  "lastMarket": "CoinMarketCap",
		  "volume24Hr": 603982528,
		  "volumeAllCurrencies": 603982528,
		  "fromCurrency": "DOGE",
		  "toCurrency": "EUR=X",
		  "fiftyDayAverage": 0.12946571,
		  "fiftyDayAverageChange": -0.010223657,
		  "fiftyDayAverageChangePercent": -0.07896806,
		  "twoHundredDayAverage": 0.18273024,
		  "twoHundredDayAverageChange": -0.063488185,
		  "twoHundredDayAverageChangePercent": -0.34744212,
		  "marketCap": 15819934720,
		  "sourceInterval": 15,
		  "exchangeDataDelayedBy": 0,
		  "tradeable": false,
		  "regularMarketChange": -0.0013814345,
		  "regularMarketChangePercent": -1.1452425,
		  "regularMarketTime": 1646291894,
		  "regularMarketPrice": 0.11924206,
		  "regularMarketDayHigh": 0.12031095,
		  "regularMarketDayRange": "0.11823941 - 0.12031095",
		  "regularMarketDayLow": 0.11823941,
		  "regularMarketVolume": 603982528,
		  "regularMarketPreviousClose": 0.12006761,
		  "fullExchangeName": "CCC",
		  "regularMarketOpen": 0.12006761,
		  "averageDailyVolume3Month": 943651179,
		  "averageDailyVolume10Day": 740768158,
		  "startDate": 1387065600,
		  "coinImageUrl": "https://s.yimg.com/uc/fin/img/reports-thumbnails/74.png",
		  "fiftyTwoWeekLowChange": 0.07925106,
		  "fiftyTwoWeekLowChangePercent": 1.9817224,
		  "fiftyTwoWeekRange": "0.039991 - 0.60639",
		  "fiftyTwoWeekHighChange": -0.48714793,
		  "fiftyTwoWeekHighChangePercent": -0.8033575,
		  "fiftyTwoWeekLow": 0.039991,
		  "fiftyTwoWeekHigh": 0.60639,
		  "symbol": "DOGE-EUR"
		},
		{
		  "language": "fr-FR",
		  "region": "FR",
		  "quoteType": "EQUITY",
		  "quoteSourceName": "Delayed Quote",
		  "triggerable": false,
		  "customPriceAlertConfidence": "LOW",
		  "currency": "EUR",
		  "marketState": "PREPRE",
		  "exchange": "PAR",
		  "shortName": "VIVENDI SE",
		  "longName": "Vivendi SE",
		  "messageBoardId": "finmb_120343_lang_fr",
		  "exchangeTimezoneName": "Europe/Paris",
		  "exchangeTimezoneShortName": "CET",
		  "gmtOffSetMilliseconds": 3600000,
		  "market": "fr_market",
		  "esgPopulated": false,
		  "firstTradeDateMilliseconds": 946886400000,
		  "priceHint": 2,
		  "sharesOutstanding": 1045400000,
		  "bookValue": 16.707,
		  "fiftyDayAverage": 11.6782,
		  "fiftyDayAverageChange": -0.5581999,
		  "fiftyDayAverageChangePercent": -0.047798455,
		  "twoHundredDayAverage": 18.241282,
		  "twoHundredDayAverageChange": -7.1212816,
		  "twoHundredDayAverageChangePercent": -0.3903937,
		  "marketCap": 11624848384,
		  "forwardPE": 15.444444,
		  "priceToBook": 0.6655892,
		  "sourceInterval": 15,
		  "exchangeDataDelayedBy": 0,
		  "prevName": "Vivendi SA",
		  "nameChangeDate": "2022-03-02",
		  "averageAnalystRating": "2.1 - Buy",
		  "tradeable": false,
		  "regularMarketChange": 0.055000305,
		  "regularMarketChangePercent": 0.49706557,
		  "regularMarketTime": 1646238928,
		  "regularMarketPrice": 11.12,
		  "regularMarketDayHigh": 11.165,
		  "regularMarketDayRange": "10.95 - 11.165",
		  "regularMarketDayLow": 10.95,
		  "regularMarketVolume": 3412554,
		  "regularMarketPreviousClose": 11.065,
		  "bid": 0,
		  "ask": 0,
		  "bidSize": 0,
		  "askSize": 0,
		  "fullExchangeName": "Paris",
		  "financialCurrency": "EUR",
		  "regularMarketOpen": 10.99,
		  "averageDailyVolume3Month": 2955256,
		  "averageDailyVolume10Day": 3397892,
		  "fiftyTwoWeekLowChange": 1.1199999,
		  "fiftyTwoWeekLowChangePercent": 0.11199999,
		  "fiftyTwoWeekRange": "10.0 - 33.48",
		  "fiftyTwoWeekHighChange": -22.36,
		  "fiftyTwoWeekHighChangePercent": -0.66786146,
		  "fiftyTwoWeekLow": 10,
		  "fiftyTwoWeekHigh": 33.48,
		  "earningsTimestamp": 1646838000,
		  "earningsTimestampStart": 1646838000,
		  "earningsTimestampEnd": 1646838000,
		  "trailingAnnualDividendRate": 0.6,
		  "trailingPE": 10.52034,
		  "trailingAnnualDividendYield": 0.05422504,
		  "epsTrailingTwelveMonths": 1.057,
		  "epsForward": 0.72,
		  "symbol": "VIV.PA"
		}
	  ],
	  "error": null
	}
  }`

func FinancialAPICall(url string, key string, method string, action string) ([]byte, error) {
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
	if err != nil {
		clog.Error("APICall", method, "%s", err)
		return nil, err
	}
	req.Header.Set("user-agent", "GoLedMatrix Agent")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", key)

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

func (S *Scenario) ComposeMessage(interface{}) string {
	return ""
}

func (S *Stock) drawLine(startX int) int {
	var mess string
	var lineLength int = 0

	for _, symbol := range S.req.Data.Result {
		S.ctx.SetHexColor("#FFFFFF")
		if symbol.QuoteType == "CRYPTOCURRENCY" {
			mess = fmt.Sprintf("%s:", symbol.FromCurrency)
		} else {
			mess = fmt.Sprintf("%s:", symbol.Symbol)
		}
		S.ctx.DrawString(mess, float64(startX+lineLength), float64(S.sprite.Pos.Y))
		lineLength += len(mess) * 5

		if symbol.RegularMarketChange < 0 {
			S.ctx.SetHexColor("#FF0000")
		} else {
			S.ctx.SetHexColor("#00FF00")
		}
		mess = fmt.Sprintf("%.3f ", symbol.RegularMarketPrice)
		S.ctx.DrawString(mess, float64(startX+lineLength), float64(S.sprite.Pos.Y))
		lineLength += len(mess) * 5
	}
	return lineLength
}

func (S *Stock) DisplaySprite(param interface{}) {
	// var this *rgbmatrix.Sprite = param.(*rgbmatrix.Sprite)
	S.sprite.Size.X = S.drawLine(S.sprite.Pos.X)
	if S.sprite.Pos.X+S.sprite.Size.X < S.sprite.ScreenSize.X {
		S.drawLine(S.sprite.Pos.X + S.sprite.Size.X)
	}
}

func (S *Scenario) Business() {
	ticker := time.NewTicker(time.Minute * time.Duration(S.conf.API.QuoteInterval))
	defer func() {
		ticker.Stop()
	}()

	stock := Stock{}

	size := S.tk.Canvas.Bounds().Max
	strHeight := 8
	stock.ctx = gg.NewContext(size.X, size.Y)
	stock.ctx.SetFontFace(bitmapfont.Gothic10r)

	stock.sprite = &rgbmatrix.Sprite{
		ID:         1,
		ScreenSize: size,
		Size:       image.Point{0, strHeight},
		Pos:        image.Point{5, strHeight},
		Style:      rgbmatrix.Restart,
		DirX:       -1,
		DirY:       1,
		Draw:       stock.DisplaySprite,
	}

	body, _ := FinancialAPICall(S.conf.API.QuoteURL, S.conf.API.QuoteKey, "GET", S.conf.API.QuoteSymbols)
	json.Unmarshal(body, &stock.req)

	for {
		select {
		case <-ticker.C:
			body, _ := FinancialAPICall(S.conf.API.QuoteURL, S.conf.API.QuoteKey, "GET", S.conf.API.QuoteSymbols)
			json.Unmarshal(body, &stock.req)
		default:
			stock.ctx.SetHexColor("#000000")
			stock.ctx.Clear()
			stock.sprite.Move()
			S.tk.PlayImage(stock.ctx.Image(), time.Millisecond*50)
		}
	}
}
