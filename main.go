package main

import (
	"machine"
	"time"
)

// // Define global I2C pin constants needed by the machine package
// const (
// 	SDA_PIN = machine.GPIO26
// 	SCL_PIN = machine.GPIO22
// )

func main() {
	led := machine.Pin(26)
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for {
		led.High()
		time.Sleep(time.Second / 2)

		led.Low()
		time.Sleep(time.Second / 2)
	}
}
