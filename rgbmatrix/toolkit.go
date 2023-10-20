package rgbmatrix

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

// ToolKit is a convinient set of function to operate with a led of Matrix
type ToolKit struct {
	// Canvas is the Canvas wrapping the Matrix, if you want to instanciate
	// a ToolKit with a custom Canvas you can use directly the struct,
	// without calling NewToolKit
	Canvas *Canvas

	// Transform function if present is applied just before draw the image to
	// the Matrix, this is a small example:
	//	tk.Transform = func(img image.Image) *image.NRGBA {
	//		return imaging.Fill(img, 64, 96, imaging.Center, imaging.Lanczos)
	//	}
	Transform func(img image.Image) *image.NRGBA
}

type MoveStyle int

const (
	Idle MoveStyle = iota
	ExternalBounce
	InternalBounce
	Continue
	Restart
)

type Sprite struct {
	ID         int
	ScreenSize image.Point
	Size       image.Point
	Pos        image.Point
	DirX       int
	DirY       int
	Style      MoveStyle
	Draw       func(interface{})
	Text       string
	BgColor    string
	FgColor    string

	startPos image.Point
}

// NewToolKit returns a new ToolKit wrapping the given Matrix
func NewToolKit(m Matrix) *ToolKit {
	return &ToolKit{
		Canvas: NewCanvas(m),
	}
}

func (tk *ToolKit) NewSprite(width, height int, x, y int) *Sprite {
	return &Sprite{
		ScreenSize: tk.Canvas.Bounds().Max,
		Size:       image.Point{width, height},
		Pos:        image.Point{x, y},
		startPos:   image.Point{x, y},
	}
}

// PlayImage draws the given image during the given delay
func (tk *ToolKit) PlayImage(i image.Image, delay time.Duration) error {
	start := time.Now()
	defer func() { time.Sleep(delay - time.Since(start)) }()

	if tk.Transform != nil {
		i = tk.Transform(i)
	}

	draw.Draw(tk.Canvas, tk.Canvas.Bounds(), i, image.ZP, draw.Over)
	return tk.Canvas.Render()
}

type Animation interface {
	Init() chan bool
	Next() (image.Image, <-chan time.Time, error)
}

// PlayAnimation play the image during the delay returned by Next, until an err
// is returned, if io.EOF is returned, PlayAnimation finish without an error
func (tk *ToolKit) PlayAnimation(a Animation) error {
	var err error
	var i image.Image
	var n <-chan time.Time
	quit := a.Init()

	for {
		select {
		case <-quit:
			return nil
		default:
			i, n, err = a.Next()
			if err != nil {
				break
			}

			if err := tk.PlayImageUntil(i, n); err != nil {
				return err
			}
		}
	}
}

// PlayImageUntil draws the given image until is notified to stop
func (tk *ToolKit) PlayImageUntil(i image.Image, notify <-chan time.Time) error {
	defer func() {
		<-notify
	}()

	if tk.Transform != nil {
		i = tk.Transform(i)
	}

	draw.Draw(tk.Canvas, tk.Canvas.Bounds(), i, image.ZP, draw.Over)
	return tk.Canvas.Render()
}

// PlayImages draws a sequence of images during the given delays, the len of
// images should be equal to the len of delay. If loop is true the function
// loops over images until a true is sent to the returned chan
func (tk *ToolKit) PlayImages(images []image.Image, delay []time.Duration, loop int) chan bool {
	quit := make(chan bool)

	go func() {
		l := len(images)
		i := 0
		for {
			select {
			case <-quit:
				return
			default:
				tk.PlayImage(images[i], delay[i])
			}

			i++
			if i >= l {
				if loop == 0 {
					i = 0
					continue
				}

				break
			}
		}
	}()

	return quit
}

// PlayGIF reads and draw a gif file from r. It use the contained images and
// delays and loops over it, until a true is sent to the returned chan
func (tk *ToolKit) PlayGIF(r io.Reader) (chan bool, error) {
	gif, err := gif.DecodeAll(r)
	if err != nil {
		return nil, err
	}

	delay := make([]time.Duration, len(gif.Delay))
	images := make([]image.Image, len(gif.Image))
	for i, image := range gif.Image {
		images[i] = image
		delay[i] = time.Millisecond * time.Duration(gif.Delay[i]) * 10
	}

	return tk.PlayImages(images, delay, gif.LoopCount), nil
}

func (tk *ToolKit) DrawText(lines []string, x, y int, fontfile string, size, spacing float64) {
	dpi := 72
	hinting := "none"

	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	// Initialize the context
	fg, bg := image.NewUniform(color.RGBA{0xff, 0xe5, 0x1E, 0xff}), image.Black
	// ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}

	rgba := image.NewRGBA(image.Rect(0, 0, tk.Canvas.w, tk.Canvas.h))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(float64(dpi))
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	switch hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Draw the text.
	pt := freetype.Pt(x, y+int(c.PointToFixed(size)>>6))
	for _, s := range lines {
		_, err = c.DrawString(s, pt)
		if err != nil {
			log.Println(err)
			return
		}
		pt.Y += c.PointToFixed(size * spacing)
	}

	tk.PlayImage(rgba, time.Second)
}

// Close close the toolkit and the inner canvas
func (tk *ToolKit) Close() error {
	return tk.Canvas.Close()
}

func (S *Sprite) ExternalBounce() {
	if S.Size.X > S.ScreenSize.X {
		S.Pos.X += S.DirX
		if S.Pos.X+S.Size.X < S.ScreenSize.X {
			S.DirX = 1
		}
		if S.Pos.X == 0 {
			S.DirX = -1
		}
	}
	if S.Size.Y > S.ScreenSize.Y {
		S.Pos.Y += S.DirY
		if S.Pos.Y+S.Size.Y < S.ScreenSize.Y {
			S.DirY = 1
		}
		if S.Pos.Y == 0 {
			S.DirY = -1
		}
	}
}

func (S *Sprite) InternalBounce() {

	S.Pos.X += S.DirX
	if S.Pos.X+S.Size.X >= S.ScreenSize.X {
		S.DirX = -1
	}
	if S.Pos.X <= 0 {
		S.DirX = 1
	}

	S.Pos.Y += S.DirY
	if S.Pos.Y+S.Size.Y >= S.ScreenSize.Y {
		S.DirY = -1
	}
	if S.Pos.Y <= 0 {
		S.DirY = 1
	}
}

func (S *Sprite) Continue() {
}

func (S *Sprite) Restart() {
	S.Pos.X += S.DirX
	if S.Pos.X >= S.ScreenSize.X {
		S.Pos.X = S.startPos.X
	}
	if S.Pos.X+S.Size.X <= 0 {
		S.Pos.X = S.startPos.X
	}
}

func (S *Sprite) Move() {
	S.Draw(S)
	if S.Style == Idle {
		return
	}
	switch S.Style {
	case Idle:
		return
	case ExternalBounce:
		S.ExternalBounce()
	case InternalBounce:
		S.InternalBounce()
	case Continue:
		S.Continue()
	case Restart:
		S.Restart()
	}
}
