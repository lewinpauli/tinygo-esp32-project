package main

import (
	"machine"
	"time"
	// "fmt" 
)

// how to use :  screen /dev/cu.usbserial-1440 115200   

func getDistanceCm(ultrasonicPin machine.Pin) float32 {
	ultrasonicPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ultrasonicPin.Low()
	time.Sleep(2 * time.Microsecond)
	ultrasonicPin.High()
	time.Sleep(10 * time.Microsecond)
	ultrasonicPin.Low()

	ultrasonicPin.Configure(machine.PinConfig{Mode: machine.PinInput})

	var startTime, endTime int64
	timeoutStart := time.Now().UnixNano() + 25000*1000

	for !ultrasonicPin.Get() {
		if time.Now().UnixNano() > timeoutStart {
			return -1.0
		}
	}
	startTime = time.Now().UnixNano()

	timeoutEnd := startTime + 25000*1000

	for ultrasonicPin.Get() {
		if time.Now().UnixNano() > timeoutEnd {
			return -2.0
		}
	}
	endTime = time.Now().UnixNano()

	pulseDurationNano := endTime - startTime
	if pulseDurationNano <= 0 {
		return -3.0
	}


	distanceCm := float32(pulseDurationNano) * 0.0000343 / 2.0

	return distanceCm
}

func main() {
	// Serielle Ausgabe konfigurieren (wichtig für println)
	uart := machine.UART0
	uart.Configure(machine.UARTConfig{BaudRate: 115200})
	// Kurze Verzögerung, damit der serielle Monitor Zeit hat, sich zu verbinden
	time.Sleep(2 * time.Second) 
	println("ESP32 Ultraschall- und Taster-Demo gestartet...")

	// Gelbe LED an Pin 26
	led := machine.Pin(26)
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// Taster an Pin 22
	button := machine.Pin(22)
	button.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// Ultraschallsensor an Pin 18
	ultrasonicPin := machine.Pin(18)

	ledState := false
	led.Low() // LED initial ausschalten

	previousButtonState := true
	
	var loopCounter uint32 = 0


	for {
		// Tasterlogik
		currentButtonState := button.Get()
		if previousButtonState && !currentButtonState {
			ledState = !ledState
			if ledState {
				led.High()
				println("LED AN")
			} else {
				led.Low()
				println("LED AUS")
			}
		}
		previousButtonState = currentButtonState

		if loopCounter%10 == 0 {
			dist := getDistanceCm(ultrasonicPin)
			if dist >= 0 {

				println("Entfernung:", int(dist), "cm")
			} else if dist == -1.0 {
				println("Sensor Fehler: Echo Start Timeout")
			} else if dist == -2.0 {
				println("Sensor Fehler: Echo Ende Timeout")
			} else if dist == -3.0 {
				println("Sensor Fehler: Ungültige Impulsdauer")
			} else {
				println("Sensor Fehler: Unbekannt")
			}
		}
		loopCounter++

		time.Sleep(time.Millisecond * 50) // Hauptschleifenverzögerung
	}
} 