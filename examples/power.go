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
