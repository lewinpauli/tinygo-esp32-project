# TinyGo ESP32 Project

## Setup

```bash
# Install required tools
brew install go
brew tap tinygo-org/tools
brew install tinygo

# Install additional dependencies for ESP32
brew install esptool

# Install Go modules/dependencies
go mod tidy
go get

# Update Go modules/dependencies
go get -u


```

## Configuration

The project uses an `.env` file to store configuration variables:
- `TINYGO_TARGET`: Specifies the target board (currently set to `esp32-coreboard-v2`)

You can modify the `.env` file to change the target board without having to update multiple files.

## Building and Flashing in VS Code

Use keyboard shortcut: `CMD + SHIFT + B`

Configuration details are in `.vscode/tasks.json`

## Manual Building and Flashing

```bash
# IMPORTANT: Use the TinyGo compiler, not the standard Go compiler

# Load environment variables
source .env

# Build your project
tinygo build -target=$TINYGO_TARGET -o output.bin main.go

# Flash to your ESP32 device
tinygo flash -target=$TINYGO_TARGET main.go
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

Find your device port on Linux:
```bash
ls /dev/ttyUSB*
```

## Troubleshooting

If you encounter errors with undefined I2C pin constants like `SCL_PIN` and `SDA_PIN`, you may need to define them in your code:

```go
// Define global I2C pin constants needed by the machine package
var (
    SCL_PIN = machine.GPIO22
    SDA_PIN = machine.GPIO21
)
```

### Linux Port Permissions

On Linux, you may need to give your user permission to access the serial port:

```bash
# Replace ttyUSB0 with your actual port
sudo chmod 777 /dev/ttyUSB0
```

Alternatively, add your user to the dialout group for permanent access:

```bash
sudo usermod -a -G dialout $USER
# Log out and back in for this to take effect
```