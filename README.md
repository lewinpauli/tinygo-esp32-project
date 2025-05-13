# TinyGo ESP32 Project

## Setup

```bash
# Install required tools
brew install go
brew tap tinygo-org/tools
brew install tinygo

# Install Go modules/dependencies
go mod tidy
go get

# Update Go modules/dependencies
go get -u
```

## Building and Flashing in VS Code

Use keyboard shortcut: `CMD + SHIFT + B`

Configuration details are in `.vscode/tasks.json`

## Manual Building and Flashing

```bash
# IMPORTANT: Use the TinyGo compiler, not the standard Go compiler

# Build your project
tinygo build -target=esp32-coreboard-v2 -o output.bin main.go

# Flash to your ESP32 device
tinygo flash -target=esp32-coreboard-v2 -o output.bin main.go
```

## Device Ports

Common ESP32 device ports:
- macOS: `/dev/cu.SLAB_USBtoUART` or `/dev/cu.usbserial-*`
- Linux: `/dev/ttyUSB0`
- Windows: `COM3` (or another COM port number)

Find your device port on macOS:
```bash
ls /dev/cu.*
```