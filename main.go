package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/io/i2c"
	"golang.org/x/net/websocket"
)

type Cargo struct {
	dev *i2c.Device
}

const (
	motorL        = 0
	motorR        = 1
	outputConfig0 = 2
	output0       = 8
)

func NewCargo() (c Cargo) {
	var err error

	c.dev, err = i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x22)
	if err != nil {
		panic(err)
	}
	c.Move(0, 0)

	// Setup servo on output 0
	c.setOutput(outputConfig0, 2)
	return
}

func (c *Cargo) setOutput(output, value byte) {
	c.dev.Write([]byte{output, value})
}

func (c *Cargo) Move(speedL, speedR int8) {
	c.setOutput(motorL, byte(speedL))
	c.setOutput(motorR, byte(speedR))
}

func (c *Cargo) SetServo(pos int8) {
	c.setOutput(output0, byte(pos))
}

func (c *Cargo) GetDistance() int64 {
	out, _ := exec.Command("python", "./sonar.py").Output()
	s := strings.TrimSpace(string(out))
	fmt.Println(s)
	distance, _ := strconv.ParseInt(s, 10, 32)
	return distance
}

func (c *Cargo) ExecuteProgram(program string) {
	for _, s := range strings.Split(strings.Replace(program, ";", "\n", -1), "\n") {
		line := strings.TrimSpace(s)
		splits := strings.Split(line, " ")

		switch splits[0] {
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

		d, _ := strconv.ParseInt(splits[1], 10, 32)

		time.Sleep(time.Duration(d) * 10 * time.Millisecond)
		c.Move(0, 0)
	}
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
			c.SetServo(90)
		case "reverse":
			c.Move(-127, -127)
			c.SetServo(-90)
		case "left":
			c.Move(127, -127)
			c.SetServo(-45)
		case "right":
			c.Move(-127, 127)
			c.SetServo(45)
		case "stop":
			c.Move(0, 0)
			c.SetServo(0)
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
