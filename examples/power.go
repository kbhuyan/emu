package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kbhuyan/emu"
)

func main() {
	// Configure command-line flags
	port := flag.String("port", "/dev/ttyACM1", "Serial port device path")
	logLevel := flag.String("log", "LOG_WARNING", "Emu logging level (LOG_ALL, LOG_INFO, LOG_WARNING, LOG_ERROR, LOG_OFF)")
	flag.Parse()

	ll, err := emu.StringToLogLevel(*logLevel)
	if err != nil {
		log.Fatalf("Bad log level: %s\n", err)
	}

	// Connect to EMU-2 device
	device, err := emu.NewEmu(*port,
		emu.WithBaudRate(115200),
		emu.WithTimeOut(15*time.Second),
		emu.WithLoggingLevel(ll))
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer device.Close()

	// Start the device communication
	device.Start()

	// Process the power consumption data
	// if power, err := device.GetInstantaneousPowerConsumption(); err == nil {
	// 	log.Printf("Power Consumption: %v kW", power.Power)
	// }
	sub := []emu.MessageName{emu.InstantaneousPower, emu.CumulativeEnergy}
	handler := processMessage
	device.Subscribe(sub, &handler)

	waitingToBeTerminate(device)
}

func waitingToBeTerminate(device emu.Emu) {
	// Create a channel to receive signals.
	sigChan := make(chan os.Signal, 1)

	// Notify the channel of specific signals.
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Program is running. Press Ctrl+C to interrupt.")

	// Block until a signal is received.
	sig := <-sigChan

	// Handle the signal.
	switch sig {
	case syscall.SIGINT:
		fmt.Println("SIGINT received. Exiting...")
		device.Close()
		os.Exit(0)
	case syscall.SIGTERM:
		fmt.Println("SIGTERM received. Exiting...")
		device.Close()
		os.Exit(0)
	default:
		fmt.Println("Unexpected signal received.")
	}
}

func processMessage(msg emu.Message) {
	name := emu.MessageName(msg.GetName())
	switch name {
	case emu.TimeCluster:
		var local, utc int64
		if value, ok := msg.GetAttrib("LocalTime"); !ok {
			log.Printf("LocalTime not found %+v\n", msg)
			return
		} else {
			local = value.(int64)
		}
		if value, ok := msg.GetAttrib("UTCTime"); !ok {
			log.Printf("UTCTime not found %+v\n", msg)
			return
		} else {
			utc = value.(int64)
		}
		log.Printf("TimeCluster: LocalTime: %s UTCTime %s\n",
			time.Unix(local, 0).In(time.UTC).Format("2006-01-02 15:04:05"),
			time.Unix(utc, 0).In(time.UTC).Format("2006-01-02 15:04:05"))
	case emu.InstantaneousPower:
		if power, ok := msg.(*emu.InstantaneousPowerDemand); ok {
			log.Printf("TimeStamp: %s Instantaneous Demand: %.3fkW\n", time.Unix(power.TimeStamp, 0), power.Power)
		} else {
			log.Printf("invalid message: expecting emu.InstantaneousPowerDemand insted got %T. %+v", msg, msg)
		}
	case emu.CumulativeEnergy:
		if energy, ok := msg.(*emu.CumulativeEnergyConsumption); ok {
			log.Printf("TimeStamp: %s Current Cumulative Energy Delivered: %.3fkWh\n", time.Unix(energy.TimeStamp, 0), energy.Energy)
		} else {
			log.Printf("invalid message: expecting emu.CumulativeEnergyConsumption insted got %T. %+v", msg, msg)
		}
	default:
		log.Printf("Message: %+v\n", msg)
	}
}
