package scenario

import (
	"GoLedMatrix/clog"
	"GoLedMatrix/rgbmatrix"
	"encoding/json"
	"fmt"
	"image"
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

type WeatherForecast struct {
	Day       int    `mapstructure:"day" json:"day"`
	Date      string `mapstructure:"datetime" json:"datetime"`
	ProbaRain int    `mapstructure:"probarain" json:"probarain"`
	TMin      int    `mapstructure:"tmin" json:"tmin"`
	TMax      int    `mapstructure:"tmax" json:"tmax"`
}

type WeatherCity struct {
	INSEE string `mapstructure:"insee" json:"insee"`
	CP    int    `mapstructure:"cp" json:"cp"`
	Name  string `mapstructure:"name" json:"name"`
}

type WeatherResponse struct {
	City     WeatherCity
	Update   string
	Forecast []WeatherForecast
}

type Weather struct {
	ctx    *gg.Context
	sprite *rgbmatrix.Sprite
	req    WeatherResponse
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
	tickerQuote := time.NewTicker(time.Minute * time.Duration(S.conf.QuoteAPI.QuoteInterval))
	defer func() {
		tickerQuote.Stop()
	}()

	stock := Stock{}
	weather := Weather{}

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

	// Quotekey := ApiToken{
	// 	Value: S.conf.QuoteAPI.QuoteKey,
	// 	Auth:  XAPIKEY,
	// }

	// body, _ := APICall(S.conf.QuoteAPI.QuoteURL, Quotekey, "GET", S.conf.QuoteAPI.QuoteSymbols)
	// json.Unmarshal(body, &stock.req)

	WeatherKey := ApiToken{
		Value: S.conf.WeatherAPI.WeatherKey,
		Auth:  PARAM,
	}

	body, _ := APICall(S.conf.WeatherAPI.WeatherURL, WeatherKey, "GET", S.conf.WeatherAPI.WeatherRoute+"?insee="+S.conf.WeatherAPI.WeatherINSEE)
	json.Unmarshal(body, &weather.req)
	clog.Test("Scenario", "Business", "%v", weather.req)

	for {
		select {
		case <-tickerQuote.C:
			// body, _ = APICall(S.conf.QuoteAPI.QuoteURL, Quotekey, "GET", S.conf.QuoteAPI.QuoteSymbols)
			// json.Unmarshal(body, &stock.req)
		default:
			stock.ctx.SetHexColor("#000000")
			stock.ctx.Clear()
			stock.sprite.Move()
			S.tk.PlayImage(stock.ctx.Image(), time.Millisecond*50)
		}
	}
}
