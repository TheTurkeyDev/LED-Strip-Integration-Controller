package main

import (
	"time"

	ws281x "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	IDLE            int = 0
	SOLID           int = 1
	RAINBOW         int = 2
	BLOCK_COLOR     int = 3
	ALTERNATE_COLOR int = 4
	POLICE          int = 5
)

const (
	ledCounts = 240
	gpioPin   = 18
	freq      = 800000
	sleepTime = 20
)

type (
	Color struct {
		r, g, b uint8
	}

	ws struct {
		ws2811     *ws281x.WS2811
		colors     []Color
		brightness int
		display    int
	}
)

func (ws *ws) init() error {
	err := ws.ws2811.Init()
	if err != nil {
		return err
	}

	return nil
}

func (ws *ws) renderAll() error {
	for i := 0; i < len(ws.ws2811.Leds(0)); i++ {
		ws.ws2811.Leds(0)[i] = rgbToColor(ws.colors[0].r, ws.colors[0].g, ws.colors[0].b)
	}

	if err := ws.ws2811.Render(); err != nil {
		return err
	}

	return nil
}

func (ws *ws) renderAllHex(color Color) {
	intColor := rgbToColor(color.r, color.g, color.b)
	for i := 0; i < len(ws.ws2811.Leds(0)); i++ {
		ws.ws2811.Leds(0)[i] = intColor
	}

	if err := ws.ws2811.Render(); err != nil {
		println("Render Error!")
	}
}

func (ws *ws) close() {
	ws.ws2811.Fini()
}

func rgbToColor(r uint8, g uint8, b uint8) uint32 {
	return uint32(uint32(r)<<16 | uint32(g)<<8 | uint32(b))
}

func (ws *ws) rainbowRGB() {
	ws.colors[0].r = 255
	err := ws.renderAll()
	if err != nil {
		println("Render Error Rainbow!")
	}
}

func (ws *ws) idleStart() {
	ws.ws2811.Leds(0)[0] = rgbToColor(255, 255, 255)
}

func (ws *ws) idle() {
	for i := 0; i < len(ws.ws2811.Leds(0)); i++ {
		index := i
		if index == 0 {
			index = len(ws.ws2811.Leds(0))
		}
		index -= 1
		ws.ws2811.Leds(0)[i] = ws.ws2811.Leds(0)[index]
	}
	err := ws.renderAll()
	if err != nil {
		println("Render Error Idle!")
	}
}

func (ws *ws) endsToMiddle() error {
	return nil
}

func (ws *ws) pulseColor() error {
	return nil
}

func lightsLoop(ws ws) {
	for {
		switch ws.display {
		case 0:
			ws.idle()
		case 1:
			ws.renderAllHex(ws.colors[0])
		case 2:
			ws.rainbowRGB()
		}
		time.Sleep(sleepTime * time.Millisecond)
	}
}

func NewLEDLights() *ws {
	opt := ws281x.DefaultOptions
	opt.Channels[0].Brightness = 128
	opt.Channels[0].LedCount = ledCounts
	opt.Channels[0].GpioPin = gpioPin
	opt.Frequency = freq

	ws2811, err := ws281x.MakeWS2811(&opt)
	if err != nil {
		panic(err)
	}

	ws := ws{
		ws2811: ws2811,
		colors: []Color{
			{
				r: 0,
				g: 0,
				b: 0,
			},
		},
		brightness: 128,
		display:    IDLE,
	}

	err = ws.init()
	if err != nil {
		panic(err)
	}
	defer ws.close()

	ws.idleStart()
	go lightsLoop(ws)

	return &ws
}
