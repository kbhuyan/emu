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

func sigHandler(device emu.Emu) {
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
func main() {
	// Configure command-line flags
	port := flag.String("port", "/dev/ttyACM1", "Serial port device path")
	baud := flag.Int("baud", 115200, "Baud rate (115200, 9600, etc)")
	timeout := flag.Duration("timeout", 15*time.Second, "Read timeout duration")
	logLevel := flag.String("log", "LOG_WARNING", "Emu logging level (LOG_ALL, LOG_INFO, LOG_WARNING, LOG_ERROR, LOG_OFF)")
	list := flag.Bool("list", false, "List available commands and exit")

	flag.Parse()

	// Handle --list flag
	if *list {
		printAvailableCommands()
		return
	}

	// Get command from positional arguments
	args := flag.Args()
	if len(args) < 1 {
		log.Fatalf("Usage: emuctl [flags] <command>\n%s", cmdList)
	}

	ll, err := emu.StringToLogLevel(*logLevel)
	if err != nil {
		log.Fatalf("Bad log level: %s\n", err)
	}
	cmdStr := args[0]
	command, err := emu.StrToCommandId(cmdStr)
	if err != nil {
		log.Fatalf("Bad command: %v\n %s\n", err, cmdList)
	}

	device, err := emu.NewEmu(*port, emu.WithBaudRate(*baud), emu.WithTimeOut(*timeout), emu.WithLoggingLevel(ll))
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer device.Close()
	device.Start()
	go func() {
		//fmt.Println("commanf: ", command)
		// Execute command
		cmd, err := emu.NewCommand(command)
		if err != nil {
			log.Fatalf("%v", err)
		}

		if msg, err := executeCommand(device, cmd); err != nil {
			log.Fatalf("%v", err)
		} else {
			processMessage(msg)
		}

		if energy, err := device.GetCumulativeEnergyConsumption(); err != nil {
			log.Fatalf("%v", err)
		} else {
			log.Printf("TimeStamp: %s Current Cumulative Energy Delivered: %.3fkWh\n", time.Unix(energy.TimeStamp, 0), energy.Energy)
		}

		if power, err := device.GetInstantaneousPowerConsumption(); err != nil {
			log.Fatalf("%v", err)
		} else {
			log.Printf("TimeStamp: %s Instantaneous Demand: %.3fkW\n", time.Unix(power.TimeStamp, 0), power.Power)
		}

		sub := []emu.MessageName{emu.TimeCluster, emu.CurrentSummationDelivered, emu.InstantaneousDemand}
		for {
			msg, err := device.GetMessage(sub)
			if err != nil {
				log.Fatalf("GetMessage failed: %v", err)
			}
			processMessage(msg)
		}
	}()

	sigHandler(device)
}

func processMessage(msg emu.Message) {
	if msg.GetName() == "InstantaneousDemand" {

		if power, err := emu.GetInstantaneousPowerConsumption(msg); err != nil {
			log.Println(err)
			return
		} else {
			log.Printf("TimeStamp: %s Instantaneous Demand: %.3fkW\n", time.Unix(power.TimeStamp, 0), power.Power)
		}

	} else if msg.GetName() == "CurrentSummationDelivered" {

		if energy, err := emu.GetCumulativeEnergyConsumption(msg); err != nil {
			log.Println(err)
			return
		} else {
			log.Printf("TimeStamp: %s Current Cumulative Energy Delivered: %.3fkWh\n", time.Unix(energy.TimeStamp, 0), energy.Energy)
		}
	} else if msg.GetName() == "TimeCluster" {
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
	} else {
		log.Printf("Message: %+v\n", msg)
	}
}

func executeCommand(device emu.Emu, cmd emu.Command) (emu.Message, error) {
	// Send command
	if err := device.SendCommand(cmd); err != nil {
		return nil, fmt.Errorf("command failed: %v", err)
	}

	// Get response
	//rawResp, err := device.ReadResponseSync()
	rsp, err := device.GetResponse()
	if err != nil {
		return nil, fmt.Errorf("response error: %v", err)
	}
	return rsp, nil
}

var cmdList = `Available EMU commands:
	RESTART				- restarts the emu-2 device
	GET_DEVICE_INFO		- gets the basic emu-2 device info HW/SW version, make/model etc.
	GET_TIME			- gets the time (local and UTC) on the emu-2 as sync with the smart energy meter
	GET_CONN_STATUS		- gets the current connection status with the smart energy meter`

func printAvailableCommands() {
	fmt.Println(cmdList)
}
