package scenario

import (
	"encoding/json"
	"fmt"
	"image"
	"time"

	"github.com/Djoulzy/GoLedMatrix/clog"
	"github.com/Djoulzy/GoLedMatrix/rgbmatrix"

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
	ctx    *gg.Context
	sprite *rgbmatrix.Sprite
	req    StockResponse
	// message string
	// active  bool
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

type Clock struct {
	ctx    *gg.Context
	sprite *rgbmatrix.Sprite
}

var WeatherDateFormat string = "2006-01-02T15:04:05-0700"
var jourDeLaSemaine [7]string = [7]string{"Dim", "Lun", "Mar", "Mer", "Jeu", "Ven", "Sam"}
var jourDeLaSemaineL [7]string = [7]string{"Dimanche", "Lundi", "Mardi", "Mercredi", "Jeudi", "Vendredi", "Samedi"}

// var temp_meteo string = `{
//     "city": {
//         "insee": "35238",
//         "cp": 35000,
//         "name": "Rennes",
//         "latitude": 48.112,
//         "longitude": -1.6819,
//         "altitude": 38
//     },
//     "update": "2020-10-29T12:40:08+0100",
//     "forecast": [
//         {
//             "insee": "35238",
//             "cp": 35000,
//             "latitude": 48.112,
//             "longitude": -1.6819,
//             "day": 3,
//             "datetime": "2020-11-01T01:00:00+0100",
//             "wind10m": 30,
//             "gust10m": 41,
//             "dirwind10m": 210,
//             "rr10": 15.5,
//             "rr1": 21,
//             "probarain": 90,
//             "weather": 11,
//             "tmin": 15,
//             "tmax": 18,
//             "sun_hours": 0,
//             "etp": 1,
//             "probafrost": 0,
//             "probafog": 0,
//             "probawind70": 20,
//             "probawind100": 0,
//             "gustx": 61
//         },
//         {
//             "insee": "35238",
//             "cp": 35000,
//             "latitude": 48.112,
//             "longitude": -1.6819,
//             "day": 4,
//             "datetime": "2020-11-02T01:00:00+0100",
//             "wind10m": 30,
//             "gust10m": 56,
//             "dirwind10m": 207,
//             "rr10": 10.2,
//             "rr1": 27,
//             "probarain": 80,
//             "weather": 211,
//             "tmin": 9,
//             "tmax": 18,
//             "sun_hours": 0,
//             "etp": 1,
//             "probafrost": 0,
//             "probafog": 0,
//             "probawind70": 30,
//             "probawind100": 0,
//             "gustx": 86
//         },
//         {
//             "insee": "35238",
//             "cp": 35000,
//             "latitude": 48.112,
//             "longitude": -1.6819,
//             "day": 5,
//             "datetime": "2020-11-03T01:00:00+0100",
//             "wind10m": 15,
//             "gust10m": 30,
//             "dirwind10m": 216,
//             "rr10": 2.8,
//             "rr1": 12.2,
//             "probarain": 60,
//             "weather": 41,
//             "tmin": 6,
//             "tmax": 14,
//             "sun_hours": 6,
//             "etp": 1,
//             "probafrost": 10,
//             "probafog": 0,
//             "probawind70": 0,
//             "probawind100": 0,
//             "gustx": 45
//         },
//         {
//             "insee": "35238",
//             "cp": 35000,
//             "latitude": 48.112,
//             "longitude": -1.6819,
//             "day": 6,
//             "datetime": "2020-11-04T01:00:00+0100",
//             "wind10m": 15,
//             "gust10m": 26,
//             "dirwind10m": 49,
//             "rr10": 0.4,
//             "rr1": 2.2,
//             "probarain": 60,
//             "weather": 40,
//             "tmin": 5,
//             "tmax": 13,
//             "sun_hours": 6,
//             "etp": 1,
//             "probafrost": 10,
//             "probafog": 10,
//             "probawind70": 0,
//             "probawind100": 0,
//             "gustx": 36
//         },
//         {
//             "insee": "35238",
//             "cp": 35000,
//             "latitude": 48.112,
//             "longitude": -1.6819,
//             "day": 7,
//             "datetime": "2020-11-05T01:00:00+0100",
//             "wind10m": 20,
//             "gust10m": 32,
//             "dirwind10m": 71,
//             "rr10": 0,
//             "rr1": 0,
//             "probarain": 20,
//             "weather": 3,
//             "tmin": 5,
//             "tmax": 14,
//             "sun_hours": 5,
//             "etp": 1,
//             "probafrost": 10,
//             "probafog": 0,
//             "probawind70": 10,
//             "probawind100": 0,
//             "gustx": 32
//         },
//         {
//             "insee": "35238",
//             "cp": 35000,
//             "latitude": 48.112,
//             "longitude": -1.6819,
//             "day": 8,
//             "datetime": "2020-11-06T01:00:00+0100",
//             "wind10m": 15,
//             "gust10m": 30,
//             "dirwind10m": 88,
//             "rr10": 0,
//             "rr1": 0,
//             "probarain": 40,
//             "weather": 3,
//             "tmin": 7,
//             "tmax": 17,
//             "sun_hours": 5,
//             "etp": 1,
//             "probafrost": 10,
//             "probafog": 0,
//             "probawind70": 0,
//             "probawind100": 0,
//             "gustx": 30
//         },
//         {
//             "insee": "35238",
//             "cp": 35000,
//             "latitude": 48.112,
//             "longitude": -1.6819,
//             "day": 9,
//             "datetime": "2020-11-07T01:00:00+0100",
//             "wind10m": 15,
//             "gust10m": 27,
//             "dirwind10m": 92,
//             "rr10": 2.4,
//             "rr1": 5.2,
//             "probarain": 60,
//             "weather": 41,
//             "tmin": 8,
//             "tmax": 17,
//             "sun_hours": 4,
//             "etp": 1,
//             "probafrost": 0,
//             "probafog": 0,
//             "probawind70": 0,
//             "probawind100": 0,
//             "gustx": 38
//         },
//         {
//             "insee": "35238",
//             "cp": 35000,
//             "latitude": 48.112,
//             "longitude": -1.6819,
//             "day": 10,
//             "datetime": "2020-11-08T01:00:00+0100",
//             "wind10m": 15,
//             "gust10m": 24,
//             "dirwind10m": 110,
//             "rr10": 6.4,
//             "rr1": 8.4,
//             "probarain": 60,
//             "weather": 40,
//             "tmin": 8,
//             "tmax": 16,
//             "sun_hours": 3,
//             "etp": 1,
//             "probafrost": 0,
//             "probafog": 0,
//             "probawind70": 0,
//             "probawind100": 0,
//             "gustx": 35
//         },
//         {
//             "insee": "35238",
//             "cp": 35000,
//             "latitude": 48.112,
//             "longitude": -1.6819,
//             "day": 11,
//             "datetime": "2020-11-09T01:00:00+0100",
//             "wind10m": 15,
//             "gust10m": 22,
//             "dirwind10m": 140,
//             "rr10": 1.6,
//             "rr1": 7,
//             "probarain": 60,
//             "weather": 41,
//             "tmin": 7,
//             "tmax": 15,
//             "sun_hours": 3,
//             "etp": 1,
//             "probafrost": 0,
//             "probafog": 0,
//             "probawind70": 0,
//             "probawind100": 0,
//             "gustx": 37
//         },
//         {
//             "insee": "35238",
//             "cp": 35000,
//             "latitude": 48.112,
//             "longitude": -1.6819,
//             "day": 12,
//             "datetime": "2020-11-10T01:00:00+0100",
//             "wind10m": 10,
//             "gust10m": 21,
//             "dirwind10m": 153,
//             "rr10": 2.2,
//             "rr1": 6.5,
//             "probarain": 60,
//             "weather": 41,
//             "tmin": 6,
//             "tmax": 15,
//             "sun_hours": 4,
//             "etp": 1,
//             "probafrost": 10,
//             "probafog": 30,
//             "probawind70": 0,
//             "probawind100": 0,
//             "gustx": 31
//         },
//         {
//             "insee": "35238",
//             "cp": 35000,
//             "latitude": 48.112,
//             "longitude": -1.6819,
//             "day": 13,
//             "datetime": "2020-11-11T01:00:00+0100",
//             "wind10m": 15,
//             "gust10m": 21,
//             "dirwind10m": 201,
//             "rr10": 3.5,
//             "rr1": 10,
//             "probarain": 60,
//             "weather": 41,
//             "tmin": 6,
//             "tmax": 15,
//             "sun_hours": 4,
//             "etp": 1,
//             "probafrost": 10,
//             "probafog": 0,
//             "probawind70": 10,
//             "probawind100": 0,
//             "gustx": 31
//         }
//     ]
// }`

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

func (W *Weather) drawWidget(day WeatherForecast, x, y float64) {
	t, err := time.Parse(WeatherDateFormat, day.Date)
	if err != nil {
		clog.Error("Business", "drawWidget", "%s", err)
	}

	W.ctx.SetHexColor("#FFFFFF")
	W.ctx.DrawString(fmt.Sprintf("%s.%02d", jourDeLaSemaine[t.Weekday()], t.Day()), x, y)
	W.ctx.DrawString(fmt.Sprintf(" %d°C", day.TMax), x, y+8)
	W.ctx.DrawString(fmt.Sprintf(" %d°C", day.TMin), x, y+16)
}

func (W *Weather) DisplaySprite(param interface{}) {
	XstartLine1 := 4
	XstartLine2 := 38

	// var this *rgbmatrix.Sprite = param.(*rgbmatrix.Sprite)
	W.ctx.SetHexColor("#0000FF")
	W.ctx.DrawRectangle(float64(W.sprite.Pos.X), float64(W.sprite.Pos.Y), float64(W.sprite.Size.X), float64(W.sprite.Size.Y))
	W.ctx.Stroke()

	W.ctx.SetHexColor("#FFFFFF")
	day := int(time.Now().Weekday())
	toWe := 6 - day
	// clog.Trace("Business", "drawWidget", "Day:%d - ToWE:%d", day, toWe)

	W.drawWidget(W.req.Forecast[0], float64(W.sprite.Pos.X+XstartLine1), float64(W.sprite.Pos.Y+8))
	W.drawWidget(W.req.Forecast[1], float64(W.sprite.Pos.X+XstartLine2), float64(W.sprite.Pos.Y+8))

	switch day {
	case 5:
		W.drawWidget(W.req.Forecast[toWe+1], float64(W.sprite.Pos.X+XstartLine1), float64(W.sprite.Pos.Y+40))
		W.drawWidget(W.req.Forecast[toWe+7], float64(W.sprite.Pos.X+XstartLine2), float64(W.sprite.Pos.Y+40))
	case 6:
		W.drawWidget(W.req.Forecast[toWe+7], float64(W.sprite.Pos.X+XstartLine1), float64(W.sprite.Pos.Y+40))
		W.drawWidget(W.req.Forecast[toWe+8], float64(W.sprite.Pos.X+XstartLine2), float64(W.sprite.Pos.Y+40))
	default:
		W.drawWidget(W.req.Forecast[toWe], float64(W.sprite.Pos.X+XstartLine1), float64(W.sprite.Pos.Y+40))
		W.drawWidget(W.req.Forecast[toWe+1], float64(W.sprite.Pos.X+XstartLine2), float64(W.sprite.Pos.Y+40))
	}
}

func (C *Clock) DisplaySprite(param interface{}) {
	C.ctx.SetHexColor("#FF0000")
	C.ctx.DrawRectangle(float64(C.sprite.Pos.X), float64(C.sprite.Pos.Y), float64(C.sprite.Size.X), float64(C.sprite.Size.Y))
	C.ctx.Stroke()

	actual := time.Now()
	heure := actual.Format("15:04:05")
	day := int(time.Now().Weekday())
	C.ctx.DrawString(heure, float64(C.sprite.Pos.X+2), float64(C.sprite.Pos.Y+14))
	C.ctx.DrawString(jourDeLaSemaineL[day], float64(C.sprite.Pos.X+2), float64(C.sprite.Pos.Y+30))
}

func (S *Scenario) Business() {
	var body []byte

	tickerQuote := time.NewTicker(time.Minute * time.Duration(S.conf.QuoteAPI.QuoteInterval))
	defer func() {
		tickerQuote.Stop()
	}()
	tickerWeather := time.NewTicker(time.Minute * time.Duration(S.conf.WeatherAPI.WeatherInterval))
	defer func() {
		tickerQuote.Stop()
	}()

	Actions := true
	Meteo := true

	stock := Stock{}
	weather := Weather{}
	clock := Clock{}

	ctx := gg.NewContext(128, 128)

	size := S.tk.Canvas.Bounds().Max
	strHeight := 8
	stock.ctx = ctx
	weather.ctx = ctx
	clock.ctx = ctx

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

	weather.sprite = &rgbmatrix.Sprite{
		ID:         2,
		ScreenSize: size,
		Size:       image.Point{64, 64},
		Pos:        image.Point{0, 32},
		Style:      rgbmatrix.Idle,
		DirX:       1,
		DirY:       1,
		Draw:       weather.DisplaySprite,
	}

	clock.sprite = &rgbmatrix.Sprite{
		ID:         2,
		ScreenSize: size,
		Size:       image.Point{64, 64},
		Pos:        image.Point{64, 32},
		Style:      rgbmatrix.Idle,
		DirX:       1,
		DirY:       1,
		Draw:       clock.DisplaySprite,
	}

	Quotekey := ApiToken{
		Value: S.conf.QuoteAPI.QuoteKey,
		Auth:  XAPIKEY,
	}
	if Actions {
		body, _ = APICall(S.conf.QuoteAPI.QuoteURL, Quotekey, "GET", S.conf.QuoteAPI.QuoteSymbols)
		json.Unmarshal(body, &stock.req)
	}

	WeatherKey := ApiToken{
		Value: S.conf.WeatherAPI.WeatherKey,
		Auth:  PARAM,
	}
	if Meteo {
		body, _ = APICall(S.conf.WeatherAPI.WeatherURL, WeatherKey, "GET", S.conf.WeatherAPI.WeatherRoute+"?insee="+S.conf.WeatherAPI.WeatherINSEE)
		json.Unmarshal(body, &weather.req)
	}

	for {
		select {
		case <-tickerQuote.C:
			if Actions {
				body, _ = APICall(S.conf.QuoteAPI.QuoteURL, Quotekey, "GET", S.conf.QuoteAPI.QuoteSymbols)
				json.Unmarshal(body, &stock.req)
			}
		case <-tickerWeather.C:
			if Meteo {
				body, _ = APICall(S.conf.WeatherAPI.WeatherURL, WeatherKey, "GET", S.conf.WeatherAPI.WeatherRoute+"?insee="+S.conf.WeatherAPI.WeatherINSEE)
				json.Unmarshal(body, &weather.req)
			}
		default:
			ctx.SetHexColor("#000000")
			ctx.Clear()
			ctx.SetFontFace(bitmapfont.Gothic10r)
			if Actions {
				stock.sprite.Move()
			}
			if Meteo {
				weather.sprite.Move()
			}
			ctx.LoadFontFace(S.conf.DefaultConf.FontDir+"fixed/Pixel_NES.otf", 12)
			clock.sprite.Move()
			S.tk.PlayImage(stock.ctx.Image(), time.Millisecond*50)
		}
	}
}
