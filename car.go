package main

import (
	"fmt"
	"github.com/mrmorphic/hwio"
	"golang.org/x/net/websocket"
	"net/http"
)

const (
	in1 = iota
	in2
	in3
	in4
	TotalPins
)

type Car struct {
	pins [TotalPins]hwio.Pin
}

func (c *Car) Motor(v1, v2, v3, v4 int) {
	hwio.DigitalWrite(c.pins[in1], v1)
	hwio.DigitalWrite(c.pins[in2], v2)
	hwio.DigitalWrite(c.pins[in3], v3)
	hwio.DigitalWrite(c.pins[in4], v4)
}

func (c *Car) Forward() {
	c.Motor(hwio.HIGH, hwio.LOW, hwio.LOW, hwio.HIGH)
}

func (c *Car) Reverse() {
	c.Motor(hwio.LOW, hwio.HIGH, hwio.HIGH, hwio.LOW)
}

func (c *Car) Left() {
	c.Motor(hwio.LOW, hwio.HIGH, hwio.LOW, hwio.HIGH)
}

func (c *Car) Right() {
	c.Motor(hwio.HIGH, hwio.LOW, hwio.HIGH, hwio.LOW)
}

func (c *Car) Stop() {
	c.Motor(hwio.LOW, hwio.LOW, hwio.LOW, hwio.LOW)
}

func (c *Car) Setup() {
	c.pins[in1], _ = hwio.GetPin("gpio24")
	c.pins[in2], _ = hwio.GetPin("gpio17")
	c.pins[in3], _ = hwio.GetPin("gpio23")
	c.pins[in4], _ = hwio.GetPin("gpio22")
	_ = hwio.PinMode(c.pins[in1], hwio.OUTPUT)
	_ = hwio.PinMode(c.pins[in2], hwio.OUTPUT)
	_ = hwio.PinMode(c.pins[in3], hwio.OUTPUT)
	_ = hwio.PinMode(c.pins[in4], hwio.OUTPUT)
}

func socketHandler(ws *websocket.Conn) {
	for {
		receivedtext := make([]byte, 100)
		n, _ := ws.Read(receivedtext)
		s := string(receivedtext[:n])
		fmt.Printf("Received: %s\n", s)

		switch s {
		case "forward":
			c.Forward()
		case "reverse":
			c.Reverse()
		case "left":
			c.Left()
		case "right":
			c.Right()
		case "stop":
			c.Stop()
		default:
			c.Stop()
		}
	}
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./control.html")
}

func setupWebserver() {
	http.Handle("/control", websocket.Handler(socketHandler))
	http.HandleFunc("/", webHandler)
	http.ListenAndServe(":9090", nil)
}

var c Car

func main() {
	c.Setup()

	setupWebserver()
	for {
	}
}
