//go:build tinygo

package main

import (
	"machine"
	"time"
	// Import für serielle Ausgabe mit Formatierung, falls benötigt (siehe Hinweis unten)
	// "fmt"
)

const (
	// Pin-Definitionen
	ledPinGPIO            = machine.Pin(26)
	buttonPinGPIO         = machine.Pin(22)
	ultrasonicSensorPinGPIO = machine.Pin(18)

	// getDistanceCm Konstanten
	ultrasonicTimeoutMicroseconds = 25000 // 25 Millisekunden
	soundSpeedCmPerNs float32     = 0.0000343 // Explizit float32

	// Fehlercodes für getDistanceCm (explizit float32)
	errEchoStartTimeout float32 = -1.0
	errEchoEndTimeout   float32 = -2.0
	errInvalidPulse     float32 = -3.0
	
	// Hauptschleifen-Konstanten
	mainLoopDelay                   = 50 * time.Millisecond
	ultrasonicReadIntervalInCycles = 10 // Alle 10 Zyklen lesen (10 * 50ms = 500ms)
)

// getDistanceCm misst die Entfernung mit einem Single-Pin Ultraschallsensor.
// ultrasonicPin: Der GPIO-Pin, der mit dem SIG-Pin des Sensors verbunden ist.
// Gibt die Distanz in cm zurück, oder negative Werte bei Fehlern/Timeouts.
func getDistanceCm(ultrasonicPin machine.Pin) float32 {
	ultrasonicPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ultrasonicPin.Low()
	time.Sleep(2 * time.Microsecond)
	ultrasonicPin.High()
	time.Sleep(10 * time.Microsecond) // Trigger-Impuls für 10 Mikrosekunden
	ultrasonicPin.Low()

	ultrasonicPin.Configure(machine.PinConfig{Mode: machine.PinInput})

	var startTime, endTime int64
	// Timeout für den Start des Echo-Impulses (in Nanosekunden)
	timeoutStartThreshold := time.Now().UnixNano() + int64(ultrasonicTimeoutMicroseconds*1000)

	// Warten auf den Beginn des Echo-Impulses
	for !ultrasonicPin.Get() {
		if time.Now().UnixNano() > timeoutStartThreshold {
			return errEchoStartTimeout // Direkt float32 zurückgeben
		}
	}
	startTime = time.Now().UnixNano()

	// Timeout für das Ende des Echo-Impulses (in Nanosekunden)
	timeoutEndThreshold := startTime + int64(ultrasonicTimeoutMicroseconds*1000)

	// Warten auf das Ende des Echo-Impulses
	for ultrasonicPin.Get() {
		if time.Now().UnixNano() > timeoutEndThreshold {
			return errEchoEndTimeout // Direkt float32 zurückgeben
		}
	}
	endTime = time.Now().UnixNano()

	pulseDurationNano := endTime - startTime
	if pulseDurationNano <= 0 {
		return errInvalidPulse // Direkt float32 zurückgeben
	}

	// Geschwindigkeit des Schalls: ca. 34300 cm/s
	// Distanz (cm) = (Dauer_ns * Geschwindigkeit_cm_pro_ns) / 2
	distanceCm := float32(pulseDurationNano) * soundSpeedCmPerNs / 2.0

	return distanceCm
}

func main() {
	// Serielle Ausgabe konfigurieren (wichtig für println)
	uart := machine.UART0
	uart.Configure(machine.UARTConfig{BaudRate: 115200})
	// Kurze Verzögerung, damit der serielle Monitor Zeit hat, sich zu verbinden
	time.Sleep(2 * time.Second) 
	println("ESP32 Ultraschall- und Taster-Demo gestartet...")

	// Gelbe LED
	led := ledPinGPIO
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// Taster
	button := buttonPinGPIO
	button.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// Ultraschallsensor
	ultrasonicPin := ultrasonicSensorPinGPIO
	// Initialkonfiguration nicht hier, da sie in getDistanceCm() dynamisch wechselt

	ledState := false
	led.Low() // LED initial ausschalten

	previousButtonState := button.Get() // Sicherstellen, dass der initiale Zustand korrekt gelesen wird
	
	var loopCounter uint32 = 0


	for {
		// Tasterlogik
		currentButtonState := button.Get()
		if previousButtonState && !currentButtonState { // Flankenerkennung: von nicht gedrückt (true bei Pullup) zu gedrückt (false)
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

		// Ultraschallsensor auslesen
		if loopCounter%ultrasonicReadIntervalInCycles == 0 {
			dist := getDistanceCm(ultrasonicPin)
			if dist >= 0 {
				// Für genauere Float-Ausgabe "fmt" importieren und fmt.Printf verwenden:
				// fmt.Printf("Entfernung: %.2f cm\n", dist)
				// Für einfache Ausgabe mit println (als Integer):
				println("Entfernung:", int(dist), "cm")
			} else {
				// Fehlerbehandlung basierend auf den definierten Fehlercodes
				switch dist { // dist ist bereits float32, kein Cast nötig
				case errEchoStartTimeout:
					println("Sensor Fehler: Echo Start Timeout")
				case errEchoEndTimeout:
					println("Sensor Fehler: Echo Ende Timeout")
				case errInvalidPulse:
					println("Sensor Fehler: Ungültige Impulsdauer")
				default:
					println("Sensor Fehler: Unbekannt")
				}
			}
		}
		loopCounter++

		time.Sleep(mainLoopDelay) // Hauptschleifenverzögerung
	}
} 