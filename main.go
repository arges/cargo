package main

import (
	"fmt"
	"net/http"

	"golang.org/x/exp/io/i2c"
	"golang.org/x/net/websocket"
)

type Cargo struct {
	dev *i2c.Device
}

const (
	motorL = 0
	motorR = 1
)

func NewCargo() (c Cargo) {
	var err error

	c.dev, err = i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x22)
	if err != nil {
		panic(err)
	}
	c.Move(0, 0)
	return
}

func (c *Cargo) Move(speedL, speedR int8) {
	c.dev.Write([]byte{motorL, byte(speedL)})
	c.dev.Write([]byte{motorR, byte(speedR)})
}

func (c *Cargo) SocketHandler(ws *websocket.Conn) {
	for {
		receivedtext := make([]byte, 16)
		n, _ := ws.Read(receivedtext)
		s := string(receivedtext[:n])
		fmt.Printf("Received: %s\n", s)

		switch s {
		case "forward":
			c.Move(127, 127)
		case "reverse":
			c.Move(-127, -127)
		case "left":
			c.Move(127, -127)
		case "right":
			c.Move(-127, 127)
		case "stop":
			c.Move(0, 0)
		default:
			c.Move(0, 0)
		}
	}
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./control.html")
}

func (c *Cargo) SetupWebserver() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", webHandler)
	mux.Handle("/control", websocket.Handler(c.SocketHandler))
	http.ListenAndServe(":9090", mux)
}

func main() {
	c := NewCargo()

	c.SetupWebserver()
}
