package main

import (
	"fmt"
	"github.com/mrmorphic/hwio"
	"golang.org/x/net/websocket"
	"net/http"
	"os"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func forward() {
	fmt.Println("forward")
	hwio.DigitalWrite(in1, hwio.HIGH)
	hwio.DigitalWrite(in2, hwio.LOW)
	hwio.DigitalWrite(in3, hwio.LOW)
	hwio.DigitalWrite(in4, hwio.HIGH)
}

func reverse() {
	fmt.Println("reverse")
	hwio.DigitalWrite(in1, hwio.LOW)
	hwio.DigitalWrite(in2, hwio.HIGH)
	hwio.DigitalWrite(in3, hwio.HIGH)
	hwio.DigitalWrite(in4, hwio.LOW)
}

func left() {
	fmt.Println("left")
	hwio.DigitalWrite(in1, hwio.LOW)
	hwio.DigitalWrite(in2, hwio.HIGH)
	hwio.DigitalWrite(in3, hwio.LOW)
	hwio.DigitalWrite(in4, hwio.HIGH)
}

func right() {
	fmt.Println("right")
	hwio.DigitalWrite(in1, hwio.HIGH)
	hwio.DigitalWrite(in2, hwio.LOW)
	hwio.DigitalWrite(in3, hwio.HIGH)
	hwio.DigitalWrite(in4, hwio.LOW)
}

func stop() {
	fmt.Println("stop")
	hwio.DigitalWrite(in1, hwio.LOW)
	hwio.DigitalWrite(in2, hwio.LOW)
	hwio.DigitalWrite(in3, hwio.LOW)
	hwio.DigitalWrite(in4, hwio.LOW)
}

var in1 hwio.Pin
var in2 hwio.Pin
var in3 hwio.Pin
var in4 hwio.Pin

func setupControls() {
	in1, _ = hwio.GetPin("gpio4")
	in2, _ = hwio.GetPin("gpio17")
	in3, _ = hwio.GetPin("gpio18")
	in4, _ = hwio.GetPin("gpio22")
	_ = hwio.PinMode(in1, hwio.OUTPUT)
	_ = hwio.PinMode(in2, hwio.OUTPUT)
	_ = hwio.PinMode(in3, hwio.OUTPUT)
	_ = hwio.PinMode(in4, hwio.OUTPUT)
}

func socketHandler(ws *websocket.Conn) {
	for {
		receivedtext := make([]byte, 100)
		n, _ := ws.Read(receivedtext)
		s := string(receivedtext[:n])
		fmt.Printf("Received: %s\n", s)

		switch s {
		case "forward":
			forward()
		case "reverse":
			reverse()
		case "left":
			left()
		case "right":
			right()
		case "stop":
			stop()
		default:
			stop()
		}
	}
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "control.html")
}

func setupWebserver() {
	http.Handle("/control", websocket.Handler(socketHandler))
	http.HandleFunc("/", webHandler)
	http.ListenAndServe(":9090", nil)
}

func main() {
	setupControls()
	setupWebserver()
	for {
	}
}
