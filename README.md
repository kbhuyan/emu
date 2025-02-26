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

    // Get instantaneous power demand
    cmd, _ := emu.NewCommand("get_instantaneous_demand")
    device.SendCommand(cmd)
    
    msg, err := device.GetResponse()
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    
    // Process the power consumption data
    if power, err := emu.GetInstantaneousPowerConsumption(msg); err == nil {
        log.Printf("Power Consumption: %v kW", power.Power)
    }
}
```

### Available Commands

- `get_device_info` - Device manufacturer and version info
- `get_network_info` - MAC addresses and connection status
- `get_current_summation` - Total energy consumption
- `get_instantaneous_demand` - Current power demand

### Configuration Options

```go
emu.NewEmu("/dev/ttyACM1",
    emu.WithBaudRate(115200),    // Default: 115200
    emu.WithTimeOut(15*time.Second), // Default: 15s
    emu.WithLogWriter(os.Stdout), // Default: os.Stdout
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
```

### CumulativeEnergyConsumption
```go
type CumulativeEnergyConsumption struct {
    TimeStamp   int64   // Unix timestamp
    Energy      float64 // Energy in kWh
    DeviceMacId string  // EMU-2 MAC address
    MeterMacId  string  // Smart meter MAC address
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
