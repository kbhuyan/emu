# EMU-2 Go Library

[![Go Reference](https://pkg.go.dev/badge/github.com/kbhuyan/emu.svg)](https://pkg.go.dev/github.com/kbhuyan/emu)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A pure Go library for interfacing with the EMU-2 (Energy Monitoring Unit) device through USB-serial communication. This library provides a clean API to query power consumption data from smart meters connected to the EMU-2.

## Features

- Serial port communication with EMU-2 device
- Real-time power consumption monitoring
- Support for multiple EMU-2 commands
- Thread-safe message handling
- Configurable logging
- Type-safe data structures for power and energy readings

## Installation

```bash
go get github.com/kbhuyan/emu
```

## Usage

### Basic Example

```go
package main

import (
    "log"
    "time"
    "github.com/kbhuyan/emu"
)

func main() {
    // Connect to EMU-2 device
    device, err := emu.NewEmu("/dev/ttyACM1",
        emu.WithBaudRate(115200),
        emu.WithTimeOut(15*time.Second))
    if err != nil {
        log.Fatalf("Connection failed: %v", err)
    }
    defer device.Close()
    
    // Start the device communication
    device.Start()
    
    // Process the power consumption data
    if power, err := device.GetInstantaneousPowerConsumption(); err == nil {
        log.Printf("Power Consumption: %v kW", power.Power)
    }
}
```

### Configuration Options

```go
emu.NewEmu("/dev/ttyACM1",
    emu.WithBaudRate(115200),    // Default: 115200
    emu.WithTimeOut(15*time.Second), // Default: 15s
    emu.WithLogWriter(os.Stdout), // Default: os.Stdout
    emu.WithLoggingLevel(emu.LOG_WARNING) //Default: emu.LOG_WARNING
)
```

## Data Structures

### InstantaneousPowerConsumption
```go
type InstantaneousPowerConsumption struct {
    TimeStamp   int64   // Unix timestamp
    Power       float64 // Power in kW
    DeviceMacId string  // EMU-2 MAC address
    MeterMacId  string  // Smart meter MAC address
}

if power, err := device.GetInstantaneousPowerConsumption(); err == nil {
    log.Printf("Power Consumption: %.3f kW", power.Power)
}
```

### CumulativeEnergyConsumption
```go
type CumulativeEnergyConsumption struct {
    TimeStamp   int64   // Unix timestamp
    Energy      float64 // Energy in kWh
    DeviceMacId string  // EMU-2 MAC address
    MeterMacId  string  // Smart meter MAC address
}
if energy, err := device.GetCumulativeEnergyConsumption(); err == nil {
    log.Printf("Current Cumulative Energy Delivered: %.3fkWh\n", energy.Energy)
}
```

### Available Commands

- `emu.RESTART`				- restarts the emu-2 device
- `emu.GET_DEVICE_INFO`		- gets the basic emu-2 device info HW/SW version, make/model etc.
- `emu.GET_TIME`			- gets the time (local and UTC) on the emu-2 as sync with the smart energy meter
- `emu.GET_CONN_STATUS`		- gets the current connection status with the smart energy meter`

```go
    if cmd, err := emu.NewCommand(emu.RESTART); err == nil {
        if err := device.SendCommand(cmd); err == nil {
            // Get response
            if rsp, err := device.GetResponse(); err == nil {
               fmt.Printf("Response %+v", rsp)
            }
	    }
    }
```

### Asyncronous Message Reception
```go
	sub := []emu.MessageName{emu.InstantaneousPower, emu.CumulativeEnergy}
	var handler = func(m emu.Message) {
		log.Printf("Received %+v", m)
	}
	device.Subscribe(sub, &handler)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
