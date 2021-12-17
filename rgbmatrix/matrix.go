// +build !arm64

package rgbmatrix

import (
	"fmt"
	"image/color"

	"GoLedMatrix/clog"
	"GoLedMatrix/emulator"
)

// DefaultConfig default WS281x configuration
var DefaultConfig = HardwareConfig{
	Rows:              32,
	Cols:              64,
	ChainLength:       4,
	Parallel:          1,
	PWMBits:           11,
	PWMLSBNanoseconds: 130,
	Brightness:        100,
	ScanMode:          Progressive,
}

// HardwareConfig rgb-led-matrix configuration
type HardwareConfig struct {
	// Rows the number of rows supported by the display, so 32 or 16.
	Rows int
	// Cols the number of columns supported by the display, so 32 or 64 .
	Cols int
	// ChainLengthis the number of displays daisy-chained together
	// (output of one connected to input of next).
	ChainLength int
	// Parallel is the number of parallel chains connected to the Pi; in old Pis
	// with 26 GPIO pins, that is 1, in newer Pis with 40 interfaces pins, that
	// can also be 2 or 3. The effective number of pixels in vertical direction is
	// then thus rows * parallel.
	Parallel int
	// Set PWM bits used for output. Default is 11, but if you only deal with
	// limited comic-colors, 1 might be sufficient. Lower require less CPU and
	// increases refresh-rate.
	PWMBits int
	// Change the base time-unit for the on-time in the lowest significant bit in
	// nanoseconds.  Higher numbers provide better quality (more accurate color,
	// less ghosting), but have a negative impact on the frame rate.
	PWMLSBNanoseconds int // the DMA channel to use
	// Brightness is the initial brightness of the panel in percent. Valid range
	// is 1..100
	Brightness int
	// ScanMode progressive or interlaced
	ScanMode ScanMode // strip color layout
	// Disable the PWM hardware subsystem to create pulses. Typically, you don't
	// want to disable hardware pulsing, this is mostly for debugging and figuring
	// out if there is interference with the sound system.
	// This won't do anything if output enable is not connected to GPIO 18 in
	// non-standard wirings.
	DisableHardwarePulsing bool

	ShowRefreshRate bool
	InverseColors   bool

	// Name of GPIO mapping used
	HardwareMapping string

	// A string describing a sequence of pixel mappers that should be applied
	// to this matrix. A semicolon-separated list of pixel-mappers with optional
	// parameter.
	PixelMapperConfig  string
	PWMDitherBits      int
	LimitRefreshRateHz int
	Multiplexing       int
}

func (c *HardwareConfig) geometry() (width, height int) {
	// return c.Cols * c.ChainLength, c.Rows * c.Parallel
	return 128, 128
}

type RuntimeOptions struct {
	GpioSlowdown int // 0 = no slowdown.          Flag: --led-slowdown-gpio

	// ----------
	// If the following options are set to disabled with -1, they are not
	// even offered via the command line flags.
	// ----------

	// Thre are three possible values here
	//   -1 : don't leave choise of becoming daemon to the command line parsing.
	//        If set to -1, the --led-daemon option is not offered.
	//    0 : do not becoma a daemon, run in forgreound (default value)
	//    1 : become a daemon, run in background.
	//
	// If daemon is disabled (= -1), the user has to call
	// RGBMatrix::StartRefresh() manually once the matrix is created, to leave
	// the decision to become a daemon
	// after the call (which requires that no threads have been started yet).
	// In the other cases (off or on), the choice is already made, so the thread
	// is conveniently already started for you.
	Daemon int // -1 disabled. 0=off, 1=on. Flag: --led-daemon

	// Drop privileges from 'root' to 'daemon' once the hardware is initialized.
	// This is usually a good idea unless you need to stay on elevated privs.
	DropPrivileges int // -1 disabled. 0=off, 1=on. flag: --led-drop-privs

	// By default, the gpio is initialized for you, but if you run on a platform
	// not the Raspberry Pi, this will fail. If you don't need to access GPIO
	// e.g. you want to just create a stream output (see content-streamer.h),
	// set this to false.
	DoGpioInit bool
}

type ScanMode int8

const (
	Progressive ScanMode = 0
	Interlaced  ScanMode = 1
)

// RGBLedMatrix matrix representation for ws281x
type RGBLedMatrix struct {
	Config *HardwareConfig

	height int
	width  int
	leds   []uint32
}

// NewRGBLedMatrix returns a new matrix using the given size and config
func NewRGBLedMatrix(configHard *HardwareConfig, configRuntime *RuntimeOptions) (c Matrix, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("error creating matrix: %v", r)
			}
		}
	}()

	clog.Warn("matrix", "NewRGBLedMatrix", "Starting Emulator")
	w, h := configHard.geometry()
	return emulator.NewEmulator(w, h, emulator.DefaultPixelPitch, true), nil
}

// Initialize initialize library, must be called once before other functions are
// called.
func (c *RGBLedMatrix) Initialize() error {
	return nil
}

// Geometry returns the width and the height of the matrix
func (c *RGBLedMatrix) Geometry() (width, height int) {
	return c.width, c.height
}

// Apply set all the pixels to the values contained in leds
func (c *RGBLedMatrix) Apply(leds []color.Color) error {
	for position, l := range leds {
		c.Set(position, l)
	}

	return c.Render()
}

// Render update the display with the data from the LED buffer
func (c *RGBLedMatrix) Render() error {
	clog.Test("rgbmatrix", "RGBLedMatrix", "Render")
	w, h := c.Config.geometry()

	// C.led_matrix_swap(
	// 	c.matrix,
	// 	c.buffer,
	// 	C.int(w), C.int(h),
	// 	(*C.uint32_t)(unsafe.Pointer(&c.leds[0])),
	// )

	c.leds = make([]uint32, w*h)
	return nil
}

func (c *RGBLedMatrix) Start() {
}

// At return an Color which allows access to the LED display data as
// if it were a sequence of 24-bit RGB values.
func (c *RGBLedMatrix) At(position int) color.Color {
	return uint32ToColor(c.leds[position])
}

// Set set LED at position x,y to the provided 24-bit color value.
func (c *RGBLedMatrix) Set(position int, color color.Color) {
	c.leds[position] = uint32(colorToUint32(color))
}

// Close finalizes the ws281x interface
func (c *RGBLedMatrix) Close() error {
	return nil
}

func colorToUint32(c color.Color) uint32 {
	if c == nil {
		return 0
	}

	// A color's RGBA method returns values in the range [0, 65535]
	red, green, blue, _ := c.RGBA()
	return (red>>8)<<16 | (green>>8)<<8 | blue>>8
}

func uint32ToColor(u uint32) color.Color {
	return color.RGBA{
		uint8(u>>16) & 255,
		uint8(u>>8) & 255,
		uint8(u>>0) & 255,
		0,
	}
}
